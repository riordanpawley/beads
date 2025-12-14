package main

import (
	"os"
	"testing"
)

func TestTruncateForBox(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		maxLen int
		want   string
	}{
		{
			name:   "short path no truncate",
			path:   "/home/user",
			maxLen: 20,
			want:   "/home/user",
		},
		{
			name:   "exact length",
			path:   "12345",
			maxLen: 5,
			want:   "12345",
		},
		{
			name:   "needs truncate",
			path:   "/very/long/path/to/somewhere/deep",
			maxLen: 15,
			want:   "...mewhere/deep",
		},
		{
			name:   "truncate to minimum",
			path:   "abcdefghij",
			maxLen: 5,
			want:   "...ij",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateForBox(tt.path, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateForBox(%q, %d) = %q, want %q", tt.path, tt.maxLen, got, tt.want)
			}
			if len(got) > tt.maxLen {
				t.Errorf("truncateForBox(%q, %d) returned %q with length %d > maxLen %d",
					tt.path, tt.maxLen, got, len(got), tt.maxLen)
			}
		})
	}
}

func TestGitRevParse(t *testing.T) {
	// Basic test - should either return a value or empty string (if not in git repo)
	result := gitRevParse("--git-dir")
	// Just verify it doesn't panic and returns a string
	if result != "" {
		// In a git repo
		t.Logf("Git dir: %s", result)
	} else {
		// Not in a git repo or error
		t.Logf("Not in git repo or error")
	}
}

func TestDaemonDisableReasonEnvVar(t *testing.T) {
	// Test that BEADS_NO_DAEMON disables daemon
	tests := []struct {
		envValue string
		expected DaemonDisableReason
	}{
		{"1", DaemonDisabledEnvVar},
		{"true", DaemonDisabledEnvVar},
		{"yes", DaemonDisabledEnvVar},
		{"on", DaemonDisabledEnvVar},
		{"TRUE", DaemonDisabledEnvVar},
		{"  1  ", DaemonDisabledEnvVar},
	}

	for _, tt := range tests {
		t.Run("BEADS_NO_DAEMON="+tt.envValue, func(t *testing.T) {
			os.Setenv("BEADS_NO_DAEMON", tt.envValue)
			defer os.Unsetenv("BEADS_NO_DAEMON")

			got := getDaemonDisableReason()
			if got != tt.expected {
				t.Errorf("getDaemonDisableReason() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsWorktreeWithoutSyncBranch(t *testing.T) {
	// Ensure env var doesn't interfere
	os.Unsetenv("BEADS_NO_DAEMON")
	os.Unsetenv("BEADS_SYNC_BRANCH")

	// This test verifies the function doesn't panic and returns a boolean
	// The actual result depends on whether we're in a worktree
	result := isWorktreeWithoutSyncBranch()
	t.Logf("isWorktreeWithoutSyncBranch: %v", result)
}
