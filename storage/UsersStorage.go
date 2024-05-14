package storage

import (
	"fmt"

	"github.com/Iyusuf40/go-auth/models"
)

// type Storage[T any] interface {
// 	Get(id string) T
// 	Create(data T) string
// 	Update(id string, data UpdateDesc) bool
// 	Delete(id string) bool
// 	GetByField(data any) []T
// 	GetAll() []T
// 	BuildSelf(obj any) T
// }

type UserStorage struct {
	DB DB_Engine
}

func (us *UserStorage) Get(id string) (models.User, error) {
	val, err := us.DB.Get(id)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	obj := us.BuildSelf(val)
	return obj, nil
}

func (us *UserStorage) Save(data models.User) string {
	if us.userWithEmailExist(data.Email) {
		return ""
	}
	id, _ := us.DB.Save(data)
	us.DB.Commit()
	return id
}

func (us *UserStorage) BuildSelf(obj any) models.User {
	user := models.User{}
	if map_rep, ok := obj.(map[string]any); ok {
		user.FirstName = map_rep["firstName"].(string)
		user.LastName = map_rep["lastName"].(string)
		user.Email = map_rep["email"].(string)
		user.Phone = int(map_rep["phone"].(float64))
	}
	return user
}

func (us *UserStorage) Update(id string, data UpdateDesc) bool {
	return false
}

func (us *UserStorage) Delete(id string) bool {
	return false
}

func (us *UserStorage) GetByField(data any) []models.User {
	var res []models.User
	return res
}

func (us *UserStorage) GetAll() []models.User {
	var res []models.User
	return res
}

func (us *UserStorage) userWithEmailExist(email string) bool {
	queryRes, _ := us.DB.GetRecordsByField("User", "email", email)
	return len(queryRes) > 0
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
