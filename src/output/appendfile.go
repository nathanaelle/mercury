package	output

import (
	"os"
)



type	AppendFile	struct {
	GenericOutput
	File	string		`json:"file"`
}

func (p *AppendFile)DriverName() string {
	return	"o_appendfile"
}




func (p *AppendFile)Run(errchan chan<- error) {
	p.end	= make(chan bool,1)
	p.source= make(chan string,100)

	f, err	:= os.OpenFile(p.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		errchan <- &OutputError { p.Driver, p.Id, "Open "+p.File, err }
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
					errchan <- &OutputError { p.Driver, p.Id, "Write "+p.File, err }
					return
				}
		}
	}
}
