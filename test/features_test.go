package test_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/cucumber/godog"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type HTTP struct {
	Client      *http.Client
	Response    *http.Response
	RequestData string
	Cookies     []http.Cookie
}

type DB struct {
	Client *sql.DB
}

type testContext struct {
	err  error
	HTTP HTTP
	DB   DB
}

func NewTestContext(psqlInfo string) (testContext, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return testContext{}, err
	}

	return testContext{
		HTTP: HTTP{
			Client: &http.Client{
				Timeout: time.Second * 1,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		},
		DB: DB{
			Client: db,
		},
	}, nil
}

// Errorf is used by the called assertion to report an error and is required to
// make testify assertions work
func (c *testContext) Errorf(format string, args ...interface{}) {
	c.err = fmt.Errorf(format, args...)
}

func (t *testContext) iSendRequestTo(method, path string) error {
	h := t.HTTP
	d := h.RequestData
	req, err := http.NewRequest(
		method,
		"http://web:8080"+path,
		bytes.NewBuffer([]byte(d)),
	)

	if err != nil {
		return err
	}

	for _, c := range t.HTTP.Cookies {
		log.Printf("adding cookie for '%s'", c.Name)
		req.AddCookie(&c)
	}

	t.HTTP.Response, err = t.HTTP.Client.Do(req)

	return err
}

func (t *testContext) theResponseCodeWillBe(code int) error {
	if t.HTTP.Response.StatusCode != code {
		return fmt.Errorf("expected response with code %d but got %d", code, t.HTTP.Response.StatusCode)
	}
	return nil
}

func (t *testContext) theBodyWillBe(expected string) error {
	buf := new(bytes.Buffer)

	_, err := buf.ReadFrom(t.HTTP.Response.Body)
	if err != nil {
		return err
	}

	assert.Equal(t, expected, buf.String())

	return t.err
}

func (t *testContext) iHaveSetDataBufferTo(file string) error {
	data, err := ioutil.ReadFile("reference/" + file)
	if err != nil {
		return err
	}

	t.HTTP.RequestData = string(data)

	return nil
}

func (t *testContext) iHaveACookieWithValue(key, value string) error {
	t.HTTP.Cookies = append(t.HTTP.Cookies, http.Cookie{
		Name:  key,
		Value: value,
	})
	return nil
}

func (t *testContext) iHaveRunDatabaseScript(file string) error {
	data, err := ioutil.ReadFile("reference/" + file)
	if err != nil {
		return err
	}

	_, err = t.DB.Client.Exec(string(data))
	if err != nil {
		return err
	}

	return t.err
}

func (t *testContext) responseCookieIsPresent(cookie string) error {
	found := false
	for _, c := range t.HTTP.Response.Cookies() {
		if c.Name == cookie {
			found = true
		}
	}

	assert.True(t, found, `Could not find cookie "%s" in response:\n%+v`, cookie, t.HTTP.Response)

	return t.err
}

func (t *testContext) theResponseHasAHeaderThatMatches(arg1, arg2 string) error {
	return godog.ErrPending
}

func (t *testContext) theBodyWillLookLike(file string) error {
	data, err := ioutil.ReadFile("reference/" + file)
	if err != nil {
		return err
	}

	r, err := ioutil.ReadAll(t.HTTP.Response.Body)
	if err != nil {
		return err
	}

	assert.JSONEq(t, string(data), string(r))
	return t.err
}

func (t *testContext) theBodyWillLookLikeJsonExcluding(file, exclude string) error {
	expectedData, err := ioutil.ReadFile("reference/" + file)
	if err != nil {
		return err
	}

	var expected map[string]json.RawMessage
	err = json.Unmarshal(expectedData, &expected)
	if err != nil {
		return err
	}

	actualData, err := ioutil.ReadAll(t.HTTP.Response.Body)
	if err != nil {
		return err
	}

	var actual map[string]json.RawMessage
	err = json.Unmarshal(actualData, &actual)
	if err != nil {
		return err
	}

	delete(expected, exclude)
	delete(actual, exclude)

	e, _ := json.Marshal(expected)
	a, _ := json.Marshal(actual)

	assert.Equal(t, e, a)
	return t.err
}

func (t *testContext) theDatabaseContainsAPollWithName(name string) error {
	row := t.DB.Client.QueryRow(
		`SELECT * from polls WHERE name=$1`,
		name,
	)
	var id string
	var pollName string
	var pollType string
	var options string
	return row.Scan(&id, &pollName, &pollType, &options)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	var testCtx testContext

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		psqlInfo := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)

		fmt.Println("# Resetting schemas...")

		var err error
		testCtx, err = NewTestContext(psqlInfo)
		if err != nil {
			return ctx, err
		}

		cmd := exec.Command("goose", "-v", "-dir", "../db/migrations", "postgres", psqlInfo, "reset")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Drop failed with: %s", string(output))
			return ctx, fmt.Errorf("unable to drop DB: %w", err)
		}

		cmd = exec.Command("goose", "-v", "-dir", "../db/migrations", "postgres", psqlInfo, "up")
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Printf("Migration failed with: %s", string(output))
			return ctx, fmt.Errorf("unable to migrate DB: %w", err)
		}

		if os.Getenv("DEBUG") == "1" {
			err = logTableNames(testCtx.DB.Client)
			if err != nil {
				return ctx, err
			}
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		testCtx.DB.Client.Close()

		return ctx, nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, testCtx.iSendRequestTo)
	ctx.Step(`^I make a "([^"]*)" request to "([^"]*)"$`, testCtx.iSendRequestTo)
	ctx.Step(`^the response code will be (\d+)$`, testCtx.theResponseCodeWillBe)
	ctx.Step(`^the body will be "([^"]*)"$`, testCtx.theBodyWillBe)
	ctx.Step(`^I have set data buffer to "([^"]*)"$`, testCtx.iHaveSetDataBufferTo)
	ctx.Step(`^I have a cookie "([^"]*)" with value "([^"]*)"$`, testCtx.iHaveACookieWithValue)
	ctx.Step(`^the body will look like "([^"]*)"$`, testCtx.theBodyWillLookLike)
	ctx.Step(`^the body will look like json "([^"]*)" excluding "([^"]*)"$`, testCtx.theBodyWillLookLikeJsonExcluding)

	ctx.Step(`^the database contains a poll with name "([^"]*)"$`, testCtx.theDatabaseContainsAPollWithName)
	ctx.Step(`^I have run the database script "([^"]*)"$`, testCtx.iHaveRunDatabaseScript)
	ctx.Step(`the response has a cookie "([^"]*)" that is present`, testCtx.responseCookieIsPresent)
	ctx.Step(`^the response has a header "([^"]*)" that matches "([^"]*)"$`, testCtx.theResponseHasAHeaderThatMatches)
}

func logTableNames(db *sql.DB) error {
	rows, err := db.Query(
		`SELECT 
			tablename
		FROM 
			pg_catalog.pg_tables
		ORDER BY
			tablename
		`,
	)
	if err != nil {
		return fmt.Errorf("unable to read DB tables: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			// handle this error
			if err != nil {
				return fmt.Errorf("unable to scan table row: %w", err)
			}
		}
		log.Printf("%s", name)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return fmt.Errorf("unable to read table row: %w", err)
	}

	return nil
}
