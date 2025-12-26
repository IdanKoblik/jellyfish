package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NotificationRequest struct {
	DeviceID string   `json:"device"`
	Message  string   `json:"message"`
	Targets  []string `json:"targets"`
}

func SendNotification(cfg *Config, title string, notification *Notification) error {
	reqURL := fmt.Sprintf("%s/api/send", cfg.WhisperURI)

	message := BuildSMSMessage(title, notification)
	reqBody := NotificationRequest{
		DeviceID: cfg.DeviceID,
		Message:  message,
		Targets:  []string{ExtractPrefix(notification.Request.RequestedByEmail)},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Token", cfg.WhisperToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non-200 response: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func ExtractPrefix(email string) string {
	if email == "" || !strings.Contains(email, "@") {
		return ""
	}
	parts := strings.Split(email, "@")
	return parts[0]
}

func BuildSMSMessage(title string, n *Notification) string {
	var b strings.Builder
	if n.Subject != "" {
		b.WriteString(fmt.Sprintf("%s — %s\n", title, n.Subject))
	} else {
		b.WriteString(title + "\n")
	}

	if n.Event != "" {
		b.WriteString(n.Event + "\n")
	}

	if n.Message != "" {
		b.WriteString(n.Message + "\n")
	}

	if n.Request != nil {
		if n.Request.RequestedByUsername != "" {
			b.WriteString("Requested by " + n.Request.RequestedByUsername + "\n")
		} else if n.Request.RequestedByEmail != "" {
			b.WriteString("Requested by " + ExtractPrefix(n.Request.RequestedByEmail) + "\n")
		}
	}

	if n.Media != nil {
		meta := []string{}

		if n.Media.MediaType != "" {
			meta = append(meta, n.Media.MediaType)
		}
		if n.Media.Status != "" {
			meta = append(meta, n.Media.Status)
		}

		if len(meta) > 0 {
			b.WriteString(strings.Join(meta, " / "))
		}

		ids := []string{}
		if n.Media.TmdbID != "" {
			ids = append(ids, "TMDB: "+n.Media.TmdbID)
		}
		if n.Media.TvdbID != "" {
			ids = append(ids, "TVDB: "+n.Media.TvdbID)
		}

		if len(ids) > 0 {
			if len(meta) > 0 {
				b.WriteString(" | ")
			}
			b.WriteString(strings.Join(ids, " "))
			b.WriteString("\n")
		}
	}

	for _, e := range n.Extra {
		b.WriteString(e.Name + ": " + e.Value + "\n")
	}

	return strings.TrimSpace(b.String())
}
