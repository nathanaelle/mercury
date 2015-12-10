package	output


import	(
	"net"
	"time"
)


type	GenericOutput struct {
	Id	string		`json:"id"`
	Driver	string		`json:"driver"`

	end	chan bool
	source	chan string
	errchan chan<- error
}


func (p *GenericOutput)DriverType() string {
	return	"OUTPUT"
}


func (p *GenericOutput)End() {
	close(p.end)
}


func (p *GenericOutput)Send(msg string) {
	p.source <- msg
}



func connect_remote(remote_host string)	(*net.TCPConn,error)  {
/*
net/OpError 0: Op string = dial
net/OpError 1: Net string = tcp
net/OpError 2: Addr net.Addr = <nil>
net/OpError 3: Err error = lookup aglog.local: no such host

if ( nOErr.Op != "write" ||(
nOErr.Err.Error() != "connection refused" &&
nOErr.Err.Error() != "broken pipe" )) {
exterminate(err)
}


net/OpError 0: Op string = dial
net/OpError 1: Net string = tcp
net/OpError 2: Addr net.Addr = 172.16.3.91:5140
net/OpError 3: Err error = no route to host

*/
	conn, err := net.Dial("tcp", remote_host)
	ErrCnt	:= 0
	for err != nil {
		ErrCnt++
		if ErrCnt == 50 {
			return nil,err
		}

		nOErr := err.(*net.OpError)
		errTxt	:= nOErr.Err.Error()

		switch {
			case nOErr.Op == "dial" && nOErr.Net == "tcp" && errTxt[0:6] == "lookup":
				time.Sleep(60 * time.Second)
				conn, err = net.Dial("tcp", remote_host)

			case nOErr.Op == "dial" && nOErr.Net == "tcp" && errTxt == "no route to host":
				time.Sleep(60 * time.Second)
				conn, err = net.Dial("tcp", remote_host)

			default:
				return nil,err
		}
	}

	tcpconn	:= conn.(*net.TCPConn)

	/*
	 *	Linger Close	= 5s
	 *	KeepAlive	= 2s
	 */

	tcpconn.SetWriteBuffer(1<<20)

	if err	= tcpconn.SetLinger( 5 ); err != nil {
		return nil,err
	}

	if err	= tcpconn.SetKeepAlivePeriod( 10 * time.Second ); err != nil {
		return nil,err
	}

	if err	= tcpconn.SetKeepAlive( true ); err != nil {
		return nil,err
	}

	return tcpconn,nil
}
