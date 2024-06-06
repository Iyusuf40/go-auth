package storage

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Iyusuf40/go-auth/models"
)

type UserStorage struct {
	DB DB_Engine
}

func (us *UserStorage) Get(id string) (models.User, error) {
	val, err := us.DB.Get(id)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	obj := us.BuildClient(val)
	return obj, nil
}

func (us *UserStorage) Save(user models.User) (msg string, success bool) {
	if !us.isValidUser(user) {
		success = false
		msg = "invalid user"
		return msg, success
	}

	if us.userWithEmailExist(user.Email) {
		success = false
		msg = fmt.Sprintf("user with email %s exists", user.Email)
		return msg, success
	}

	user.HashPassword()
	id, err := us.DB.Save(user)

	if err != nil {
		success = false
		msg = err.Error()
		return msg, success
	}

	us.DB.Commit()
	success = true
	msg = id

	return msg, success
}

func (us *UserStorage) Update(id string, data UpdateDesc) bool {
	field := data.Field

	// check if field exists on User struct
	exists := fieldExistsOnUser(field)
	canRebuild := us.userIsRebuildableWithUpdatedData(id, field, data.Value)

	if !exists || !canRebuild {
		return false
	}

	res := us.DB.Update(id, data)
	us.DB.Commit()
	return res
}

func (us *UserStorage) Delete(id string) {
	us.DB.Delete(id)
	us.DB.Commit()
}

func (us *UserStorage) GetByField(field string, value any) []models.User {
	var users []models.User
	var retrievedUsers []map[string]any
	retrievedUsers, _ = us.DB.GetRecordsByField(field, value)
	users = us.buildManyUsers(retrievedUsers)
	return users
}

func (us *UserStorage) GetIdByField(field string, value any) string {
	return us.DB.GetIdByFieldAndValue(field, value)
}

func (us *UserStorage) GetAll() []models.User {
	var users []models.User
	retrievedUsers := us.DB.GetAllOfRecords()
	users = us.buildManyUsers(retrievedUsers)
	return users
}

func (us *UserStorage) buildManyUsers(retrievedUsers []map[string]any) []models.User {
	var users []models.User

	for _, userDesc := range retrievedUsers {
		user := us.BuildClient(userDesc)
		users = append(users, user)
	}

	return users
}

func (us *UserStorage) BuildClient(obj any) models.User {

	// after recovery, zero value of enclosing function
	// will be returned
	defer RecoverFromPanic()

	user := models.User{}
	if map_rep, ok := obj.(map[string]any); ok {
		user.FirstName = map_rep["firstName"].(string)
		user.LastName = map_rep["lastName"].(string)
		user.Email = map_rep["email"].(string)
		user.Password = map_rep["password"].(string)
		phoneFloatVal, _ := getFloat64Equivalent(map_rep["phone"])
		user.Phone = int(phoneFloatVal)
	}
	return user
}

func (us *UserStorage) userWithEmailExist(email string) bool {
	queryRes, _ := us.DB.GetRecordsByField("email", email)
	return len(queryRes) > 0
}

func (us *UserStorage) isValidUser(user models.User) bool {
	if user.Email == "" ||
		user.FirstName == "" ||
		user.LastName == "" ||
		user.Phone == 0 ||
		user.Password == "" {
		return false
	}
	return true
}

// try to rebuild user with updated data and return
// true if possible else return false
func (us *UserStorage) userIsRebuildableWithUpdatedData(id, field string, value any) bool {
	prevDesc, err := us.DB.Get(id)
	if err != nil {
		return false
	}
	if concDesc, ok := prevDesc.(map[string]any); ok {
		var copyUserDesc = map[string]any{}
		// copy concDesc
		for key, value := range concDesc {
			copyUserDesc[key] = value
		}
		copyUserDesc[field] = value
		user := us.BuildClient(copyUserDesc)
		return us.isValidUser(user)
	}
	return false
}

func fieldExistsOnUser(field string) bool {
	// map rep of user was used because of json reps of
	// User struct has it fields having json tagged keys
	// and requests come in using this lower cased keys
	// which do not match with the Capitalised exported
	// struct keys otherwise an approach similar to
	//
	// _, exists := reflect.TypeOf(User{}).FieldByName(field)
	//
	// would have been used
	_, exists := getMapRepOfUser()[field]
	return exists
}

func getMapRepOfUser() map[string]any {
	var mapRep map[string]any
	jsonBytes, _ := json.Marshal(models.User{})
	json.Unmarshal(jsonBytes, &mapRep)
	return mapRep
}

func RecoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("UserStorage.BuildClient:", r)
	}
}

func MakeUserStorage(db_path, recordsName string) Storage[models.User] {
	if recordsName == "" {
		recordsName = reflect.TypeOf(models.User{}).Name()
	}

	// STORAGE, err := MakeFileDb(db_path, recordsName)
	STORAGE, err := MakePostgresEngine(db_path,
		recordsName,
		SQL_TABLE_COLUMN_FIELD_AND_DESC{`email`, "VARCHAR(128)"},
		SQL_TABLE_COLUMN_FIELD_AND_DESC{`firstName`, "VARCHAR(128)"},
		SQL_TABLE_COLUMN_FIELD_AND_DESC{`lastName`, "VARCHAR(128)"},
		SQL_TABLE_COLUMN_FIELD_AND_DESC{`phone`, "integer"},
		SQL_TABLE_COLUMN_FIELD_AND_DESC{`password`, "VARCHAR(128)"},
	)

	if err != nil {
		panic(err)
	}

	var DB DB_Engine = STORAGE
	US := new(UserStorage)
	US.DB = DB
	return US
}
