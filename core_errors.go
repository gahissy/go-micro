package micro

type BindingError struct {
	Err error
}

func (e *BindingError) Unwrap() error { return e.Err }
