package handlers

import (
	"fmt"
	controllers "jwt-practice/Controllers"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type userHandler struct {
	userController controllers.UserController
}

func NewUserHandler(c controllers.UserController) userHandler {
	return userHandler{userController: c}
}

func (handler userHandler) Signup(c *fiber.Ctx) error {
	request := controllers.SignupRequest{}
	// แปลง Body ให้เป็น SignupRequest
	err := c.BodyParser(&request)
	if err != nil {
		// CODE 406
		return c.SendStatus(fiber.ErrNotAcceptable.Code)
	}

	response, err := handler.userController.Signup(request)
	if err != nil {
		// CODE 422
		return c.SendStatus(fiber.ErrUnprocessableEntity.Code)
	}
	// CODE 201
	c.Status(fiber.StatusCreated)
	return c.JSON(response)
}

func (handler userHandler) Login(c *fiber.Ctx) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")

	request := controllers.LoginRequest{}
	err := c.BodyParser(&request)
	if err != nil {
		// CODE 406
		return c.SendStatus(fiber.ErrNotAcceptable.Code)
	}

	// เช็คข้อมูลที่รับกับ db
	user, err := handler.userController.Login(request)
	// ถ้าไม่ตรงกับ db
	if err != nil {
		// CODE 404
		fmt.Println("user not found")
		return c.SendStatus(fiber.ErrNotFound.Code)
	}

	claims := jwt.MapClaims{
		"iss": "issuer",
		"exp": time.Now().Local().Add(time.Hour * 1).Unix(),
		"data": map[string]string{
			"id":       strconv.Itoa(user.Id),
			"username": user.Username,
			"role":     user.Role,
		},
	}

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)

	// // ใส่ Secret key
	token, err := jwtToken.SignedString([]byte(os.Getenv("JWTSECRET")))
	if err != nil {
		return c.SendStatus(fiber.ErrInternalServerError.Code)
	}

	return c.JSON(fiber.Map{
		"jwtToken": token,
		// "ExpiresAt": cliams.ExpiresAt,
	})
}

func (handler userHandler) View(c *fiber.Ctx) error {
	fmt.Println("welcome")
	// รับข้อมูลที่ผ่านการตรวจสอบ token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	data := claims["data"]
	// ดึงข้อมูล role มาเช็คว่าเป็น admin หรือไม่
	role := (data.(map[string]interface{})["role"]).(string)
	role = strings.ToUpper(role)
	if role != "ADMIN" {
		return c.Redirect("/viewuser")
		// return fiber.NewError(fiber.StatusUnauthorized, "not admin")
	}
	return c.JSON("welcome " + role)
}

func (handler userHandler) ViewUser(c *fiber.Ctx) error {

	return c.JSON("welcome user")
}
