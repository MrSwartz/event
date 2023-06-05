package main

import (
	"fmt"
	"os"

	"github.com/MrSwartz/event/pkg/eventservice"
)

func main() {
	if err := eventservice.Run(); err != nil {
		fmt.Fprintf(os.Stdout, "service stopped with error: %v", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "service successfully stopped")
	os.Exit(0)
}
