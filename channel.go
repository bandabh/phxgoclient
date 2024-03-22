package phxgoclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	Url    url.URL
	Socket *websocket.Conn
}

type Channel struct {
	State State
	Topic string

	InitialPayload interface{}

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
	Topic   string  `json:"topic"`
	Ref     int     `json:"ref"`
	Payload Payload `json:"payload"`
	Event   string  `json:"event"`
}

type Payload struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
}

func (channel *Channel) Reconnect() {
	if channel.State == ErroredState {
		channel.Join(channel.InitialPayload)
	}
}

func (channel *Channel) ForceReconnect() error {
	client, err := Connect(channel.Client.Url)

	channel.Client = *client

	if err != nil {
		return err
	}

	channel.Join(channel.InitialPayload)

	return nil
}

func (channel *Channel) incrementRef() int64 {
	channel.RefCount = channel.RefCount + 1
	return channel.RefCount
}

// To avoid a connection timeout the client needs to send the server a heartbeat event.
func (channel *Channel) heartbeat() {

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

func unwrapResponse(b []byte) (bool, MessageResponse) {
	var data MessageResponse

	json.Unmarshal(b, &data)

	if data.Payload.Status != "" {
		return true, data
	}

	return false, data
}

func checkIfSocketClose(err error) bool {
	return strings.Contains(err.Error(), "websocket: close 1000")
}

func (channel *Channel) Read() MessageResponse {
	var resp Message

	err := channel.Client.Socket.ReadJSON(&resp)

	prepared, _ := json.Marshal(resp)

	status, response := unwrapResponse(prepared)

	if err != nil {
		println(err.Error())
		if checkIfSocketClose(err) {
			channel.ForceReconnect()
		}
	}

	if status {
		switch response.Payload.Status {
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
		return response
	}

	return MessageResponse{resp.Topic, int(resp.Ref), Payload{"", resp.Payload}, resp.Event}

}

func (channel *Channel) CanPush() bool {
	return channel.State == JoinedState
}

func (channel *Channel) Join(payload interface{}) *Channel {
	channel.State = JoiningState

	channel.InitialPayload = payload

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
		nil,
		client,
		0,
		make(map[Event]ChannelCallbackFunc),
		make(map[string]interface{}),
	}

	return channel
}

func (client Client) MakeAndJoinAChannel(topic string, payload interface{}) Channel {
	ch := client.MakeChannel(topic)
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
		url,
		socket,
	}

	if err != nil {
		return nil, errors.New("dial:" + err.Error())
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
