package app

import (
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/controllers/polo"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", polo.Polo)
	router.POST("/repositories", repositories.CreateRepo)
}
