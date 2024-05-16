package storage

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type FileDb struct {
	path                   string
	inMemoryStore          map[string]any
	OBJ_TYPE_KEY_SEPARATOR string
}

func (db *FileDb) New(db_path string) (*FileDb, error) {
	db.path = db_path
	db.OBJ_TYPE_KEY_SEPARATOR = "-"
	err := db.Reload()
	return db, err
}

func (db *FileDb) Reload() error {
	db.inMemoryStore = make(map[string]any)
	content, _ := os.ReadFile(db.path)

	err := json.Unmarshal(content, &(db.inMemoryStore))
	if err != nil {
		return err
	}
	return nil
}

func (db *FileDb) AllRecordsCount() int {
	return len(db.inMemoryStore)
}

func (db *FileDb) Save(obj any) (string, error) {
	uid := uuid.NewString()
	id := reflect.TypeOf(obj).Name() + db.OBJ_TYPE_KEY_SEPARATOR + uid
	json_rep, err := json.Marshal(obj) // test if it can be jsoned
	if err != nil {
		return "", err
	}

	// save map[string]any rep
	var saved_version map[string]any
	json.Unmarshal(json_rep, &saved_version)

	db.inMemoryStore[id] = saved_version

	return id, nil
}

// returns objects with any type so users can rebuild
// objects with their type builders
func (db *FileDb) Get(id string) (any, error) {
	stored, found := db.inMemoryStore[id]
	if found {
		return stored, nil
	}
	return nil, errors.New("FileDb: Get: failed to get object with id: " + id)
}

func (db *FileDb) GetRecordsByField(objTypeName, field string, value any) ([]map[string]any, error) {
	if objTypeName == "" {
		return nil, errors.New("FileDb: GetRecordsByField: no records found for " + objTypeName)
	}

	var listOfRecordsOfSameType = db.GetAllOfType(objTypeName)

	var listOfMatchedRecords []map[string]any
	var compValue any

	// convert value number to float64 if value is a number
	numberVal, ok := getFloat64Equivalent(value)
	if ok {
		compValue = numberVal
	} else {
		compValue = value
	}

	for _, record := range listOfRecordsOfSameType {
		if record[field] == compValue {
			listOfMatchedRecords = append(listOfMatchedRecords, record)
		}
	}

	return listOfMatchedRecords, nil
}

func (db *FileDb) GetAllOfType(objTypeName string) []map[string]any {
	var listOfRecordsOfSameType []map[string]any

	for key, val := range db.inMemoryStore {
		if strings.HasPrefix(key, objTypeName) {
			concVal, ok := val.(map[string]any)
			if !ok {
				panic(`FileDb: GetRecordsByField: records found for is not of  
					map[string]any type` + objTypeName)
			}
			listOfRecordsOfSameType = append(listOfRecordsOfSameType, concVal)
		}
	}

	return listOfRecordsOfSameType
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

	if concrete_obj, ok := obj.(map[string]any); ok {
		concrete_obj[data.Field] = data.Value
		db.inMemoryStore[id] = concrete_obj
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

func getFloat64Equivalent(value any) (float64, bool) {
	if concVal, ok := value.(int); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int8); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int16); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int32); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int64); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(float32); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(float64); ok {
		return float64(concVal), true
	}

	return 0, false
}

var Db, err = new(FileDb).New("db.json")

func MakeFileDb(db_path string) (*FileDb, error) {
	return new(FileDb).New(db_path)
}

var GLOBAL_FILE_DB, _ = MakeFileDb("file_db.json")
