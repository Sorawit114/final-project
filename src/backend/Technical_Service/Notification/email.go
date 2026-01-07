package Notification

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type NotificationService struct {
	dialer *gomail.Dialer
}

func NewNotificationService() *NotificationService {
	// ⚠️ Using hardcoded credentials for DEMO. In production, use Environment Variables.
	// You need to replace these with real credentials to send actual emails.
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	email := "your_email@gmail.com"
	password := "your_app_password"

	d := gomail.NewDialer(smtpHost, smtpPort, email, password)

	return &NotificationService{
		dialer: d,
	}
}

func (s *NotificationService) SendEmail(to []string, stockName string, pattern string, confidence float64, time string) error {
	m := gomail.NewMessage()
	// Set generic sender
	m.SetHeader("From", "admin@stockapp.com")
	m.SetHeader("To", to...)
	m.SetHeader("Subject", fmt.Sprintf("Stock Alert: %s - %s Detected!", stockName, pattern))

	body := fmt.Sprintf(`
	<h2>Stock Pattern Alert</h2>
	<p>Hello,</p>
	<p>A new pattern has been detected for your monitored stock.</p>
	<ul>
		<li><strong>Stock:</strong> %s</li>
		<li><strong>Pattern:</strong> %s</li>
		<li><strong>Confidence:</strong> %.2f%%</li>
		<li><strong>Time:</strong> %s</li>
	</ul>
	<p>Please check your dashboard for more details.</p>
	`, stockName, pattern, confidence*100, time)

	m.SetBody("text/html", body)

	// if err := s.dialer.DialAndSend(m); err != nil {
	// 	// For development without real creds, just log it.
	// 	// fmt.Println("❌ Email failed (expected without creds):", err)
	// 	// return err
	// }

	// For demo purposes, we just print to console to verify logic workflow
	fmt.Printf("\n📧 [MOCK EMAIL] To: %v | Subject: Stock Alert %s - %s\n", to, stockName, pattern)
	return nil
}
