package route

const (
	pathURLAssets               = "/assets"
	pathURLRepositoryList       = "/"
	pathURLBuildsByRepositoryID = "/builds"
	pathURLBindList             = "/binds"

	pathURLBuildAPI         = "/api/build"
	pathURLTestBindAPI      = "/api/bind"
	pathURLRepositoryAPI    = "/api/repository"
	pathURLGithubWebhookAPI = "/api/webhook/github"

	pathTokenOwner        = "owner_name"
	pathTokenBuildID      = "build_id"
	pathTokenRepositoryID = "repository_id"
	pathTokenFileName     = "file_name"

	templatesDirPathFromGOPATH = "/src/github.com/dpolansky/grader-ci/backend/static/tmpl"
	assetsDirRelToGOPATH       = "/src/github.com/dpolansky/grader-ci/backend/static/assets/"
)
