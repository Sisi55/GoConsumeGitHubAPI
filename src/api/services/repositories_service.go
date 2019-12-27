package services

import (
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/config"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/domain/github"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/domain/repositories"
	"github.com/Sisi55/GoConsumeGitHubAPI/src/api/provider/github_provider"
	errors "github.com/Sisi55/GoConsumeGitHubAPI/src/api/utils"
	"strings"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, errors.NewBadRequestError("invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Homepage:    "",
		Private:     false,
		HasIssues:   false,
		HasProjects: false,
		HasWiki:     false,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Owner: response.Owner.Login,
		Name:  response.Name,
	}
	return &result, nil
}
