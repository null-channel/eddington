package git

import "github.com/go-git/go-git/v5"

// Eventually clone should support cloning a specific branch
func Clone(repo, path string) (*git.Repository, error) {
	return git.PlainClone(path, false, &git.CloneOptions{
		URL: repo,

		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}
