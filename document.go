package main

// Document represents a searchable document
type Document struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	URL      string            `json:"url,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewDocument creates a new document
func NewDocument(id, title, content string) *Document {
	return &Document{
		ID:       id,
		Title:    title,
		Content:  content,
		Metadata: make(map[string]string),
	}
}

// SearchableText returns the full text for indexing
func (d *Document) SearchableText() string {
	return d.Title + " " + d.Content
}

// DocStats stores document statistics for BM25 ranking
type DocStats struct {
	ID               string         `json:"id"`
	Length           int            `json:"length"`
	TermFrequencies  map[string]int `json:"term_frequencies"`
}

// NewDocStats creates new document statistics
func NewDocStats(id string, length int) *DocStats {
	return &DocStats{
		ID:              id,
		Length:          length,
		TermFrequencies: make(map[string]int),
	}
}
