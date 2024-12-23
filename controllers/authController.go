package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	gen "math/rand/v2"
	"os"
	"strconv"
	"time"

	"QuickPicsAuth/database"
	"QuickPicsAuth/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Pix   string `json:"pix"`
	CPF   string `json:"cpf"`
}

func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, World!"})
}

var validate = validator.New()

func Register(c *fiber.Ctx) error {

	fmt.Println("Received")

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User

	if err := database.DB.Where("email = ? OR cpf = ? OR pix = ?", data["email"], data["cpf"], data["pix"]).First(&user).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot hash password",
		})
	}

	userModel := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Pix:      data["pix"],
		CPF:      data["cpf"],
		Password: string(hashedPassword),
	}

	if err := validate.Struct(userModel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error on email, cpf or pix",
		})
	} else {
		if err = database.DB.Create(&userModel).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to register user",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "User has been registered",
	})

}

func Login(c *fiber.Ctx) error {

	fmt.Println("login request")

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		fmt.Println("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		fmt.Println("Incorrect password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("erro ao gerar o token", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	cookie := fiber.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Login successful",
	})
}

func User(c *fiber.Ctx) error {

	fmt.Println("User request")

	cookie := c.Cookies("session")

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthenticated",
		})
	}

	claims, ok := token.Claims.(*jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthenticated",
		})
	}

	id, _ := strconv.Atoi((*claims)["sub"].(string))
	user := models.User{ID: uint(id)}
	var userResponse UserResponse

	if err := database.DB.Model(&user).Select("id", "name", "email", "pix", "cpf").Scan(&userResponse).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(userResponse)
}

func ForgotPassword(c *fiber.Ctx) error {

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User

	randomNumber := func(min, max int) int {
		return gen.IntN(max-min) + min
	}

	// generate := randomNumber(1000, 9999)

	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		time.Sleep(time.Millisecond * time.Duration(randomNumber(1800, 3600)))
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": "Password reset link has been sent to your email",
		})
	}

	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Erro ao gerar bytes aleatórios:", err)
	}

	token := hex.EncodeToString(randomBytes)

	database.DB.Model(&user).Update("password_reset_token", token)
	database.DB.Model(&user).Update("password_reset_expires", time.Now().Add(time.Hour*1))

	message := gomail.NewMessage()
	message.SetHeader("From", "no-reply@quickpics.ai")
	message.SetHeader("To", data["email"])
	message.SetHeader("Subject", "Reset Password Notification")

	url := "http://localhost:3000/reset-password/" + token
	htmlFilePath := "./resources/mail/auth/forgot_password.html"
	htmlContent, err := os.ReadFile(htmlFilePath)
	if err != nil {
		panic(fmt.Sprintf("Erro ao ler o arquivo HTML: %v", err))
	}

	message.SetBody("text/html", fmt.Sprintf(string(htmlContent), user.Name, url, 1))

	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), 2525, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

	return c.JSON(fiber.Map{
		"message": "Password reset link has been sent to your email",
	})
}

func Logout(c *fiber.Ctx) error {

	cookie := fiber.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
