package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	dataDir string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "simplefts",
		Short: "Simple Full-Text Search Engine",
		Long:  "A full-text search engine with BM25 ranking, HTTP API, and persistent storage",
	}

	rootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "./data/search.db", "Data directory for storage")

	// Serve command
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP API server",
		Run:   runServe,
	}
	serveCmd.Flags().StringP("host", "H", "127.0.0.1", "Server host")
	serveCmd.Flags().IntP("port", "p", 3000, "Server port")

	// Insert command
	insertCmd := &cobra.Command{
		Use:   "insert",
		Short: "Insert a document",
		Run:   runInsert,
	}
	insertCmd.Flags().StringP("id", "i", "", "Document ID (required)")
	insertCmd.Flags().StringP("title", "t", "", "Document title (required)")
	insertCmd.Flags().StringP("content", "c", "", "Document content (required)")
	insertCmd.Flags().StringP("url", "u", "", "Document URL")
	insertCmd.MarkFlagRequired("id")
	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("content")

	// Search command
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for documents",
		Run:   runSearch,
	}
	searchCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchCmd.Flags().IntP("limit", "l", 10, "Maximum results")
	searchCmd.Flags().BoolP("ranked", "r", true, "Use BM25 ranking")
	searchCmd.Flags().String("mode", "and", "Search mode: and or or")
	searchCmd.MarkFlagRequired("query")

	// Get command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a document by ID",
		Run:   runGet,
	}
	getCmd.Flags().StringP("id", "i", "", "Document ID (required)")
	getCmd.MarkFlagRequired("id")

	// Delete command
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a document",
		Run:   runDelete,
	}
	deleteCmd.Flags().StringP("id", "i", "", "Document ID (required)")
	deleteCmd.MarkFlagRequired("id")

	// Stats command
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show index statistics",
		Run:   runStats,
	}

	rootCmd.AddCommand(serveCmd, insertCmd, searchCmd, getCmd, deleteCmd, statsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServe(cmd *cobra.Command, args []string) {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetInt("port")

	log.Printf("Starting search engine with data: %s", dataDir)

	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Server listening on http://%s", addr)
	log.Println("API Documentation:")
	log.Println("  GET    /health              - Health check")
	log.Println("  POST   /documents           - Insert a document")
	log.Println("  POST   /documents/batch     - Batch insert documents")
	log.Println("  GET    /documents/:id       - Get a document")
	log.Println("  PUT    /documents/:id       - Update a document")
	log.Println("  DELETE /documents/:id       - Delete a document")
	log.Println("  GET    /search?query=...    - Search documents")
	log.Println("  GET    /stats               - Get index statistics")

	api := NewAPI(engine)
	if err := api.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runInsert(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	title, _ := cmd.Flags().GetString("title")
	content, _ := cmd.Flags().GetString("content")
	url, _ := cmd.Flags().GetString("url")

	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	doc := NewDocument(id, title, content)
	doc.URL = url

	if err := engine.UpsertDocument(doc); err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	fmt.Printf("✓ Document '%s' inserted successfully\n", id)
}

func runSearch(cmd *cobra.Command, args []string) {
	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt("limit")
	ranked, _ := cmd.Flags().GetBool("ranked")
	modeStr, _ := cmd.Flags().GetString("mode")

	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	options := DefaultSearchOptions()
	options.Limit = limit
	options.UseRanking = ranked

	if modeStr == "or" {
		options.Mode = SearchModeOR
	}

	start := time.Now()
	result, err := engine.Search(query, options)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	duration := time.Since(start)

	fmt.Printf("\n🔍 Search Results for: \"%s\"\n", query)
	fmt.Printf("Found %d documents in %v\n\n", result.Total, duration)

	for i, doc := range result.Documents {
		if len(result.Scores) > 0 {
			fmt.Printf("%d. [Score: %.4f] %s\n", i+1, result.Scores[i], doc.Title)
		} else {
			fmt.Printf("%d. %s\n", i+1, doc.Title)
		}
		fmt.Printf("   ID: %s\n", doc.ID)
		if doc.URL != "" {
			fmt.Printf("   URL: %s\n", doc.URL)
		}
		contentPreview := doc.Content
		if len(contentPreview) > 100 {
			contentPreview = contentPreview[:100] + "..."
		}
		fmt.Printf("   Content: %s\n\n", contentPreview)
	}
}

func runGet(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")

	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	doc, err := engine.GetDocument(id)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	if doc == nil {
		fmt.Printf("❌ Document '%s' not found\n", id)
		return
	}

	fmt.Println("\n📄 Document")
	fmt.Printf("ID:      %s\n", doc.ID)
	fmt.Printf("Title:   %s\n", doc.Title)
	if doc.URL != "" {
		fmt.Printf("URL:     %s\n", doc.URL)
	}
	fmt.Printf("Content: %s\n\n", doc.Content)
}

func runDelete(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")

	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	if err := engine.DeleteDocument(id); err != nil {
		log.Fatalf("Failed to delete document: %v", err)
	}

	fmt.Printf("✓ Document '%s' deleted successfully\n", id)
}

func runStats(cmd *cobra.Command, args []string) {
	engine, err := NewSearchEngine(dataDir)
	if err != nil {
		log.Fatalf("Failed to create search engine: %v", err)
	}
	defer engine.Close()

	stats := engine.Stats()

	fmt.Println("\n📊 Index Statistics")
	fmt.Printf("Total Documents:       %d\n", stats.TotalDocuments)
	fmt.Printf("Total Unique Tokens:   %d\n", stats.TotalTokens)
	fmt.Printf("Avg Docs per Token:    %.2f\n\n", stats.AvgDocsPerToken)
}
