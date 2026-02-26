package git

import (
	"testing"
)

func TestParseGitLog(t *testing.T) {
	raw := `abc1234567890§§첫 번째 커밋§§wook§§2026-02-26T15:00:00+09:00

 3 files changed, 100 insertions(+), 20 deletions(-)
def7890123456§§두 번째 커밋§§wook§§2026-02-26T16:00:00+09:00

 1 file changed, 10 insertions(+)`

	commits, err := parseGitLog(raw)
	if err != nil {
		t.Fatal(err)
	}

	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	// 첫 번째 커밋
	c := commits[0]
	if c.Hash != "abc1234" {
		t.Errorf("hash = %q, want abc1234", c.Hash)
	}
	if c.Message != "첫 번째 커밋" {
		t.Errorf("message = %q", c.Message)
	}
	if c.Author != "wook" {
		t.Errorf("author = %q", c.Author)
	}
	if c.Files != 3 {
		t.Errorf("files = %d, want 3", c.Files)
	}

	// 두 번째 커밋
	c2 := commits[1]
	if c2.Hash != "def7890" {
		t.Errorf("hash = %q, want def7890", c2.Hash)
	}
	if c2.Files != 1 {
		t.Errorf("files = %d, want 1", c2.Files)
	}
}

func TestParseGitLog_Empty(t *testing.T) {
	commits, err := parseGitLog("")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 0 {
		t.Errorf("expected 0 commits, got %d", len(commits))
	}
}

func TestParseFileCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"3 files changed, 100 insertions(+), 20 deletions(-)", 3},
		{"1 file changed, 10 insertions(+)", 1},
		{"", 0},
	}

	for _, tt := range tests {
		if got := parseFileCount(tt.input); got != tt.expected {
			t.Errorf("parseFileCount(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}
