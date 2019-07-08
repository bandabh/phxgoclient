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



//func main() {
//	var addr = flag.String("addr", "localhost:4000", "http service address")
//
//	//conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000/socket", http.Header{})
//	flag.Parse()
//
//	log.SetFlags(0)
//
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	u := url.URL{Scheme: "ws", Host: *addr, Path: "/socket/websocket"}
//	log.Printf("connecting to %s", u.String())
//
//	//c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//
//	c, err := Connect(u)
//
//	if err != nil {
//		log.Fatal("dial:", err)
//	}
//
//	ch := c.MakeAndJoinAChannel("room:E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C", nil)
//	println(ch.CanPush())
//
//	ch.Leave()
//
//	defer c.Close()
//
//	done := make(chan struct{})
//
//	go func() {
//		defer close(done)
//
//		ch.Push("ping", nil) //Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
//
//		//beat := Message{
//		//	"room:E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C", JoinEvent.EventToString(), Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}, 0,
//		//}
//		//
//		//c.Socket.WriteJSON(beat)
//
//		//ch.Read()
//		//
//		//var k Response
//		//
//		//types := ch.Client.Socket.ReadJSON(&k)
//		ch.On(CloseEvent, func(response MessageResponse) (data interface{}, err error) {
//
//			println("WORKSSSSS!!!!!")
//			println(response.Event)
//			println(response.Payload.Status)
//			println(response.Topic)
//			println(response.Ref)
//
//			return response, err
//		})
//		//ch.On(ToEvent("after_join"), func(response MessageResponse) (data interface{}, err error) {
//		//
//		//	//println("testeststsetse")
//		//	//println(response.Ref)
//		//	//println(response.Topic)
//		//	//println(response.Event)
//		//	//println(response.Payload.Response)
//		//	//println("testeststsetse")
//		//	return nil, nil
//		//})
//		ch.Push("ping", nil) //Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
//		ch.On(ErrorEvent, func(response MessageResponse) (updated interface{}, err error) {
//
//			//println(types)
//
//			//println(response.Event)
//			//println(response.Topic)
//			//println(response.Ref)
//
//			//println(kresponse.Payload.Response.Reason)
//
//			return response, nil
//		})
//		ch.Push("ping", nil) //Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
//		ch.On(ReplyEvent, func(response MessageResponse) (updated interface{}, err error) {
//
//			////println(types)
//			//println("11111111111111111111111111111e")
//			//println(response.Event)
//			//println(response.Topic)
//			//println(response.Ref)
//			//println("11111111111111111111111111111e")
//			//println(response.Payload.Status)
//			//println("11111111111111111111111111111e")
//			return response, nil
//		})
//
//		//
//		//if err != nil {
//		//	log.Println("read:", err)
//		//	return
//		//}
//
//		println(ch.CanPush())
//		//log.Printf("recv: %s", message)
//
//	}()
//
//	go func() {
//		for now := range time.Tick(time.Minute) {
//			fmt.Println(now, "test")
//		}
//	}()
//
//	timer2 := time.NewTimer(time.Second)
//
//	go func() {
//		<-timer2.C
//		for true {
//			println("here")
//			ch.Client.Socket.PingHandler()
//			ch.heartbeath()
//			println("pass")
//			n := screenshot.NumActiveDisplays()
//
//			for i := 0; i < n; i++ {
//				bounds := screenshot.GetDisplayBounds(i)
//
//				img, err := screenshot.CaptureRect(bounds)
//				if err != nil {
//					panic(err)
//				}
//
//				//fileName := fmt.Sprintf("%d_%dx%d.out",  i, bounds.Dx(), bounds.Dy())
//				//file, _ := os.Create(fileName)
//				//defer file.Close()
//				//png.Encode(file, img)
//
//				frameBuffer := new(bytes.Buffer)
//
//				//frameBuffer := make([]byte, imgFrame.ImageSize())
//				err1 := png.Encode(frameBuffer, img)
//				if err1 != nil {
//					panic(err)
//				}
//
//				//header := []byte("data:image/png;base64,")
//
//				// convert the buffer bytes to base64 string - use buf.Bytes() for new image
//				imgBase64Str := base64.StdEncoding.EncodeToString(frameBuffer.Bytes())
//
//				//_ = imgBase64Str
//				//file.Write([]byte(imgBase64Str))
//				//file.Close()
//				time.Sleep(300 * time.Microsecond)
//				ch.Push("screen_send", Image{imgBase64Str})
//
//			}
//		}
//	}()
//
//	ticker := time.NewTicker(time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		//case <-done:
//		//	return
//		case t := <-ticker.C:
//
//			//beat := Message{
//			//	"room:E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C", "echo", Payload{"stuff"}, 0,
//			//}
//			//
//			//c.Socket.WriteJSON(beat)
//			//println(t.String())
//
//			//ch.Push("ping", nil) //Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
//			ch.On(ErrorEvent, func(response MessageResponse) (updated interface{}, err error) {
//
//				//println(types)
//
//				println(response.Event)
//				println(response.Topic)
//				println(response.Ref)
//
//				//println(kresponse.Payload.Response.Reason)
//
//				return response, nil
//			})
//			//ch.Push("ping", nil) //Payload{"E7AD-49CB-9010-47F7-A06E-9D2D-3351-DD1C"}
//
//			t.String()
//			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
//			if err != nil {
//				log.Println("write:", err)
//				return
//			}
//		case <-interrupt:
//			log.Println("interrupt")
//
//			// Cleanly close the connection by sending a close message and then
//			// waiting (with timeout) for the server to close the connection.
//			err := c.Socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//			if err != nil {
//				log.Println("write close:", err)
//				return
//			}
//			select {
//			case <-done:
//			case <-time.After(time.Second):
//			}
//			return
//		}
//	}
//
//	for true {
//		_ = ""
//	}
//}
