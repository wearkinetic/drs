package drs

type Command struct {
	Key    string      `json:"key,omitempty"`
	Action string      `json:"action,omitempty"`
	Body   interface{} `json:"body,omitempty"`
}

func (this *Command) Map() map[string]interface{} {
	return this.Body.(map[string]interface{})
}

const (
	ERROR     = "drs.error"
	RESPONSE  = "drs.response"
	EXCEPTION = "drs.exception"
)
