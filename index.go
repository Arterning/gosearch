package main

import "sync"

// Index is an improved inverted index with CRUD support
type Index struct {
	mu        sync.RWMutex
	index     map[string][]string // token -> []docID
	docCount  int
}

// NewIndex creates a new inverted index
func NewIndex() *Index {
	return &Index{
		index: make(map[string][]string),
	}
}

// AddDocument adds or updates a document in the index
func (idx *Index) AddDocument(docID string, tokens []string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Get unique tokens
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token] = true
	}

	// Add to index
	for token := range uniqueTokens {
		docList := idx.index[token]

		// Check if doc already exists
		found := false
		for _, id := range docList {
			if id == docID {
				found = true
				break
			}
		}

		if !found {
			idx.index[token] = append(docList, docID)
		}
	}

	idx.docCount++
}

// RemoveDocument removes a document from the index
func (idx *Index) RemoveDocument(docID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	for token, docList := range idx.index {
		newList := make([]string, 0, len(docList))
		for _, id := range docList {
			if id != docID {
				newList = append(newList, id)
			}
		}

		if len(newList) == 0 {
			delete(idx.index, token)
		} else {
			idx.index[token] = newList
		}
	}

	idx.docCount--
	if idx.docCount < 0 {
		idx.docCount = 0
	}
}

// UpdateDocument updates a document (remove old, add new)
func (idx *Index) UpdateDocument(docID string, tokens []string) {
	idx.RemoveDocument(docID)
	idx.AddDocument(docID, tokens)
}

// SearchAND returns documents containing ALL tokens
func (idx *Index) SearchAND(tokens []string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if len(tokens) == 0 {
		return nil
	}

	// Start with first token's documents
	result := make(map[string]bool)
	firstDocs, ok := idx.index[tokens[0]]
	if !ok {
		return nil
	}

	for _, docID := range firstDocs {
		result[docID] = true
	}

	// Intersect with remaining tokens
	for i := 1; i < len(tokens); i++ {
		docs, ok := idx.index[tokens[i]]
		if !ok {
			return nil
		}

		// Create new result with intersection
		newResult := make(map[string]bool)
		for _, docID := range docs {
			if result[docID] {
				newResult[docID] = true
			}
		}
		result = newResult

		if len(result) == 0 {
			return nil
		}
	}

	// Convert to slice
	resultSlice := make([]string, 0, len(result))
	for docID := range result {
		resultSlice = append(resultSlice, docID)
	}

	return resultSlice
}

// SearchOR returns documents containing ANY token
func (idx *Index) SearchOR(tokens []string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	result := make(map[string]bool)

	for _, token := range tokens {
		if docs, ok := idx.index[token]; ok {
			for _, docID := range docs {
				result[docID] = true
			}
		}
	}

	// Convert to slice
	resultSlice := make([]string, 0, len(result))
	for docID := range result {
		resultSlice = append(resultSlice, docID)
	}

	return resultSlice
}

// DocFrequency returns number of documents containing the token
func (idx *Index) DocFrequency(token string) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if docs, ok := idx.index[token]; ok {
		return len(docs)
	}
	return 0
}

// TotalDocuments returns total number of indexed documents
func (idx *Index) TotalDocuments() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.docCount
}

// Stats returns index statistics
func (idx *Index) Stats() IndexStats {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	totalDocs := 0
	for _, docs := range idx.index {
		totalDocs += len(docs)
	}

	avgDocsPerToken := 0.0
	if len(idx.index) > 0 {
		avgDocsPerToken = float64(totalDocs) / float64(len(idx.index))
	}

	return IndexStats{
		TotalDocuments:  idx.docCount,
		TotalTokens:     len(idx.index),
		AvgDocsPerToken: avgDocsPerToken,
	}
}

// IndexStats contains index statistics
type IndexStats struct {
	TotalDocuments  int     `json:"total_documents"`
	TotalTokens     int     `json:"total_tokens"`
	AvgDocsPerToken float64 `json:"avg_docs_per_token"`
}
