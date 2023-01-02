package main

import (
	"fmt"
	controllers "jwt-practice/Controllers"
	handlers "jwt-practice/Handlers"
	models "jwt-practice/Models"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// https://www.sohamkamani.com/golang/jwt-authentication/

// https://www.bacancytechnology.com/blog/golang-jwt

func iniConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// viper.AutomaticEnv()
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func initDatabase() *sqlx.DB {

	DSN := fmt.Sprintf("%v:%v@/%v",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.database"),
	)
	db, err := sqlx.Open(viper.GetString("db.driver"), DSN)

	if err != nil {
		panic(err)
	}

	// db.SetConnMaxIdleTime(3 * time.Minute)
	// db.SetMaxOpenConns(3)
	// db.SetMaxIdleConns(3)

	return db
}

func main() {

	initTimeZone()
	iniConfig()
	db := initDatabase()

	userRepository := models.NewUserRepositoryDB(db)
	userService := controllers.NewUserController(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	app := fiber.New()

	app.Post("/signup", userHandler.Signup)
	app.Post("/login", userHandler.Login)

	// middleware
	app.Use("/view", LoginAuth())
	app.Get("/view", userHandler.View)
	app.Get("/viewuser", userHandler.ViewUser)
	app.Listen(":8000")
}

func LoginAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWTSECRET")),
		// ถ้าผ่าน
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		// ถ้าไม่ผ่าน
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// CODE 401
			return fiber.ErrUnauthorized
		},
	})
}

func isAuth2(c *fiber.Ctx) error {
	jwtware.New(jwtware.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(os.Getenv("JWTSECRET")),
		// ถ้าผ่าน
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		// ถ้าไม่ผ่าน
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// CODE 401
			return fiber.ErrUnauthorized
		},
	})
	return nil
}
