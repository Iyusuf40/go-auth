package storage

type UpdateDesc struct {
	Field string
	Value any
}

// Every table or collection must implement the Storage interface
type Storage[T any] interface {
	Get(id string) (T, error)
	Save(data T) (msg string, success bool)
	Update(id string, data UpdateDesc) bool
	Delete(id string)
	GetByField(field string, value any) []T
	GetAll() []T
	BuildClient(obj any) T
}

type DB_Engine interface {
	Get(id string) (any, error)
	Save(data any) (string, error)
	// for document stores using noSql, Callers of this function
	// must validate both fields passed in data else, unwanted fields
	// may be added to the records on disc and values of
	// an inappropriate type might be added, causing errors in
	// rebuilding objects
	Update(id string, data UpdateDesc) bool
	Delete(id string)
	// if FileDb is the Engine, field is the json tag if it
	// is defined on the obj
	GetRecordsByField(objTypeName, field string, value any) ([]map[string]any, error)
	GetAllOfType(objTypeName string) []map[string]any
	Commit() error
}

// for every transaction these processes
//		receive connection to db
//		setup transaction to be for particular table
//		operations
//			create common executions
