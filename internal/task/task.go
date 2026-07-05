package task

type Task struct {
	Name      string
	Desc      string
	Completed bool
}

func MakeTask(name, desc string) Task {
	return Task{
		Name:      name,
		Desc:      desc,
		Completed: false,
	}
}
