package storage

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/google/uuid"
)

type FileDb struct {
	path          string
	inMemoryStore map[string]any
}

func (db *FileDb) New(db_path string) (*FileDb, error) {
	db.path = db_path
	db.inMemoryStore = make(map[string]any)
	content, err := os.ReadFile(db.path)

	if err != nil { // file doesnt exist
		return db, nil
	}

	json.Unmarshal(content, &(db.inMemoryStore))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *FileDb) Save(obj any) (string, error) {
	uid := uuid.NewString()
	id := reflect.TypeOf(obj).Name() + "-" + uid
	_, err := json.Marshal(obj) // test if it can be jsoned
	if err != nil {
		return "", err
	}

	db.inMemoryStore[id] = obj

	return id, nil
}

// returns objects with any type so users can build
// objects using type assertions
func (db *FileDb) Get(id string) (any, error) {
	stored, found := db.inMemoryStore[id]
	if found {
		return stored, nil
	}
	return nil, errors.New("FileDb: Get: failed to get object with id: " + id)
}

func (db *FileDb) Delete(id string) {
	delete(db.inMemoryStore, id)
}

func (db *FileDb) Update(id string, data UpdateDesc) bool {
	_, exists := db.inMemoryStore[id]
	if !exists {
		return false
	}

	obj := db.inMemoryStore[id]

	if conc_obj, ok := obj.(map[string]any); ok {
		conc_obj[data.field] = data.value
	} else {
		panic("typeof inMemoryStore[id] is not map[string]any")
	}

	return true
}

func (db *FileDb) Commit() error {
	json_rep, err := json.Marshal(db.inMemoryStore)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json_rep, 0644)
	return err
}

func (db *FileDb) DeleteDb() error {
	err = os.Remove(db.path)
	return err
}

var Db, err = new(FileDb).New("db.json")

func MakeFileDb(db_path string) (*FileDb, error) {
	return new(FileDb).New(db_path)
}
