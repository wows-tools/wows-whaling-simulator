package common

import (
	"errors"
)

var (
	ErrDraw       = errors.New("did not found a proper interval when drawing")
	ErrCastCopyLB = errors.New("Error casting when copying the LB data")
	ErrNoMoreItem = errors.New("No more collectable items, should have been substituted")
)
