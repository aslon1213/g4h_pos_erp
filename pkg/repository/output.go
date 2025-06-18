package models

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Output struct {
	Data  interface{} `json:"data"`
	Error []Error     `json:"error"`
}

func NewOutput(data interface{}, errors ...Error) Output {
	return Output{
		Data:  data,
		Error: errors,
	}
}
