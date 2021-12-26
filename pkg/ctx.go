package pkg

import (
	"fmt"
)

const (
	contextKeyPrefix = "yandex-practicum-devops-"
)

type ContextKey string

func (c ContextKey) String() string {
	return fmt.Sprintf("%s%s", contextKeyPrefix, string(c))
}
