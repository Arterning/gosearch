package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API represents the HTTP API server
type API struct {
	engine *SearchEngine
	router *gin.Engine
}

// NewAPI creates a new API server
func NewAPI(engine *SearchEngine) *API {
	router := gin.Default()

	api := &API{
		engine: engine,
		router: router,
	}

	api.setupRoutes()
	return api
}

// setupRoutes sets up API routes
func (api *API) setupRoutes() {
	api.router.GET("/health", api.handleHealth)
	api.router.POST("/documents", api.handleInsertDocument)
	api.router.POST("/documents/batch", api.handleBatchInsert)
	api.router.GET("/documents/:id", api.handleGetDocument)
	api.router.PUT("/documents/:id", api.handleUpdateDocument)
	api.router.DELETE("/documents/:id", api.handleDeleteDocument)
	api.router.GET("/search", api.handleSearch)
	api.router.GET("/stats", api.handleStats)
}

// Run starts the API server
func (api *API) Run(addr string) error {
	return api.router.Run(addr)
}

// Response types
type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Request types
type insertDocumentRequest struct {
	ID      string `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	URL     string `json:"url"`
}

type batchInsertRequest struct {
	Documents []insertDocumentRequest `json:"documents" binding:"required"`
}

type searchResponse struct {
	Documents []*Document `json:"documents"`
	Total     int         `json:"total"`
	Query     string      `json:"query"`
	Scores    []float64   `json:"scores,omitempty"`
}

// Handlers

func (api *API) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Data:    "OK",
	})
}

func (api *API) handleInsertDocument(c *gin.Context) {
	var req insertDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	doc := NewDocument(req.ID, req.Title, req.Content)
	doc.URL = req.URL

	if err := api.engine.UpsertDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Message: "Document inserted successfully",
	})
}

func (api *API) handleBatchInsert(c *gin.Context) {
	var req batchInsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	for _, docReq := range req.Documents {
		doc := NewDocument(docReq.ID, docReq.Title, docReq.Content)
		doc.URL = docReq.URL

		if err := api.engine.UpsertDocument(doc); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Message: "Documents inserted successfully",
	})
}

func (api *API) handleGetDocument(c *gin.Context) {
	id := c.Param("id")

	doc, err := api.engine.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if doc == nil {
		c.JSON(http.StatusNotFound, errorResponse{
			Success: false,
			Error:   "Document not found",
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Data:    doc,
	})
}

func (api *API) handleUpdateDocument(c *gin.Context) {
	id := c.Param("id")

	var req insertDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	doc := NewDocument(id, req.Title, req.Content)
	doc.URL = req.URL

	if err := api.engine.UpsertDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Message: "Document updated successfully",
	})
}

func (api *API) handleDeleteDocument(c *gin.Context) {
	id := c.Param("id")

	if err := api.engine.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Message: "Document deleted successfully",
	})
}

func (api *API) handleSearch(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, errorResponse{
			Success: false,
			Error:   "query parameter is required",
		})
		return
	}

	// Parse options
	options := DefaultSearchOptions()

	if mode := c.Query("mode"); mode == "or" {
		options.Mode = SearchModeOR
	}

	if ranked := c.Query("ranked"); ranked == "false" {
		options.UseRanking = false
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			options.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			options.Offset = offset
		}
	}

	// Perform search
	result, err := api.engine.Search(query, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Data: searchResponse{
			Documents: result.Documents,
			Total:     result.Total,
			Query:     query,
			Scores:    result.Scores,
		},
	})
}

func (api *API) handleStats(c *gin.Context) {
	stats := api.engine.Stats()

	c.JSON(http.StatusOK, successResponse{
		Success: true,
		Data:    stats,
	})
}
