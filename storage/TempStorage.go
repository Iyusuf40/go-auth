package storage

type TempStore interface {
	GetVal(key string) string
	SetKeyToVal(key string, value string) bool
	SetKeyToValWIthExpiry(key string, value string, expiry float64) bool
	ChangeKeyEpiry(key string, newExpiry float64) bool
	DelKey(key string) bool
}

func GET_TempStore(database, recordsName string) TempStore {
	return MakeTempStoreFileDbImpl(database, recordsName)
}
