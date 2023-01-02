package controllers

import (
	"database/sql"
	model "jwt-practice/Models"
	"jwt-practice/logs"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type userController struct {
	userRepository model.UserRepository
}

func NewUserController(r model.UserRepository) UserController {
	return userController{userRepository: r}
}

func (c userController) Signup(req SignupRequest) (*SignupResponse, error) {
	// Validation
	if req.Username == "" || req.Password == "" {
		logs.Error("Invalid Username or Password")
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid Username or Password")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		logs.Error(err)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	user := model.User{
		Username: req.Username,
		Password: string(password),
		Role:     "User",
	}

	newUser, err := c.userRepository.Create(user)
	if err != nil {
		logs.Error(err)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	res := SignupResponse{
		Username: newUser.Username,
		Role:     newUser.Role,
	}

	return &res, nil
}

func (c userController) Login(req LoginRequest) (*LoginResponse, error) {
	// Validation
	if req.Username == "" || req.Password == "" {
		logs.Error("Invalid Username or Password")
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	checkUser := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := c.userRepository.Login(checkUser)

	if err != nil {
		logs.Error(err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Incorrect Username or Password")
	}

	res := LoginResponse{
		Id:       user.Id,
		Username: user.Username,
		Role:     user.Role,
	}

	return &res, nil
}

func (c userController) View(id int) (*LoginResponse, error) {
	user, err := c.userRepository.Find(id)

	if err != nil {
		if err == sql.ErrNoRows {
			logs.Error(err)
			return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "User not found")
		}
		logs.Error(err)
		return nil, fiber.NewError(fiber.StatusExpectationFailed)
	}

	res := LoginResponse{
		Id:       user.Id,
		Username: user.Username,
		Role:     user.Role,
	}

	return &res, nil
}
