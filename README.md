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

## 실습 순서

### 1. 패키지 구조를 위한 디렉토리 생성

먼저 프로젝트 디렉터리를 설정하고 필요한 디렉터리들을 생성합니다.

```bash
mkdir golang-net
cd golang-net
go mod init server

mkdir -p cmd/server
mkdir -p cmd/client
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

### 30.1 gnet을 이용해서 echo 서버 제작
 gnet은 TCP/UDP와 같은 기본적인 네트워크 프로토콜을 이용한 서버 프로그램을 만들 수 있게 도와주는 오픈 소스 패키지 입니다.<br/>
 gnet은 사용이 쉽고 고성능을 자랑합니다.

 gnet 설치하기
```bash
go get -u github.com/panjf2000/gnet/v2@v2.3.0
```

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

// 1. OnOpen은 새로운 클라이언트가 접속할 때 호출됩니다.
func (es *echoServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("client connected. address:%s", c.RemoteAddr().String())
	return nil, gnet.None
}

// 2. OnClose는 클라이언트 접속을 해제할 때 호출됩니다.
func (es *echoServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("clinet disconnected. address:%s",
		c.RemoteAddr().String())
	return gnet.None
}

// 3. 서버가 시작될 때 호출됩니다.
func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n",
		es.multicore, es.addr)
	return gnet.None
}

// 4. 서버가 데이터를 네트워크를 통해서 클라이언트로부터 수신할 때 호출됩니다.
func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

func main() {
	var port int
	var multicore bool

	// 5. flag 패키지를 통해 실행 인수를 읽어오게 됩니다.
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	// 6. gnet을 통해서 서버를 실행하는 부분입니다.
	echo := &echoServer{
		addr:      fmt.Sprintf("tcp://:%d", port),
		multicore: multicore,
	}
	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))
}
```

이 코드를 통해 프로그램을 실행하면 `TCP port:9000`을 사용하는 서버가 실행됩니다.

6. gnet.Run() 함수를 호출해서 서버를 싱행합니다. Run()함수는 다음과 같이 세가지 인수를 받습니다.
```go
Run(eventHandler EventHandler, protoAddr string, opts ...Option)
```
첫번째 인수는 eventHandler입니다. 이것은 EventHandler 인터페이스 타입으로 gnet 내부에서 이벤트들이 클라이언트 접속, 연결 해제, 데이터 수신 등이 발생하면 그게 맞는 메서드를 호출해주게 됩니다.<br/>
따라서 EventHandler 인터페이스를 통해서 서버에서 데이터를 수신하거나 클라이언트 접속/해제를 알 수 있습니다.<br/>
두번째 인수는 protoAddr 입니다. 이것은 서버가 어떤 프로토콜을 통해서 통신하게 되는지, 또 그 주소는 어디에 바인딩하는지를 나타냅니다.<br/>
다음과 같이 하면 TCP프로토콜을 사용하고 IP주소 192.168.0.10 주소에 9851 포트 번호에 바인딩하게 됩니다.
```go
"tcp://192.168.0.10:9851"
```
다음과 같이 설정하면 서버 머신에서 바인딩 가능한 모든 네트워크 인터페이스의 9851 포트에 바인딩하게 됩니다.<br/>
클라이언트는 서버의 네트워크 인터페이스 중 어떤 곳에서라도 9851 포트를 통해서 통신을 할 수 있습니다.
```go
"tcp://:9851"
```
아래와 같이 설정하면 UDP 프로토콜과 IPv4 프로토콜을 사용하고 9851 포트에 바인딩 하게 됩니다.
```go
"udp4://:9851"
```
세번째는 옵션 리스트를 정하게 됩니다. gnet은 다양한 옵션을 제공하는데 여기서는 Multicore 옵션만 적용했습니다.<br/>
Multicore 옵션을 적용하면 서버가 여러 CPU 코어를 사용해서 동작하게 됩니다. <br/>
더 빠른 성능을 가지게 되지만, 멀티 쓰레딩 환경에서 동작하기 때문에 메모리 자원 점유에 대해서 주의해야 합니다.

### 30.2 클라이언트 제작

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

	// 1. 실행 인수를 통해 접속하려는 ip, port를 설정합니다.
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.StringVar(&addr, "address", "localhost", "--address localhost")
	flag.Parse()

	// 2. net 패키지를 이용해서 tcp 연결을 맺습니다.
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed:", err)
	}

	// 3. 연결된 conn을 통해 데이터를 읽어서 출력합니다.
	go func() {
		scan := bufio.NewScanner(conn)
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			fmt.Println(scan.Text())
		}
	}()

	// 4. 키보드로부터 텍스트를 입력받아 데이터를 전송합니다.
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
이 코드를 통해 프로그램을 실행하면 서버에 메세지를 보내는 클라이언트가 실행됩니다.<br/>
1. flag 패키기를 이용해서 실행 인수를 통해 접속하려는 곳의 ip 주소와 port 번호를 설정합니다.<br/>
2. net 패키기를 이용해서 서버에 접속합니다. net.ResolveTCPAddr 함수는 ip주소와 port번호를 통해서 *TcpAddr 객체를 반환합니다.<br/>
이렇게 만든 net.DialTCP 함수를 통해서 앞서만든 *TCPAddr가 가리키는 주소로 접속하고 접속된 연결을 반환합니다.<br/>
3. 연결된 net.Coon 객체에서 데이터를 읽어서 출력하는 고루틴을 실행합니다.<br/>
4. 키보드로 텍스트를 한 줄 입력받아 연결된 net.Conn 객체를 통해 데이터를 송신합니다.<br/>


### 30.3 채팅 서버 제작

echo 서버를 기반으로 채팅 서버를 만들겠습니다.<br/>
채팅 서버는 연결된 모든 클라이언트들에게 모두 전송을 합니다.<br/>
이렇게 여러 클라이언트에게 전송하는 것을 방송한다고 말하고 브로드캐스트(broadcast)라고 말하기도 합니다.<br/>
클라이언트가 보낸 데이터를 모든 클라이언트들에게 브로드캐스트하는 채팅 서버를 만들어 보겠습니다.

```go
// cmd/server/server.go
package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	gnet "github.com/panjf2000/gnet/v2"
)

type chatServer struct {
	gnet.BuiltinEventEngine

	// 1. 연결된 커넥션을 보관하는 맵
	cliMap sync.Map
}

func (cs *chatServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("client connected. address:%s", c.RemoteAddr().String())
	// 2. 새로운 연결이 되면 커넥션을 맵에 보관한다.
	cs.cliMap.Store(c, true)
	return nil, gnet.None
}

func (cs *chatServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("clinet disconnected. address:%s", c.RemoteAddr().String())
	// 3. 연결이 해제되면 맵에서 삭제한다.
	if _, ok := cs.cliMap.LoadAndDelete(c); ok {
		log.Printf("connection removed")
	}
	return gnet.None
}

func (cs *chatServer) OnBoot(eng gnet.Engine) gnet.Action {
	log.Printf("chat server is listening\n")
	return gnet.None
}

func (cs *chatServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	// 4. 데이터 수신 시 모든 커넥션에 데이터를 전송한다.
	cs.cliMap.Range(func(key, value any) bool {
		if conn, ok := key.(gnet.Conn); ok {
			conn.AsyncWrite(buf, nil)
		}
		return true
	})
	return gnet.None
}

func main() {
	var port int
	var multicore bool

	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()

	echo := &chatServer{}
	log.Fatal(gnet.Run(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore)))
}
```
1. 연결된 net.Conn 객제들을 보관하는 맵입니다. 멀티 코어로 동작할수 있기 때문에 동시성 프로그래밍에 적합한 sync.Map 객체로 생성했습니다.<br/>
2. 새로운 연결이 맺어지면 맵에 보관합니다. sync.Map이기 때문에 별도로 Lock을 잡을 필요가 없습니다.<br/>
3. 연결이 해제되면 맵에서 삭제합니다.<br/>
4. 데이터 수신시 모든 연결에게 데이터를 브로드캐스트합니다. 이렇게 해서 하나의 클라이언트가 보낸 메시지가 모든 클라이언트에게 보여질 수 있는 채팅 서버를 만들게 됩니다.<br/>


위 코드를 작성 후, build 후에 서버에 접속을 합니다.
```bash
cd cmd/client
go build client.go
```
### 실행화면
![스크린샷 2025-01-12 오후 7 14 16](https://github.com/user-attachments/assets/ac37d156-a4ee-4b01-992d-296be51f953e)
