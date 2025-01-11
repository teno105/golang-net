아래는 실습 순서에 맞춰 다시 작성한 `README.md`입니다.

---

# echoserver

`echoserver`는 Golang으로 작성된 간단한 Echo 애플리케이션으로, 기본적인 client-server 통신하는 방법을 익히기 위한 실습입니다.


## 프로젝트 구조

```plaintext
golang-net/
│
├── cmd/
│   └── server/
│       └── server.go
│
├── cmd/
│   └── client/
│       └── clinet.go
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
mkdir golang-net
cd golang-net
go mod init server

mkdir -p cmd/server
```

### 2. `Makefile` 작성

이제 프로젝트의 빌드 및 실행을 자동화하기 위한 `Makefile`을 프로젝트 루트에 작성합니다.

```makefile
# Go 관련 변수 설정
APP_NAME := server
CMD_DIR := ./cmd/server
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

4.`server.go`을 작성 후, `Makefile`을 이용하여 코드를 빌드하고 실행할 수 있습니다.

```bash
make run
```

### 3. gnet 설치
 gnet은 TCP/UDP와 같은 기본적인 네트워크 프로토콜을 이용한 서버 프로그램을 만들 수 있게 도와주는 오픈 소스 패키지 입니다.
 gnet은 사용이 쉽고 고성능을 자랑합니다.

 gnet 설치하기
```bash
go get -u github.com/panjf2000/gnet/v2@v2.3.0
```

### 4. `server.go` 생성

`cmd/server/` 디렉터리 아래에 `server.go` 파일을 생성하고, gnet 을 사용한 Echo 서버 코드를 작성합니다.

```go
// cmd/server/server.go
package main

import (
	"flag"
	"fmt"
	"log"

	gnet "github.com/panjf2000/gnet/v2"
)

type echoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *echoServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("client connected. address:%s", c.RemoteAddr().String())
	return nil, gnet.None
}

func (es *echoServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("clinet disconnected. address:%s",
		c.RemoteAddr().String())
	return gnet.None
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n",
		es.multicore, es.addr)
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

func main() {
	var port int
	var multicore bool

	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	echo := &echoServer{
		addr:      fmt.Sprintf("tcp://:%d", port),
		multicore: multicore,
	}
	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))
}
```

이 코드를 통해 프로그램을 실행하면 `TCP port:9000`을 사용하는 서버가 실행됩니다.




### 5. `client.go` 생성

`cmd/client/` 디렉터리 아래에 `client.go` 파일을 생성하고, net 을 사용한 클라이언트 코드를 작성합니다.

```go
// cmd/client/client.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	var port int
	var addr string

	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.StringVar(&addr, "address", "localhost", "--address localhost")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	go func() {
		scan := bufio.NewScanner(conn)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
	}()

	for {
		inputScan := bufio.NewScanner(os.Stdin)
		inputScan.Split(bufio.ScanLines)
		for inputScan.Scan() {
			if inputScan.Text() == "exit" {
				return
			}
			conn.Write([]byte(fmt.Sprintf("%s\n", inputScan.Text())))
		}
	}
}
```
```bash
go run client.go
```
이 코드를 통해 프로그램을 실행하면 서버에 메세지를 보내는 클라이언트가 실행됩니다.
