package	main

import	(
	"strings"
)


type	stringList	[]string

func (sl *stringList)Set(data string) error {
	for _, dt := range strings.Split(data, ",") {
		*sl = append(*sl, dt)
	}

	return nil
}


func (sl stringList)String() string {
	return strings.Join(sl,",")
}
