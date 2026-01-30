package repository

import "errors"

var ErrCodeAlreadyExists = errors.New("code already exists")
var ErrNotFound = errors.New("url not found")
