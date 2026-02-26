package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/wook/gitday/internal/git"
)

// ToMarkdown은 리포트를 마크다운 문자열로 변환한다.
func ToMarkdown(results []git.RepoResult, since, until time.Time, summary string) string {
	var sb strings.Builder

	weekday := weekdayKo(since.Weekday())
	sb.WriteString(fmt.Sprintf("# 📅 %s (%s)\n\n", since.Format("2006-01-02"), weekday))

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

		sb.WriteString(fmt.Sprintf("## %s (%d commits)\n\n", r.Name, commitCount))
		for _, c := range r.Commits {
			if c.Files > 0 {
				sb.WriteString(fmt.Sprintf("- `%s` %s (%d files)\n", c.Hash, c.Message, c.Files))
			} else {
				sb.WriteString(fmt.Sprintf("- `%s` %s\n", c.Hash, c.Message))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("---\n\n📊 **총 %d commits | %d개 프로젝트 | %d files changed**\n",
		totalCommits, len(results), totalFiles))

	if summary != "" {
		sb.WriteString(fmt.Sprintf("\n## 📝 요약\n\n%s\n", summary))
	}

	return sb.String()
}
