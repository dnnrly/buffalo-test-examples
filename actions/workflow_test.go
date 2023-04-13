package actions

import (
	"buffalo_test_examples/models"
	"net/http"
)

func (as *ActionSuite) Test_Workflow_HomepageLogin() {
	as.LoadFixture("base-data")

	{
		res := as.HTML("/").Get()
		as.Equal(http.StatusFound, res.Code)
		as.Equal("/auth/new", res.Location())
	}

	{
		res := as.HTML("/auth/new").Get()
		as.Equal(http.StatusOK, res.Code)
		as.Contains(res.Body.String(), "Sign In")
	}

	{
		res := as.HTML("/auth").Post(&models.User{
			Email:    "user1@example.com",
			Password: "password",
		})
		as.Equal(http.StatusFound, res.Code)
		as.Equal("/", res.Location())
	}

	{
		res := as.HTML("/").Get()
		as.Equal(http.StatusOK, res.Code)
	}
}
