package main

import (
	"os"

	"github.com/marcuwynu23/git-share/internal/cli"
)

func main() {
	opts := cli.Parse(os.Args)
	cli.Run(opts)
}
