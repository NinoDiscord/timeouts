package pkg

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type WSServer struct {
	client     *Client
	queue      []Timeout
	queueMutex *sync.Mutex
	upgrader   websocket.Upgrader
}

func (s *WSServer) HasClient() bool {
	return s.client != nil
}

func (s *WSServer) Queue(t Timeout) {
	s.queueMutex.Lock()
	s.queue = append(s.queue, t)
	s.queueMutex.Unlock()
}

var (
	Server *WSServer
	header = http.CanonicalHeaderKey("Authorization")
)

func NewServer() {
	if Server != nil {
		panic("Attempt to initialise another server instance!")
	}
	Server = &WSServer{
		client:     nil,
		queue:      []Timeout{},
		queueMutex: &sync.Mutex{},
		upgrader:   websocket.Upgrader{},
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	con, err := Server.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Socket upgrade error: %v", err)
		return
	}
	if req.Header.Get(header) == "" || req.Header.Get(header) != GetAuth() {
		logrus.Warnf("Got bad auth %s from client!", req.Header.Get(header))
		_ = con.Close()
		return
	}
	Server.client = &Client{
		Conn:      con,
		writeLock: &sync.Mutex{},
	}
	Server.client.WriteMessage(WSMessage{Op: Ready})
	go func() {
		for {
			if Server.client == nil {
				break
			}
			var msg WSMessage
			err := con.ReadJSON(&msg)
			if err != nil {
				Server.client = nil
				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr) {
					logrus.Info("Bot disconnected from server, when it reconnects we'll play back lost events...")
					break
				}
			}
			if len(Server.queue) > 0 {
				// Replay any lost events
				for _, event := range Server.queue {
					Server.client.WriteMessage(WSMessage{
						Op: Apply,
						D:  event,
					})
				}
				// Drop old queue since we replayed the events we needed to
				Server.queue = []Timeout{}
			}
			go Server.client.HandleMessage(msg)
		}
	}()
}

func StartServer() {
	http.HandleFunc("/", handleRequest)
	logrus.Info("Listening on port 4025!")
	if err := http.ListenAndServe("0.0.0.0:4025", nil); err != nil {
		logrus.Fatal(err)
	}
}
