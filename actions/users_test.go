package actions

import (
	"net/http"

	"buffalo_test_examples/models"
)

func (as *ActionSuite) Test_Users_New() {
	res := as.HTML("/users/new").Get()
	as.Equal(http.StatusOK, res.Code)
}

func (as *ActionSuite) Test_Users_Create() {
	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(0, count)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	res := as.HTML("/users").Post(u)
	as.Equal(http.StatusFound, res.Code)

	found := as.getUserDetails("mark@example.com")
	as.Empty(found.Password) // Make sure the password isn't persisted
}

func (as *ActionSuite) Test_Users_Show() {
	as.LoadFixture("base-data")
	u := as.getUserDetails("user1@example.com")

	res := as.HTML("/users/" + u.ID.String()).Get()

	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "This is user1@example.com")
}

func (as *ActionSuite) Test_Users_Me() {
	as.LoadFixture("base-data")
	as.setUserSession("user1@example.com")

	res := as.HTML("/users/me").Get()

	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "You are user1@example.com")
}
