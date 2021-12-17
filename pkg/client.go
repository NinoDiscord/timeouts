package pkg

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"strings"
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
	t := Timeout{
		Type:        item["type"].(string),
		GuildId:     item["guild_id"].(string),
		UserId:      item["user_id"].(string),
		IssuedAt:    int64(item["issued_at"].(float64)),
		ExpiresAt:   int64(item["expires_at"].(float64)),
		ModeratorId: item["moderator"].(string),
	}

	if item["reason"] != nil {
		t.Reason = item["reason"].(string)
	}

	return t
}

func mapAll(toMap interface{}) []Timeout {
	var timeouts []Timeout
	bytes, _ := json.Marshal(toMap)
	_ = json.Unmarshal(bytes, &timeouts)
	return timeouts
}

func (c *Client) WriteMessage(msg Message) {
	c.writeLock.Lock()
	err := c.Conn.WriteJSON(msg)
	c.writeLock.Unlock()

	if err != nil {
		logrus.Errorf("Unable to write %s to client: %v", marshalToString(msg), err)
	} else {
		logrus.Debugf("Wrote data to client: %s", marshalToString(msg))
	}
}

func (c *Client) HandleTimeout(t Timeout) {
	logrus.Debugf("Told to handle timeout (type=%s; guild=%s; user=%s)", t.Type, t.GuildId, t.UserId)
	go func() {
		select {
		case <-time.After(time.Duration(t.ExpiresAt-t.IssuedAt) * time.Millisecond):
			{
				if !Server.HasClient() {
					Server.QueueIn(t)
					logrus.Warnf("Client has been disconnected, added pending timeout to replay soon.")
					return
				}

				c.WriteMessage(Message{
					OP:   Apply,
					Data: t,
				})
			}
		}
	}()
}

func (c *Client) HandleMessage(msg Message) {
	switch msg.OP {
	case RequestAll:
		{
			data, err := Redis.Connection.HGetAll(context.TODO(), "nino:timeouts").Result()
			if err != nil {
				logrus.Warnf("Unable to retrieve all timeouts, are we connected?\n%v", err)
				c.WriteMessage(Message{
					OP:   RequestAllBack,
					Data: []Timeout{},
				})

				return
			}

			// serialize output to map[string]Timeout{}
			mappedData := map[string]Timeout{}
			for key, value := range data {
				timeout := Timeout{}
				if err := json.NewDecoder(strings.NewReader(value)).Decode(&timeout); err != nil {
					logrus.Warnf("Unable to decode packet %s, skipping", value)
					continue
				}

				mappedData[key] = timeout
			}

			c.WriteMessage(Message{
				OP:   RequestAll,
				Data: mappedData,
			})
		}

	case Request:
		{
			c.HandleTimeout(toTimeout(msg.Data.(map[string]interface{})))
		}
	}
}
