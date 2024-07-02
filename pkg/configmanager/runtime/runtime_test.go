package runtime

import (
	"nestos-kubernetes-deployer/pkg/api"
	"nestos-kubernetes-deployer/pkg/constants"
	"testing"

	"github.com/pkg/errors"
)

func TestGetRuntime(t *testing.T) {
	tests := []struct {
		name     string
		runtime  string
		expected api.Runtime
		err      error
	}{
		{
			name:     "Empty string defaults to Isulad",
			runtime:  "",
			expected: &isuladRuntime{},
			err:      nil,
		},
		{
			name:     "Isulad",
			runtime:  constants.Isulad,
			expected: &isuladRuntime{},
			err:      nil,
		},
		{
			name:     "Docker",
			runtime:  constants.Docker,
			expected: &dockerRuntime{},
			err:      nil,
		},
		{
			name:     "Crio",
			runtime:  constants.Crio,
			expected: &crioRuntime{},
			err:      nil,
		},
		{
			name:     "Containerd",
			runtime:  constants.Containerd,
			expected: &containerdRuntime{},
			err:      nil,
		},
		{
			name:     "Unsupported",
			runtime:  "unsupported",
			expected: nil,
			err:      errors.New("unsupported runtime"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRuntime(tt.runtime)

			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Errorf("expected error %v, got %v", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if _, ok := got.(*isuladRuntime); tt.runtime == constants.Isulad && !ok {
					t.Errorf("expected IsuladRuntime, got %T", got)
				}
				if _, ok := got.(*dockerRuntime); tt.runtime == constants.Docker && !ok {
					t.Errorf("expected DockerRuntime, got %T", got)
				}
				if _, ok := got.(*crioRuntime); tt.runtime == constants.Crio && !ok {
					t.Errorf("expected CrioRuntime, got %T", got)
				}
				if _, ok := got.(*containerdRuntime); tt.runtime == constants.Containerd && !ok {
					t.Errorf("expected ContainerdRuntime, got %T", got)
				}
			}
		})
	}
}
