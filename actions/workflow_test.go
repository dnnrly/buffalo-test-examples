package actions

import (
	"buffalo_test_examples/models"
	"net/http"

	"github.com/gobuffalo/httptest"
)

func (as *ActionSuite) Test_Workflow_HomepageLogin() {
	as.LoadFixture("base-data")

	req := httptest.New(as.App)

	{
		res := req.HTML("/").Get()
		as.Equal(http.StatusFound, res.Code)
		as.Equal("/auth/new", res.Location())
	}

	{
		res := req.HTML("/auth/new").Get()
		as.Equal(http.StatusOK, res.Code)
		as.Contains(res.Body.String(), "Sign In")
	}

	{
		res := req.HTML("/auth").Post(&models.User{
			Email:    "user1@example.com",
			Password: "password",
		})
		as.Equal(http.StatusFound, res.Code)
		as.Equal("/", res.Location())
	}

	{
		res := req.HTML("/").Get()
		as.Equal(http.StatusOK, res.Code)
	}
}
