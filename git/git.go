// Package git will manage all requirements to clone repository
package git

import (
	"os"

	"github.com/Lord-Y/cypress-parallel-cli/logger"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Repository string // HTTP(s) git repository
	Username   string // Username to use to fetch repository if required
	Password   string // Password to use to fetch repository if required
	Ref        string // Ref in which branch e.g test or refs/head/test
}

func init() {
	logger.SetLoggerLogLevel()
}

// Clone permit to clone git repository
func (c *Repository) Clone() (z string, err error) {
	var (
		targetRef string
		result    *git.Repository
	)

	if c.Ref != "" {
		if c.Ref == "master" || c.Ref == "refs/heads/master" {
			targetRef = ""
		} else {
			targetRef = c.Ref
		}
	} else {
		targetRef = ""
	}
	log.Debug().Msgf("Branch or tag %s", targetRef)

	z, err = os.MkdirTemp(os.TempDir(), fake.CharactersN(10))
	if err != nil {
		return
	}

	if targetRef != "" {
		if c.Username != "" {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
				ReferenceName: plumbing.ReferenceName(targetRef),
				Depth:         1,
			})
			if err != nil {
				return z, err
			}
			_, err := result.Head()
			if err != nil {
				return z, err
			}
		} else {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL:           c.Repository,
				ReferenceName: plumbing.ReferenceName(targetRef),
				SingleBranch:  true,
				Depth:         1,
			})
			if err != nil {
				log.Error().Err(err).Msg("putain")
				return
			}
			_, err := result.Head()
			if err != nil {
				log.Error().Err(err).Msg("zzz")
				return z, err
			}
		}
	} else {
		if c.Username != "" {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
				Depth: 1,
			})
		} else {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL:   c.Repository,
				Depth: 1,
			})
		}
		if err != nil {
			return z, err
		}
	}
	return
}
