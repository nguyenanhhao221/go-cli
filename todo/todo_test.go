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
func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := [3]string{"foo", "bar", "baz"}
	for _, taskName := range tasks {
		l.Add(taskName)
	}

	_ = l.Delete(2)
	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d", 2, len(l))
	}

	if tasks[2] != l[1].Task {
		t.Errorf("Expected %q, got %q", tasks[2], l[1].Task)
	}

}
