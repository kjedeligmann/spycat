package handlers

import (
	"net/http"
	"strconv"

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

	// (Optional) Validate breed with TheCatAPI here:
	// if !isValidBreed(cat.Breed) { ... }

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

func (h *SpyCatHandler) UpdateSpyCat(c *gin.Context) {
	// Convert path param (id) to int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Bind incoming JSON to a SpyCat struct
	var cat models.SpyCat
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure the ID in the path matches the one we're updating
	cat.ID = id

	// (Optional) Validate breed with TheCatAPI here before updating:
	// if !isValidBreed(cat.Breed) {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat breed"})
	//     return
	// }

	// Call the repository to update the record
	if err := h.repo.Update(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cat)
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
