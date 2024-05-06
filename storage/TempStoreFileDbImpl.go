package storage

import (
	"reflect"
)

type TempStoreFileDbImpl struct {
	db       *FileDb
	MapStore map[string]any
	Init     bool
}

func (TS *TempStoreFileDbImpl) New(dbFile string) TempStore {
	fileDb, _ := MakeFileDb(dbFile)
	TS.db = fileDb
	TS.reload()
	return TS
}

func (TS *TempStoreFileDbImpl) reload() {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		TS.MapStore = map[string]any{}
		TS.Init = true
		TS.db.Save(*TS)
		TS.commit()
	}
}

func (TS *TempStoreFileDbImpl) GetVal(key string) string {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		return ""
	}
	val, _ := mapStore[key].(string)
	return val
}

func (TS *TempStoreFileDbImpl) SetKeyToVal(key string, value string) bool {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		return false
	}
	mapStore[key] = value
	TS.commit()
	return true
}

func (TS *TempStoreFileDbImpl) SetKeyToValWIthExpiry(key string, value string, expiry int) bool {
	return false
}

func (TS *TempStoreFileDbImpl) ExtendKeyEpiry(key string, newExpiry int) bool {
	return false
}

func (TS *TempStoreFileDbImpl) DelKey(key string) bool {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		return false
	}
	delete(mapStore, key)
	TS.commit()
	return true
}

func (TS *TempStoreFileDbImpl) commit() {
	TS.db.Commit()
}

func (TS *TempStoreFileDbImpl) DeleteDb() {
	TS.db.DeleteDb()
}

func (TS *TempStoreFileDbImpl) getMapStore() map[string]any {
	TS_name := reflect.TypeOf(*TS).Name()
	mapStoreList, _ := TS.db.GetRecordsByField(TS_name, "Init", true)
	if mapStoreList == nil {
		return nil
	}

	if len(mapStoreList) > 1 {
		panic(`TempStoreFileDbImpl: getMapStore: there should only be one 
			instance in existenc`)
	}
	mapStore := mapStoreList[0].(map[string]any)
	return mapStore
}

func MakeTempStoreFileDbImpl(db_path string) TempStore {
	return new(TempStoreFileDbImpl).New(db_path)
}
