package storage

type UpdateDesc struct {
	field string
	value any
}

// Every table or collection must implement the Storage interface
type Storage[T any] interface {
	Get(id string) T
	Create(data T) string
	Update(id string, data UpdateDesc) bool
	Delete(id string) bool
	GetByField(data any) []T
	GetAll() []T
}

// for every transaction these processes
//		connection to db
//		assumption of transaction to be for particular table
//		operations
//			create common executions
