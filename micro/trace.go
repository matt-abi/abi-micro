package micro

import (
	"strings"

	"github.com/google/uuid"
)

func NewTrace() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
