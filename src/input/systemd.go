package	input

import	(
	"os"
	"net"
	"bytes"
	"strings"
	"unicode"

	"../message"
)


type	JournalReader struct {
	GenericInput
	Journald	string		`json:"journald"`
}


func (jrnl *JournalReader)DriverName() string {
	return	"i_journald"
}


func (jrnl *JournalReader)Run(dest chan<- Message, errchan chan<- error) {
	jrnl.end	= make(chan bool,1)
	conn, err	:= net.ListenUnixgram("unixgram",  &net.UnixAddr { jrnl.Journald, "unixgram" } )

	for err != nil {
		switch err.(type) {
			case *net.OpError:
				if err.(*net.OpError).Err.Error() != "bind: address already in use" {
					errchan <- &InputError{ jrnl.Driver, jrnl.Id,"listen "+jrnl.Journald , err }
					return
				}

			default:
				errchan <- &InputError{ jrnl.Driver, jrnl.Id,"listen "+jrnl.Journald , err }
				return
		}

		if _, r_err := os.Stat(jrnl.Journald); r_err != nil {
			errchan <- &InputError{ jrnl.Driver, jrnl.Id,"lstat "+jrnl.Journald , err }
			return
		}
		os.Remove(jrnl.Journald)

		conn, err = net.ListenUnixgram("unixgram",  &net.UnixAddr { jrnl.Journald, "unixgram" } )
	}

	defer	conn.Close()

	for {
		select {
			case <-jrnl.end:
				return

			default:
				buffer	:= make([]byte, 65536)
				_, _, err := conn.ReadFrom(buffer)
				if err != nil {
					errchan <- &InputError{ jrnl.Driver, jrnl.Id,"ReadFrom "+jrnl.Journald, err }
					return
				}

				line	:= string(bytes.TrimRight(buffer,"\t \n\r\000"))
				if (line == "" ){
					continue
				}
				pos	:= strings.Index(line, ">")

				if  pos > 0 && unicode.IsDigit( rune(line[pos+1]) ) {
					dest <- packmsg(jrnl.Id, *message.ParseMessage_5424( line ))
				} else {
					dest <- packmsg(jrnl.Id, *message.ParseMessage_3164( line ))
				}
		}
	}
}
