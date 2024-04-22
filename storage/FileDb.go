package storage

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type FileDb struct {
	path          string
	inMemoryStore map[string]any
}

func (db *FileDb) New() (*FileDb, error) {
	db.path = "db.json"
	content, err := os.ReadFile(db.path)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(content, &(db.inMemoryStore))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *FileDb) save(obj any) (string, error) {
	id := uuid.NewString()
	json_rep, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	// TO-DO
	//		check type and do name-spacing?
	db.inMemoryStore[id] = json_rep
	return id, nil
}

func (db *FileDb) get(id string) any {
	val := db.inMemoryStore[id]
	return val
}

func (db *FileDb) delete(id string) {
	delete(db.inMemoryStore, id)
}

func (db *FileDb) update(id string, data UpdateDesc) bool {
	_, exists := db.inMemoryStore[id]
	if !exists {
		return false
	}

	obj := db.inMemoryStore[id]

	if conc_obj, ok := obj.(map[string]any); ok {
		conc_obj[data.field] = data.value
	} else {
		panic("cannot covert object in inMemoryStore to map")
	}

	return true
}

func (db *FileDb) commit() error {
	json_rep, err := json.Marshal(db.inMemoryStore)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json_rep, 0444)
	return err
}

var db, err = new(FileDb).New()
