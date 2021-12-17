package main

import (
	"fmt"
	"time"
)

// This script exists for a cross-platform version of `date`
// instead of trying to try to mesh it from the Makefile.
func main() {
	fmt.Println(time.Now().Format(time.RFC1123))
}
