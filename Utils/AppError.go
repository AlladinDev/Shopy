package utils

type AppError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"message"`
	Reason     string `json:"reason"`
}

func (err *AppError) Error() string {
	return err.Reason
}

func ReturnAppError(err error, msg string, statusCode int) error {

	return &AppError{
		StatusCode: statusCode,
		Msg:        msg,
		Reason:     err.Error(),
	}
}
