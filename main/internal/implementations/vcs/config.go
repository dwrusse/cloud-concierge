package vcs

// Config contains the values that allow for authentication and the specific repo
// traits needed.
type Config struct {
	// VCSBaseBranch is the name of the base branch within the version control into which
	// new PRs should be opened.
	VCSBaseBranch string `required:"true"`

	// VCSToken is the auth token needed to read code and open pull requests within a customer's VCS
	// environment.
	VCSToken string `required:"true"`

	// VCSUser is the name of the user with whom VCSToken is associated.
	VCSUser string `required:"true"`

	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	VCSRepo string `required:"true"`

	// VCSSystem is the name of the version control system chosen.
	// At the moment only GitHub is supported.
	VCSSystem string `required:"true"`

	// PullReviewers is the name of the pull request reviewer who will be tagged on the opened pull request.
	PullReviewers []string `default:"NoReviewer"`
}
