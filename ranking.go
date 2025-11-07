package main

import (
	"math"
	"sort"
)

// BM25 implements the BM25 ranking algorithm
type BM25 struct {
	K1 float64 // Term frequency saturation parameter
	B  float64 // Length normalization parameter
}

// NewBM25 creates a new BM25 ranker with default parameters
func NewBM25() *BM25 {
	return &BM25{
		K1: 1.5,
		B:  0.75,
	}
}

// Score calculates BM25 score for a document
func (bm25 *BM25) Score(queryTerms []string, docStats *DocStats, idx *Index, avgDocLength float64) float64 {
	score := 0.0
	docLength := float64(docStats.Length)
	totalDocs := float64(idx.TotalDocuments())

	for _, term := range queryTerms {
		// Get term frequency in document
		tf := float64(docStats.TermFrequencies[term])
		if tf == 0 {
			continue
		}

		// Calculate IDF (Inverse Document Frequency)
		docFreq := float64(idx.DocFrequency(term))
		idf := 0.0
		if docFreq > 0 {
			idf = math.Log((totalDocs-docFreq+0.5)/(docFreq+0.5) + 1.0)
		}

		// Calculate BM25 score component
		normalizedTF := (tf * (bm25.K1 + 1.0)) / (tf + bm25.K1*(1.0-bm25.B+bm25.B*(docLength/avgDocLength)))

		score += idf * normalizedTF
	}

	return score
}

// ScoredDocument represents a document with its relevance score
type ScoredDocument struct {
	DocID string
	Score float64
}

// RankDocuments ranks documents using BM25
func RankDocuments(queryTerms []string, candidateDocs []string, docStatsMap map[string]*DocStats, idx *Index, avgDocLength float64) []ScoredDocument {
	bm25 := NewBM25()
	scored := make([]ScoredDocument, 0, len(candidateDocs))

	for _, docID := range candidateDocs {
		if stats, ok := docStatsMap[docID]; ok {
			score := bm25.Score(queryTerms, stats, idx, avgDocLength)
			scored = append(scored, ScoredDocument{
				DocID: docID,
				Score: score,
			})
		}
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}
