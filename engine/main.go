package main

import (
	"fmt"
	"os"

	"github.com/b7r-dev/lyt/engine/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "lyt: %v\n", err)
		os.Exit(1)
	}
}
