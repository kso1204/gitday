package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "이번 주 Git 활동 요약",
	RunE:  runWeek,
}

func init() {
	rootCmd.AddCommand(weekCmd)
}

func runWeek(cmd *cobra.Command, args []string) error {
	now := time.Now()

	// 이번 주 월요일 00:00
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 일요일
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	since := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())

	return runReport(since, now)
}
