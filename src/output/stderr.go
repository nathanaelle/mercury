package	output

import (
	"os"
)


type	StdErr	struct {
	GenericOutput
}


func (p *StdErr) DriverName() string {
	return	"o_stderr"
}


func (p *StdErr) Configure(errchan chan<- error) {
	p.end		= make(chan bool)
	p.source	= make(chan string, 100)
	p.errchan	= errchan
}


func (p *StdErr) Run() {

	for {
		select {
			case <- p.end:
				return

			case text := <- p.source:
				_, err := os.Stderr.WriteString(text+"\n");
				if err != nil {
					p.errchan <- &OutputError { p.Driver, p.Id, "Write STDErr", err }
					return
				}
		}
	}
}
