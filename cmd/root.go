package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "gitday",
	Short:   "Git 데일리 활동 요약 도구",
	Long:    `멀티레포 Git 로그를 스캔해서 오늘의 작업 내용을 보여주고, AI로 자연어 요약하는 CLI 도구.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 기본 동작 = today (순환 참조 방지를 위해 init에서 설정)
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runToday(cmd, args)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "설정 파일 경로 (기본: ~/.gitday.yaml)")
	rootCmd.PersistentFlags().String("author", "", "Git 저자 필터")
	rootCmd.PersistentFlags().Bool("summary", false, "AI 요약 포함")
	rootCmd.PersistentFlags().Bool("compact", false, "간략 출력 모드")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("output.compact", rootCmd.PersistentFlags().Lookup("compact"))
	viper.BindPFlag("summary", rootCmd.PersistentFlags().Lookup("summary"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".gitday")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("GITDAY")
	viper.AutomaticEnv()

	// 기본값
	viper.SetDefault("scan_paths", []string{"."})
	viper.SetDefault("exclude", []string{"node_modules", "vendor", ".cache", ".venv"})
	viper.SetDefault("ai.provider", "claude")
	viper.SetDefault("ai.ollama_url", "http://localhost:11434")
	viper.SetDefault("output.color", true)
	viper.SetDefault("output.compact", false)

	viper.ReadInConfig()
}
