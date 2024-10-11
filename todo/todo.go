package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Item struct {
	CompletedAt time.Time
	CreatedAt   time.Time
	Task        string
	Done        bool
}

// List represents a list of ToDo items and Verbose mode
type List struct {
	Items        []Item
	VerboseMode  bool
	HideComplete bool
}

// Add creates a ToDo item and append it to the List
func (l *List) Add(task string) {
	t := Item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	l.Items = append(l.Items, t)
}

// Complete method mark a ToDo item as completed by settings Done = true
// And CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls.Items) {
		return fmt.Errorf("item %d does not exist", i)
	}

	ls.Items[i-1].Done = true
	ls.Items[i-1].CompletedAt = time.Now()

	return nil
}

// Delete method deletes a ToDo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls.Items) {
		return fmt.Errorf("item %d does not exist", i)
	}

	l.Items = append(ls.Items[:i-1], ls.Items[i:]...)
	return nil

}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l.Items)
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
	return json.Unmarshal(file, &l.Items)
}

// Implements the fmt.Stringer interface
func (l *List) String() string {
	formated := ""
	for k, t := range l.Items {
		if l.HideComplete && t.Done {
			continue
		}
		prefix := " "
		if t.Done {
			prefix = "X "
		}
		if l.VerboseMode {
			formated += fmt.Sprintf("%s%d: %s, Created at:%s\n", prefix, k+1, t.Task, t.CreatedAt)
		} else {
			formated += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
		}
	}
	return formated
}
