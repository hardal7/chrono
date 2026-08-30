package user

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hardal7/chrono/internal/db"
	"github.com/hardal7/chrono/internal/dto"
	"github.com/hardal7/chrono/internal/util/config"
	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/mailgun/mailgun-go/v5"
)

const otpExpirationMinutes = 5

func ResetPassword(ctx context.Context, r dto.ResetUserPasswordRequest) error {
	if r.Email == "" && r.Username == "" {
		return fmt.Errorf("Both email and username fields cannot be empty")
	}

	var err error
	email := r.Email
	if r.Email == "" {
		email, err = getEmail(ctx, r.Username)
		if err != nil {
			return err
		}
	}

	err = sendResetEmail(email)

	return err
}

func getEmail(ctx context.Context, username string) (string, error) {
	user, err := db.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}

func sendResetEmail(email string) error {
	mg := mailgun.NewMailgun(config.App.MAIL_API_KEY)
	err := mg.SetAPIBase(mailgun.APIBaseEU)
	if err != nil {
		logger.Warn("Failed to set API Base for email")
		return err
	}

	token, err := generateResetToken(email)
	if err != nil {
		return err
	}
	domain := config.App.Domain

	recipient := email
	sender := config.App.MailAddress
	subject := "Password Reset"
	body := "Click this link to reset your password: https://" + domain + "/reset-password?token=" + token

	data, err := os.ReadFile("static/reset-password.html")
	if err != nil {
		return err
	}
	html := string(data)
	html = strings.ReplaceAll(html, "RESET_TOKEN", token)
	html = strings.ReplaceAll(html, "DOMAIN_NAME", domain)

	message := mailgun.NewMessage(domain, sender, subject, body, recipient)
	message.SetHTML(html)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err = mg.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func generateResetToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     email,
		"exp":     time.Now().Add(time.Minute * time.Duration(otpExpirationMinutes)).Unix(),
		"purpose": "password_reset",
	})
	tokenString, err := token.SignedString([]byte(config.App.HashSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO: Invalidate past session tokens on password reset
// Also implement in account endpoint

// TODO: Expire token after consuming, not only after expiration
