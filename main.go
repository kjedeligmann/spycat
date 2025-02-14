package main

import (
	"log"

	"github.com/kjedeligmann/spycat/config"
	"github.com/kjedeligmann/spycat/internal/handlers"
	"github.com/kjedeligmann/spycat/internal/repos"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Spy Cat Agency CRUD")

	config.ConnectDB()

	spyCatRepo := repos.NewSpyCatRepo(config.DB)
	spyCatHandler := handlers.NewSpyCatHandler(spyCatRepo)

	missionRepo := repos.NewMissionRepo(config.DB)
	missionHandler := handlers.NewMissionHandler(missionRepo)

	targetRepo := repos.NewTargetRepo(config.DB)
	targetHandler := handlers.NewTargetHandler(targetRepo)

	// Default router uses Logger and Recovery middleware
	// by default, so no need to do it explicitly
	r := gin.Default()

	r.POST("/spy-cats", spyCatHandler.CreateSpyCat)
	r.GET("/spy-cats", spyCatHandler.ListSpyCats)
	r.GET("/spy-cats/:id", spyCatHandler.GetSpyCat)
	r.PATCH("/spy-cats/:id", spyCatHandler.UpdateSpyCatSalary)
	r.DELETE("/spy-cats/:id", spyCatHandler.DeleteSpyCat)

	r.POST("/missions", missionHandler.CreateMission) // along with targets
	r.GET("/missions", missionHandler.ListMissions)
	r.GET("/missions/:id", missionHandler.GetMission)
	r.PATCH("/missions/:id/assign", missionHandler.AssignSpyCat)
	r.PATCH("/missions/:id/completed", missionHandler.MarkMissionCompleted) // mark it as completed
	r.DELETE("/missions/:id", missionHandler.DeleteMission)                 // cannot be deleted if already assigned to the cat

	r.POST("/missions/:id/newtarget", missionHandler.AddTarget) // cannot be added if the mission is already completed
	r.POST("/missions/:id/:targetid/completed", targetHandler.MarkTargetCompleted)
	r.PATCH("/missions/:id/:targetid/notes", targetHandler.UpdateTargetNotes) // cannot be updated if either the target or the mission is completed
	r.DELETE("/missions/:id/:targetid", targetHandler.DeleteTarget)           // cannot be deleted if it's already completed

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
