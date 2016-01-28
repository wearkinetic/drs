package drs

type DRSError struct {
	Message string `json:"message"`
}

func (this *DRSError) Error() string {
	return this.Message
}

func Error(message string) *DRSError {
	return &DRSError{
		Message: message,
	}
}
