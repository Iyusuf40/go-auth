package tests

import (
	"os"
	"testing"

	"github.com/Iyusuf40/go-auth/storage"
)

var temp_test_db_path = "temp_store_test_db.json"
var TS storage.TempStore

func beforeEachTSF() {
	TS = storage.MakeTempStoreFileDbImpl(temp_test_db_path)
}

func afterEachTSF() {
	os.Remove(temp_test_db_path)
}

func TestSetKeyToValAndGetVal(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key"
	got := TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"

	ok := TS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestSetKeyToValAndGetVal: failed to set value")
	}
	got = TS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be " + val + " got " + got)
	}
}

func TestDelKey(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key"
	val := "value"

	ok := TS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestDelKey: failed to set value")
	}
	got := TS.GetVal(key)
	// key exists
	if got != val {
		t.Fatal("TestDelKey: expected value to be " + val + " got " + got)
	}

	ok = TS.DelKey(key)
	if !ok {
		t.Fatal("TestDelKey: failed to delete value")
	}

	got = TS.GetVal(key)
	// key does not exists
	if got != "" {
		t.Fatal("TestDelKey: expected value to be empty")
	}
}
