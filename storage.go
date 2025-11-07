package main

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	docsBucket  = []byte("documents")
	statsBucket = []byte("doc_stats")
	indexBucket = []byte("index")
	metaBucket  = []byte("metadata")
)

// Storage handles persistent storage using BoltDB
type Storage struct {
	db *bolt.DB
}

// NewStorage creates a new storage instance
func NewStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create buckets
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{docsBucket, statsBucket, indexBucket, metaBucket}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create buckets: %w", err)
	}

	return &Storage{db: db}, nil
}

// Close closes the database
func (s *Storage) Close() error {
	return s.db.Close()
}

// SaveDocument saves a document
func (s *Storage) SaveDocument(doc *Document) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(docsBucket)
		data, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		return b.Put([]byte(doc.ID), data)
	})
}

// GetDocument retrieves a document by ID
func (s *Storage) GetDocument(id string) (*Document, error) {
	var doc *Document

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(docsBucket)
		data := b.Get([]byte(id))
		if data == nil {
			return nil
		}

		doc = &Document{}
		return json.Unmarshal(data, doc)
	})

	return doc, err
}

// DeleteDocument deletes a document
func (s *Storage) DeleteDocument(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(docsBucket)
		return b.Delete([]byte(id))
	})
}

// GetAllDocuments retrieves all documents
func (s *Storage) GetAllDocuments() ([]*Document, error) {
	var docs []*Document

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(docsBucket)
		return b.ForEach(func(k, v []byte) error {
			var doc Document
			if err := json.Unmarshal(v, &doc); err != nil {
				return err
			}
			docs = append(docs, &doc)
			return nil
		})
	})

	return docs, err
}

// CountDocuments returns total number of documents
func (s *Storage) CountDocuments() (int, error) {
	var count int

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(docsBucket)
		count = b.Stats().KeyN
		return nil
	})

	return count, err
}

// SaveDocStats saves document statistics
func (s *Storage) SaveDocStats(stats *DocStats) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(statsBucket)
		data, err := json.Marshal(stats)
		if err != nil {
			return err
		}
		return b.Put([]byte(stats.ID), data)
	})
}

// GetDocStats retrieves document statistics
func (s *Storage) GetDocStats(id string) (*DocStats, error) {
	var stats *DocStats

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(statsBucket)
		data := b.Get([]byte(id))
		if data == nil {
			return nil
		}

		stats = &DocStats{}
		return json.Unmarshal(data, stats)
	})

	return stats, err
}

// GetAllDocStats retrieves all document statistics
func (s *Storage) GetAllDocStats() (map[string]*DocStats, error) {
	statsMap := make(map[string]*DocStats)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(statsBucket)
		return b.ForEach(func(k, v []byte) error {
			var stats DocStats
			if err := json.Unmarshal(v, &stats); err != nil {
				return err
			}
			statsMap[stats.ID] = &stats
			return nil
		})
	})

	return statsMap, err
}

// DeleteDocStats deletes document statistics
func (s *Storage) DeleteDocStats(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(statsBucket)
		return b.Delete([]byte(id))
	})
}

// SaveIndex saves the inverted index
func (s *Storage) SaveIndex(idx *Index) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(indexBucket)

		// Serialize index data
		data, err := json.Marshal(struct {
			Index    map[string][]string `json:"index"`
			DocCount int                 `json:"doc_count"`
		}{
			Index:    idx.index,
			DocCount: idx.docCount,
		})

		if err != nil {
			return err
		}

		return b.Put([]byte("main_index"), data)
	})
}

// LoadIndex loads the inverted index
func (s *Storage) LoadIndex() (*Index, error) {
	idx := NewIndex()

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(indexBucket)
		data := b.Get([]byte("main_index"))
		if data == nil {
			return nil
		}

		var indexData struct {
			Index    map[string][]string `json:"index"`
			DocCount int                 `json:"doc_count"`
		}

		if err := json.Unmarshal(data, &indexData); err != nil {
			return err
		}

		idx.index = indexData.Index
		idx.docCount = indexData.DocCount
		return nil
	})

	return idx, err
}

// SaveMetadata saves metadata
func (s *Storage) SaveMetadata(key, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(metaBucket)
		return b.Put([]byte(key), []byte(value))
	})
}

// GetMetadata retrieves metadata
func (s *Storage) GetMetadata(key string) (string, error) {
	var value string

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(metaBucket)
		data := b.Get([]byte(key))
		if data != nil {
			value = string(data)
		}
		return nil
	})

	return value, err
}
