package handler

import (
	"net/http"
	"strconv"

	"github.com/didikurnia/api-quest/internal/model"
	"github.com/didikurnia/api-quest/internal/store"
	"github.com/gin-gonic/gin"
)

// BookHandler contains handlers for all book-related endpoints.
type BookHandler struct {
	store *store.BookStore
}

// NewBookHandler creates a new BookHandler with the given store.
func NewBookHandler(s *store.BookStore) *BookHandler {
	return &BookHandler{store: s}
}

// Create handles POST /books — Level 3.
func (h *BookHandler) Create(c *gin.Context) {
	var req model.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	book := h.store.Create(req)
	c.JSON(http.StatusCreated, book)
}

// List handles GET /books — Level 3, 6 (search & paginate).
func (h *BookHandler) List(c *gin.Context) {
	// Level 6: Search by author
	if author := c.Query("author"); author != "" {
		books := h.store.Search(author)
		c.JSON(http.StatusOK, books)
		return
	}

	// Level 6: Pagination
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	if pageStr != "" || limitStr != "" {
		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		result := h.store.Paginate(page, limit)
		c.JSON(http.StatusOK, result)
		return
	}

	// Default: return all books
	books := h.store.GetAll()
	c.JSON(http.StatusOK, books)
}

// GetByID handles GET /books/:id — Level 3.
func (h *BookHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	book, ok := h.store.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// Update handles PUT /books/:id — Level 4.
func (h *BookHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	book, ok := h.store.Update(id, req)
	if !ok {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// Delete handles DELETE /books/:id — Level 4.
func (h *BookHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if !h.store.Delete(id) {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "book not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
