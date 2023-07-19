package models

import "testing"

func TestValidateNullApplicationService(t *testing.T) {
	tests := []struct {
		name    string
		service *NullApplicationService
		wantErr bool
	}{
		{
			name: "valid container image",
			service: &NullApplicationService{
				Name: "test",
				Type: ContainerImage,
			},
			wantErr: false,
		},
		{
			name: "valid git repo",
			service: &NullApplicationService{
				Name:    "test",
				Type:    Go,
				GitRepo: "https://github.com",
				GitSha:  "1234",
			},
			wantErr: false,
		},
		{
			name: "invalid git repo",
			service: &NullApplicationService{
				Name:    "test",
				Type:    ContainerImage,
				GitRepo: "https://github.com",
			},
			wantErr: true,
		},
		{
			service: &NullApplicationService{
				Name:    "no name test",
				Type:    ContainerImage,
				GitRepo: "https://github.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNullApplicationService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
