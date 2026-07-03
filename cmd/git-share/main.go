package main

import (
	"os"

	"github.com/markwayne/git-share/internal/cli"
)

func main() {
	opts := cli.Parse(os.Args)
	cli.Run(opts)
}
