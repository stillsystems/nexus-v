package cli

import (
	"flag"
	"os"

	"github.com/stillsystems/nexus-v/internal/web"
)

func runServe(args []string) {
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	port := serveCmd.Int("port", 8080, "Port to run the visual scaffolder on")
	serveCmd.Parse(args)

	if err := web.Start(*port); err != nil {
		Error("Failed to start visual scaffolder: " + err.Error())
		os.Exit(1)
	}
}

