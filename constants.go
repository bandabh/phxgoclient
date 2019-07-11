package phxgoclient

type Status string
type State string
type Event string
type Transport string

const (
	PhxSocket   Transport = "/websocket"
	PhxLongPool Transport = "/longpoll"
)

func (transport Transport) ToPath() string {
	return string(transport)
}

const (
	OkStatus      Status = "ok"
	ErrorStatus   Status = "error"
	TimeoutStatus Status = "timeout"
)

const (
	ClosedState  State = "closed"
	ErroredState State = "errored"
	JoinedState  State = "joined"
	JoiningState State = "joining"
	LeavingState State = "leaving"
)

const (
	// MessageEvent represents a regular message on a topic.
	MessageEvent Event = "phx_message"
	// JoinEvent represents a successful join on a channel.
	JoinEvent Event = "phx_join"
	// CloseEvent represents the closing of a channel.
	CloseEvent Event = "phx_close"
	// ErrorEvent represents an error.
	ErrorEvent Event = "phx_error"
	// ReplyEvent represents a reply to a message sent on a topic.
	ReplyEvent Event = "phx_reply"
	// LeaveEvent represents leaving a channel and unsubscribing from a topic.
	LeaveEvent Event = "phx_leave"
)

func (status Status) StatusToString() string {
	return string(status)
}

func (state State) StateToString() string {
	return string(state)
}

func (event Event) EventToString() string {
	return string(event)
}

func ToEvent(str string) Event {
	var event Event = Event(str)
	return event
}
