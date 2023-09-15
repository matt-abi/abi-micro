package smtp

import (
	"fmt"

	"github.com/ability-sh/abi-micro/micro"
)

type SMTPService interface {
	micro.Service
	Send(to []string, subject string, body string, contentType string) error
	AsyncSend(to []string, subject string, body string, contentType string) error
}

func GetSMTPService(ctx micro.Context, name string) (SMTPService, error) {

	s, err := ctx.GetService(name)

	if err != nil {
		return nil, err
	}

	ss, ok := s.(SMTPService)

	if !ok {
		return nil, fmt.Errorf("service %s not instanceof SMTPService", name)
	}

	return ss, nil
}
