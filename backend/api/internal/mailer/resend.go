package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

// resendMailer sends email through the Resend HTTP API.
// https://resend.com/docs/api-reference/emails/send-email
type resendMailer struct {
	apiKey     string
	from       string
	httpClient *http.Client
}

// NewResendMailer constructs a Mailer backed by Resend. from should be a
// verified sender address (e.g. "Caretree <invites@yourdomain.com>");
// Resend's shared "onboarding@resend.dev" sandbox sender works without
// domain verification for testing.
func NewResendMailer(apiKey, from string) Mailer {
	return &resendMailer{
		apiKey:     apiKey,
		from:       from,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

type resendErrorResponse struct {
	Message string `json:"message"`
}

func (m *resendMailer) SendFamilyInvite(ctx context.Context, email FamilyInviteEmail) error {
	subject := fmt.Sprintf("%s invited you to join %s on Caretree", email.InviterName, email.FamilyName)

	body := resendRequest{
		From:    m.from,
		To:      []string{email.To},
		Subject: subject,
		HTML:    familyInviteHTML(email),
		Text:    familyInviteText(email),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding resend request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("building resend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		var errResp resendErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
			return fmt.Errorf("resend returned %d: %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("resend returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (m *resendMailer) SendTaskAssigned(ctx context.Context, email TaskAssignedEmail) error {
	subject := fmt.Sprintf("You've been assigned a task on Caretree: %s", email.TaskTitle)

	body := resendRequest{
		From:    m.from,
		To:      []string{email.To},
		Subject: subject,
		HTML:    taskAssignedHTML(email),
		Text:    taskAssignedText(email),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding resend request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("building resend request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		var errResp resendErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
			return fmt.Errorf("resend returned %d: %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("resend returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func taskAssignedText(email TaskAssignedEmail) string {
	return fmt.Sprintf(
		"%s assigned you a task on Caretree: %q",
		email.AssignedBy, email.TaskTitle,
	)
}

func taskAssignedHTML(email TaskAssignedEmail) string {
	return fmt.Sprintf(`<p>%s assigned you a task on Caretree:</p>
<p style="font-size:18px;font-weight:600;">%s</p>`,
		email.AssignedBy, email.TaskTitle,
	)
}

func familyInviteText(email FamilyInviteEmail) string {
	return fmt.Sprintf(
		"%s invited you to join %q on Caretree.\n\nUse invite code %s to join: open the app, choose \"Join a family\", and enter this code.",
		email.InviterName, email.FamilyName, email.InviteCode,
	)
}

func familyInviteHTML(email FamilyInviteEmail) string {
	return fmt.Sprintf(`<p>%s invited you to join <strong>%s</strong> on Caretree.</p>
<p>Open the app, choose "Join a family", and enter this invite code:</p>
<p style="font-size:24px;font-weight:600;letter-spacing:2px;">%s</p>`,
		email.InviterName, email.FamilyName, email.InviteCode,
	)
}
