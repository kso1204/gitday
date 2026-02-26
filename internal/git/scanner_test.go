package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanRepos(t *testing.T) {
	// 임시 디렉토리에 가짜 git 레포 생성
	tmp := t.TempDir()

	// 레포 2개 생성
	for _, name := range []string{"project-a", "project-b"} {
		gitDir := filepath.Join(tmp, name, ".git")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	// 제외 대상
	excluded := filepath.Join(tmp, "node_modules", ".git")
	os.MkdirAll(excluded, 0755)

	// 숨김 디렉토리
	hidden := filepath.Join(tmp, ".hidden-repo", ".git")
	os.MkdirAll(hidden, 0755)

	// 일반 디렉토리 (git 아님)
	os.MkdirAll(filepath.Join(tmp, "not-a-repo"), 0755)

	repos, err := ScanRepos([]string{tmp}, []string{"node_modules"})
	if err != nil {
		t.Fatal(err)
	}

	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d: %v", len(repos), repos)
	}

	// 이름 확인
	names := make(map[string]bool)
	for _, r := range repos {
		names[filepath.Base(r)] = true
	}
	if !names["project-a"] || !names["project-b"] {
		t.Errorf("expected project-a and project-b, got %v", names)
	}
}

func TestScanRepos_SelfIsGitRepo(t *testing.T) {
	tmp := t.TempDir()
	os.MkdirAll(filepath.Join(tmp, ".git"), 0755)

	repos, err := ScanRepos([]string{tmp}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(repos) != 1 {
		t.Errorf("expected 1 repo (self), got %d", len(repos))
	}
}

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	result := expandHome("~/test")
	expected := filepath.Join(home, "test")
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	// 절대 경로는 그대로
	abs := "/tmp/test"
	if expandHome(abs) != abs {
		t.Errorf("absolute path should not change")
	}
}

func TestShouldExclude(t *testing.T) {
	excludes := []string{"node_modules", "vendor"}

	tests := []struct {
		name     string
		expected bool
	}{
		{"node_modules", true},
		{"vendor", true},
		{".hidden", true},
		{"my-project", false},
	}

	for _, tt := range tests {
		if got := shouldExclude(tt.name, excludes); got != tt.expected {
			t.Errorf("shouldExclude(%q) = %v, want %v", tt.name, got, tt.expected)
		}
	}
}
