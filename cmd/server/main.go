package main

import (
	"fmt"
	"os"

	"github.com/atozi-ai/gateway/internal/app"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Print version if requested
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("gateway version %s (built at %s)\n", version, buildTime)
		os.Exit(0)
	}

	app.Start()
}
