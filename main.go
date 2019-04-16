package main

import (
	"fmt"
	"os"

	"github.com/gernest/go-wasm-server/server"
	"github.com/urfave/cli"
)

func main() {
	a := cli.NewApp()
	a.Name = "go-wasm-server"
	a.Usage = "helper for working with wasm in Go"
	a.Commands = []cli.Command{
		server.Serve(),
	}
	if err := a.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
