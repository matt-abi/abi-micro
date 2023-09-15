package smtp

import (
	"log"
	"strings"
	"sync"

	"gopkg.in/gomail.v2"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-micro/micro"
)

type smtpService struct {
	config interface{}
	name   string

	From        string `json:"from"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Concurrency int    `json:"concurrency"`

	waitGroup sync.WaitGroup
	e         chan *gomail.Message
}

func newSmtpService(name string, config interface{}) SMTPService {
	return &smtpService{name: name, config: config}
}

/**
* 服务名称
**/
func (s *smtpService) Name() string {
	return s.name
}

/**
* 服务配置
**/
func (s *smtpService) Config() interface{} {
	return s.config
}

/**
* 初始化服务
**/
func (s *smtpService) OnInit(ctx micro.Context) error {

	dynamic.SetValue(s, s.config)

	if s.Concurrency <= 0 {
		s.Concurrency = 1
	}

	s.e = make(chan *gomail.Message, 204800)

	p := ctx.Payload()

	for i := 0; i < s.Concurrency; i++ {

		s.waitGroup.Add(1)

		go func() {

			defer s.waitGroup.Done()

			d := gomail.NewDialer(s.Host, s.Port, s.User, s.Password)

			for {

				m, ok := <-s.e

				if !ok {
					break
				}

				if m == nil {
					break
				}

				ctx, err := p.NewContext("__smtp__", micro.NewTrace())

				if err != nil {
					log.Println(err)
					continue
				}

				err = d.DialAndSend(m)

				if err != nil {
					ctx.Println("[err:1]", err)
				}

				ctx.Recycle()
			}
		}()

	}

	return nil
}

/**
* 校验服务是否可用
**/
func (s *smtpService) OnValid(ctx micro.Context) error {
	return nil
}

func (s *smtpService) Recycle() {

	if s.e != nil {
		close(s.e)
		s.e = nil
		s.waitGroup.Wait()
	}

}

func (s *smtpService) Send(to []string, subject string, body string, contentType string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", s.From)
	m.SetHeader("To", strings.Join(to, ";"))
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)

	d := gomail.NewDialer(s.Host, s.Port, s.User, s.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *smtpService) AsyncSend(to []string, subject string, body string, contentType string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", s.From)
	m.SetHeader("To", strings.Join(to, ";"))
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)

	s.e <- m

	return nil
}
