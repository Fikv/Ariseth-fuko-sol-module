package service

import (
	"context"
	"errors"
	"regexp"
	"sort"
	"strings"
	"time"

	"ariseth-fuko-sol-module/internal/domain"
)

var (
	wordRegex   = regexp.MustCompile(`[A-Za-z0-9_]+`)
	numberRegex = regexp.MustCompile(`\b\d+(?:[.,]\d+)?\b`)
)

type DataAnalysisService domain.DataAnalysisService

func NewDataAnalysisService() *DataAnalysisService {
	return &DataAnalysisService{}
}

func (s *DataAnalysisService) Analyze(ctx context.Context, req domain.DataAnalysisRequest) (domain.DataAnalysisResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.DataAnalysisResult{}, err
	}
	if strings.TrimSpace(req.Content) == "" {
		return domain.DataAnalysisResult{}, errors.New("content is required")
	}

	content := req.Content
	lines := 1 + strings.Count(content, "\n")
	if content == "" {
		lines = 0
	}

	wordMatches := wordRegex.FindAllString(content, -1)
	numberMatches := numberRegex.FindAllString(content, -1)
	termCount := make(map[string]int, len(wordMatches))

	for _, word := range wordMatches {
		if err := ctx.Err(); err != nil {
			return domain.DataAnalysisResult{}, err
		}

		normalized := strings.ToLower(word)
		if len(normalized) < 3 {
			continue
		}
		termCount[normalized]++
	}

	topTerms := make([]domain.TermFrequency, 0, len(termCount))
	for term, count := range termCount {
		topTerms = append(topTerms, domain.TermFrequency{
			Term:  term,
			Count: count,
		})
	}

	sort.Slice(topTerms, func(i, j int) bool {
		if topTerms[i].Count == topTerms[j].Count {
			return topTerms[i].Term < topTerms[j].Term
		}
		return topTerms[i].Count > topTerms[j].Count
	})

	if len(topTerms) > 10 {
		topTerms = topTerms[:10]
	}

	return domain.DataAnalysisResult{
		Source:      req.Source,
		Bytes:       len(content),
		Lines:       lines,
		Words:       len(wordMatches),
		UniqueWords: len(termCount),
		Numbers:     len(numberMatches),
		TopTerms:    topTerms,
		AnalyzedAt:  time.Now().UTC(),
	}, nil
}
