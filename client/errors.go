package client

import "fmt"

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "not found"
}

type BadRequestError struct {
	Problem string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.Problem)
}

type ConflictError struct {
	Problem string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Problem)
}

type UnknownError struct {
	Problem string
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("unknown error: %s", e.Problem)
}
