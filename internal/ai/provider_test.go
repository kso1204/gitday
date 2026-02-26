package ai

import (
	"strings"
	"testing"

	"github.com/kso1204/gitday/internal/git"
)

func TestNewProvider_Claude(t *testing.T) {
	p, err := NewProvider("claude", "test-key", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name() != "Claude" {
		t.Errorf("name = %q, want Claude", p.Name())
	}
}

func TestNewProvider_OpenAI(t *testing.T) {
	p, err := NewProvider("openai", "test-key", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name() != "OpenAI" {
		t.Errorf("name = %q, want OpenAI", p.Name())
	}
}

func TestNewProvider_Ollama(t *testing.T) {
	p, err := NewProvider("ollama", "", "", "http://localhost:11434")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name() != "Ollama" {
		t.Errorf("name = %q, want Ollama", p.Name())
	}
}

func TestNewProvider_NoKey(t *testing.T) {
	_, err := NewProvider("claude", "", "", "")
	if err == nil {
		t.Error("expected error for missing API key")
	}
}

func TestNewProvider_Unknown(t *testing.T) {
	_, err := NewProvider("unknown", "", "", "")
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}

func TestBuildPrompt(t *testing.T) {
	results := []git.RepoResult{
		{
			Name: "rpg",
			Commits: []git.Commit{
				{Message: "전투 시스템 수정"},
				{Message: "인벤토리 UI 개선"},
			},
		},
		{
			Name: "petition",
			Commits: []git.Commit{
				{Message: "API 에러 핸들링"},
			},
		},
	}

	prompt := BuildPrompt(results, "2026-02-26")

	if !strings.Contains(prompt, "rpg") {
		t.Error("prompt should contain repo name")
	}
	if !strings.Contains(prompt, "전투 시스템 수정") {
		t.Error("prompt should contain commit message")
	}
	if !strings.Contains(prompt, "한국어") {
		t.Error("prompt should request Korean")
	}
}
