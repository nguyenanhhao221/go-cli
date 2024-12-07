package pomodoro_test

import (
	"os"
	"testing"

	"haonguyen.tech/interactiveTools/pomo/pomodoro"
	"haonguyen.tech/interactiveTools/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	tf, err := os.CreateTemp("", "pomo")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	dbRepo, err := repository.NewSqlite3Repo(tf.Name())
	if err != nil {
		t.Fatal(err)
	}

	return dbRepo, func() {
		os.Remove(tf.Name())
	}
}
