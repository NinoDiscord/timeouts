package pkg

import "os"

func GetAuth() string {
	return os.Getenv("AUTH")
}
