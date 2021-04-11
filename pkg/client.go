package pkg

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Client struct {
	Conn      *websocket.Conn
	writeLock *sync.Mutex
}

func marshalToString(d interface{}) string {
	bytes, _ := json.Marshal(d)
	return string(bytes)
}

func toTimeout(item map[string]interface{}) Timeout {
	return Timeout{
		Type:    PunishmentType(item["type"].(string)),
		Guild:   item["guild"].(string),
		User:    item["user"].(string),
		Issued:  int64(item["issued"].(float64)),
		Expired: int64(item["expired"].(float64)),
	}
}

func mapAll(toMap interface{}) []Timeout {
	var timeouts []Timeout
	bytes, _ := json.Marshal(toMap)
	_ = json.Unmarshal(bytes, &timeouts)
	return timeouts
}

func (c *Client) WriteMessage(msg WSMessage) {
	c.writeLock.Lock()
	err := c.Conn.WriteJSON(msg)
	c.writeLock.Unlock()
	if err != nil {
		logrus.Errorf("Failed to write %s to client!", marshalToString(msg))
	} else {
		logrus.Debugf("Wrote %s to client!", marshalToString(msg))
	}
}

func (c *Client) HandleTimeout(t Timeout) {
	logrus.Debugf("Received a request to handle timeout: %s", marshalToString(t))
	go func() {
		select {
		case <-time.After(time.Duration(t.Expired-t.Issued) * time.Millisecond):
			{
				if !Server.HasClient() {
					Server.Queue(t)
					logrus.Warnf("Client has disconnected, adding a pending timeout to the replay queue.")
					return
				}
				c.WriteMessage(WSMessage{
					Op: Apply,
					D:  t,
				})
			}
		}
	}()
}

func (c *Client) HandleMessage(msg WSMessage) {
	switch msg.Op {
	case Request:
		{
			c.HandleTimeout(toTimeout(msg.D.(map[string]interface{})))
		}
	case Acknowledged:
		{
			timeouts := mapAll(msg.D)
			for _, timeout := range timeouts {
				c.HandleTimeout(timeout)
			}
			logrus.Info("Client acknowledged ready!")
		}
	}
}
