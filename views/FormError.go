package views

type FormError struct {
	Value string
	Error error
}

func NewFormError(v string, err error) FormError {
	return FormError{Value: v, Error: err}
}

func (e FormError) Invalid() string {
	if e.Error == nil {
		return "false"
	}
	return "true"
}

func (e FormError) Vlu() string {
	return e.Value
}

func (e FormError) Err() string {
	if e.Error == nil {
		return ""
	}
	return e.Error.Error()
}
