package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wook/gitday/internal/ai"
	"github.com/wook/gitday/internal/git"
	"github.com/wook/gitday/internal/output"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "오늘의 활동을 로컬 파일에 자동 저장",
	Long:  `커밋 로그 + AI 요약을 ~/.gitday/logs/YYYY-MM-DD.md에 저장합니다.`,
	RunE:  runLog,
}

func init() {
	logCmd.Flags().String("period", "today", "기간: today, week")
	rootCmd.AddCommand(logCmd)
}

func runLog(cmd *cobra.Command, args []string) error {
	period, _ := cmd.Flags().GetString("period")
	now := time.Now()
	var since time.Time

	switch period {
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := now.AddDate(0, 0, -(weekday - 1))
		since = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
	default:
		since = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}

	scanPaths := viper.GetStringSlice("scan_paths")
	excludes := viper.GetStringSlice("exclude")
	author := viper.GetString("author")

	repos, err := git.ScanRepos(scanPaths, excludes)
	if err != nil {
		return fmt.Errorf("레포 스캔 실패: %w", err)
	}

	results, err := git.CollectLogs(repos, since, now, author)
	if err != nil {
		return fmt.Errorf("커밋 로그 수집 실패: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("저장할 커밋이 없습니다.")
		return nil
	}

	// AI 요약 시도
	var summaryText string
	summaryText = getSummary(results, since)

	// 마크다운 생성
	md := output.ToMarkdown(results, since, now, summaryText)

	// 저장 경로
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(home, ".gitday", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	filename := since.Format("2006-01-02") + ".md"
	if period == "week" {
		filename = since.Format("2006-01-02") + "_week.md"
	}
	logPath := filepath.Join(logDir, filename)

	if err := os.WriteFile(logPath, []byte(md), 0644); err != nil {
		return fmt.Errorf("로그 저장 실패: %w", err)
	}

	// 터미널에도 출력
	compact := viper.GetBool("output.compact")
	output.PrintReport(results, since, now, compact)

	if summaryText != "" {
		output.PrintSummary(summaryText)
	}

	fmt.Printf("\n✓ 저장됨: %s\n", logPath)
	return nil
}

func getSummary(results []git.RepoResult, since time.Time) string {
	providerName := viper.GetString("ai.provider")
	apiKey := viper.GetString("ai.api_key")
	model := viper.GetString("ai.model")
	ollamaURL := viper.GetString("ai.ollama_url")

	if envKey := os.Getenv("GITDAY_API_KEY"); envKey != "" {
		apiKey = envKey
	} else if envKey := os.Getenv("ANTHROPIC_API_KEY"); envKey != "" && providerName == "claude" {
		apiKey = envKey
	} else if envKey := os.Getenv("OPENAI_API_KEY"); envKey != "" && providerName == "openai" {
		apiKey = envKey
	}

	provider, err := ai.NewProvider(providerName, apiKey, model, ollamaURL)
	if err != nil {
		return ""
	}

	prompt := ai.BuildPrompt(results, since.Format("2006-01-02"))

	fmt.Printf("📝 AI 요약 생성 중 (%s)...\n", provider.Name())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	text, err := provider.Summarize(ctx, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ AI 요약 실패: %v\n", err)
		return ""
	}

	return text
}
