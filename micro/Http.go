package micro

import "fmt"

type HttpContent struct {
	Code    int               `json:"code"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

func (s *HttpContent) Error() string {
	return fmt.Sprintf("HttpContent: %d", s.Code)
}
