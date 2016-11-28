package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	webrtc "github.com/keroserene/go-webrtc"
)

func main() {
	var upgrader = socketUpgrader{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},

		OnChannel: func(channel *webrtc.DataChannel) {
			channel.Send([]byte("Hello"))
		},
	}

	http.HandleFunc("/connect", upgrader.Handler)

	http.Handle("/", http.RedirectHandler("/static", http.StatusSeeOther))
	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	bind := "127.0.0.1:3000"
	fmt.Println("listening on:", bind)
	http.ListenAndServe(bind, nil)
}
