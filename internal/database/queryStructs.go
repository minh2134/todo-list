package database

type EnumCompleted int

const (
	INCOMPLETE EnumCompleted = iota
	COMPLETED
	ALL
)

type ListQuery struct {
	Name      string
	Completed EnumCompleted
}

type EditQuery struct {
	Id        int
	Name      string
	Desc      string
	Completed EnumCompleted
}
