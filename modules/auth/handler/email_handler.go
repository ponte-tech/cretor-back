package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/ponte-tech/cretor-back/shared/response"
	"github.com/ponte-tech/cretor-back/shared/validator"
	"go.uber.org/zap"
)

type EmailHandler struct {
	logger *zap.Logger
}

func NewEmailHandler(logger *zap.Logger) *EmailHandler {
	return &EmailHandler{logger: logger}
}

type SendEmailRequest struct {
	To         string `json:"to" validate:"required,email"`
	Subject    string `json:"subject" validate:"required"`
	Body       string `json:"body" validate:"required"`
	Attachment string `json:"attachment"`
	FileName   string `json:"file_name"`
}

// POST /email/send
func (h *EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	from := "daniel@danielkrammes.com"

	if req.Attachment != "" {
		if err := h.sendWithAttachment(ctx, from, req); err != nil {
			h.logger.Error("failed to send email with attachment", zap.Error(err))
			response.Error(w, "failed to send email", http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.sendSimple(ctx, from, req); err != nil {
			h.logger.Error("failed to send email", zap.Error(err))
			response.Error(w, "failed to send email", http.StatusInternalServerError)
			return
		}
	}

	h.logger.Info("email sent", zap.String("to", req.To), zap.String("subject", req.Subject))
	response.Message(w, "email sent", http.StatusOK)
}

func (h *EmailHandler) getSESClient(ctx context.Context) (*sesv2.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}
	return sesv2.NewFromConfig(cfg), nil
}

func (h *EmailHandler) sendSimple(ctx context.Context, from string, req SendEmailRequest) error {
	client, err := h.getSESClient(ctx)
	if err != nil {
		return err
	}

	replyTo := "danielkrammes27@gmail.com"
	_, err = client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress:   &from,
		ReplyToAddresses:   []string{replyTo},
		Destination: &types.Destination{
			ToAddresses: []string{req.To},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: &req.Subject},
				Body: &types.Body{
					Html: &types.Content{Data: &req.Body},
				},
			},
		},
	})
	return err
}

func (h *EmailHandler) sendWithAttachment(ctx context.Context, from string, req SendEmailRequest) error {
	client, err := h.getSESClient(ctx)
	if err != nil {
		return err
	}

	attachData := req.Attachment
	if idx := strings.Index(attachData, ","); idx >= 0 {
		attachData = attachData[idx+1:]
	}

	fileName := req.FileName
	if fileName == "" {
		fileName = "proposta.pdf"
	}

	boundary := "----=_Part_boundary_cretor"

	replyTo := "danielkrammes27@gmail.com"

	rawMsg := fmt.Sprintf("From: %s\r\nTo: %s\r\nReply-To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n--%s\r\nContent-Type: text/html; charset=UTF-8\r\nContent-Transfer-Encoding: 7bit\r\n\r\n%s\r\n\r\n--%s\r\nContent-Type: application/pdf; name=\"%s\"\r\nContent-Disposition: attachment; filename=\"%s\"\r\nContent-Transfer-Encoding: base64\r\n\r\n%s\r\n\r\n--%s--\r\n",
		from, req.To, replyTo, req.Subject, boundary,
		boundary, req.Body,
		boundary, fileName, fileName, wrapBase64(attachData),
		boundary,
	)

	rawBytes := []byte(rawMsg)

	_, err = client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: &from,
		Destination: &types.Destination{
			ToAddresses: []string{req.To},
		},
		Content: &types.EmailContent{
			Raw: &types.RawMessage{
				Data: rawBytes,
			},
		},
	})
	return err
}

func wrapBase64(data string) string {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data
	}
	encoded := base64.StdEncoding.EncodeToString(decoded)
	var wrapped strings.Builder
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		wrapped.WriteString(encoded[i:end])
		wrapped.WriteString("\r\n")
	}
	return wrapped.String()
}
