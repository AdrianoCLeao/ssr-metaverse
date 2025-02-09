package error

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func Error(code int, message string) APIError {
    return APIError{
        Code:    code,
        Message: message,
    }
}
