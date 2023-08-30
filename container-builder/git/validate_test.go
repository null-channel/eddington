package git_test

import (
	"testing"

	"github.com/null-channel/eddington/container-builder/git"
)

// TestValidateGithubUrl tests the ValidateGithubUrl function
func TestValidateGitHubURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "https://github.com/null-channel/eddington",
			wantErr: false,
		},
		{
			name:    "invalid url format",
			url:     "htps://github.com//null-channel/eddington",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := git.ValidateGitHubURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test case %q: got error %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
