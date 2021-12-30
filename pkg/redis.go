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
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

var Redis *RedisClient

type RedisClient struct {
	Connection *redis.Client
}

func NewRedis() error {
	if Redis != nil {
		panic(errors.New("tried to create new redis client, no thank you"))
	}

	logrus.Info("Now connecting to Redis...")
	index, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}

	password := os.Getenv("REDIS_PASSWORD")

	var connection *redis.Client
	if sentinels := os.Getenv("REDIS_SENTINELS"); sentinels != "" {
		hosts := strings.Split(os.Getenv("REDIS_SENTINELS"), ";")
		connection = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: hosts,
			MasterName:    os.Getenv("REDIS_MASTER"),
			Password:      password,
			DB:            index,
			DialTimeout:   10 * time.Second,
			ReadTimeout:   15 * time.Second,
			WriteTimeout:  15 * time.Second,
		})
	} else {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		connection = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Password:     password,
			DB:           index,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		})
	}

	if err := connection.Ping(context.TODO()).Err(); err != nil {
		return err
	} else {
		logrus.Info("Connected to Redis!")
	}

	Redis = &RedisClient{
		Connection: connection,
	}

	return nil
}

func (r RedisClient) Connect() error {
	logrus.Info("Now connecting to Redis...")
	index, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}

	password := os.Getenv("REDIS_PASSWORD")
	if sentinels := os.Getenv("REDIS_SENTINELS"); sentinels != "" {
		hosts := strings.Split(os.Getenv("REDIS_SENTINELS"), ";")
		r.Connection = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: hosts,
			MasterName:    os.Getenv("REDIS_MASTER"),
			Password:      password,
			DB:            index,
			DialTimeout:   10 * time.Second,
			ReadTimeout:   15 * time.Second,
			WriteTimeout:  15 * time.Second,
		})
	} else {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		r.Connection = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", host, port),
			Password:     password,
			DB:           index,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		})
	}

	if err := r.Connection.Ping(context.TODO()).Err(); err != nil {
		return err
	} else {
		logrus.Info("Connected to Redis!")
		return nil
	}
}
