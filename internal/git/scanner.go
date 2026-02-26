package git

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanRepos는 주어진 경로들에서 .git 디렉토리를 찾아 레포 경로 목록을 반환한다.
// 1단계 깊이만 탐색 (scan_path/project/.git)
func ScanRepos(scanPaths []string, excludes []string) ([]string, error) {
	var repos []string
	seen := make(map[string]bool)

	for _, sp := range scanPaths {
		expanded := expandHome(sp)

		// 해당 경로 자체가 git 레포인 경우
		if isGitRepo(expanded) {
			abs, _ := filepath.Abs(expanded)
			if !seen[abs] {
				repos = append(repos, abs)
				seen[abs] = true
			}
			continue
		}

		// 하위 디렉토리 탐색
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			name := entry.Name()

			if shouldExclude(name, excludes) {
				continue
			}

			fullPath := filepath.Join(expanded, name)
			abs, _ := filepath.Abs(fullPath)

			if isGitRepo(abs) && !seen[abs] {
				repos = append(repos, abs)
				seen[abs] = true
			}
		}
	}

	return repos, nil
}

func isGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && info.IsDir()
}

func shouldExclude(name string, excludes []string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}
	for _, ex := range excludes {
		if name == ex {
			return true
		}
	}
	return false
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
