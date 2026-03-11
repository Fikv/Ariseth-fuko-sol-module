package domain

import "time"

type TermFrequency struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type DataAnalysisRequest struct {
	Source  string `json:"source,omitempty"`
	Content string `json:"content"`
}

type DataAnalysisResult struct {
	Source      string          `json:"source,omitempty"`
	Bytes       int             `json:"bytes"`
	Lines       int             `json:"lines"`
	Words       int             `json:"words"`
	UniqueWords int             `json:"uniqueWords"`
	Numbers     int             `json:"numbers"`
	TopTerms    []TermFrequency `json:"topTerms"`
	AnalyzedAt  time.Time       `json:"analyzedAt"`
}
