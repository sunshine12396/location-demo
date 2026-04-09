package guid

import (
	"github.com/google/uuid"
)

var (
	uuidFunc = uuid.NewString
)

type impl struct {
	uuidFn func() string
}

func New() Guid {
	return &impl{
		uuidFn: uuidFunc,
	}
}
