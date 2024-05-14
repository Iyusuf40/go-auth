package tests

import (
	"os"
	"reflect"
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

	id := US.Save(user)
	if id == "" {
		t.Fatal("TestSaveAndGetUser: id should not be empty")
	}

	retrievedUser, _ := US.Get(id)

	if reflect.DeepEqual(retrievedUser, user) == false {
		t.Fatal("TestSaveAndGetUser: retrievedUser should be equal to saved")
	}

	// test user with similar mail cannot be duplicated
	if US.Save(user) != "" {
		t.Fatal("TestSaveAndGetUser: id should be empty")
	}
}
