package todo_test

import (
	"testing"
	"todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New Task should not be completed")
	}

	_ = l.Complete(1)
	if !l[0].Done {
		t.Errorf("This task should be completed")
	}
}
