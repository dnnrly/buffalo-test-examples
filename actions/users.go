package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"buffalo_test_examples/models"
)

// UsersNew renders the users form
func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
	}

	c.Session().Set("current_user_id", u.ID)

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Flash().Add("success", "Welcome to buffalo-test-examples!")
		return c.Redirect(http.StatusFound, "/")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Redirect(http.StatusFound, "/users/"+u.ID.String())
	}).Respond(c)
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				c.Logger().Warnf("user attempted to access with current_user_id '%v' that is not found: %v", uid, err)

				c.Session().Delete("current_user_id")
				c.Session().Set("redirectURL", c.Request().URL.String())
				c.Flash().Add("danger", "You must be authorized with a correct user to see that page")
				return c.Redirect(http.StatusFound, "/auth/new")
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// UsersShow default implementation.
func UsersShow(c buffalo.Context) error {
	uid := uuid.Must(uuid.FromString(c.Param("id")))

	u := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Find(u, uid)
	if err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("email", u.Email)
		return c.Render(http.StatusOK, r.HTML("users/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(u))
	}).Respond(c)
}

// UsersMe default implementation.
func UsersMe(c buffalo.Context) error {
	uid := c.Session().Get("current_user_id")
	if uid == nil {
		return c.Redirect(http.StatusFound, "/auth/new")
	}

	u := &models.User{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Find(u, uid)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("email", u.Email)
	return c.Render(http.StatusOK, r.HTML("users/me.html"))
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			return responder.Wants("html", func(c buffalo.Context) error {
				c.Flash().Add("danger", "You must be authorized to see that page")
				return c.Redirect(http.StatusFound, "/auth/new")
			}).Wants("json", func(c buffalo.Context) error {
				return c.Render(200, r.JSON(errors.New("not authorized")))
			}).Respond(c)
		}
		return next(c)
	}
}
