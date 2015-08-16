package message


import	(
	"errors"
	"fmt"
)



const	(
	Err_Invalid_Priority	int = iota
)


type MessageErr struct {
	EType	int

	Driver	string
	Key	string
	File	string
	ExpType	string

	Err	error
}


func error_invalid_prio(prio string) *MessageErr {
	return &MessageErr { Err_Invalid_Priority, "", "", prio, "", errors.New("") }
}


func (e MessageErr)Error()string  {
	err	:= ""
	switch e.EType {
		case	Err_Invalid_Priority:
			err	= fmt.Sprintf("invalid syslog priority. prio: %s", e.File )
	}
	return	err + e.Err.Error()
}
