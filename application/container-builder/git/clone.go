package git

import "github.com/go-git/go-git/v5"

// Eventually clone should support cloning a specific branch
func Clone(repoName string) (*git.Repository, error) {
	return git.PlainClone(repoName, false, &git.CloneOptions{
		URL:               repoName,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}
