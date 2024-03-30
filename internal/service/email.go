package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yosakoo/task-traker/pkg/rabbitmq"
)

type EmailService struct {
	queueConn *rabbitmq.Connection
}

func NewEmailService(queueConn *rabbitmq.Connection) *EmailService {
	return &EmailService{
		queueConn: queueConn,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, email *Email) error {
	data, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %w", err)
	}
	fmt.Print(data)
	if err := s.queueConn.PublishMessage(ctx, "application/json", data); err != nil {
		return fmt.Errorf("failed to send email message to queue: %w", err)
	}

	return nil
}
