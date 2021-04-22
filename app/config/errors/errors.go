package errors

import "errors"

var ErrBlankParam = errors.New("cannot be blank")
var ErrConflictingBlocks = errors.New("one of the blocks must be defined")
