package router

import (
	"net/http"

	"github.com/didikurnia/api-quest/internal/config"
	"github.com/didikurnia/api-quest/internal/handler"
	"github.com/didikurnia/api-quest/internal/middleware"
	"github.com/didikurnia/api-quest/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/pakornv/scalar-go"
)

// Setup creates and configures the Gin engine with all routes.
func Setup(cfg *config.Config, bookStore *store.BookStore) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), CORSMiddleware())

	// --- Handlers ---
	bookHandler := handler.NewBookHandler(bookStore)
	authHandler := handler.NewAuthHandler(cfg)

	// --- Public routes ---
	r.GET("/ping", handler.Ping)
	r.POST("/echo", handler.Echo)
	r.POST("/auth/token", authHandler.Token)

	// --- Book routes ---
	// Default to Optional auth: Level 3 & 4 tests make requests without tokens.
	books := r.Group("/books")
	books.Use(middleware.OptionalJWTAuth(cfg))
	{
		books.GET("", bookHandler.List)
		books.GET("/:id", bookHandler.GetByID)
		books.POST("", bookHandler.Create)
		books.PUT("/:id", bookHandler.Update)
		books.DELETE("/:id", bookHandler.Delete)
	}

	// --- API Docs (Scalar) ---
	scalarRef, err := scalar.New("docs/openapi.yaml", &scalar.Config{
		Title: "API Quest — Documentation",
		Theme: scalar.ThemeKepler,
	})
	if err == nil {
		r.GET("/docs", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			html, renderErr := scalarRef.RenderHTML()
			if renderErr != nil {
				http.Error(w, "Failed to render docs", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(html))
		})))
	}

	return r
}

// CORSMiddleware allows Cross-Origin requests from the browser.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
