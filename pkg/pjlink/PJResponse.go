package pjlink

import (
	"errors"
	"strings"
)

type Response struct {
	Class    string   `json:"class"`
	Command  string   `json:"command"`
	Response []string `json:"response"`
}

func NewPJResponse() *Response {
	return &Response{}
}

func (res *Response) Parse(raw string) error {
	// If password is wrong, response will be 'PJLINK ERRA'
	if strings.Contains(raw, ERRA) {
		return errors.New("incorrect password")
	}
	if len(raw) == 0 {
		return errors.New("empty Response")
	}

	tokens := strings.Split(raw, " ")

	token0 := tokens[0]
	param1 := []string{token0[7:]}
	paramsN := tokens[1:]
	params := append(param1, paramsN...)

	res.Class = token0[1:2]
	res.Command = token0[2:6]
	res.Response = params

	return nil
}

// Checks if a Command was a success
func (res *Response) Success() bool {
	return res.Response[0] == OK
}
