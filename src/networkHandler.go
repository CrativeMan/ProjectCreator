package main

import (
	"fmt"

	"github.com/pkg/sftp"
)

type SFTPError struct {
	Code    int
	Message string
}

func (e *SFTPError) Error() string {
	return fmt.Sprintf("SFTPError - Code: %d, Message: %s", e.Code, e.Message)
}

func NewSFTPError(code int, message string) error {
	return &SFTPError{
		Code:    code,
		Message: message,
	}
}

func initSFTPClient() *sftp.Client{
	
}