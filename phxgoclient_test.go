package phxgoclient

import (
	"testing"
	"time"
)

func TestNewPheonixWebsocket(t *testing.T) {

	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	if socket.Status != PhxGoClosed {
		t.Errorf("Socket was created in incorrect state")
	}
}

func TestPheonixGoSocket_Listen(t *testing.T){
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	err := socket.Listen()

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestPheonixGoSocket_JoinChannel(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	socket.OpenChannel("room:lobby")
	err := socket.JoinChannel("room:lobby", nil)

	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestPheonixGoSocket_OpenChannel(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	channel, err := socket.OpenChannel("room:lobby")

	if err != nil {
		t.Errorf(err.Error())
	}

	if channel == nil {
		t.Errorf("Channel should exist!")
	}
}

func TestPheonixGoSocket_SetCustomTimeout(t *testing.T){

	timeOut := 5 * time.Second

	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)
	socket.SetCustomTimeout(timeOut)

	if socket.Timeout != timeOut {
		t.Errorf("Socket timeout mismatch!")
	}
}

func TestPheonixGoSocket_GetChannel(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	socket.OpenChannel("room:lobby")
	socket.JoinChannel("room:lobby", nil)

	channel, err := socket.GetChannel("room:lobby")

	if err != nil {
		if err.Error() != "channel does not exist or is already closed" {
			t.Errorf("unknown error: " + err.Error())
		}
	}

	if channel == nil {
		t.Errorf("Channel should exist!")
	}
}

func TestPheonixGoSocket_GetChannel_DoesNotExist(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	channel, err := socket.GetChannel("room:lobby")

	if err != nil {
		if err.Error() != "cannot get channel that does not exist" {
			t.Errorf("unknown error: " + err.Error())
		}
	}

	if channel != nil {
		t.Error("Channel should not exist")
	}
}

func TestPheonixGoSocket_CloseChannel(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	socket.OpenChannel("room:lobby")

	err := socket.CloseChannel("room:lobby")

	if err != nil {
		if err.Error() != "channel does not exist or is already closed" {
			t.Errorf("unknown error: " + err.Error())
		}
	}
}

func TestPheonixGoSocket_CloseChannel_No_Channels_Exists(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	err := socket.CloseChannel("NOT EXISTING")

	if err != nil {
		if err.Error() != "channel does not exist or is already closed" {
			t.Errorf("unknown error: " + err.Error())
		}
	}
}

func TestPheonixGoSocket_ClosePheonixWebsocket(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)
	socket.Listen()
	socket.ClosePheonixWebsocket()

	if socket.Status != PhxGoClosed {
		t.Errorf("Failed to close socket!")
	}
}

func TestPheonixGoSocket_ClosePheonixWebsocket_Never_Started(t *testing.T) {
	socket := NewPheonixWebsocket("localhost:4000", "/socket", "ws", false)

	socket.ClosePheonixWebsocket()

	if socket.Status != PhxGoClosed {
		t.Errorf("Socket in unknown state, should be close since it was never started!")
	}
}