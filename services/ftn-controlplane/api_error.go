package controlplane

type APIError struct{Code string `json:"code"`;Message string `json:"message"`;Retryable bool `json:"retryable"`}
func NewAPIError(code,message string,retryable bool)APIError{return APIError{Code:code,Message:message,Retryable:retryable}}
