package main

import (
	"dev.floofy.nino/timeouts/pkg"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
}

func main() {
	pkg.NewServer()
	pkg.StartServer()
}
