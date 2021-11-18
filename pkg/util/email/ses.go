package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	// DefaultCharSet The default character encoding for the email.
	DefaultCharSet = "UTF-8"
)

// Input simplifies aws-ses input
type Input struct {
	From     string
	To       []string
	Bcc      []string
	Subject  string
	HTMLBody string
	TextBody string
	CharSet  string
}

// Attachments supporting attachments input
type Attachments struct {
	Content     []byte
	ContentType string
	Name        string
}

// SendEmail sends email using SES service
func (s *Email) SendEmail(input Input) error {
	if input.CharSet == "" {
		input.CharSet = DefaultCharSet
	}
	if input.From == "" {
		input.From = s.cfg.Sender
	}
	toAddr := make([]*string, 0, len(input.To))
	for _, val := range input.To {
		toAddr = append(toAddr, &val)
	}

	// Assemble the email.
	awsInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddr,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(input.CharSet),
					Data:    aws.String(input.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(input.CharSet),
					Data:    aws.String(input.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(input.CharSet),
				Data:    aws.String(input.Subject),
			},
		},
		Source: aws.String(input.From),
	}

	// Attempt to send the email.
	_, err := s.ses.SendEmail(awsInput)
	if err != nil {
		return err
	}
	return nil
}

// SendRaw sends raw email using SES service
func (s *Email) SendRaw(input Input, attachments []*Attachments) error {
	if input.From == "" {
		input.From = s.cfg.Sender
	}

	// Parse input and attachments to MIME type
	data, err := s.writeDataToRaw(input, attachments)
	if err != nil {
		return err
	}

	// Assemble the email.
	awsInput := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: data,
		},
	}

	// Attempt to send the email.
	_, err = s.ses.SendRawEmail(awsInput)
	if err != nil {
		return err
	}
	return nil
}

func (s *Email) writeDataToRaw(input Input, attachments []*Attachments) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// email main header:
	h := make(textproto.MIMEHeader)
	h.Set("From", input.From)
	h.Set("To", strings.Join(input.To, ","))
	h.Set("Bcc", strings.Join(input.Bcc, ","))
	h.Set("Return-Path", input.From)
	h.Set("Subject", input.Subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	// body:
	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/html; charset=utf-8")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(input.HTMLBody))
	if err != nil {
		return nil, err
	}

	// file attachments:
	for _, a := range attachments {
		attachment := base64.StdEncoding.EncodeToString(a.Content)
		h = make(textproto.MIMEHeader)
		h.Add("Content-Transfer-Encoding", "base64")
		h.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s; size=%d", a.Name, len(attachment)))
		h.Add("Content-Type", fmt.Sprintf("%s; name=%s", a.ContentType, a.Name))
		part, err = writer.CreatePart(h)
		if err != nil {
			return nil, err
		}
		_, err = part.Write([]byte(attachment))
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Strip boundary line before header (doesn't work with it present)
	result := buf.String()
	if strings.Count(result, "\n") < 2 {
		panic("Invalid e-mail content")
	}
	result = strings.SplitN(result, "\n", 2)[1]

	return []byte(result), nil
}
