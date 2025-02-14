package handlers

import (
	"net/http"
	"strconv"

	"github.com/kjedeligmann/spycat/internal/repos"

	"github.com/gin-gonic/gin"
)

type TargetHandler struct {
	repo *repos.TargetRepo
}

func NewTargetHandler(repo *repos.TargetRepo) *TargetHandler {
	return &TargetHandler{repo: repo}
}

// ListTargetsByMission lists all targets associated with a given mission.
// func (h *TargetHandler) ListTargetsByMission(c *gin.Context) {
// 	missionID, err := strconv.Atoi(c.Param("mission_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
// 		return
// 	}
// 	targets, err := h.repo.ListByMission(c.Request.Context(), missionID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, targets)
// }

// MarkTargetCompleted marks a target as completed.
func (h *TargetHandler) MarkTargetCompleted(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	targetID, err := strconv.Atoi(c.Param("targetid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}
	if err := h.repo.MarkTargetCompleted(c.Request.Context(), missionID, targetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Target marked as completed"})
}

// UpdateTargetNotes updates a target's notes if allowed.
func (h *TargetHandler) UpdateTargetNotes(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	targetID, err := strconv.Atoi(c.Param("targetid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}
	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.repo.UpdateTargetNotes(c.Request.Context(), missionID, targetID, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Target notes updated"})
}

// DeleteTarget deletes a target if permitted.
func (h *TargetHandler) DeleteTarget(c *gin.Context) {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	targetID, err := strconv.Atoi(c.Param("targetid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}
	if err := h.repo.DeleteTarget(c.Request.Context(), missionID, targetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}
