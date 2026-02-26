package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/kso1204/gitday/internal/git"
)

// Provider는 AI 요약 프로바이더 인터페이스이다.
type Provider interface {
	Summarize(ctx context.Context, prompt string) (string, error)
	Name() string
}

// NewProvider는 설정에 따라 적절한 AI 프로바이더를 생성한다.
func NewProvider(providerName, apiKey, model, ollamaURL string) (Provider, error) {
	switch strings.ToLower(providerName) {
	case "claude":
		if apiKey == "" {
			return nil, fmt.Errorf("Claude API 키가 필요합니다 (GITDAY_API_KEY 환경변수 또는 설정 파일)")
		}
		return NewClaude(apiKey, model), nil
	case "openai":
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API 키가 필요합니다 (GITDAY_API_KEY 환경변수 또는 설정 파일)")
		}
		return NewOpenAI(apiKey, model), nil
	case "ollama":
		return NewOllama(ollamaURL, model), nil
	default:
		return nil, fmt.Errorf("지원하지 않는 AI 프로바이더: %s (claude/openai/ollama)", providerName)
	}
}

// BuildPrompt는 커밋 데이터로 요약 프롬프트를 생성한다.
func BuildPrompt(results []git.RepoResult, since string) string {
	var sb strings.Builder
	sb.WriteString("다음은 개발자의 Git 커밋 로그입니다. 이 내용을 바탕으로 오늘 한 일을 자연어로 간결하게 요약해주세요.\n")
	sb.WriteString("- 프로젝트별로 핵심 작업을 1-2문장으로 요약\n")
	sb.WriteString("- 마지막에 전체적인 한줄 요약 추가\n")
	sb.WriteString("- 한국어로 작성\n\n")

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("## %s (%d commits)\n", r.Name, len(r.Commits)))
		for _, c := range r.Commits {
			sb.WriteString(fmt.Sprintf("- %s\n", c.Message))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
