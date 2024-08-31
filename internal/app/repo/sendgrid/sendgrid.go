package sendgrid

import (
	"context"
	"fmt"
	"log/slog"
	"mahir-trade-be/internal/app/models"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/dig"
)

type (
	SendgridRepo interface {
		SendEmail(ctx context.Context, req models.SendgridSendEmailRequest) (err error)
	}

	SendgridRepoImpl struct {
		dig.In
	}
)

func NewSendgridRepo(impl SendgridRepoImpl) SendgridRepo {
	return &impl
}

func (s *SendgridRepoImpl) SendEmail(ctx context.Context, req models.SendgridSendEmailRequest) (err error) {
	var (
		mailer *mail.SGMailV3
	)
	apiKey := os.Getenv("SENDGRID_API_KEY")

	if apiKey == "" {
		slog.ErrorContext(ctx, "[SendgridRepoImpl.SendEmail] failed to send email", "empty api key", err)
		err = fmt.Errorf("failed to send email")
		return
	}

	from := mail.NewEmail(req.SenderName, req.From)
	to := mail.NewEmail(req.RecepientName, req.To)

	if req.Attach.Content != "" {
		attachment := mail.NewAttachment()
		attachment.SetContent(req.Attach.Content)
		attachment.SetType(req.Attach.Type)
		attachment.SetFilename(req.Attach.Filename)
		attachment.SetDisposition(req.Attach.Disposition)
		attachment.SetContentID(req.Attach.ContentID)
		mailer = mail.NewSingleEmail(from, req.Subject, to, "", req.Body)
		mailer.AddAttachment(attachment)
	} else {
		mailer = mail.NewSingleEmail(from, req.Subject, to, "", req.Body)
	}

	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(mailer)
	if err != nil {
		slog.ErrorContext(ctx, "[SendgridRepoImpl.ClientSend] failed to send email", err)
		err = fmt.Errorf("failed to send email")
		return
	}

	if response.StatusCode != 202 {
		slog.ErrorContext(ctx, "[SendgridRepoImpl][response][code] failed to send email", fmt.Sprintf("%d", response.StatusCode), err)
		err = fmt.Errorf("failed to send email")
		return
	}

	return
}
