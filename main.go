package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting reservation service")
	notificationsBaseURL := "http://localhost:8010"
	notificationClient := NewNotificationClient(notificationsBaseURL)
	reservationService := NewReservationService(notificationClient)
	err := reservationService.Reserve("test@wp.pl")
	if err != nil {
		fmt.Printf("Failed to make reservation: %v", err)
	}
	return
}

type ReservationService struct {
	notificationClient *NotificationClient
}

func NewReservationService(notificationClient *NotificationClient) *ReservationService {
	return &ReservationService{
		notificationClient: notificationClient,
	}
}

func (rs *ReservationService) Reserve(email string) error {
	// here should be logic representing reservation

	err := rs.notificationClient.SendEmail(EmailMessage{
		Address: email,
		Title:   "Reservation",
		Body:    "You have successfully reserved a table",
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Println(fmt.Sprintf("Reservation made for %s", email))
	return err
}

type NotificationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewNotificationClient(baseURL string) *NotificationClient {
	return &NotificationClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *NotificationClient) SendEmail(message EmailMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/emails", c.baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	fmt.Println(fmt.Sprintf("Reservation Email sent to %s", message.Address))
	return nil
}

type EmailMessage struct {
	Address string `json:"address"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}
