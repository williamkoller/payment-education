package infra_email

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResendClient struct {
	mock.Mock
}

func (m *MockResendClient) SendEmail(to, subject, html string) error {
	args := m.Called(to, subject, html)
	return args.Error(0)
}

func TestNewResendEmailNotifier(t *testing.T) {
	mockClient := new(MockResendClient)

	notifier := NewResendEmailNotifier(mockClient)

	assert.NotNil(t, notifier)
	assert.IsType(t, &ResendEmailNotifier{}, notifier)
}

func TestResendEmailNotifier_SendWelcomeEmail_Success(t *testing.T) {
	mockClient := new(MockResendClient)
	notifier := NewResendEmailNotifier(mockClient)

	name := "John Doe"
	email := "john.doe@example.com"

	mockClient.On("SendEmail",
		email,
		"Bem-vindo á Plataforma",
		mock.MatchedBy(func(html string) bool {
			return len(html) > 0
		}),
	).Return(nil)

	err := notifier.SendWelcomeEmail(name, email)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestResendEmailNotifier_SendWelcomeEmail_EmailContainsUserName(t *testing.T) {
	mockClient := new(MockResendClient)
	notifier := NewResendEmailNotifier(mockClient)

	name := "Maria Silva"
	email := "maria.silva@example.com"

	var capturedHTML string
	mockClient.On("SendEmail",
		email,
		"Bem-vindo á Plataforma",
		mock.MatchedBy(func(html string) bool {
			capturedHTML = html
			return true
		}),
	).Return(nil)

	err := notifier.SendWelcomeEmail(name, email)

	assert.NoError(t, err)
	assert.Contains(t, capturedHTML, name, "HTML should contain the user's name")
	assert.Contains(t, capturedHTML, "Olá, Maria Silva!", "HTML should have personalized greeting")
	mockClient.AssertExpectations(t)
}

func TestResendEmailNotifier_SendWelcomeEmail_ClientError(t *testing.T) {
	mockClient := new(MockResendClient)
	notifier := NewResendEmailNotifier(mockClient)

	name := "John Doe"
	email := "john.doe@example.com"
	expectedError := errors.New("resend API error")

	mockClient.On("SendEmail",
		email,
		"Bem-vindo á Plataforma",
		mock.Anything,
	).Return(expectedError)

	err := notifier.SendWelcomeEmail(name, email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send welcome email")
	assert.Contains(t, err.Error(), email)
	mockClient.AssertExpectations(t)
}

func TestResendEmailNotifier_SendWelcomeEmail_HTMLStructure(t *testing.T) {
	mockClient := new(MockResendClient)
	notifier := NewResendEmailNotifier(mockClient)

	name := "Test User"
	email := "test@example.com"

	var capturedHTML string
	mockClient.On("SendEmail",
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(html string) bool {
			capturedHTML = html
			return true
		}),
	).Return(nil)

	err := notifier.SendWelcomeEmail(name, email)

	assert.NoError(t, err)
	assert.Contains(t, capturedHTML, "<!DOCTYPE html", "Should be valid HTML")
	mockClient.AssertExpectations(t)
}

func TestResendEmailNotifier_SendWelcomeEmail_CorrectSubject(t *testing.T) {
	mockClient := new(MockResendClient)
	notifier := NewResendEmailNotifier(mockClient)

	name := "Test User"
	email := "test@example.com"
	expectedSubject := "Bem-vindo á Plataforma"

	mockClient.On("SendEmail",
		email,
		expectedSubject,
		mock.Anything,
	).Return(nil)

	err := notifier.SendWelcomeEmail(name, email)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockClient.AssertCalled(t, "SendEmail", email, expectedSubject, mock.Anything)
}
