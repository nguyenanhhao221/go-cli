package main

import (
	"fmt"
	"os"

	"haonguyen.tech/interactiveTools/pomo/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
