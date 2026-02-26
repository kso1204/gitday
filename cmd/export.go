package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/kso1204/gitday/internal/git"
	"github.com/kso1204/gitday/internal/output"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "리포트를 마크다운으로 출력/저장",
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().StringP("output", "o", "", "출력 파일 경로 (미지정 시 stdout)")
	exportCmd.Flags().String("period", "today", "기간: today, week")
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	period, _ := cmd.Flags().GetString("period")
	outputPath, _ := cmd.Flags().GetString("output")

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
		fmt.Println("내보낼 커밋이 없습니다.")
		return nil
	}

	md := output.ToMarkdown(results, since, now, "")

	if outputPath == "" {
		fmt.Print(md)
		return nil
	}

	if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
		return fmt.Errorf("파일 저장 실패: %w", err)
	}

	fmt.Printf("✓ 리포트 저장됨: %s\n", outputPath)
	return nil
}
