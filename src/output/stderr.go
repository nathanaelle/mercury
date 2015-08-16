package	output

import (
	"os"
)


type	StdErr	struct {
	GenericOutput
}


func (p *StdErr)DriverName() string {
	return	"o_stderr"
}


func (p *StdErr)Run(errchan chan<- error) {
	p.end	= make(chan bool,1)
	p.source= make(chan string,100)

	for {
		select {
			case <- p.end:
				return

			case text := <- p.source:
				_, err := os.Stderr.WriteString(text+"\n");
				if err != nil {
					errchan <- &OutputError { p.Driver, p.Id, "Write STDErr", err }
					return
				}
		}
	}
}
