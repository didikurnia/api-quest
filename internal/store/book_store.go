package store

import (
	"strings"
	"sync"

	"github.com/didikurnia/api-quest/internal/model"
	"github.com/google/uuid"
)

// BookStore is a thread-safe, in-memory book storage.
type BookStore struct {
	mu    sync.RWMutex
	books map[string]model.Book
	order []string // preserves insertion order
}

// NewBookStore creates a new empty BookStore.
func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]model.Book),
		order: make([]string, 0),
	}
}

// Create adds a new book and returns it with a generated ID.
func (s *BookStore) Create(req model.CreateBookRequest) model.Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := model.Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
	}
	s.books[book.ID] = book
	s.order = append(s.order, book.ID)
	return book
}

// GetAll returns all books in insertion order.
func (s *BookStore) GetAll() []model.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Book, 0, len(s.order))
	for _, id := range s.order {
		if b, ok := s.books[id]; ok {
			result = append(result, b)
		}
	}
	return result
}

// GetByID returns a single book by ID.
func (s *BookStore) GetByID(id string) (model.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.books[id]
	return b, ok
}

// Update replaces a book by ID. Returns the updated book and whether it existed.
func (s *BookStore) Update(id string, req model.UpdateBookRequest) (model.Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return model.Book{}, false
	}

	book := model.Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
	}
	s.books[id] = book
	return book, true
}

// Delete removes a book by ID. Returns whether it existed.
func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return false
	}

	delete(s.books, id)
	// Remove from order slice
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
	return true
}

// Search returns books matching the given author (case-insensitive substring match).
func (s *BookStore) Search(author string) []model.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Book, 0)
	authorLower := strings.ToLower(author)
	for _, id := range s.order {
		b := s.books[id]
		if strings.Contains(strings.ToLower(b.Author), authorLower) {
			result = append(result, b)
		}
	}
	return result
}

// Paginate returns a page of books with metadata.
func (s *BookStore) Paginate(page, limit int) model.PaginatedResponse {
	all := s.GetAll()
	total := len(all)

	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	data := all[start:end]
	if data == nil {
		data = []model.Book{}
	}

	return model.PaginatedResponse{
		Data:  data,
		Page:  page,
		Limit: limit,
		Total: total,
	}
}
