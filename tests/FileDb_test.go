package tests

import (
	"testing"

	"github.com/Iyusuf40/go-auth/storage"
)

var DB, _ = storage.MakeFileDb("test_db.json")

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGet(t *testing.T) {
	user := User{"test user", 20}
	if id, err := DB.Save(user); err == nil {
		var obj, _ = DB.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}

		saved_user, ok := obj.(User)

		if !ok {
			t.Fatal("TestGet: object doesnt impl User interface")
		}

		if saved_user.Name != user.Name {
			t.Fatal(
				"TestGet: name of user not equal: user.name = ",
				user.Name, "got =", saved_user.Name)
		}

		if saved_user.Age != user.Age {
			t.Fatal(
				"TestGet: age of user not equal: user.age =",
				user.Age, "got =", saved_user.Age)
		}

	} else {
		t.Fatal("TestGet: failed to Save")
	}
}
