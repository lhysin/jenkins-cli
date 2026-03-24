# Jenkins CLI

<p align="center">
  <img src="https://img.shields.io/github/stars/lhysin/jenkins-cli?style=flat-square" alt="GitHub Stars">
  <img src="https://img.shields.io/github/license/lhysin/jenkins-cli?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/v/release/lhysin/jenkins-cli?style=flat-square" alt="Release">
</p>

> Jenkins REST API용 명령줄 인터페이스입니다. 터미널에서 jobs, builds, nodes 등을 관리하세요.

[English](README.md) | **한국어**

## ✨ 주요 기능

- 🔗 **멀티호스트 관리** - 여러 Jenkins 서버에 로그인하고 전환
- 📋 **Job 관리** - 목록, 생성, 수정, 삭제
- 🚀 **Build 관리** - 빌드 목록, 로그 확인, 빌드 실행
- 👁️ **View 관리** - View 목록 및[View]]내 job 확인
- 🖥️ **Node 관리** - 에이전트/노드 목록
- 🔓 **익명 접근** - 공개 Jenkins 인스턴스의 경우 인증 없이 사용 가능

## 📦 설치

### 빠른 설치 (macOS/Linux)

#### 방법 1: 시스템 전역 (sudo 필요)

```bash
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/install.sh | sudo sh
```

#### 방법 2: 사용자 디렉토리 (sudo 불필요)

```bash
mkdir -p ~/.local/bin
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/install.sh | DEST=~/.local/bin/jenkins-cli sh
```

셸 프로파일(`.zshrc`, `.bashrc` 등)에 추가:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### 제거

```bash
curl -fsSL https://raw.githubusercontent.com/lhysin/jenkins-cli/main/uninstall.sh | sh
```

### 바이너리 다운로드

```bash
# 시스템 전역 (sudo 필요)
curl -fsSL https://github.com/lhysin/jenkins-cli/releases/latest/download/jenkins-cli_darwin_arm64 -o /usr/local/bin/jenkins-cli
chmod +x /usr/local/bin/jenkins-cli

# 사용자 디렉토리 (sudo 불필요)
mkdir -p ~/.local/bin
curl -fsSL https://github.com/lhysin/jenkins-cli/releases/latest/download/jenkins-cli_darwin_arm64 -o ~/.local/bin/jenkins-cli
chmod +x ~/.local/bin/jenkins-cli
```

[GitHub Releases](https://github.com/lhysin/jenkins-cli/releases)에서 모든 플랫폼 다운로드

### 소스코드 빌드

```bash
git clone https://github.com/lhysin/jenkins-cli.git
cd jenkins-cli
go build -o jenkins-cli ./cmd/
```

## 🚀 빠른 시작

```bash
# Jenkins 서버에 로그인
jenkins-cli login local http://localhost:8080 admin
# Token: your-api-token

# 호스트 전환
jenkins-cli use prod

# 호스트 목록
jenkins-cli hosts

# Job 목록
jenkins-cli jobs list

# 빌드 실행
jenkins-cli builds trigger my-job

# 빌드 로그 확인
jenkins-cli builds logs my-job 42
```

## 📚 명령어

### 🔗 호스트 관리

| 명령어 | 설명 |
|--------|------|
| `login <이름> <url> [user]` | Jenkins에 로그인 (토큰 입력 요청, 익명 가능) |
| `use <이름>` | 호스트 선택 |
| `logout <이름>` | 호스트 삭제 |
| `hosts` | 호스트 목록 |

### 📋 Jobs

| 명령어 | 설명 |
|--------|------|
| `jobs list` | Job 목록 |
| `jobs info <이름>` | Job 상세 |
| `jobs create <이름> [-c 설정]` | Job 생성 |
| `jobs delete <이름>` | Job 삭제 |
| `jobs update <이름> [-c 설정]` | Job 수정 |

### 🚀 Builds

| 명령어 | 설명 |
|--------|------|
| `builds list <job>` | 빌드 목록 |
| `builds info <job> <번호>` | 빌드 상세 |
| `builds logs <job> <번호>` | 빌드 로그 |
| `builds trigger <job>` | 빌드 실행 |

### 👁️ Views

| 명령어 | 설명 |
|--------|------|
| `views list` | View 목록 |
| `views info <이름>` | View 상세 및[View]]내 job |

### 🖥️ 기타

| 명령어 | 설명 |
|--------|------|
| `nodes list` | 노드 목록 |
| `status` | 연결 상태 확인 |

## ⚙️ 설정

설정 파일: `~/.jenkins-cli/config.yaml`:

```yaml
current: local
hosts:
  local:
    url: http://localhost:8080
    user: admin
    token: your-api-token
  prod:
    url: https://jenkins.prod.com
    user: admin
    token: prod-token
```

## 🔑 API Token 생성

1. Jenkins에 로그인
2. 사용자 이름 클릭 → 설정(Configure)
3. "API Token" 섹션에서 "Add new Token" 클릭
4. 생성된 토큰 복사

## 🧪 테스트

```bash
# 단위 테스트
go test ./...

# 통합 테스트 (Docker Jenkins 필요)
./scripts/integration-test.sh
```

## 🤝 기여하기

### 커밋 메시지 형식

버전이 커밋 메시지에 따라 자동 증가됩니다:

| 타입 | 설명 | 버전 변경 |
|------|------|-----------|
| `fix:` | 버그 수정 | Patch +1 |
| `feat:` | 새 기능 | Minor +1 |
| `perf:` | 성능 개선 | Patch +1 |
| `BREAKING CHANGE:` | BREAKING 변경 | Major +1 |

예시:
```bash
git commit -m "fix: 로그인 타임아웃 문제 해결"
git commit -m "feat: 파이프라인 지원 추가"
git commit -m "feat!: Node.js 16 지원 중단"
```

### 작업 흐름

1. Feature 브랜치 생성: `git checkout -b feature/my-feature`
2. 올바른 접두사로 커밋
3. 푸시 후 PR 생성 → main으로 머지
4. main 머지 → 자동 릴리즈

## 📄 라이선스

MIT License - 자세한 내용은 [LICENSE](LICENSE)를 참조하세요.
