package am

import (
	"github.com/google/uuid"

	"github.com/adriancrafter/todoapp/internal/am/errors"
)

var (
	NoSlugPrefixErr = errors.NewError("no slug prefix provided")
)

var (
	ZeroUUID = uuid.UUID{}
)
