package drs

type Command struct {
	Key    string      `json:"key"`
	Action string      `json:"action"`
	Body   interface{} `json:"body"`
}

func (this *Command) Map() map[string]interface{} {
	return this.Body.(map[string]interface{})
}
