package actions

import (
	"buffalo_test_examples/models"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/antonmedv/expr"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
)

func ln() string {
	_, _, line, _ := runtime.Caller(1)
	return fmt.Sprintf("line %d", line)
}

func (as *ActionSuite) Test_API() {
	as.LoadFixture("base-data")

	uid := func(email string) string {
		return as.getUserDetails(email).ID.String()
	}

	tcases := []struct {
		Name        string
		UserEmail   string
		Method      string
		URL         string
		RequestBody any
		Status      int
		Location    string
		Assertions  []string
	}{
		// Anon user
		{ln(), "", "GET", "/users/" + uid("user1@example.com"), nil, http.StatusOK, "", []string{}},
		{ln(), "", "GET", "/users/me", nil, http.StatusUnauthorized, "", []string{}},
		{ln(), "", "POST", "/users/", &models.User{Email: "user99@example.com"}, http.StatusBadRequest, "", []string{`"current_user_id" not in session.Session.Values`}},
		{ln(), "", "POST", "/users/", &models.User{Email: "user99@example.com", Password: "A", PasswordConfirmation: "A"}, 302, "/users/.+", []string{`"current_user_id" in session.Session.Values`}},

		// // // Logged in user
		{ln(), "user2@example.com", "GET", "/users/" + uid("user1@example.com"), nil, 200, "", []string{}},
		{ln(), "user1@example.com", "GET", "/users/me", nil, 200, "", []string{`resp.Body.String() contains "user1"`}},
	}

	getResult := func(m, url string, body any) *httptest.JSONResponse {
		j := as.JSON(url)
		switch m {
		case "POST":
			return j.Post(body)
		case "PUT":
			return j.Put(body)
		case "PATCH":
			return j.Patch(body)
		case "GET":
			return j.Get()
		case "DELETE":
			return j.Delete()
		}

		panic("Unmapped method: " + m)
	}

	for _, tt := range tcases {
		as.Run(tt.Name, func() {
			_ = as.CleanDB()
			as.LoadFixture("base-data")
			as.Session.Clear()

			var user = &models.User{}
			if tt.UserEmail != "" {
				as.setUserSession(tt.UserEmail)
				user = as.getUserDetails(tt.UserEmail)
			}

			response := getResult(tt.Method, tt.URL, tt.RequestBody)

			testDetails := printTestEnv(as, as.Session, response, tt.RequestBody, user)

			as.NotNil(response, "%s", testDetails)
			as.Equal(tt.Status, response.Code, "%s", testDetails)
			as.Regexp("^"+tt.Location+"$", response.Location(), "%s", testDetails)

			if response.Code != 204 && (response.Code < 300 || response.Code >= 400) {
				as.Contains(response.HeaderMap.Get("Content-Type"), "application/json", testDetails)
			}

			env := map[string]any{
				"session":  as.Session,
				"resp":     response,
				"respBody": response.Body.String(),
				"body":     tt.RequestBody,
				"user":     user,
			}

			for _, a := range tt.Assertions {
				result, err := expr.Eval(a, env)
				as.NoError(err, "%s", testDetails)

				evaluation, ok := result.(bool)
				as.True(ok, "assertion result not a boolean for \"%s\"\n%s", a, testDetails)
				as.True(evaluation, "assertion failed: %s\n%s", a, testDetails)
			}
		})
	}
}

func printTestEnv(as *ActionSuite, session *buffalo.Session, response *httptest.JSONResponse, requestBody any, user *models.User) string {
	result := ""

	{
		result += "\nSession:\n  ID: " + session.Session.ID + "\n  Name: " + session.Session.Name() + "\n  Values:\n"
		for k, v := range session.Session.Values {
			result += fmt.Sprintf("    %s: %s\n", k, v)
		}
	}

	{
		b, err := json.MarshalIndent(response, "", "  ")
		as.NoError(err)
		result += "\nResponse:\n" + string(b)
	}

	{
		b, err := json.MarshalIndent(requestBody, "", "  ")
		as.NoError(err)
		result += "\nRequest Body:\n" + string(b)
	}

	{
		result += "\nResponse Body:\n" + response.Body.String()
	}

	{
		b, err := json.MarshalIndent(user, "", "  ")
		as.NoError(err)
		result += "\nUser:\n" + string(b)
	}

	return result
}
