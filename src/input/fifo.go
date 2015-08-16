package	input

import	(
	"io"
	"os"
	"syscall"
	"../message"
)


type	FIFOReader	struct {
	GenericInput
	Source		string		`json:"source"`
	AppName		string		`json:"appname"`
	Priority	string		`json:"priority"`

	prio		int
	fi_des		*os.File
}


func (fifo *FIFOReader)DriverName() string {
	return	"i_fifo"
}


func (fifo *FIFOReader) Read(p []byte) (n int, err error) {
	n,err	= fifo.fi_des.Read(p)
	if err == io.EOF {
		err=nil
	}

	return
}


func (fifo *FIFOReader)Run(dest chan<- Message, errchan chan<- error) {
	var err error
	fifo.end	= make(chan bool,1)
	fifo.prio,err	= message.PriorityDecode(fifo.Priority)
	if err != nil {
		errchan <- &InputError{ fifo.Driver, fifo.Id,"Priority "+fifo.Priority, err }
		return
	}

	syscall.Mkfifo(fifo.Source, 0644)
	fifo.fi_des,err	= os.OpenFile( fifo.Source, os.O_RDONLY, 0644 )
	if err != nil {
		errchan <- &InputError{ fifo.Driver, fifo.Id,"FIFO "+fifo.Source, err }
		return
	}

	data	:= make(chan string)
	defer	fifo.fi_des.Close()

	go reader_to_channel( fifo.fi_des , data )

	for {
		select{
			case line := <- data:
				dest <- packmsg(fifo.Id, *message.CreateMessage(line, fifo.AppName, fifo.prio))

			case <- fifo.end:
				return
		}
	}
}
