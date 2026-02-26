package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/kso1204/gitday/internal/git"
	"github.com/kso1204/gitday/internal/notify"
	"github.com/kso1204/gitday/internal/output"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "리포트를 Slack으로 전송",
	RunE:  runSend,
}

func init() {
	sendCmd.Flags().Bool("slack", false, "Slack 웹훅으로 전송")
	sendCmd.Flags().String("period", "today", "기간: today, week")
	rootCmd.AddCommand(sendCmd)
}

func runSend(cmd *cobra.Command, args []string) error {
	useSlack, _ := cmd.Flags().GetBool("slack")
	if !useSlack {
		return fmt.Errorf("전송 대상을 지정하세요 (예: --slack)")
	}

	webhookURL := viper.GetString("slack.webhook_url")
	if webhookURL == "" {
		return fmt.Errorf("Slack webhook URL이 설정되지 않았습니다. ~/.gitday.yaml에서 설정하세요")
	}

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
		fmt.Println("전송할 커밋이 없습니다.")
		return nil
	}

	md := output.ToMarkdown(results, since, now, "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := notify.SendSlack(ctx, webhookURL, md); err != nil {
		return err
	}

	fmt.Println("✓ Slack 전송 완료!")
	return nil
}
