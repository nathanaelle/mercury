package	output

import	(
	"net"
	"time"
)


type	TCP5424	struct {
	GenericOutput
	Remote	string		`json:"remote"`
}


func (p *TCP5424) DriverName() string {
	return	"o_tcp5424"
}


func (p *TCP5424) Configure(errchan chan<- error) {
	p.end		= make(chan bool)
	p.source	= make(chan string, 100)
	p.errchan	= errchan

	if p.Remote == "" {
		panic("Remote mandatory")
	}
}


func (p *TCP5424) Run() {
	conn,err:= connect_remote( p.Remote )
	if err != nil {
		p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.Remote, err }
		return
	}

	for {
		select {
			case <- p.end:
				return

			case text := <- p.source:
				_, err := conn.Write( []byte(text+"\n") )

				for err != nil {
					nOErr := err.(*net.OpError)
					if ( nOErr.Op != "write" ||(
					nOErr.Err.Error() != "connection refused" &&
					nOErr.Err.Error() != "broken pipe" &&
					nOErr.Err.Error() != "connection timed out" )) {
						p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.Remote, err }
						return
					}
					time.Sleep(10 * time.Second)
					conn,err= connect_remote( p.Remote )
					if err != nil {
						p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.Remote, err }
						return
					}
					_, err = conn.Write( []byte(text+"\n") )
				}
		}
	}
}
