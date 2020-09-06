package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"chat-test/network"
)

var closed = false

type testAgent struct {
	*network.WSConn
}

func (t testAgent) Run() {
	defer t.Close()
	for {
		data, err := t.ReadMsg()
		if err != nil {
			break
		}

		msg := string(data)
		fmt.Println(msg)
	}
}

func (t testAgent) OnClose() {
	closed = true
}

func main() {
	var t *testAgent

	wsc := &network.WSClient{
		Addr:             "ws://localhost:8080",
		ConnNum:          1,
		ConnectInterval:  3 * time.Second,
		PendingWriteNum:  10,
		MaxMsgLen:        1024,
		HandshakeTimeout: 1 * time.Second,
		NewAgent: func(conn *network.WSConn) network.Agent {
			t = &testAgent{conn}
			return t
		},
	}
	wsc.Start()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./test/chat.html")
		})
		_ = http.ListenAndServe(":8888", nil)
	}()

	for !closed {
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		t.WriteMsg([]byte(strings.TrimSpace(line)))
	}

	time.Sleep(time.Hour)
	wsc.Close()
}
