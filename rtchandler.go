package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type socketUpgrader struct {
	websocket.Upgrader

	OnChannel ChannelHandler
}

func (upgrader *socketUpgrader) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	peer, err := Peer(func(signal string) {
		fmt.Println("Sending signal to client")
		conn.WriteMessage(websocket.TextMessage, []byte(signal))
	}, upgrader.OnChannel)

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	peer.CreateDataChannel("test")

	fmt.Println("Client subscribed")

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		switch msgType {
		case websocket.CloseMessage:
			return
		case websocket.TextMessage:
			fmt.Println("Signal received from client:", string(msg))
			err = peer.Connect(string(msg))
		default:
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
