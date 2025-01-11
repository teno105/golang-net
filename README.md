아래는 실습 순서에 맞춰 다시 작성한 `README.md`입니다.

---

# echoserver

`echoserver`는 Golang으로 작성된 간단한 애플리케이션으로, 기본적인 Application의 구조와 테스트 방법을 익히기 위한 실습입니다.


## 프로젝트 구조

```plaintext
echoserver/
│
├── cmd/
│   └── echoserver/
│       └── main.go
│
├── go.mod
├── Makefile
└── README.md
```

## 요구 사항

- **Go 언어**가 설치되어 있어야 합니다. [Go 설치 가이드](https://golang.org/doc/install)
- **Git**이 설치되어 있어야 합니다.

## 실습 순서

### 1. 패키지 구조를 위한 디렉토리 생성

먼저 프로젝트 디렉터리를 설정하고 필요한 디렉터리들을 생성합니다.

```bash
mkdir echoserver
cd echoserver
go mod init echoserver

mkdir -p cmd/echoserver
```

### 2. `main.go` 생성

`cmd/echoserver/` 디렉터리 아래에 `main.go` 파일을 생성하고, 간단한 `Hello, World!` 메시지를 출력하는 코드를 작성합니다.

```go
// cmd/echoserver/main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

이 코드를 통해 프로그램을 실행하면 `Hello, World!`가 출력됩니다.

### 3. `Makefile` 작성

이제 프로젝트의 빌드 및 실행을 자동화하기 위한 `Makefile`을 프로젝트 루트에 작성합니다.

```makefile
# Go 관련 변수 설정
APP_NAME := echoserver
CMD_DIR := ./cmd/echoserver
BUILD_DIR := ./build

.PHONY: all clean build run test fmt vet install

all: build

# 빌드 명령어
build:
	@echo "==> Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# 실행 명령어
run: build
	@echo "==> Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

# 코드 포맷팅
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

# 코드 분석
vet:
	@echo "==> Running go vet..."
	go vet ./...

# 의존성 설치
install:
	@echo "==> Installing dependencies..."
	go mod tidy

# 테스트 실행
test:
	@echo "==> Running tests..."
	go test -v ./...

# 빌드 정리
clean:
	@echo "==> Cleaning build directory..."
	rm -rf $(BUILD_DIR)
```

`Makefile`을 이용하여 코드를 빌드하고 실행할 수 있습니다.

```bash
make run
```

이 명령어를 통해 `main.go`에서 작성한 `Hello, World!` 메시지를 확인할 수 있습니다.
