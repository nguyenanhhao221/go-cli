package pomodoro_test

import (
	"testing"

	"haonguyen.tech/interactiveTools/pomo/pomodoro"
	"haonguyen.tech/interactiveTools/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
