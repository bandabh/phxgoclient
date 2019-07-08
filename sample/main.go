package main

import (
	"github.com/bandabh/phxgoclient"
)

func main() {
	socket := phxgoclient.NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	socket.Listen()

	channel, err := socket.OpenChannel("room:lobby")
	socket.JoinChannel("room:lobby", nil)

	channel.Register(phxgoclient.ReplyEvent, func(response phxgoclient.MessageResponse) (data interface{}, err error) {

		println("---------")
		println(response.Event)
		println(response.Topic)
		println("---------")
		return response, nil
	})

	channel.Register(phxgoclient.ErrorEvent, func(response phxgoclient.MessageResponse) (data interface{}, err error) {

		println(response.Event)

		return response, nil
	})

	if err != nil {
		println(err)
	}

	channel.Push("ping", nil)
	channel.Push("ping", nil)

	channel.Observe()

	channel.Push("ping", nil)

	for {
		_ = ""
	}

}