package	input

import	(
	"io"
	"os"
	"syscall"

	"github.com/nathanaelle/syslog5424"
)


type	FIFOReader	struct {
	GenericInput
	Source		string		`json:"source"`
	AppName		string		`json:"appname"`
	Priority	string		`json:"priority"`

	prio		*syslog5424.Priority
	fi_des		*os.File
}


func (fifo *FIFOReader) DriverName() string {
	return	"i_fifo"
}


func (fifo *FIFOReader) Configure(errchan chan<- error) {
	fifo.end	= make(chan bool,1)
	fifo.errchan	= errchan

	if fifo.Source == "" {
		panic("Source mandatory")
	}

	if fifo.AppName == "" {
		panic("AppName mandatory")
	}

	if fifo.Priority == "" {
		panic("Priority mandatory")
	}

	fifo.prio	= new(syslog5424.Priority)
	err		:= fifo.prio.Set(fifo.Priority)
	if err != nil {
		fifo.errchan <- &InputError{ fifo.Driver, fifo.Id,"Priority "+fifo.Priority, err }
		return
	}

}



func (fifo *FIFOReader) Read(p []byte) (n int, err error) {
	n,err	= fifo.fi_des.Read(p)
	if err == io.EOF {
		err=nil
	}

	return
}


func (fifo *FIFOReader) Run(dest chan<- Message) {
	var err error

	syscall.Mkfifo(fifo.Source, 0644)
	fifo.fi_des,err	= os.OpenFile( fifo.Source, os.O_RDONLY, 0644 )
	if err != nil {
		fifo.errchan <- &InputError{ fifo.Driver, fifo.Id,"FIFO "+fifo.Source, err }
		return
	}

	data	:= make(chan string,100)
	defer	fifo.fi_des.Close()

	go reader_to_channel( fifo.fi_des , data )

	for {
		select{
			case line := <- data:
				dest <- packmsg(fifo.Id, syslog5424.CreateMessage(fifo.AppName, *fifo.prio, line))

			case <- fifo.end:
				return
		}
	}
}
