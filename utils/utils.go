package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaxIvanyshen/budget-tracker/types"
	"github.com/golang-jwt/jwt/v5"
)

func SendTelegramMessage(msg types.SupportMsg) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_TOKEN"))

	markdownMsg, err := buildMessage(msg)
	if err != nil {
		return err
	}
	reqBody, err := json.Marshal(map[string]string{
		"chat_id":    os.Getenv("TELEGRAM_CHAT_ID"),
		"text":       markdownMsg,
		"parse_mode": "Markdown",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func buildMessage(msg types.SupportMsg) (string, error) {
	return fmt.Sprintf("Name: *%s*\nEmail: *%s*\nSubject: *%s*\nMessage: *%s*", msg.Name, msg.Email, msg.Subject, msg.Message), nil
}

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userEmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": userEmail,
			"exp":   time.Now().Add(time.Hour * 24 * 14).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("password must contain at least one digit")
	}
	return nil
}
