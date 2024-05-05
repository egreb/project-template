package errors

type NotFound struct {
	Err error
}

func (e NotFound) Error() string {
	return e.Err.Error()
}

func (e NotFound) Unwrap() error {
	return e.Err
}

type BadRequest struct {
	Err error
}

func (e BadRequest) Error() string {
	return e.Err.Error()
}

func (e BadRequest) Unwrap() error {
	return e.Err
}

type InternalError struct {
	Err error
}

func (e InternalError) Error() string {
	return e.Err.Error()
}

func (e InternalError) Unwrap() error {
	return e.Err
}
