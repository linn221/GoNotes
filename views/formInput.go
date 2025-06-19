package views

import "strconv"

type FormInput struct {
	Value string
	Error error
}

func NewFormInput(input string, err error) FormInput {
	return FormInput{Value: input, Error: err}
}

func (e FormInput) Invalid() string {
	if e.Error == nil {
		return "false"
	}
	return "true"
}

func (e FormInput) Vlu() string {
	return e.Value
}

func (e FormInput) VluInt() int {
	i, _ := strconv.Atoi(e.Value)
	return i
}

func (e FormInput) Err() string {
	if e.Error == nil {
		return ""
	}
	return e.Error.Error()
}
