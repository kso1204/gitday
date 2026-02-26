package git

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Commit은 단일 커밋 정보를 나타낸다.
type Commit struct {
	Hash    string
	Message string
	Author  string
	Date    time.Time
	Files   int // 변경된 파일 수
}

// RepoResult는 단일 레포의 커밋 수집 결과이다.
type RepoResult struct {
	Name    string
	Path    string
	Commits []Commit
}

// CollectLogs는 여러 레포에서 병렬로 커밋 로그를 수집한다.
func CollectLogs(repos []string, since, until time.Time, author string) ([]RepoResult, error) {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []RepoResult
	)

	for _, repo := range repos {
		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()

			commits, err := getCommits(repoPath, since, until, author)
			if err != nil || len(commits) == 0 {
				return
			}

			mu.Lock()
			results = append(results, RepoResult{
				Name:    filepath.Base(repoPath),
				Path:    repoPath,
				Commits: commits,
			})
			mu.Unlock()
		}(repo)
	}

	wg.Wait()
	return results, nil
}

const separator = "§§"

func getCommits(repoPath string, since, until time.Time, author string) ([]Commit, error) {
	format := "%H" + separator + "%s" + separator + "%an" + separator + "%aI"
	args := []string{
		"log",
		"--all",
		"--format=" + format,
		"--since=" + since.Format(time.RFC3339),
		"--until=" + until.Format(time.RFC3339),
		"--shortstat",
	}

	if author != "" {
		args = append(args, "--author="+author)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseGitLog(string(out))
}

func parseGitLog(raw string) ([]Commit, error) {
	var commits []Commit
	lines := strings.Split(strings.TrimSpace(raw), "\n")

	var current *Commit
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 커밋 라인: hash§§message§§author§§date
		if strings.Contains(line, separator) {
			parts := strings.Split(line, separator)
			if len(parts) < 4 {
				continue
			}

			date, _ := time.Parse(time.RFC3339, parts[3])
			c := Commit{
				Hash:    truncate(parts[0], 7),
				Message: parts[1],
				Author:  parts[2],
				Date:    date,
			}
			commits = append(commits, c)
			current = &commits[len(commits)-1]
			continue
		}

		// shortstat 라인: " 3 files changed, 45 insertions(+), 12 deletions(-)"
		if current != nil && strings.Contains(line, "file") {
			current.Files = parseFileCount(line)
			current = nil
		}
	}

	return commits, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func parseFileCount(stat string) int {
	parts := strings.Fields(stat)
	if len(parts) > 0 {
		n, err := strconv.Atoi(parts[0])
		if err == nil {
			return n
		}
	}
	return 0
}
