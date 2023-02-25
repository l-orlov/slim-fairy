package service

import "github.com/pkg/errors"

var (
	ErrWrongPassword = errors.Errorf("wrong password")
)
