package	input

import	(
	"os"
	"net"
)


type	JournalReader struct {
	GenericInput
	Journald	string		`json:"journald"`
}


func (jrnl *JournalReader) DriverName() string {
	return	"i_journald"
}

func (jrnl *JournalReader) Configure(errchan chan<- error) {
	jrnl.end	= make(chan bool,1)
	jrnl.errchan	= errchan

	if jrnl.Journald == "" {
		jrnl.Journald = "/run/systemd/journal/syslog"
	}
}



func (jrnl *JournalReader) Run(dest chan<- Message) {

	conn, err	:= net.ListenUnixgram("unixgram",  &net.UnixAddr { jrnl.Journald, "unixgram" } )

	for err != nil {
		switch err.(type) {
			case *net.OpError:
				if err.(*net.OpError).Err.Error() != "bind: address already in use" {
					jrnl.errchan <- &InputError{ jrnl.Driver, jrnl.Id,"listen "+jrnl.Journald , err }
					return
				}

			default:
				jrnl.errchan <- &InputError{ jrnl.Driver, jrnl.Id,"listen "+jrnl.Journald , err }
				return
		}

		if _, r_err := os.Stat(jrnl.Journald); r_err != nil {
			jrnl.errchan <- &InputError{ jrnl.Driver, jrnl.Id,"lstat "+jrnl.Journald , err }
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
			jrnl.cope_with(conn, make([]byte, 65536), dest )
		}
	}
}


func (jrnl *JournalReader) cope_with(conn *net.UnixConn, buffer []byte, dest chan<- Message) {
	_,_,err := conn.ReadFrom(buffer)
	if err != nil {
		jrnl.errchan <- &InputError{ jrnl.Driver, jrnl.Id,"ReadFrom "+jrnl.Journald, err }
		return
	}

	line	:= rtrim_blank(buffer)
	if (len(line) == 0 ){
		return
	}

	l, err := parse_3164_or_5424( jrnl.Id, line )
	if err != nil {
		jrnl.errchan <- &InputError{ jrnl.Driver, jrnl.Id,"parse_3164_or_5424 "+jrnl.Journald, err }
		return
	}

	dest <- l
}
