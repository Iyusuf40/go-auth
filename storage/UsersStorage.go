package storage

import (
	"encoding/json"
	"fmt"

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
	retrievedUsers, _ = us.DB.GetRecordsByField("User", field, value)
	users = us.buildManyUsers(retrievedUsers)
	return users
}

func (us *UserStorage) GetAll() []models.User {
	var users []models.User
	retrievedUsers := us.DB.GetAllOfType("User")
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("UserStorage.BuildClient:", r)
		}
	}()

	user := models.User{}
	if map_rep, ok := obj.(map[string]any); ok {
		user.FirstName = map_rep["firstName"].(string)
		user.LastName = map_rep["lastName"].(string)
		user.Email = map_rep["email"].(string)
		phoneFloatVal, _ := getFloat64Equivalent(map_rep["phone"])
		user.Phone = int(phoneFloatVal)
	}
	return user
}

func (us *UserStorage) userWithEmailExist(email string) bool {
	queryRes, _ := us.DB.GetRecordsByField("User", "email", email)
	return len(queryRes) > 0
}

func (us *UserStorage) isValidUser(user models.User) bool {
	if user.Email == "" || user.FirstName == "" || user.LastName == "" || user.Phone == 0 {
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

func MakeUserStorage(db_path string) Storage[models.User] {
	var DB DB_Engine = GLOBAL_FILE_DB
	if db_path != "" {
		DB, _ = MakeFileDb(db_path)
	}
	US := new(UserStorage)
	US.DB = DB
	return US
}
