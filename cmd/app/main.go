package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ariseth-fuko-sol-module/internal/scanner"
	"ariseth-fuko-sol-module/internal/sources/mock"
	"ariseth-fuko-sol-module/internal/ui"
)

type menuItem struct {
	label  string
	action func() error
}

func runMenu(title string, items []menuItem) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("==", title, "==")
		for i, item := range items {
			fmt.Printf("%d) %s\n", i+1, item.label)
		}
		fmt.Println("0) Exit")
		fmt.Print("Choose option: ")

		raw, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read input: %w", err)
		}

		choice := strings.TrimSpace(raw)
		if choice == "0" {
			fmt.Println("Bye.")
			return nil
		}

		handled := false
		for i, item := range items {
			if choice == fmt.Sprintf("%d", i+1) {
				handled = true
				if err := item.action(); err != nil {
					fmt.Println("Error:", err)
				}
				break
			}
		}

		if !handled {
			fmt.Println("Invalid option. Try again.")
		}
	}
}

func main() {
	items := []menuItem{
		{
			label: "Run token scanner + web UI",
			action: func() error {
				return runScannerUI()
			},
		},
		{
			label: "Show scanner architecture notes",
			action: func() error {
				fmt.Println("Source -> dedup store -> enrichment -> UI/WebSocket/SSE")
				fmt.Println("Current source: mock stream (replace with Raydium/Orca/Meteora watchers)")
				fmt.Println("UI endpoint: http://localhost:8080")
				return nil
			},
		},
	}

	if err := runMenu("Ariseth CLI", items); err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(1)
	}
}

func runScannerUI() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	store := scanner.NewStore()
	logger := log.New(os.Stdout, "[scanner] ", log.LstdFlags)
	service := scanner.NewService(
		store,
		logger,
		mock.New("solana", "raydium"),
	)

	uiServer := ui.NewServer(store)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: uiServer.Handler(),
	}

	errCh := make(chan error, 2)

	go func() {
		if err := service.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	fmt.Println("Scanner is running.")
	fmt.Println("Open UI: http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop and return to menu.")

	select {
	case <-sigCh:
		cancel()
	case err := <-errCh:
		cancel()
		return err
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	_ = httpServer.Shutdown(shutdownCtx)
	return nil
}
