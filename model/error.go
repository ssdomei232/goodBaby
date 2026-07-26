package model

// ValidationError 表示一个可以直接展示给用户的参数校验错误
type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string { return e.Msg }

// ErrValidation 构造一个参数校验错误
func ErrValidation(msg string) error { return ValidationError{Msg: msg} }
