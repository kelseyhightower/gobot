package main

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"code.google.com/p/go.net/websocket"
	"github.com/daneharrigan/hipchat"
)

var (
	client   *hipchat.Client
	reply    = make(chan string, 2)
	hub      = make(map[string]chan string)
	fullName = "Go Bot"
	password = ""
	resource = "bot"
	room     = ""
	username = ""
)

func init() {
	var err error
	client, err = hipchat.NewClient(username, password, resource)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Register(name string, c chan string) {
	hub[name] = c
}

func sendMessage(message string) {
	origin := "http://localhost/"
	url := "ws://localhost:8080/bot"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if err := websocket.Message.Send(ws, message); err != nil {
		log.Println(err.Error())
	}
	for {
		var msg string
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
		}
		if msg == "done" {
			break
		}
		reply <- msg
	}
}

func responder() {
	for {
		message := <-reply
		client.Say(room, fullName, message)
	}
}

func main() {
	r := regexp.MustCompile(`GoBot`)
	client.Status("chat")
	client.Join(room, fullName)
	go client.KeepAlive()
	go responder()
	m := client.Messages()
	for {
		message := <-m
		fmt.Println(message)
		if r.MatchString(message.Body) {
			go sendMessage(message.Body)
		}
	}
}
