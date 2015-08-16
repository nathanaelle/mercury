package	input

import	(
	"net"
	"os"
	"strings"
	"unicode"

	"../message"
)
//	/run/systemd/journal/syslog
//	/run/systemd/journal/syslog


type	DevLogReader struct {
	GenericInput
	Devlog	string		`json:"devlog"`
}


func (devlog *DevLogReader)DriverName() string {
	return	"i_devlog"
}


func (devlog *DevLogReader)Run(dest chan<- Message, errchan chan<- error) {
	devlog.end	= make(chan bool,1)
	conn, err := net.ListenUnixgram("unixgram",  &net.UnixAddr { devlog.Devlog, "unixgram" } )
	for err != nil {
		switch err.(type) {
			case *net.OpError:
				if err.(*net.OpError).Err.Error() != "bind: address already in use" {
					errchan <- &InputError{ devlog.Driver, devlog.Id,"Listen "+devlog.Devlog , err }
					return
				}

			default:
				errchan <- &InputError{ devlog.Driver, devlog.Id,"Listen "+devlog.Devlog , err }
				return
		}

		if _, r_err := os.Stat(devlog.Devlog); r_err != nil {
			errchan <- &InputError{ devlog.Driver, devlog.Id,"lstat "+devlog.Devlog , err }
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
				buffer	:= make([]byte, 65536)
				_, _, err := conn.ReadFrom(buffer)
				if err != nil {
					errchan <- &InputError{ devlog.Driver, devlog.Id,"ReadFrom "+devlog.Devlog , err }
					return
				}

				line	:= strings.TrimRight(string(buffer),"\t \n\r\000")
				if (line == "" ){
					continue
				}
				pos	:= strings.Index(line, ">")

				if  pos > 0 && unicode.IsDigit( rune(line[pos+1]) ) {
					dest <- packmsg(devlog.Id, *message.ParseMessage_5424( line ) )
				} else {
					dest <- packmsg(devlog.Id, *message.ParseMessage_3164( line ) )
				}
		}
	}

}
