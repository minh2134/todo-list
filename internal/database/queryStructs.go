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
