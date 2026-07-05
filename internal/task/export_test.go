package task

import "fmt"

// this file exports functions that only available in tests

func (t Task) String() string {
	return fmt.Sprintf("Name: %s\nDesc: %s\nCompleted: %t\n", t.name, t.desc, t.completed)
}
