package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"code.google.com/p/go.net/websocket"
)

func BotServer(ws *websocket.Conn) {
	var message string
	websocket.Message.Receive(ws, &message)
	io.WriteString(os.Stdout, message)
	websocket.Message.Send(ws, "Starting task.")
	time.Sleep(5 * time.Second)
	websocket.Message.Send(ws, "Task complete.")
	websocket.Message.Send(ws, "done")
}

func main() {
	http.Handle("/bot", websocket.Handler(BotServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
