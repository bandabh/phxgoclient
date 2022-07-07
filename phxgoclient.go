package phxgoclient

import (
	"errors"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type PhxGoSocketStatus string

const (
	PhxGoOpen   PhxGoSocketStatus = "open"
	PhxGoError  PhxGoSocketStatus = "error"
	PhxGoClosed PhxGoSocketStatus = "closed"
)

type PheonixGoSocket struct {
	Host     string
	Schema   string
	Path     string
	RawQuery string

	Status PhxGoSocketStatus

	Timeout time.Duration

	CustomAbsoultePath bool

	Transport Transport

	HeartbeatWorker *Worker

	Channels map[string]*Channel
}

// Set your timeout intervel and heartbeat interval in format (interval * Duration) e.g (30 * time.Seconds)
func (phx *PheonixGoSocket) SetCustomTimeout(interval time.Duration) {
	phx.Timeout = interval
	return
}

func (phx *PheonixGoSocket) ClosePheonixWebsocket() {
	if phx.HeartbeatWorker != nil {
		phx.HeartbeatWorker.Shutdown()
	}
	phx.Status = PhxGoClosed
	phx.Channels = nil
}

// Creates New Pheonix Websocket connection
func NewPheonixWebsocket(Host string, Path string, Schema string, CustomAbsoultePath bool, RawQuery string) PheonixGoSocket {
	return PheonixGoSocket{
		Host,
		Schema,
		Path,
		RawQuery,
		PhxGoClosed,
		30 * time.Second,
		CustomAbsoultePath,
		PhxSocket,
		nil,
		make(map[string]*Channel),
	}
}

func (client *Client) heartbeat() {
	message := Message{
		"phoenix", "heartbeat", nil, 0,
	}

	client.Socket.WriteJSON(message)
}

// Starts the Pheonix Websocket
func (phx PheonixGoSocket) Listen() error {
	phx.Status = PhxGoOpen

	path := phx.Path

	if !phx.CustomAbsoultePath {
		path = path + phx.Transport.ToPath()
	}

	u := url.URL{Scheme: phx.Schema, Host: phx.Host, Path: path, RawQuery: phx.RawQuery}

	client, err := Connect(u)

	if err != nil {
		return err
	}

	phx.HeartbeatWorker = NewWorker(phx.Timeout, func() {
		client.heartbeat()
	})

	go phx.HeartbeatWorker.Run()

	return nil
}

// Raw access for a channel
func (phx *PheonixGoSocket) GetChannel(topic string) (*Channel, error) {
	channel, ok := phx.Channels[topic]

	if ok {
		return channel, nil
	} else {
		return nil, errors.New("cannot get channel that does not exist")
	}
}

func (phx *PheonixGoSocket) JoinChannel(topic string, payload interface{}) error {
	channel, ok := phx.Channels[topic]

	if ok {
		channel := channel.Join(payload)
		channel.Read()

		phx.Channels[topic] = channel
		return nil
	} else {
		return errors.New("channel does not exist, failed to join")
	}
}

func (phx *PheonixGoSocket) OpenChannel(topic string) (*Channel, error) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	path := phx.Path

	if !phx.CustomAbsoultePath {
		path = path + phx.Transport.ToPath()
	}

	u := url.URL{Scheme: phx.Schema, Host: phx.Host, Path: path, RawQuery: phx.RawQuery}

	client, err := Connect(u)

	if err != nil {
		phx.Status = PhxGoError
		return nil, errors.New("dial:" + err.Error())
	}

	ch := client.MakeChannel(topic)

	cha, ok := phx.Channels[topic]

	if !ok {
		phx.Channels[topic] = &ch
	} else {
		return cha, nil
	}

	return &ch, nil
}

func (phx *PheonixGoSocket) CloseChannel(topic string) error {

	channel, ok := phx.Channels[topic]

	if ok {
		channel.Leave()
		if !channel.IsOpen() {
			delete(phx.Channels, topic)
		} else {
			delete(phx.Channels, topic)
			return errors.New("force closing channel, failed to gracefully leave")
		}

		return nil
	}

	return errors.New("channel does not exist or is already closed")
}
