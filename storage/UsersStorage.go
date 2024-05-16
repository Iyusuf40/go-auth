package storage

import (
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

func MakeUserStorage(db_path string) Storage[models.User] {
	var DB DB_Engine = GLOBAL_FILE_DB
	if db_path != "" {
		DB, _ = MakeFileDb(db_path)
	}
	US := new(UserStorage)
	US.DB = DB
	return US
}
