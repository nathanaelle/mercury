package	output

import	(
	"fmt"
)


type	OutputError struct {
	Driver	string
	Id	string
	Action	string
	Err	error
}




func (e OutputError)Error() string {
	return fmt.Sprintf("%s [%s] %s got %v", e.Id, e.Driver, e.Action, e.Err )
}
