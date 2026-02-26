package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/wook/gitday/internal/git"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	dateStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	repoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))

	hashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3"))

	msgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	statStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	summaryBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true)
)

func PrintReport(results []git.RepoResult, since, until time.Time, compact bool) {
	// 헤더
	weekday := weekdayKo(since.Weekday())
	header := fmt.Sprintf("📅 %s (%s)", since.Format("2006-01-02"), weekday)
	fmt.Println(titleStyle.Render(header))
	fmt.Println()

	totalCommits := 0
	totalFiles := 0

	for _, r := range results {
		commitCount := len(r.Commits)
		totalCommits += commitCount

		fileCount := 0
		for _, c := range r.Commits {
			fileCount += c.Files
		}
		totalFiles += fileCount

		// 레포 헤더
		repoHeader := fmt.Sprintf("━━ %s (%d commits) ", r.Name, commitCount)
		padding := 50 - len(repoHeader)
		if padding < 3 {
			padding = 3
		}
		repoHeader += strings.Repeat("━", padding)
		fmt.Println(repoStyle.Render(repoHeader))

		// 커밋 목록
		if compact {
			// 간략: 첫 3개만
			limit := 3
			if commitCount < limit {
				limit = commitCount
			}
			for _, c := range r.Commits[:limit] {
				printCommit(c)
			}
			if commitCount > 3 {
				fmt.Printf("  %s\n", statStyle.Render(fmt.Sprintf("... +%d more", commitCount-3)))
			}
		} else {
			for _, c := range r.Commits {
				printCommit(c)
			}
		}
		fmt.Println()
	}

	// 하단 통계
	bar := fmt.Sprintf("📊 총 %d commits | %d개 프로젝트 | %d files changed",
		totalCommits, len(results), totalFiles)
	fmt.Println(summaryBarStyle.Render(bar))
}

func printCommit(c git.Commit) {
	hash := hashStyle.Render(c.Hash)
	msg := msgStyle.Render(c.Message)

	if c.Files > 0 {
		stat := statStyle.Render(fmt.Sprintf("(%d files)", c.Files))
		fmt.Printf("  %s %s %s\n", hash, msg, stat)
	} else {
		fmt.Printf("  %s %s\n", hash, msg)
	}
}

var summaryTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("14"))

var summaryHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("13"))

func PrintSummary(text string) {
	fmt.Println()
	fmt.Println(summaryHeaderStyle.Render("📝 오늘의 요약"))
	fmt.Println(summaryTextStyle.Render(text))
}

func weekdayKo(w time.Weekday) string {
	days := [...]string{"일", "월", "화", "수", "목", "금", "토"}
	return days[w]
}
