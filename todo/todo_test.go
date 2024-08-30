package todo_test

import (
	"os"
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

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New Task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l1[0].Task)
	}

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error when creating temp file %s", err)
	}

	defer os.Remove(tmpFile.Name())

	if err := l1.Save(tmpFile.Name()); err != nil {
		t.Fatalf("Error saving list to file %s", err)
	}
	if err := l2.Get(tmpFile.Name()); err != nil {
		t.Fatalf("Error saving list to file %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q, should match %q task", l1[0].Task, l2[0].Task)
	}

}
