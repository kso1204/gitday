# gitday

멀티레포 Git 로그를 스캔해서 오늘의 작업 내용을 보여주고, AI로 자연어 요약하는 CLI 도구.

```
$ gitday

📅 2026-02-26 (목)

━━ rpg (3 commits) ━━━━━━━━━━━━━━━━━━━━━━━━━━━
  dbc7067 기능 잠금 + 뽑기 숨기기: 배포 준비 (16 files)
  7ced58e 정수 시스템 + 전투/타워 세션 DB 영속화 (11 files)
  5558aab 보스 클리어 기록을 localStorage → DB로 이전 (4 files)

━━ petition (2 commits) ━━━━━━━━━━━━━━━━━━━━━━
  d06e651 테스트: 오디오 파이프라인 테스트 스크립트 추가 (1 files)
  41f57c1 문서: 프로젝트 리뷰 및 분석 보고서 추가 (8 files)

📊 총 5 commits | 2개 프로젝트 | 40 files changed

📝 오늘의 요약
RPG에서 전투 시스템과 데이터 영속화 작업을 집중적으로 했고,
petition은 테스트와 문서화 위주로 진행했습니다.
```

## 설치

### Homebrew (macOS/Linux)

```bash
brew install wook/tap/gitday
```

### Go Install

```bash
go install github.com/wook/gitday@latest
```

### 바이너리 다운로드

[GitHub Releases](https://github.com/wook/gitday/releases)에서 플랫폼에 맞는 바이너리를 다운로드하세요.

## 시작하기

```bash
# 설정 파일 생성
gitday init

# ~/.gitday.yaml에서 scan_paths 수정
# scan_paths:
#   - ~/Documents/projects
#   - ~/work

# 오늘 커밋 보기
gitday today

# AI 요약 포함
gitday today --summary

# 이번 주 커밋
gitday week
```

## 사용법

```bash
# 기본 (= today)
gitday                          # 오늘 커밋 로그
gitday today                    # 동일
gitday today --summary          # + AI 요약
gitday today --compact          # 간략 모드

# 기간
gitday week                     # 이번 주
gitday week --summary           # + AI 요약

# 필터
gitday --author "wook"          # 특정 저자만

# 내보내기
gitday export                   # 마크다운으로 stdout
gitday export -o report.md      # 파일 저장
gitday export --period week     # 주간 리포트

# 전송
gitday send --slack             # Slack 웹훅 전송

# 설정
gitday init                     # ~/.gitday.yaml 초기화
```

## 설정

`~/.gitday.yaml`:

```yaml
# 스캔 대상 디렉토리
scan_paths:
  - ~/Documents/home
  - ~/work

# 제외 패턴
exclude:
  - node_modules
  - vendor
  - .cache

# Git 저자 (비워두면 전체)
author: ""

# AI 설정
ai:
  provider: claude    # claude | openai | ollama
  api_key: ""         # 환경변수 GITDAY_API_KEY 우선
  model: ""           # 비워두면 기본값 사용
  ollama_url: "http://localhost:11434"

# Slack
slack:
  webhook_url: ""

# 출력
output:
  color: true
  compact: false
```

### AI 프로바이더

| 프로바이더 | 기본 모델 | API 키 |
|-----------|----------|--------|
| Claude | claude-haiku-4-5 | `GITDAY_API_KEY` 또는 설정 파일 |
| OpenAI | gpt-4o-mini | `GITDAY_API_KEY` 또는 설정 파일 |
| Ollama | llama3.2 | 불필요 (로컬) |

```bash
# 환경변수로 API 키 설정
export GITDAY_API_KEY="your-api-key"

# 또는 설정 파일에 직접 입력
# ai:
#   api_key: "your-api-key"
```

## 라이선스

MIT
