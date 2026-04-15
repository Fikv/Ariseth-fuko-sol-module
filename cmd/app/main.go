package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"ariseth-fuko-sol-module/internal/domain"
	"ariseth-fuko-sol-module/internal/execution"
	"ariseth-fuko-sol-module/internal/simulator"
	"ariseth-fuko-sol-module/internal/solana"
	"ariseth-fuko-sol-module/internal/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	svc := simulator.NewService()
	solanaClient := solana.NewClient("")
	execClient := execution.NewClient()
	if execClient.UsingDefaultRPC() {
		log.Printf("warning: using default public Solana RPC %s. Auto-copy can hit rate limits. Set SOLANA_RPC_URL in .env for a private RPC.", execClient.RPCURL())
	} else {
		log.Printf("using Solana RPC %s", execClient.RPCURL())
	}
	go svc.Start(ctx)

	buyQueue := make(chan domain.AutoCopyOrder, autoCopyQueueSize())
	go startAutoBuyWorker(ctx, svc, execClient, buyQueue)

	go solana.StartWalletWatcher(ctx, solanaClient, svc.WatchedWallets, func(trades []domain.RealtimeTrade) {
		orders := svc.IngestRealtimeTrades(trades)
		if len(orders) == 0 || !execClient.BotWalletEnabled() {
			return
		}

		for _, order := range orders {
			select {
			case buyQueue <- order:
			default:
				err := errors.New("auto-buy queue full; throttling live buys")
				log.Printf(
					"ME BUY FAILED sourceWallet=%s symbol=%s mint=%s amountSOL=%.4f error=%v",
					order.WalletName,
					order.TokenSymbol,
					order.TokenAddress,
					order.AmountSOL,
					err,
				)
				svc.RecordAutoCopyExecutionFailure(order, err)
			}
		}
	})
	go startPendingSellWorker(ctx, svc, execClient)
	go startStrandedHoldingReconciler(ctx, svc, execClient)

	server := &http.Server{
		Addr:              serverAddr(),
		Handler:           web.NewHandler(svc, execClient),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown error: %v", err)
		}
	}()

	log.Printf("solana copy-trader simulator listening on http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			log.Printf("server could not start on %s because the port is already in use. Stop the other running app instance first.", server.Addr)
		}
		log.Fatal(err)
	}
}

func serverAddr() string {
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}

func startAutoBuyWorker(ctx context.Context, svc *simulator.Service, execClient *execution.Client, queue <-chan domain.AutoCopyOrder) {
	minInterval := autoCopyMinInterval()

	for {
		select {
		case <-ctx.Done():
			return
		case order := <-queue:
			if !execClient.BotWalletEnabled() {
				continue
			}
			settings := svc.ExecutionSettings()
			requestCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
			result, err := execClient.ExecuteBotSwap(requestCtx, execution.PrepareSwapRequest{
				OutputMint:          order.TokenAddress,
				OutputSymbol:        order.TokenSymbol,
				AmountSOL:           order.AmountSOL,
				AmountLamports:      int64(order.AmountSOL * 1_000_000_000),
				SlippageBps:         order.SlippageBps,
				MinWalletReserveSOL: settings.MinWalletReserveSOL,
			})
			cancel()
			if err != nil {
				log.Printf(
					"ME BUY FAILED sourceWallet=%s symbol=%s mint=%s amountSOL=%.4f error=%v",
					order.WalletName,
					order.TokenSymbol,
					order.TokenAddress,
					order.AmountSOL,
					err,
				)
				svc.RecordAutoCopyExecutionFailure(order, err)
			} else {
				log.Printf(
					"ME BUY SUBMITTED sourceWallet=%s symbol=%s mint=%s amountSOL=%.4f backend=%s signatures=%s",
					order.WalletName,
					order.TokenSymbol,
					order.TokenAddress,
					order.AmountSOL,
					result.Backend,
					strings.Join(result.Signatures, ","),
				)
				svc.RecordAutoCopyExecutionSuccess(order, result.Backend, result.Signatures)
			}
			if minInterval > 0 {
				_ = sleepWithContext(ctx, minInterval)
			}
		}
	}
}

func startPendingSellWorker(ctx context.Context, svc *simulator.Service, execClient *execution.Client) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	maxPerTick := autoSellMaxPerTick()
	minInterval := autoSellMinInterval()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		orders := svc.ConsumePendingSells()
		if len(orders) == 0 || !execClient.BotWalletEnabled() {
			continue
		}
		if maxPerTick > 0 && len(orders) > maxPerTick {
			for _, overflow := range orders[maxPerTick:] {
				svc.RequeuePendingSell(overflow, 2*time.Second)
			}
			orders = orders[:maxPerTick]
		}

		for _, order := range orders {
			slippage := order.SlippageBps
			if order.RetryCount > 0 {
				slippage = minInt(order.SlippageBps+order.RetryCount*200, 2000)
			}
			requestCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
			result, err := execClient.ExecuteBotSellMint(requestCtx, order.TokenAddress, slippage)
			cancel()
			if err != nil {
				if shouldTreatSellAsSuccess(err) {
					log.Printf(
						"ME SELL CONFIRMED sourceWallet=%s symbol=%s mint=%s reason=%s slippage=%d note=%v",
						order.WalletName,
						order.TokenSymbol,
						order.TokenAddress,
						order.Reason,
						slippage,
						err,
					)
					svc.RecordAutoSellExecutionSuccess(order, "pumpportal", nil)
				} else {
					log.Printf(
						"ME SELL FAILED sourceWallet=%s symbol=%s mint=%s reason=%s slippage=%d error=%v",
						order.WalletName,
						order.TokenSymbol,
						order.TokenAddress,
						order.Reason,
						slippage,
						err,
					)
					svc.RecordAutoSellExecutionFailure(order, err)
					if order.RetryCount < 8 {
						delay := time.Duration(order.RetryCount+1) * 15 * time.Second
						svc.RequeuePendingSell(order, delay)
					}
				}
			} else {
				log.Printf(
					"ME SELL SUBMITTED sourceWallet=%s symbol=%s mint=%s reason=%s slippage=%d backend=%s signatures=%s",
					order.WalletName,
					order.TokenSymbol,
					order.TokenAddress,
					order.Reason,
					slippage,
					result.Backend,
					strings.Join(result.Signatures, ","),
				)
				svc.RecordAutoSellExecutionSuccess(order, result.Backend, result.Signatures)
			}
			if minInterval > 0 {
				_ = sleepWithContext(ctx, minInterval)
			}
		}
	}
}

func startStrandedHoldingReconciler(ctx context.Context, svc *simulator.Service, execClient *execution.Client) {
	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	retryAfter := make(map[string]time.Time)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		if !execClient.BotWalletEnabled() {
			continue
		}

		requestCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		holdings, err := execClient.ListBotTokenHoldings(requestCtx)
		cancel()
		if err != nil {
			log.Printf("ME RECONCILE FAILED stage=list-holdings error=%v", err)
			continue
		}

		trackedMints := svc.ActiveTrackedMints()
		now := time.Now()
		for _, holding := range holdings {
			mintKey := strings.ToLower(strings.TrimSpace(holding.Mint))
			if mintKey == "" || execution.IsProtectedMint(holding.Mint) || holding.UIAmount <= 0 {
				continue
			}
			if _, tracked := trackedMints[mintKey]; tracked {
				continue
			}
			if until, waiting := retryAfter[mintKey]; waiting && until.After(now) {
				continue
			}

			requestCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
			result, err := execClient.ExecuteBotSellMint(requestCtx, holding.Mint, 1400)
			cancel()
			if err != nil {
				retryAfter[mintKey] = now.Add(1 * time.Minute)
				log.Printf(
					"ME STRAY SELL FAILED wallet=%s symbol=%s mint=%s uiAmount=%.6f error=%v",
					execClient.Config().BotWalletAddress,
					holding.SymbolHint,
					holding.Mint,
					holding.UIAmount,
					err,
				)
				continue
			}

			delete(retryAfter, mintKey)
			signatures := append([]string(nil), result.Signatures...)
			slices.Sort(signatures)
			log.Printf(
				"ME STRAY SELL SUBMITTED wallet=%s symbol=%s mint=%s uiAmount=%.6f backend=%s signatures=%s",
				result.Wallet,
				holding.SymbolHint,
				holding.Mint,
				holding.UIAmount,
				result.Backend,
				strings.Join(signatures, ","),
			)
		}
	}
}

func autoCopyQueueSize() int {
	return readEnvIntDefault("AUTOCOPY_BUY_QUEUE_SIZE", 200)
}

func autoCopyMinInterval() time.Duration {
	ms := readEnvIntDefault("AUTOCOPY_MIN_INTERVAL_MS", 750)
	if ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

func autoSellMinInterval() time.Duration {
	ms := readEnvIntDefault("AUTOSELL_MIN_INTERVAL_MS", 600)
	if ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}

func autoSellMaxPerTick() int {
	return readEnvIntDefault("AUTOSELL_MAX_PER_TICK", 4)
}

func readEnvIntDefault(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func shouldTreatSellAsSuccess(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "already been processed") ||
		strings.Contains(message, "account not found") ||
		strings.Contains(message, "insufficient funds") ||
		strings.Contains(message, "no prior credit") ||
		strings.Contains(message, "no token account")
}
