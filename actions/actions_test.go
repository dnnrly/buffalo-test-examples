package actions

import (
	"buffalo_test_examples/models"
	"os"
	"testing"

	"github.com/gobuffalo/suite/v4"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), os.DirFS("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func (as *ActionSuite) getUserDetails(email string) *models.User {
	u := &models.User{}
	err := as.DB.Where("email = ?", email).First(u)
	if err != nil {
		as.FailNowf("error getting user %s: %v", email, err)
	}
	return u
}

func (as *ActionSuite) setUserSession(email string) *models.User {
	u := &models.User{}
	err := as.DB.Where("email = ?", email).First(u)
	if err != nil {
		as.FailNowf("failed to create user session", "error creating session for user %s: %v", email, err)
	}
	as.Session.Set("current_user_id", u.ID)

	return u
}
