package auth

import (
	"github.com/Iyusuf40/go-auth/storage"
	"github.com/google/uuid"
)

type AuthHandler struct {
	ts storage.TempStore
}

var userStore = storage.MakeUserStorage("")
var DEFAULT_SESSION_TIMEOUT = 86400.0

func (auth_h *AuthHandler) HandleLogin(email, password string) string {
	retrievedUsers := userStore.GetByField("email", email)

	if len(retrievedUsers) == 0 {
		return ""
	}

	user := retrievedUsers[0]
	if !user.IsCorrectPassword(password) {
		return ""
	}

	sessionId := uuid.NewString()
	userId := userStore.GetIdByField("email", email)

	auth_h.ts.SetKeyToValWIthExpiry(sessionId, userId, DEFAULT_SESSION_TIMEOUT)

	return sessionId
}

func (auth_h *AuthHandler) HandleLogout(sessionId string) {
	// receives sessId
	// removes session from store in tempStore
}

func (auth_h *AuthHandler) IsLoggedIn(sessionId string) {
	// check if logged in
	// extend ttl if logged in
	// return object message {res: isLoggedin}
}

func (auth_h *AuthHandler) ExtendSession(sessionId string, duration float64) {
	// receives sessId or email??
	// removes session from store in tempStore
}

func MakeAUthHandler(db_path string) *AuthHandler {
	auth_h := new(AuthHandler)
	auth_h.ts = storage.MakeTempStoreFileDbImpl(db_path)
	return auth_h
}
