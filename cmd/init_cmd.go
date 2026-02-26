package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "설정 파일 초기화 (~/.gitday.yaml)",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

const defaultConfig = `# gitday 설정 파일
# https://github.com/wook/gitday

# 스캔 대상 디렉토리
scan_paths:
  - ~/Documents/home

# 제외 패턴
exclude:
  - node_modules
  - vendor
  - .cache
  - .venv

# Git 저자 (비워두면 git config user.name 사용)
author: ""

# AI 설정
ai:
  provider: claude  # claude | openai | ollama
  api_key: ""       # 환경변수 GITDAY_API_KEY 우선
  model: ""         # 비워두면 기본값 사용
  ollama_url: "http://localhost:11434"

# Slack
slack:
  webhook_url: ""

# 출력
output:
  color: true
  compact: false
`

func runInit(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("홈 디렉토리 확인 실패: %w", err)
	}

	configPath := filepath.Join(home, ".gitday.yaml")

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠ 이미 설정 파일이 존재합니다: %s\n", configPath)
		fmt.Print("덮어쓰시겠습니까? (y/N): ")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" {
			fmt.Println("취소되었습니다.")
			return nil
		}
	}

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("설정 파일 생성 실패: %w", err)
	}

	fmt.Printf("✓ 설정 파일 생성됨: %s\n", configPath)
	fmt.Println("  scan_paths를 수정하여 스캔할 디렉토리를 지정하세요.")
	return nil
}
