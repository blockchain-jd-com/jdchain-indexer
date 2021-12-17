package types

type Task interface {
	ID() string
	Do() error
	Status() error
}
