package coleteonline

import "fmt"

type Error struct {
	Parameter string `json:"parameter"`
	Message   string `json:"message"`
}

type ResponseError struct {
	Message string  `json:"message"`
	Code    int     `json:"code"`
	Errors  []Error `json:"errors"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf(`%d: "%s"`, e.Code, e.Message)
}

type AuthResponseError struct {
	Name        string `json:"error"`
	Description string `json:"error_description"`
}

func (e *AuthResponseError) Error() string {
	return fmt.Sprintf(`%s: "%s"`, e.Name, e.Description)
}
