package tests

import (
	"fmt"
	"os"
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

var test_db_path = "test_db.json"
var DB *storage.FileDb

func beforeEach() {
	DB, _ = storage.MakeFileDb(test_db_path)
}

func afterEach() {
	DB.DeleteDb()
}

func TestSaveAndGet(t *testing.T) {

	beforeEach()
	defer afterEach()

	user := User{"test user", 20}
	if id, err := DB.Save(user); err == nil {
		var obj, _ = DB.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}

		saved_user := new(User).buildUser(obj)

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
		DB.Reload()
		obj, _ = DB.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}

		saved_user = new(User).buildUser(obj)

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

func TestGetRecordsByField(t *testing.T) {

	beforeEach()
	defer afterEach()

	user := User{"user", 20}
	DB.Save(user)
	var age float32 = 20

	records, err := DB.GetRecordsByField("User", "age", age)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	nonExistentAge := 10
	records, err = DB.GetRecordsByField("User", "age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}

	records, err = DB.GetRecordsByField("User", "name", "user")

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	records, err = DB.GetRecordsByField("NoRecord", "age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}
}

func TestUpdate(t *testing.T) {

	beforeEach()
	defer afterEach()

	user := User{"test", 20}
	id, _ := DB.Save(user)
	updated_name := "updated_name"

	resp := DB.Update(id, storage.UpdateDesc{
		Field: "name",
		Value: updated_name})

	if !resp {
		t.Fatal("TestUpdate: failed to update")
	}

	obj, _ := DB.Get(id)
	saved_user := new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}

	// test after reload
	DB.Commit()
	DB.Reload()

	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestUpdate: all records count should be 1")
	}

	obj, _ = DB.Get(id)
	saved_user = new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}
}

func TestDelete(t *testing.T) {

	beforeEach()
	defer afterEach()

	user := User{"test", 20}
	id, _ := DB.Save(user)

	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestReload: records in inMemoryStore should be 1")
	}

	DB.Delete(id)

	if DB.AllRecordsCount() != 0 {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}
}

func TestAllRecordsCount(t *testing.T) {

	beforeEach()
	defer afterEach()

	if DB.AllRecordsCount() != 0 {
		t.Fatal("TestAllRecordsCount: records in inMemoryStore should be 0")
	}

	noToSave := 10
	user := User{"test", 20}

	for i := 1; i <= noToSave; i++ {
		DB.Save(user)
		if DB.AllRecordsCount() != i {
			t.Fatal("TestAllRecordsCount: records in inMemoryStore should be " + fmt.Sprint(i))
		}
	}
}

func TestReload(t *testing.T) {

	beforeEach()
	defer afterEach()

	zero := 0
	// nothing in inMemoryStore
	if DB.AllRecordsCount() != zero {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}

	user := User{"test", 20}
	DB.Save(user)

	// Reload should only load committed transactions
	DB.Reload()
	if DB.AllRecordsCount() != zero {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}

	id, _ := DB.Save(user)
	DB.Commit()
	DB.Reload()
	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestReload: 1 record should be in inMemoryStore")
	}

	got, _ := DB.Get(id)
	saved_user := new(User).buildUser(got)

	if saved_user.Name != user.Name {
		t.Fatal(
			"TestReload: name of user not equal: user.name = ",
			user.Name, "got =", saved_user.Name)
	}

	if saved_user.Age != user.Age {
		t.Fatal(
			"TestReload: age of user not equal: user.age =",
			user.Age, "got =", saved_user.Age)
	}

}

func TestCommit_DeleteDb(t *testing.T) {

	beforeEach()
	defer afterEach()

	// test db_file does not exist
	_, err := os.Stat(test_db_path)
	if err == nil {
		t.Fatal("TestCommit_DeleteDb: db_file should not exist")
	}

	// test db_file should exist after commit
	DB.Commit()
	_, err = os.Stat(test_db_path)
	if err != nil {
		t.Fatal("TestCommit_DeleteDb: db_file should exist")
	}

	// test db_file should not exist after DeleteDb
	DB.DeleteDb()
	_, err = os.Stat(test_db_path)
	if err == nil {
		t.Fatal("TestCommit_DeleteDb: db_file should not exist")
	}
}
