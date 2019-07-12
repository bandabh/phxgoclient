package phxgoclient

import "testing"

func setupPhxGo() PheonixGoSocket {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)
	socket.Listen()
	socket.OpenChannel("room:lobby")
	socket.JoinChannel("room:lobby", nil)
	return socket
}

func TestChannel_Push(t *testing.T) {
	socket := setupPhxGo()
	channel, _ := socket.GetChannel("room:lobby")

	if !channel.CanPush() {
		t.Errorf("channel should be able to push")
	}

	channel.Push("ping", nil)
	readEvent := channel.Read()

	if readEvent.Payload.Status != "" {
		if readEvent.Payload.Status != "ok" {
			t.Errorf("failed to push to channel")
		}
	}
}

func TestChannel_Read(t *testing.T) {
	socket := setupPhxGo()
	channel, _ := socket.GetChannel("room:lobby")

	if !channel.CanPush() {
		t.Errorf("channel should be able to push")
	}

	channel.Push("ping", nil)
	readEvent := channel.Read()

	if readEvent.Event == "" {
		t.Errorf("failed to read!")
	}
}

//func TestChannel_Observe(t *testing.T) {
//	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)
//	socket.Listen()
//
//	socket.OpenChannel("room:lobby")
//
//	channel, _ := socket.GetChannel("room:lobby")
//	channel.Register(ReplyEvent, func(response MessageResponse) (data interface{}, err error) {
//
//		if response.Payload.Status != "ok" {
//			t.Errorf("failed to connect properly")
//		}
//
//		return response, nil
//	})
//	socket.JoinChannel("room:lobby", nil)
//
//	_, err := channel.Observe()
//
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//
//}

func TestChannel_IsOpen(t *testing.T) {

	socket := setupPhxGo()
	channel, _ := socket.GetChannel("room:lobby")

	if !channel.IsOpen() {
		t.Errorf("channel should be open")
	}
}

func TestChannel_CanPush(t *testing.T) {
	socket := setupPhxGo()
	channel, _ := socket.GetChannel("room:lobby")

	if !channel.CanPush() {
		t.Errorf("channel should be able to push")
	}
}
