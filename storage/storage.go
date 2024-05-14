package storage

type UpdateDesc struct {
	Field string
	Value any
}

// Every table or collection must implement the Storage interface
type Storage[T any] interface {
	Get(id string) (T, error)
	Save(data T) string
	Update(id string, data UpdateDesc) bool
	Delete(id string) bool
	GetByField(data any) []T
	GetAll() []T
	BuildSelf(obj any) T
}

type DB_Engine interface {
	Get(id string) (any, error)
	Save(data any) (string, error)
	Update(id string, data UpdateDesc) bool
	Delete(id string)
	// if FileDb is the Engine, field is the json tag if it
	// is defined on the obj
	GetRecordsByField(objTypeName, field string, value any) ([]any, error)
	Commit() error
}

// for every transaction these processes
//		receive connection to db
//		setup transaction to be for particular table
//		operations
//			create common executions
