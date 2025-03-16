package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MaxIvanyshen/budget-tracker/types"
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
