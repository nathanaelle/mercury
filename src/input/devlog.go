package	input

import	(
	"os"
	"net"
)

type	DevLogReader struct {
	GenericInput
	Devlog	string		`json:"devlog"`
}


func (devlog *DevLogReader)DriverName() string {
	return	"i_devlog"
}



func (devlog *DevLogReader)cope_with(conn *net.UnixConn, buffer []byte, dest chan<- Message) {
	_,_,err := conn.ReadFrom(buffer)
	if err != nil {
		devlog.errchan <- &InputError{ devlog.Driver, devlog.Id,"ReadFrom "+devlog.Devlog , err }
		return
	}

	line	:= rtrim_blank(buffer)
	if (len(line) == 0 ){
		return
	}

	l, err := parse_3164_or_5424( devlog.Id, line )
	if err != nil {
		devlog.errchan <-  &InputError{ devlog.Driver, devlog.Id,"parse_3164_or_5424 "+devlog.Devlog , err }
		return
	}

	dest <- l
}



func (devlog *DevLogReader)Run(dest chan<-Message, errchan chan<- error) {
	devlog.end	= make(chan bool,1)
	devlog.errchan	= errchan
	conn, err := net.ListenUnixgram("unixgram",  &net.UnixAddr { devlog.Devlog, "unixgram" } )
	for err != nil {
		switch err.(type) {
			case *net.OpError:
				if err.(*net.OpError).Err.Error() != "bind: address already in use" {
					devlog.errchan <- &InputError{ devlog.Driver, devlog.Id,"Listen "+devlog.Devlog , err }
					return
				}

			default:
				devlog.errchan <- &InputError{ devlog.Driver, devlog.Id,"Listen "+devlog.Devlog , err }
				return
		}

		if _, r_err := os.Stat(devlog.Devlog); r_err != nil {
			devlog.errchan <- &InputError{ devlog.Driver, devlog.Id,"lstat "+devlog.Devlog , err }
			return
		}
		os.Remove(devlog.Devlog)

		conn, err = net.ListenUnixgram("unixgram",  &net.UnixAddr { devlog.Devlog, "unixgram" } )
	}
	defer	conn.Close()


	for {
		select {
		case <-devlog.end:
			return

		default:
			devlog.cope_with(conn, make([]byte, 65536), dest )
		}
	}

}
