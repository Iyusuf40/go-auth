package storage

type TempStore interface {
	GetVal(key string) string
	SetKeyToVal(key string, value string) bool
	SetKeyToValWIthExpiry(key string, value string, expiry int) bool
	ChangeKeyEpiry(key string, newExpiry int) bool
	DelKey(key string) bool
}
