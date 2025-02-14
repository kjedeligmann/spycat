package handlers

import (
	"net/http"
	"strconv"

	"github.com/kjedeligmann/spycat/internal/models"
	"github.com/kjedeligmann/spycat/internal/repos"

	"github.com/gin-gonic/gin"
)

type MissionHandler struct {
	repo *repos.MissionRepo
}

func NewMissionHandler(repo *repos.MissionRepo) *MissionHandler {
	return &MissionHandler{repo: repo}
}

// Create mission along with targets
func (h *MissionHandler) CreateMission(c *gin.Context) {
	var mission models.Mission
	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Each mission must have between 1 and 3 targets"})
		return
	}

	// Default mission status is ongoing if not provided.
	if mission.Status == "" {
		mission.Status = "ongoing"
	}

	err := h.repo.Create(c.Request.Context(), &mission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mission"})
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// Get a single mission
func (h *MissionHandler) GetMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	mission, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}

	c.JSON(http.StatusOK, mission)
}

// List all missions
func (h *MissionHandler) ListMissions(c *gin.Context) {
	missions, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve missions"})
		return
	}

	c.JSON(http.StatusOK, missions)
}

// Assign spy cat to a mission
func (h *MissionHandler) AssignSpyCat(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	var req struct {
		CatID int `json:"cat_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = h.repo.AssignSpyCat(c.Request.Context(), missionID, req.CatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign spy cat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Spy cat assigned successfully"})
}

// Mark mission as completed
func (h *MissionHandler) MarkMissionCompleted(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	err = h.repo.MarkMissionCompleted(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark mission as completed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission marked as completed"})
}

// Delete mission (only if not assigned to a cat)
func (h *MissionHandler) DeleteMission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Mission cannot be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission deleted successfully"})
}

// Add a new target to a mission
func (h *MissionHandler) AddTarget(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.AddTarget(c.Request.Context(), missionID, &target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, target)
}
