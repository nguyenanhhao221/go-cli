package main

import (
	"fmt"
	"os"
	"time"

	"haonguyen.tech/interactiveTools/pomo/app"
	"haonguyen.tech/interactiveTools/pomo/pomodoro"
	"haonguyen.tech/interactiveTools/pomo/pomodoro/repository"
)

func main() {
	// TODO: Replace with actual user input config
	repo := repository.NewInMemoryRepo()
	config := pomodoro.NewConfig(repo, 2*time.Second, 2*time.Second, 5*time.Second)

	a, err := app.New(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := a.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
