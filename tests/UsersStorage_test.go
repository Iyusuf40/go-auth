package tests

import (
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/Iyusuf40/go-auth/models"
	"github.com/Iyusuf40/go-auth/storage"
)

var users_storage_test_db_path = "users_store_test_db.json"
var US storage.Storage[models.User]

func beforeEachUST() {
	US = storage.MakeUserStorage(users_storage_test_db_path)
}

func afterEachUST() {
	os.Remove(users_storage_test_db_path)
}

func TestSaveAndGetUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
	}

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestSaveAndGetUser: id should not be empty")
	}

	retrievedUser, _ := US.Get(id)

	if reflect.DeepEqual(retrievedUser, user) == false {
		t.Fatal("TestSaveAndGetUser: retrievedUser should be equal to saved")
	}

	// test user with similar mail cannot be duplicated
	if _, success = US.Save(user); success {
		t.Fatal("TestSaveAndGetUser: success should be false")
	}
}

func TestUpdateUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
	}

	updateField := "phone"
	updateValue := 9000

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestUpdateUser: id should not be empty")
	}

	retrievedUser, _ := US.Get(id)

	if reflect.DeepEqual(retrievedUser, user) == false {
		t.Fatal("TestUpdateUser: retrievedUser should be equal to saved")
	}

	US.Update(id, storage.UpdateDesc{Field: updateField, Value: updateValue})

	retrievedUser, _ = US.Get(id)

	if retrievedUser.Phone != updateValue {
		t.Fatal("TestUpdateUser: retrievedUser.Phone should be equal", updateValue)
	}
}

func TestDeleteUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
	}

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestDeleteUser: id should not be empty")
	}

	retrievedUser, _ := US.Get(id)

	if reflect.DeepEqual(retrievedUser, user) == false {
		t.Fatal("TestDeleteUser: retrievedUser should be equal to saved")
	}

	US.Delete(id)
	retrievedUser, err := US.Get(id)

	if reflect.DeepEqual(retrievedUser, models.User{}) != true {
		t.Fatal("TestDeleteUser: retrievedUser should be empty")
	}

	if err == nil {
		t.Fatal("TestDeleteUser: getting nonexistent user should return error")
	}
}

func TestGetUserByField(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	email := "mail@mail"
	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
	}

	_, success := US.Save(user)
	if !success {
		t.Fatal("TestGetUserByField: id should not be empty")
	}

	retrievedUser := US.GetByField("email", email)[0]

	if reflect.DeepEqual(retrievedUser, user) == false {
		t.Fatal("TestGetUserByField: retrievedUser should be equal to saved")
	}
}

func TestGetAllUsers(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	emails := []string{"mail1@mail.com", "mail2@mail.com", "mail3@mail.com"}

	for _, email := range emails {
		user := models.User{
			Email:     email,
			FirstName: "f_name",
			LastName:  "l_name",
			Phone:     8000,
		}
		_, success := US.Save(user)
		if !success {
			t.Fatal("TestGetUserByField: id should not be empty")
		}
	}

	retrievedUsers := US.GetAll()

	if len(retrievedUsers) != len(emails) {
		t.Fatal("TestGetAllUsers: retrievedUsers should have length equal to saved users")
	}

	for _, user := range retrievedUsers {
		if slices.Contains(emails, user.Email) == false {
			t.Fatal("TestGetAllUsers:", user.Email, "should be in", emails)
		}
	}
}
