package tests

import (
	"os"
	"testing"

	"github.com/Iyusuf40/go-auth/auth"
	"github.com/Iyusuf40/go-auth/models"
	"github.com/Iyusuf40/go-auth/storage"
)

// var AuthHandler_userstorage_test_db_path = "auth_users_store_test_db.json"
// var AuthHandler_tempDBstorage_test_db_path = "auth_users_store_test_db.json"
// var Auth_Users_RecordsName = "Users"
var AuthHandler_userstorage_test_db_path = "test"
var AuthHandler_tempDBstorage_test_db_path = "test"
var Auth_Users_RecordsName = "users"

var AUTH_HANDLER *auth.AuthHandler
var AUTH_US storage.Storage[models.User]

func beforeEachAUTH_TEST() {
	AUTH_US = storage.MakeUserStorage(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
	AUTH_HANDLER = auth.MakeAuthHandler(AuthHandler_tempDBstorage_test_db_path,
		AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
}

func afterEachAUTH_TEST() {
	os.Remove(AuthHandler_userstorage_test_db_path)
	os.Remove(AuthHandler_tempDBstorage_test_db_path)
	// storage.RemoveDbSingleton(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
	// storage.RemoveDbSingleton(AuthHandler_tempDBstorage_test_db_path, Auth_Users_RecordsName)
	storage.RemovePostgressEngineSingleton(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName, true)
	storage.RemovePostgressEngineSingleton(AuthHandler_tempDBstorage_test_db_path, Auth_Users_RecordsName, true)
}

func TestHandleLogin(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleLogin: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestHandleLogin: expected a sessionId got empty string")
	}

	// should fail to get a session Token
	sessId = AUTH_HANDLER.HandleLogin(email, password+"a")
	if sessId != "" {
		t.Fatal("TestHandleLogin: expected a sessionId to be empty string")
	}

	sessId = AUTH_HANDLER.HandleLogin(email+email, password)
	if sessId != "" {
		t.Fatal("TestHandleLogin: expected a sessionId to be empty string")
	}
}

func TestIsLoggedIn(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestIsLoggedIn: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestIsLoggedIn: expected a sessionId got empty string")
	}

	is_logged_in := AUTH_HANDLER.IsLoggedIn(sessId)

	if !is_logged_in {
		t.Fatal("TestIsLoggedIn: expected is_logged_in to be true")
	}

	is_logged_in = AUTH_HANDLER.IsLoggedIn(sessId + "abc")

	if is_logged_in {
		t.Fatal("TestIsLoggedIn: expected is_logged_in to be false")
	}
}

func TestHandleLogout(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleLogout: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestHandleLogout: expected a sessionId got empty string")
	}

	is_logged_in := AUTH_HANDLER.IsLoggedIn(sessId)

	if !is_logged_in {
		t.Fatal("TestHandleLogout: expected is_logged_in to be true")
	}

	// logout
	AUTH_HANDLER.HandleLogout(sessId)
	is_logged_in = AUTH_HANDLER.IsLoggedIn(sessId)

	if is_logged_in {
		t.Fatal("TestHandleLogout: expected is_logged_in to be false")
	}
}
