package phxgoclient

import (
	"net/url"
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"errors"
)

type Client struct {
	Socket *websocket.Conn
}

type Channel struct {
	State State
	Topic string

	Client Client

	RefCount int64

	Events     map[Event]ChannelCallbackFunc
	DataBuffer map[string]interface{}
}

type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
	Ref     int64       `json:"ref"`
}

type MessageResponse struct {
	Topic string `json:"topic"`
	Ref   int    `json:"ref"`
	Payload struct {
		Status   string      `json:"status"`
		Response interface{} `json:"response"`
	} `json:"payload"`
	Event string `json:"event"`
}


func (channel *Channel) Reconnect(payload interface{}) {
	if channel.State == ErroredState {
		channel.Join(payload)
	}
}

func (channel *Channel) incrementRef() int64 {
	channel.RefCount = channel.RefCount + 1
	return channel.RefCount
}

// To avoid a connection timeout the client needs to send the server a heartbeat event.
func (channel *Channel) heartbeath() {

	if !channel.IsOpen() {
		return
	}

	message := Message{
		"phoenix", "heartbeat", nil, 0,
	}

	channel.Client.Socket.WriteJSON(message)
}

type ChannelCallbackFunc func(response MessageResponse) (data interface{}, err error)

func (channel *Channel) UnRegister(event Event) error {

	_, ok := channel.Events[event]

	if ok {
		delete(channel.Events, event)
		return nil
	}

	return errors.New("event " + event.EventToString() + " exists")
}

func (channel *Channel) IsOpen() bool {
	return channel.State != ClosedState && channel.State != LeavingState
}

func (channel *Channel) Register(event Event, callback ChannelCallbackFunc) error {
	_, ok := channel.Events[event]

	if !ok {
		channel.Events[event] = callback
		return nil
	}

	return errors.New("event listner already exists")
}

// Observes all registered events on a channel
func (channel *Channel) Observe() (data interface{}, err error) {

	if !channel.IsOpen() {
		return nil, errors.New("channel is closed or in leaving state")
	}

	resp := channel.Read()

	callbackMap, ok := channel.Events[ToEvent(resp.Event)]
	if ok {
		return callbackMap(resp)
	}

	return nil, errors.New("event not registered, or error while reading response")
}

func (channel *Channel) Read() MessageResponse {
	var resp MessageResponse

	err := channel.Client.Socket.ReadJSON(&resp)

	if err != nil{
		println(err)
	}

	switch resp.Payload.Status {
	case OkStatus.StatusToString():
		channel.State = JoinedState
		break
	case ErrorStatus.StatusToString():
		channel.State = ErroredState
		break
	case TimeoutStatus.StatusToString():
		channel.State = ErroredState
		break
	}

	return resp

}

func (channel *Channel) CanPush() bool {
	return channel.State == JoinedState
}

func (channel *Channel) Join(payload interface{}) *Channel {
	channel.State = JoiningState

	msg_push(channel, payload, JoinEvent)

	return channel
}

// RawPush, unsafe push that does not check if you can push on a channel
func (channel *Channel) RawPush(event string, payload interface{}) {
	if event == "" {
		event = MessageEvent.EventToString()
	}

	msg_push(channel, payload, ToEvent(event))
}

// Push event with payload to channel
func (channel *Channel) Push(event string, payload interface{}) error {

	// Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
	if channel.CanPush() {
		channel.RawPush(event, payload)
		return nil
	} else {
		channel.DataBuffer[event] = payload
		return errors.New("cannot push payload to channel")
	}
}

func msg_push(channel *Channel, payload interface{}, event Event) {

	message := Message{
		channel.Topic, event.EventToString(), payload, channel.incrementRef(),
	}

	channel.Client.Socket.WriteJSON(message)
}

func (channel *Channel) Leave() {
	channel.ChannelLeave()
	channel.State = ClosedState
}

func (client Client) MakeChannel(topic string) Channel {
	channel := Channel{
		ClosedState,
		topic,
		client,
		0,
		make(map[Event]ChannelCallbackFunc),
		make(map[string]interface{}),
	}

	return channel
}

func (client Client) MakeAndJoinAChannel(topic string, payload interface{}) Channel {
	ch := client.MakeChannel("room:E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C")
	ch.Join(payload)
	return ch
}

func Connect(url url.URL) (*Client, error) {

	strUrl, errConv := checkUrl(url)

	if errConv != nil {
		return nil, errors.New("failed to parse url, no proper url format")
	}

	socket, _, err := websocket.DefaultDialer.Dial(strUrl, nil)
	client := Client{
		socket,
	}

	if err != nil {
		log.Fatal("dial:", err)
		return nil, errors.New("failed to connect")
	}

	return &client, nil
}

func (channel *Channel) ChannelLeave() error {
	channel.State = LeavingState
	return channel.Client.Close()
}

func (client *Client) Close() error {
	return client.Socket.Close()
}

func checkUrl(url url.URL) (string, error) {
	if url.Scheme == "" {
		url.Scheme = "ws"
		fmt.Errorf("scheme not defined, defaulting to ws")
	}

	if url.Path == "" {
		url.Path = "/socket/websocket"
		fmt.Errorf("path not defined, defaulting to /socket/websocket")
	}

	if url.Host == "" {
		return url.String(), errors.New("no host specified")
	}

	return url.String(), nil
}
