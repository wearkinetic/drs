package drs

type DRSError struct {
	Message string `json:"message"`
	Kind    string `json:"kind"`
}

func (this *DRSError) Error() string {
	return this.Message
}

func Error(message string) *DRSError {
	return &DRSError{
		Message: message,
		Kind:    "error",
	}
}

func Exception(message string) *DRSError {
	return &DRSError{
		Message: message,
		Kind:    "exception",
	}
}
