package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lassejlv/go-echo/utils"
	_ "github.com/lib/pq"
)

type (
	ValidateUser struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	ValidatePost struct {
		Title   string `json:"title" form:"title" validate:"required"`
		Content string `json:"content" form:"content" validate:"required"`
		UserId  int    `json:"user_id" form:"user_id" validate:"required"`
	}

	User struct {
		ID           int    `json:"id" db:"id"`
		Username     string `json:"username" db:"username"`
		PasswordHash string `json:"password_hash" db:"password_hash"`
		CreatedAt    string `json:"created_at" db:"created_at"`
		UpdatedAt    string `json:"updated_at" db:"updated_at"`
	}

	Post struct {
		Title string `json:"title" form:"title" validate:"required"`
		Body  string `json:"body" form:"body" validate:"required"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	godotenv.Load()

	// Setup db
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("Open error: %v", err)
		log.Fatalln(err)
	}

	// Test the connection with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Ping error: %v", err)
		log.Fatalln(err)
	}

	app := echo.New()

	app.Validator = &CustomValidator{validator: validator.New()}

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("OS: %s, ARCH: %s", runtime.GOOS, runtime.GOARCH))
	})

	app.GET("/health", func(c echo.Context) error {
		is_db_okay := db.Ping()

		if is_db_okay != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, is_db_okay.Error())
		}

		return c.String(http.StatusOK, "OK")
	})

	app.GET("/stats", func(c echo.Context) error {

		db_stats := db.Stats()

		return c.JSON(http.StatusOK, db_stats)
	})

	app.GET("/users", func(c echo.Context) error {

		var users []User
		err := db.Select(&users, "SELECT * FROM users")

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	})

	app.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user User
		err := db.Get(&user, "SELECT * FROM users WHERE id = $1", id)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	})

	app.POST("/users", func(c echo.Context) error {

		var user ValidateUser
		if err := c.Bind(&user); err != nil {
			return err
		}

		if err := c.Validate(&user); err != nil {
			return err
		}

		hashed_pass, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		var createdUser User
		err = db.QueryRowx(
			"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *",
			user.Username,
			hashed_pass,
		).StructScan(&createdUser)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, createdUser)

	})

	app.Logger.Fatal(app.Start(":8080"))
}
