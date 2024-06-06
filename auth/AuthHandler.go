package auth

import (
	"github.com/Iyusuf40/go-auth/models"
	"github.com/Iyusuf40/go-auth/storage"
	"github.com/google/uuid"
)

type AuthHandler struct {
	temp_store  storage.TempStore
	users_store storage.Storage[models.User]
}

var DEFAULT_SESSION_TIMEOUT = 86400.0

func (auth_h *AuthHandler) HandleLogin(email, password string) string {
	retrievedUsers := auth_h.users_store.GetByField("email", email)

	if len(retrievedUsers) == 0 {
		return ""
	}

	user := retrievedUsers[0]
	if !user.IsCorrectPassword(password) {
		return ""
	}

	sessionId := uuid.NewString()
	userId := auth_h.users_store.GetIdByField("email", email)

	auth_h.temp_store.SetKeyToValWIthExpiry(sessionId, userId, DEFAULT_SESSION_TIMEOUT)

	return sessionId
}

func (auth_h *AuthHandler) HandleLogout(sessionId string) {
	auth_h.temp_store.DelKey(sessionId)
}

func (auth_h *AuthHandler) IsLoggedIn(sessionId string) bool {
	return auth_h.temp_store.GetVal(sessionId) != ""
}

func (auth_h *AuthHandler) ExtendSession(sessionId string, duration float64) {
	auth_h.temp_store.ChangeKeyEpiry(sessionId, duration)
}

func MakeAuthHandler(temp_store_db, users_store_db, recordsName string) *AuthHandler {
	auth_h := new(AuthHandler)
	auth_h.temp_store = storage.GET_TempStore(temp_store_db, recordsName)
	auth_h.users_store = storage.MakeUserStorage(users_store_db, recordsName)
	return auth_h
}
