package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Iyusuf40/go-auth/api/controllers"
	"github.com/Iyusuf40/go-auth/models"
	"github.com/Iyusuf40/go-auth/storage"
	"github.com/labstack/echo/v4"
)

var users_api_test_db_path = "users_router_store_test_db.json"

func beforeEachUAPIT() {
	controllers.UserStorage = storage.MakeUserStorage(users_api_test_db_path)
}

func afterEachUAPIT() {
	os.Remove(users_api_test_db_path)
}

func TestPOSTUser(t *testing.T) {
	// Setup
	beforeEachUAPIT()
	defer afterEachUAPIT()

	e := echo.New()
	userJSON := `{"data": {"firstName":"John", "lastName": "Doe","email":"mail@mail.com", 
	"password": "xx", "phone": 90543434}}`

	// test successfully saving a user
	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}
	rec, c := SetupRequest(e, http.MethodPost, "/api/users", userJSON, headers)
	controllers.SaveUser(c)

	if http.StatusCreated != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /api/users: expected:", http.StatusCreated, "got:", rec.Code)
	}

	// test failed saving of user with email already existing
	rec, c = SetupRequest(e, http.MethodPost, "/api/users", userJSON, headers)
	controllers.SaveUser(c)

	if http.StatusBadRequest != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /api/users: expected:", http.StatusBadRequest, "got:", rec.Code)
	}

	// test failed saving of user with missing userfield
	userJSON = `{"data": {"firstName":"John", "lastName": "Doe","email":"", "phone": 90543434}}`
	rec, c = SetupRequest(e, http.MethodPost, "/api/users", userJSON, headers)
	controllers.SaveUser(c)

	if http.StatusBadRequest != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /api/users: expected:", http.StatusBadRequest, "got:", rec.Code)
	}
}

func TestGETUser(t *testing.T) {
	// Setup
	beforeEachUAPIT()
	defer afterEachUAPIT()

	user := models.User{Email: "testmail2@mail.com",
		FirstName: "fname",
		LastName:  "lname",
		Phone:     999,
		Password:  "xxx",
	}

	id, saved := controllers.UserStorage.Save(user)

	if !saved {
		t.Fatal("GET /api/user/:id: expected: true got:", saved)
	}

	e := echo.New()
	rec, c := SetupRequest(e, http.MethodPost, "/api/users", "", nil)
	c.SetParamNames("id")
	c.SetParamValues(id)
	controllers.GetUser(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("GET /api/users/:id : expected:", http.StatusOK, "got:", rec.Code)
	}

	// test non existent id
	rec, c = SetupRequest(e, http.MethodPost, "/api/users", "", nil)
	c.SetParamNames("id")
	c.SetParamValues("non-existent")
	controllers.GetUser(c)

	if http.StatusNotFound != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("GET /api/users/:id : expected:", http.StatusNotFound, "got:", rec.Code)
	}
}

func TestPUTUser(t *testing.T) {
	// Setup
	beforeEachUAPIT()
	defer afterEachUAPIT()

	user := models.User{Email: "testmail@mail.com",
		FirstName: "fname",
		LastName:  "lname",
		Phone:     999,
		Password:  "xxx",
	}

	id, saved := controllers.UserStorage.Save(user)

	if !saved {
		t.Fatal("PUT /api/user/:id: expected: true got:", saved)
	}

	newPhone := 99
	upadateJSON := fmt.Sprintf(`{"data": {"field":"phone", "value": %d}}`, newPhone)

	// test successfully saving a user
	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}

	e := echo.New()
	rec, c := SetupRequest(e, http.MethodPut, "/api/users/:id", upadateJSON, headers)
	c.SetParamNames("id")
	c.SetParamValues(id)

	controllers.UpdateUser(c)

	if rec.Code != http.StatusOK {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("PUT /api/users/:id : expected:", http.StatusOK, "got:", rec.Code)
	}

	upadatedUser, _ := controllers.UserStorage.Get(id)

	if upadatedUser.Phone != newPhone {
		t.Fatal("PUT /api/users/:id : expected retrieved user.Phone:", newPhone,
			"got:", upadatedUser.Phone)
	}
}

func TestDELETEUser(t *testing.T) {
	// Setup
	beforeEachUAPIT()
	defer afterEachUAPIT()

	user := models.User{Email: "testmail@mail.com",
		FirstName: "fname",
		LastName:  "lname",
		Phone:     999,
		Password:  "xxx",
	}

	id, saved := controllers.UserStorage.Save(user)

	if !saved {
		t.Fatal("DELETE /api/user/:id: expected: true got:", saved)
	}

	// attempt to get user
	_, err := controllers.UserStorage.Get(id)

	if err != nil {
		t.Fatal("DELETE /api/user/:id: expected error to be nil got:", err)
	}

	e := echo.New()
	rec, c := SetupRequest(e, http.MethodPut, "/api/users/:id", "", nil)
	c.SetParamNames("id")
	c.SetParamValues(id)

	controllers.DeleteUser(c)

	if rec.Code != http.StatusOK {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("DELETE /api/users/:id : expected:", http.StatusOK, "got:", rec.Code)
	}

	// attempt to get deleted user
	_, err = controllers.UserStorage.Get(id)

	if err == nil {
		t.Fatal("DELETE /api/user/:id: expected error to be non-nil got:", err)
	}
}

func SetupRequest(
	e *echo.Echo,
	httpMethod,
	route,
	body string,
	httpHeadersAndValues map[string]string,
) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(httpMethod, route, strings.NewReader(body))
	for headerKey, headerValue := range httpHeadersAndValues {
		req.Header.Set(headerKey, headerValue)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return rec, c
}
