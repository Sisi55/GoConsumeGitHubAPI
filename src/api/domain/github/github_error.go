package github

type GithubErrorResponse struct {
	StatusCode       int           `json:"status_code"`
	Message          string        `json:""`
	DocumentationUrl string        `json:""`
	Errors           []GithubError `json:""`
}

type GithubError struct {
	Resource string `json:"resource"`
	Code     string `json:"code"`
	Field    string `json:"field"`
	Message  string `json:"message"`
}
