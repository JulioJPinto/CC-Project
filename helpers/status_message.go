package helpers

import (
	"fmt"
	"strings"
)


type StatusMessage struct {
	Messages []string
	Errors   []error
}

func NewStatusMessage()StatusMessage {
	msg:=StatusMessage{}
	msg.Messages = []string{}
	msg.Errors = []error{}
	return msg
}

func (m *StatusMessage) ShowMessages() string {
	var ret string
	ret += "Messages:\n"
	for _, msg := range m.Messages {
		ret += "\t" + msg + "\n"
	}
	return ret
}
func (m *StatusMessage) ShowErrors() string {
	var ret string
	ret += "Errors:\n"
	for _, err := range m.Errors {
		ret += "\t" + err.Error() + "\n"
	}
	return ret
}

func (m *StatusMessage) Show() string {
	var ret string
	ret += "Messages:\n"
	for _, msg := range m.Messages {
		ret += "\t" + msg + "\n"
	}
	ret += "Errors:\n"
	for _, err := range m.Errors {
		ret += "\t" + err.Error() + "\n"
	}
	return ret
}

func (m *StatusMessage) Error() error {
	if len(m.Errors) == 0 {return nil}
	err_strings := make([]string, len(m.Errors))
	for i, e := range m.Errors {
		err_strings[i] = e.Error()
	}
	
	return fmt.Errorf(strings.Join(err_strings, ";"))
}
func (m *StatusMessage) AddError(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *StatusMessage) AddMessage(err error, success_message string) {
	if err == nil {
		m.Messages = append(m.Messages, success_message)
	} else {
		m.AddError(err)
	}
}
