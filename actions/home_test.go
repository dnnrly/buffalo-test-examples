package actions

import (
	"net/http"
)

func (as *ActionSuite) Test_HomeHandler() {
	res := as.HTML("/").Get()
	as.Equal(http.StatusFound, res.Code)
	as.Equal(res.Location(), "/auth/new")
}

func (as *ActionSuite) Test_HomeHandler_LoggedIn() {
	as.LoadFixture("base-data")
	as.setUserSession("user1@example.com")

	res := as.HTML("/auth").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Sign Out")

	as.Session.Clear()
	res = as.HTML("/auth").Get()
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Sign In")
}
