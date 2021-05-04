// Package git will manage all requirements to clone repository
package git

import (
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/icrowley/fake"
)

type Repository struct {
	Repository string // HTTP(s) git repository
	Username   string // Username to use to fetch repository if required
	Password   string // Password to use to fetch repository if required
	Branch     string // Branch in which specs are hold
}

// Clone permit to clone git repository
func (c *Repository) Clone() (z string, err error) {
	var (
		targetBranch string
		result       *git.Repository
	)

	if c.Branch != "" {
		if c.Branch == "master" {
			targetBranch = ""
		} else {
			targetBranch = c.Branch
		}
	} else {
		targetBranch = ""
	}

	z, err = os.MkdirTemp(os.TempDir(), fake.CharactersN(10))
	if err != nil {
		return z, err
	}

	if targetBranch != "" {
		if c.Username != "" {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
			})
		} else {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
			})
		}
		if err != nil {
			return z, err
		}
		w, err := result.Worktree()
		if err != nil {
			return z, err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(targetBranch),
		})
		if err != nil {
			return z, err
		}
	} else {
		if c.Username != "" {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
			})
		} else {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
			})
		}
		if err != nil {
			return z, err
		}
	}
	return
}
