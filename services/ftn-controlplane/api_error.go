package controlplane
import "strings"
type APIError struct{Code string `json:"code"`;Message string `json:"message"`;Retryable bool `json:"retryable"`}
func NewAPIError(code,message string,retryable bool)APIError{return APIError{Code:strings.TrimSpace(code),Message:strings.TrimSpace(message),Retryable:retryable}}
