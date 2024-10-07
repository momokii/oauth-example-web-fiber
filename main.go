package main

import (
	"net/http"
	"os"
	"time"
	"try-oauth/controllers"
	"try-oauth/db"
	"try-oauth/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/login/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/login/github/callback"),
	)
	engine := html.New("./public", ".html")
	db.InitDB()

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).Render("errorPage", fiber.Map{
				"Error": err.Error(),
				"Code":  code,
			})
		},
	})
	middlewares.InitSession()
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:               60,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{}, // sliding window rate limiter,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).Render("errorPage", fiber.Map{
				"Error": "Too many requests, please try again later.",
				"Code":  fiber.StatusTooManyRequests,
			})
		},
	}))
	// app.Use(recover.New())
	app.Static("/public", "./public")

	// * ENDPOINTS

	// * auth
	app.Get("/", middlewares.IsAuth, controllers.DashboardView)

	app.Get("/login/:provider/callback", controllers.LoginOAuthCallback)
	app.Get("/login/google", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) { controllers.LoginOAuth(w, r, "google") }))
	app.Get("/login/github", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) { controllers.LoginOAuth(w, r, "github") }))
	// app.Get("/login/:provider", controllers.LoginOAuth)
	app.Get("/login", controllers.LoginView)
	app.Post("/login", controllers.LoginPost)

	app.Get("/signup", controllers.SignupView)
	app.Post("/signup", controllers.SignupPost)

	app.Post("/logout", middlewares.IsAuth, controllers.LogoutPost)

	app.Listen(":3000")

}
