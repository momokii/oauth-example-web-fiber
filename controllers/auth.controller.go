package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"try-oauth/middlewares"
	"try-oauth/models"
	"try-oauth/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

func AuthSessChecker(c *fiber.Ctx) error {
	cu, err := middlewares.CheckSession(c, "username")
	if err != nil {
		return utils.ErrorJSON(c, 500, err.Error())
	}

	if cu != nil {
		return c.Redirect("/")
	}

	return nil
}

func SignupView(c *fiber.Ctx) error {
	AuthSessChecker(c)
	return c.Render("signup", fiber.Map{})
}

func LoginView(c *fiber.Ctx) error {
	AuthSessChecker(c)
	return c.Render("login", fiber.Map{})
}

func SignupPost(c *fiber.Ctx) error {
	AuthSessChecker(c)

	username := c.FormValue("username")
	password := c.FormValue("password")

	userNew := models.User{}
	userCheck := models.User{
		Username: username,
		Password: password,
	}

	err := utils.ValidateStruct(userCheck)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Password":
				return utils.ErrorJSON(c, 400, "Password minimal 6 karakter, mengandung angka dan huruf besar")
			case "Username":
				return utils.ErrorJSON(c, 400, "Username minimal 5 karakter, merupakan alphanumerik")
			default:
				return utils.ErrorJSON(c, 400, err.Error())
			}
		}
	}

	checkUser := userNew.CheckUserByUsername(username)
	if checkUser.Username != "" {
		return utils.ErrorJSON(c, 400, "Username sudah terdaftar")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return utils.ErrorJSON(c, 500, err.Error())
	}

	userNew.Password = string(hashedPass)
	userNew.Username = username

	err = userNew.CreateUser()
	if err != nil {
		return utils.ErrorJSON(c, 500, err.Error())
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Berhasil buat akun baru",
	})
}

func LoginPost(c *fiber.Ctx) error {
	AuthSessChecker(c)

	username := c.FormValue("username")
	password := c.FormValue("password")
	user := models.User{}

	checkUser := user.CheckUserByUsername(username)
	if checkUser.Username == "" {
		return utils.ErrorJSON(c, 400, "Username/Password tidak salah")
	}

	err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(password))
	if err != nil {
		return utils.ErrorJSON(c, 400, "Username/Password salah")
	}

	middlewares.CreateSession(c, "username", username)

	return c.JSON(fiber.Map{
		"error": false,
	})
}

func LoginOAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	ctx := context.WithValue(c.Context(), "provider", provider)

	// Membuat http.Request dari fiber.Ctx
	r, err := http.NewRequest(c.Method(), c.OriginalURL(), bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Gagal membuat request",
		})
	}

	// Menyalin header dari Fiber request ke http.Request
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	r = r.WithContext(ctx)      // Menggunakan context yang telah dimodifikasi
	w := httptest.NewRecorder() // Membuat response writer

	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return utils.ErrorJSON(c, 500, err.Error())
	}

	user := models.User{}
	var checkUser *models.User
	var usernameSess string
	if provider == "google" {
		checkUser = user.CheckUserByUsername(u.Email) // google using email
	} else {
		fmt.Println("nick: ", u.NickName)
		checkUser = user.CheckUserByUsername(u.NickName) // if github get the nickname
	}

	if checkUser.Username == "" {
		if provider == "google" {
			user.Username = u.Email
		} else {
			user.Username = u.NickName
		}
		err = user.CreateUser()
		if err != nil {
			return utils.ErrorJSON(c, 500, err.Error())
		}
	}

	usernameSess = user.Username

	gothic.Logout(w, r)
	middlewares.CreateSession(c, "username", usernameSess)

	return c.Redirect("/", fiber.StatusSeeOther)
}

func LoginOAuth(w http.ResponseWriter, r *http.Request, provider string) {
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		fmt.Println("error sini: ", gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func LogoutPost(c *fiber.Ctx) error {
	middlewares.DeleteSession(c)

	return c.Redirect("/login")
}
