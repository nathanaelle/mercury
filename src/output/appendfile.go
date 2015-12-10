package	output

import (
	"os"
)



type	AppendFile	struct {
	GenericOutput
	File	string		`json:"file"`
}


func (p *AppendFile) DriverName() string {
	return	"o_appendfile"
}


func (p *AppendFile) Configure(errchan chan<- error) {
	p.end		= make(chan bool)
	p.source	= make(chan string, 100)
	p.errchan	= errchan

	if p.File == "" {
		panic("File mandatory")
	}
}


func (p *AppendFile) Run() {
	f, err	:= os.OpenFile(p.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.File, err }
		return
	}
	defer	f.Close()

	for {
		select {
			case <- p.end:
				return

			case text := <- p.source:
				_, err = f.WriteString(text+"\n");
				if err != nil {
					p.errchan <- &OutputError { p.Driver, p.Id, "Write "+p.File, err }
					return
				}
		}
	}
}
