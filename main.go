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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"nino.sh/timeouts/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Setup logging
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})

	// Setup config
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err := godotenv.Load("./.env")
		if err != nil {
			panic(err)
		}
	}

	logrus.Infof("Using v%s (commit: %s, built at: %s)", pkg.Version, pkg.CommitHash, pkg.BuildDate)
}

func fallbackEnv(envString string, fallback string) string {
	if envString == "" {
		return fallback
	} else {
		return envString
	}
}

func main() {
	if err := pkg.NewRedis(); err != nil {
		panic(err)
	}

	// Create a new `Server` instance
	pkg.NewServer()

	enableMetrics := pkg.SetupMetrics()

	http.HandleFunc("/", pkg.HandleRequest)

	if enableMetrics {
		http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", fallbackEnv(os.Getenv("PORT"), "4025")),
		Handler: nil,
	}

	// Setup syscall signals for Docker
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		// Run the server
		logrus.Infof("Now listening at 0.0.0.0:%s", fallbackEnv(os.Getenv("PORT"), "4025"))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error has occured while listening to server: %v", err)
		}
	}()

	<-sig

	logrus.Warn("Closing off timeouts service due to signal...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Wait for all connections to slowly die off
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			logrus.Fatal("Graceful shutdown timed out, forcing exit!")
		}
	}()

	// :spin:
	defer func() {
		// Save server queue before closing
		data, err := json.Marshal(pkg.Server.Queue)
		if err != nil {
			logrus.Fatalf("Unable to deserialize server queue into JSON: %v", err)
		} else {
			logrus.Info("Saving server queue...")
			err := pkg.Redis.Connection.Set(context.TODO(), "nino:timeouts", string(data), 0).Err()
			if err != nil {
				logrus.Fatalf("Unable to save server queue: %s", err)
			}

			logrus.Info("Saved server queue!")
		}

		err = pkg.Redis.Connection.Close()
		if err != nil {
			logrus.Fatalf("Unable to close Redis server: %v", err)
		}

		// Now we cancel! ^w^
		cancel()
	}()

	// Kill off the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Fatal("Unable to shutdown server: %v", err)
		os.Exit(1)
	} else {
		logrus.Info("Goodbye...")
	}
}
