package location

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/location-demo/internal/domain"
)

// Handler holds the HTTP handlers for location endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new handler with the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes wires all location routes to the Gin engine.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/locations/search", h.Search)
		api.GET("/locations/trending", h.GetTrending)
		api.GET("/locations/autocomplete", h.Autocomplete)
		api.GET("/locations/:id", h.GetByID)
		api.GET("/locations/:id/children", h.GetChildren)
		api.POST("/locations/hydrate", h.Hydrate)

		api.POST("/posts", h.CreatePost)
		api.GET("/posts", h.GetPosts)
		api.GET("/locations/:id/posts", h.GetPostsByLocation)
	}
}

// ──────────────────────────────────────────────
// Location Handlers
// ──────────────────────────────────────────────

// Search handles GET /api/v1/locations/search?q=sai+gon&lang=vi&type=city
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	lang := c.DefaultQuery("lang", "en")
	locType := c.Query("type")

	results, err := h.service.LocalSearch(c.Request.Context(), query, lang, locType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// Autocomplete handles GET /api/v1/locations/autocomplete?q=sai&lang=vi
func (h *Handler) Autocomplete(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	lang := c.DefaultQuery("lang", "en")

	results, err := h.service.Autocomplete(c.Request.Context(), query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

// GetByID handles GET /api/v1/locations/:id?lang=en
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	lang := c.DefaultQuery("lang", "en")

	detail, err := h.service.GetByID(c.Request.Context(), id, lang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": detail})
}

// GetChildren handles GET /api/v1/locations/:id/children?lang=en
func (h *Handler) GetChildren(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	lang := c.DefaultQuery("lang", "en")

	children, err := h.service.GetChildren(c.Request.Context(), id, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": children})
}

// GetTrending handles GET /api/v1/locations/trending?lang=en&limit=10
func (h *Handler) GetTrending(c *gin.Context) {
	lang := c.DefaultQuery("lang", "en")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	trending, err := h.service.GetTrending(c.Request.Context(), lang, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trending})
}

type hydrateRequest struct {
	ExternalID string `json:"external_id" binding:"required"`
	Lang       string `json:"lang"`
}

// Hydrate handles POST /api/v1/locations/hydrate
func (h *Handler) Hydrate(c *gin.Context) {
	var req hydrateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.EnsureLocationByExternalID(c.Request.Context(), req.ExternalID, req.Lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
}

// ──────────────────────────────────────────────
// Post Handlers
// ──────────────────────────────────────────────

type createPostRequest struct {
	Content    string `json:"content" binding:"required"`
	MediaType  string `json:"media_type"`
	LocationID int64  `json:"location_id"`
	ExternalID string `json:"external_id"`
	Lang       string `json:"lang"`
}

// CreatePost handles POST /api/v1/posts
func (h *Handler) CreatePost(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		Content:    req.Content,
		MediaType:  req.MediaType,
		LocationID: req.LocationID,
	}

	createdPost, err := h.service.CreatePost(c.Request.Context(), post, req.ExternalID, req.Lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdPost})
}

// GetPosts handles GET /api/v1/posts?lang=en&limit=20&offset=0&location_id=123
func (h *Handler) GetPosts(c *gin.Context) {
	lang := c.DefaultQuery("lang", "en")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	var locID *int64
	locationIDStr := c.Query("location_id")
	if locationIDStr != "" {
		id, err := strconv.ParseInt(locationIDStr, 10, 64)
		if err == nil {
			locID = &id
		}
	}

	posts, err := h.service.GetPosts(c.Request.Context(), locID, lang, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// GetPostsByLocation handles GET /api/v1/locations/:id/posts?lang=en&limit=20&offset=0
func (h *Handler) GetPostsByLocation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location ID"})
		return
	}

	lang := c.DefaultQuery("lang", "en")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	posts, err := h.service.GetPostsByLocation(c.Request.Context(), id, lang, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}
