package main

import (
	"fmt"
	"sync"
)

// SearchMode defines how to combine query terms
type SearchMode int

const (
	SearchModeAND SearchMode = iota
	SearchModeOR
)

// SearchOptions contains search parameters
type SearchOptions struct {
	Mode       SearchMode
	UseRanking bool
	Limit      int
	Offset     int
}

// DefaultSearchOptions returns default search options
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Mode:       SearchModeAND,
		UseRanking: true,
		Limit:      10,
		Offset:     0,
	}
}

// SearchResult contains search results
type SearchResult struct {
	Documents []*Document
	Total     int
	Scores    []float64
}

// SearchEngine is the main search engine
type SearchEngine struct {
	storage       *Storage
	index         *Index
	docStats      map[string]*DocStats
	avgDocLength  float64
	mu            sync.RWMutex
}

// NewSearchEngine creates a new search engine
func NewSearchEngine(storagePath string) (*SearchEngine, error) {
	storage, err := NewStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	// Load or create index
	index, err := storage.LoadIndex()
	if err != nil {
		return nil, fmt.Errorf("failed to load index: %w", err)
	}

	// Load document statistics
	docStatsMap, err := storage.GetAllDocStats()
	if err != nil {
		return nil, fmt.Errorf("failed to load doc stats: %w", err)
	}

	// Calculate average document length
	avgDocLength := 0.0
	if len(docStatsMap) > 0 {
		totalLength := 0
		for _, stats := range docStatsMap {
			totalLength += stats.Length
		}
		avgDocLength = float64(totalLength) / float64(len(docStatsMap))
	}

	return &SearchEngine{
		storage:      storage,
		index:        index,
		docStats:     docStatsMap,
		avgDocLength: avgDocLength,
	}, nil
}

// Close closes the search engine
func (e *SearchEngine) Close() error {
	return e.storage.Close()
}

// UpsertDocument inserts or updates a document
func (e *SearchEngine) UpsertDocument(doc *Document) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Analyze document text
	tokens := analyze(doc.SearchableText())

	// Calculate term frequencies
	termFreqs := make(map[string]int)
	for _, token := range tokens {
		termFreqs[token]++
	}

	// Create document statistics
	docStats := &DocStats{
		ID:              doc.ID,
		Length:          len(tokens),
		TermFrequencies: termFreqs,
	}

	// Update index
	e.index.UpdateDocument(doc.ID, tokens)

	// Update doc stats
	e.docStats[doc.ID] = docStats

	// Recalculate average document length
	e.recalculateAvgLength()

	// Save to storage
	if err := e.storage.SaveDocument(doc); err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	if err := e.storage.SaveDocStats(docStats); err != nil {
		return fmt.Errorf("failed to save doc stats: %w", err)
	}

	if err := e.storage.SaveIndex(e.index); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}

// DeleteDocument deletes a document
func (e *SearchEngine) DeleteDocument(docID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Remove from index
	e.index.RemoveDocument(docID)

	// Remove from doc stats
	delete(e.docStats, docID)

	// Recalculate average document length
	e.recalculateAvgLength()

	// Remove from storage
	if err := e.storage.DeleteDocument(docID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if err := e.storage.DeleteDocStats(docID); err != nil {
		return fmt.Errorf("failed to delete doc stats: %w", err)
	}

	if err := e.storage.SaveIndex(e.index); err != nil {
		return fmt.Errorf("failed to save index: %w", err)
	}

	return nil
}

// GetDocument retrieves a document by ID
func (e *SearchEngine) GetDocument(docID string) (*Document, error) {
	return e.storage.GetDocument(docID)
}

// Search searches for documents
func (e *SearchEngine) Search(query string, options SearchOptions) (*SearchResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Tokenize query
	queryTokens := analyze(query)
	if len(queryTokens) == 0 {
		return &SearchResult{
			Documents: []*Document{},
			Total:     0,
		}, nil
	}

	// Find matching documents
	var candidateIDs []string
	switch options.Mode {
	case SearchModeAND:
		candidateIDs = e.index.SearchAND(queryTokens)
	case SearchModeOR:
		candidateIDs = e.index.SearchOR(queryTokens)
	default:
		candidateIDs = e.index.SearchAND(queryTokens)
	}

	total := len(candidateIDs)

	// Rank documents if requested
	var sortedIDs []string
	var scores []float64

	if options.UseRanking && total > 0 {
		scoredDocs := RankDocuments(queryTokens, candidateIDs, e.docStats, e.index, e.avgDocLength)

		sortedIDs = make([]string, len(scoredDocs))
		scores = make([]float64, len(scoredDocs))

		for i, sd := range scoredDocs {
			sortedIDs[i] = sd.DocID
			scores[i] = sd.Score
		}
	} else {
		sortedIDs = candidateIDs
	}

	// Apply pagination
	start := options.Offset
	if start > len(sortedIDs) {
		start = len(sortedIDs)
	}

	end := start + options.Limit
	if end > len(sortedIDs) {
		end = len(sortedIDs)
	}

	pageIDs := sortedIDs[start:end]
	var pageScores []float64
	if len(scores) > 0 {
		pageScores = scores[start:end]
	}

	// Fetch documents
	documents := make([]*Document, 0, len(pageIDs))
	for _, docID := range pageIDs {
		doc, err := e.storage.GetDocument(docID)
		if err != nil {
			continue
		}
		if doc != nil {
			documents = append(documents, doc)
		}
	}

	return &SearchResult{
		Documents: documents,
		Total:     total,
		Scores:    pageScores,
	}, nil
}

// Stats returns index statistics
func (e *SearchEngine) Stats() IndexStats {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.index.Stats()
}

// recalculateAvgLength recalculates average document length
func (e *SearchEngine) recalculateAvgLength() {
	if len(e.docStats) == 0 {
		e.avgDocLength = 0
		return
	}

	totalLength := 0
	for _, stats := range e.docStats {
		totalLength += stats.Length
	}
	e.avgDocLength = float64(totalLength) / float64(len(e.docStats))
}
