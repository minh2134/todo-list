package task_test

import (
	"fmt"
	"testing"
	"todo-list/internal/task"
)

func TestMakeTask(t *testing.T) {
	templ := "Name: %s\nDesc: %s\nCompleted: false\n"
	tests := []struct {
		name string
		desc string
	}{
		{"", ""},
		{"Task1", "example task"},
		{"Task2: Hello world!", "hello there!"},
	}

	for _, test := range tests {
		t.Run("Test", func(t *testing.T) {
			t.Parallel()
			tsk := task.MakeTask(test.name, test.desc)
			exptected := fmt.Sprintf(templ, test.name, test.desc)
			if tsk.String() != exptected {
				t.Errorf("got:\n%s,\nexpected:\n%s", tsk.String(), exptected)
			}

		})
	}
}
