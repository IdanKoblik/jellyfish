package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	NotificationType string       `json:"notification_type"`
	Event            string       `json:"event"`
	Subject          string       `json:"subject"`
	Message          string       `json:"message"`
	Image            string       `json:"image"`
	Media            *Media       `json:"media"`
	Request          *RequestBody `json:"request"`
	Extra            []ExtraData  `json:"extra"`
}

type Media struct {
	MediaType string `json:"media_type"`
	TmdbID    string `json:"tmdbId"`
	TvdbID    string `json:"tvdbId"`
	Status    string `json:"status"`
}

type RequestBody struct {
	RequestID           string `json:"request_id"`
	RequestedByEmail    string `json:"requestedBy_email"`
	RequestedByUsername string `json:"requestedBy_username"`
}

type ExtraData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func HandleWebhook(cfg *Config, c *gin.Context) {
	rawBody, err := c.GetRawData()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read body"})
		return
	}

	var notification Notification
	if err := json.Unmarshal(rawBody, &notification); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		fmt.Printf("Invalid JSON, details: %s\n", err.Error())
		return
	}

	switch strings.ToUpper(notification.NotificationType) {
	case "MEDIA_PENDING":
		if notification.Request != nil {
			for _, admin := range cfg.Admins {
				if err := SendNotification(cfg, "New Media Request", &notification); err != nil {
					fmt.Printf("Failed to notify admin %s: %v", admin, err)
					c.JSON(400, err.Error())
				}
			}
		}

	case "MEDIA_APPROVED":
		if notification.Request != nil {
			if err := SendNotification(cfg, "Request Approved", &notification); err != nil {
				log.Printf("Failed to notify requester: %v", err)
				c.JSON(400, err.Error())
			}
		}

	case "MEDIA_DECLINED":
		if notification.Request != nil {
			if err := SendNotification(cfg, "Request Declined", &notification); err != nil {
				fmt.Printf("Failed to notify requester: %v", err)
				c.JSON(400, err.Error())
			}
		}

	case "MEDIA_AVAILABLE":
		if notification.Media != nil {
			if err := SendNotification(cfg, "Now Available", &notification); err != nil {
				fmt.Printf("Failed to notify requester: %v", err)
				c.JSON(400, err.Error())
			}
		}

	case "MEDIA_FAILED":
		if err := SendNotification(cfg, "Download Failed", &notification); err != nil {
			fmt.Printf("Failed to notify requester: %v", err)
			c.JSON(400, err.Error())
		}
	}

	c.JSON(200, gin.H{"status": "ok"})
}
