package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "오늘의 활동을 AI 요약 + 로컬 저장 (= gitday --summary)",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("summary", true)
		return runToday(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
