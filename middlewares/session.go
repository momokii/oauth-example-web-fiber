package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func InitSession() {
	Store = session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
		// CookieSameSite: "strict", // bcs using oauth not using it
	})
	fmt.Println("Session initialized")
}

func CreateSession(c *fiber.Ctx, key string, val interface{}) error {
	session, err := Store.Get(c)
	if err != nil {
		return err
	}
	defer session.Save()

	session.Set(key, val)

	return nil
}

func DeleteSession(c *fiber.Ctx) error {
	session, err := Store.Get(c)
	if err != nil {
		return err
	}

	session.Destroy()

	return nil
}

func CheckSession(c *fiber.Ctx, key string) (interface{}, error) {
	session, err := Store.Get(c)
	if err != nil {
		return nil, err
	}

	return session.Get(key), nil
}

func IsAuth(c *fiber.Ctx) error {
	user, err := CheckSession(c, "username")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if user == nil {
		return c.Redirect("/login")
	}

	// with username search user data in db

	// pass data to next middleware
	c.Locals("username", user)

	// renew session
	CreateSession(c, "username", user)

	return c.Next()
}
