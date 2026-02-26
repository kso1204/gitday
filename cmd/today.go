package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wook/gitday/internal/ai"
	"github.com/wook/gitday/internal/git"
	"github.com/wook/gitday/internal/output"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "오늘의 Git 활동 요약",
	RunE:  runToday,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

func runToday(cmd *cobra.Command, args []string) error {
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	until := now

	return runReport(since, until)
}

func runReport(since, until time.Time) error {
	scanPaths := viper.GetStringSlice("scan_paths")
	excludes := viper.GetStringSlice("exclude")
	author := viper.GetString("author")

	// 1. 레포 스캔
	repos, err := git.ScanRepos(scanPaths, excludes)
	if err != nil {
		return fmt.Errorf("레포 스캔 실패: %w", err)
	}

	if len(repos) == 0 {
		fmt.Println("스캔된 Git 레포가 없습니다. gitday init으로 scan_paths를 설정하세요.")
		return nil
	}

	// 2. 커밋 로그 수집
	results, err := git.CollectLogs(repos, since, until, author)
	if err != nil {
		return fmt.Errorf("커밋 로그 수집 실패: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("📭 %s ~ %s 기간에 커밋이 없습니다.\n",
			since.Format("2006-01-02"),
			until.Format("2006-01-02 15:04"))
		return nil
	}

	// 3. 터미널 출력
	compact := viper.GetBool("output.compact")
	output.PrintReport(results, since, until, compact)

	// 4. AI 요약 (--summary 플래그)
	summary := viper.GetBool("summary")
	if summary {
		if err := runSummary(results, since); err != nil {
			fmt.Fprintf(os.Stderr, "\n⚠ AI 요약 실패: %v\n", err)
		}
	}

	return nil
}

func runSummary(results []git.RepoResult, since time.Time) error {
	providerName := viper.GetString("ai.provider")
	apiKey := viper.GetString("ai.api_key")
	model := viper.GetString("ai.model")
	ollamaURL := viper.GetString("ai.ollama_url")

	// 환경변수 우선
	if envKey := os.Getenv("GITDAY_API_KEY"); envKey != "" {
		apiKey = envKey
	}

	provider, err := ai.NewProvider(providerName, apiKey, model, ollamaURL)
	if err != nil {
		return err
	}

	prompt := ai.BuildPrompt(results, since.Format("2006-01-02"))

	fmt.Printf("\n📝 AI 요약 생성 중 (%s)...\n", provider.Name())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	text, err := provider.Summarize(ctx, prompt)
	if err != nil {
		return err
	}

	output.PrintSummary(text)
	return nil
}
