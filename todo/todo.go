package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	CompletedAt time.Time
	CreatedAt   time.Time
	Task        string
	Done        bool
}

// List represents a list of ToDo items
type List []item

// Add creates a ToDo item and append it to the List
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

// Complete method mark a ToDo item as completed by settings Done = true
// And CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete method deletes a ToDo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)
	return nil

}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)

}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return json.Unmarshal(file, l)
}

// Implements the fmt.Stringer interface
func (l *List) String() string {
	formated := ""
	for k, t := range *l {
		prefix := " "
		if t.Done {
			prefix = "X "
		}

		formated += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}
	return formated
}
