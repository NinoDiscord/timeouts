// Copyright (c) 2021 Nino
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pkg

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
)

type WebSocketServer struct {
	upgrader websocket.Upgrader
	mutex    *sync.Mutex
	Queue    []Timeout
	client   *Client
}

func (s *WebSocketServer) HasClient() bool {
	return s.client == nil
}

func (s *WebSocketServer) QueueIn(t Timeout) {
	s.mutex.Lock()
	s.Queue = append(s.Queue, t)
	s.mutex.Unlock()
}

var (
	Server     *WebSocketServer
	authHeader = http.CanonicalHeaderKey("Authorization")
)

func NewServer() {
	if Server != nil {
		panic("Attempt to initialise another server instance!")
	}

	Server = &WebSocketServer{
		upgrader: websocket.Upgrader{},
		Queue:    []Timeout{},
		mutex:    &sync.Mutex{},
		client:   nil,
	}

	// Overwrite queue based off Redis
	data, err := Redis.Connection.Get(context.TODO(), "nino:timeouts").Result()
	if err != nil {
		if err == redis.Nil {
			logrus.Warnf("We have retrieved `null` when fetching for timeouts, skipping as empty queue.")

			// Set it to an empty array
			err = Redis.Connection.HMSet(context.TODO(), "nino:timeouts", []Timeout{}).Err()
			if err != nil {
				logrus.Errorf("Unable to set data as an empty array. :<")
				return
			}

			logrus.Infof("Queue has been set in Redis, now starting fresh!")
			Server.Queue = make([]Timeout, 0)
			return
		} else {
			logrus.Warnf("Unable to retrieve all timeouts, are we connected?\n%v", err)
			return
		}
	}

	// serialize output to []Timeout{}
	mappedData := make([]Timeout, 0)
	err = json.Unmarshal([]byte(data), &mappedData)
	if err != nil {
		panic(err)
	}

	Server.Queue = mappedData
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	conn, err := Server.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logrus.Errorf("Unable to upgrade WebSocket: %v", err)
		return
	}

	if req.Header.Get(authHeader) == "" || req.Header.Get(authHeader) != os.Getenv("AUTH") {
		logrus.Warnf("Tried to connect client, but received bad authentication key (header=%s)", req.Header.Get(authHeader))
		_ = conn.Close()
		return
	}

	Server.client = &Client{
		Conn:      conn,
		writeLock: &sync.Mutex{},
	}

	Server.client.WriteMessage(Message{OP: Ready})
	go func() {
		for {
			if Server.client == nil {
				break
			}

			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				Server.client = nil
				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway, websocket.CloseInternalServerErr) {
					logrus.Info("Received disconnect from bot, will replay events once it is back...")
					break
				}
			}

			if len(Server.Queue) > 0 {
				// Replay lost events
				for _, event := range Server.Queue {
					Server.client.WriteMessage(Message{
						OP:   Apply,
						Data: event,
					})
				}

				// Drop old queue
				Server.Queue = []Timeout{}
			}

			go Server.client.HandleMessage(message)
		}
	}()
}
