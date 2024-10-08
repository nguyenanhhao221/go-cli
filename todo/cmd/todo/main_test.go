package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName = "todo"
)

func TestMain(m *testing.M) {
	os.Setenv("TODO_FILENAME", ".test_todo.json")
	fmt.Println("Building tool....")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()
	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(os.Getenv("TODO_FILENAME"))

	os.Exit(result)

}

func TestTodoCLI(t *testing.T) {
	t.Setenv("TODO_FILENAME", ".test_todo.json")
	task := "test task number 1"

	dir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)
	t.Run("Add New Task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.WriteString(cmdStdIn, task2); err != nil {
			t.Fatal(err)
		}
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n 2: %s\n", task, task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %s instead \n", expected, out)
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		var cmd *exec.Cmd
		cmd = exec.Command(cmdPath, "-delete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n", task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %s instead \n", expected, out)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListWithCompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X 1: %s\n", task2)

		if expected != string(out) {
			t.Errorf("Expected %q, got %s instead \n", expected, out)
		}
	})

	t.Run("ListWithCompleteTaskWithHideComplete", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "--hide-complete")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := ""

		if expected != string(out) {
			t.Errorf("Expected %q, got %s instead \n", expected, out)
		}
	})
}
