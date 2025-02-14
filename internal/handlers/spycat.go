package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/kjedeligmann/spycat/internal/models"
	"github.com/kjedeligmann/spycat/internal/repos"

	"github.com/gin-gonic/gin"
)

type SpyCatHandler struct {
	repo *repos.SpyCatRepo
}

func NewSpyCatHandler(repo *repos.SpyCatRepo) *SpyCatHandler {
	return &SpyCatHandler{repo: repo}
}

func (h *SpyCatHandler) CreateSpyCat(c *gin.Context) {
	var cat models.SpyCat
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate breed with TheCatAPI
	if !isValidBreed(cat.Breed) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid breed"})
		return
	}

	err := h.repo.Create(c.Request.Context(), &cat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *SpyCatHandler) GetSpyCat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cat, err := h.repo.Read(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SpyCat not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *SpyCatHandler) ListSpyCats(c *gin.Context) {
	cats, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

func (h *SpyCatHandler) UpdateSpyCatSalary(c *gin.Context) {
	// Convert path param (id) to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Extract salary field from incoming JSON
	var s = struct {
		Salary float64 `json:"salary"`
	}{}
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the repository to update the record
	if err := h.repo.UpdateSalary(c.Request.Context(), id, s.Salary); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *SpyCatHandler) DeleteSpyCat(c *gin.Context) {
	// Convert path param (id) to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Call the repository to delete the record
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return 204 No Content to indicate successful deletion
	c.Status(http.StatusNoContent)
}

// isValidBreed checks if the provided breed exists in the list of breeds from TheCatAPI
func isValidBreed(breed string) bool {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var breeds = []struct {
		Name string `json:"name"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false
	}

	// Check if the provided breed matches any of the breeds (case-insensitive)
	for _, b := range breeds {
		if strings.EqualFold(b.Name, breed) {
			return true
		}
	}
	return false
}
