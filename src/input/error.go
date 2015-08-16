package	input

import	(
	"fmt"
)


type	InputError struct {
	Driver	string
	Id	string
	Action	string
	Err	error
}




func (e InputError)Error() string {
	return fmt.Sprintf("%s [%s] %s got %v", e.Id, e.Driver, e.Action, e.Err )
}
