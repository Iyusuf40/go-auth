package tests

import (
	"testing"

	"github.com/Iyusuf40/go-auth/storage"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) buildUser(obj any) *User {
	if map_rep, ok := obj.(map[string]any); ok {
		usr_ob := new(User)
		usr_ob.Name = map_rep["name"].(string)
		usr_ob.Age = int(map_rep["age"].(float64))
		if usr_ob.Name == "" {
			return nil
		}
		return usr_ob
	}
	return nil
}

func TestGet(t *testing.T) {

	test_db_path := "test_db.json"
	var DB, _ = storage.MakeFileDb(test_db_path)

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

		// test load from file
		DB.Commit()
		DB2, _ := storage.MakeFileDb(test_db_path)
		obj, _ = DB2.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}
		saved_user = *new(User).buildUser(obj)
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

		DB.DeleteDb()

	} else {
		t.Fatal("TestGet: failed to Save")
	}
}
