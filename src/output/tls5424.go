package	output


import	(
	"net"
	"time"
	"crypto/tls"
)



type	TLS5424	struct {
	GenericOutput
	Remote	string		`json:"remote"`
}

func (p *TLS5424) DriverName() string {
	return	"o_tls5424"
}


func (p *TLS5424) Configure(errchan chan<- error) {
	p.end		= make(chan bool)
	p.source	= make(chan string, 100)
	p.errchan	= errchan

	if p.Remote == "" {
		panic("Remote mandatory")
	}
}


func connect_tls_remote(remote_host string, tlsConfig *tls.Config)	(*tls.Conn,error)  {
	conn, err := connect_remote(remote_host)
	if err != nil {
		return nil, err
	}

	return tls.Client( conn, tlsConfig ),nil
}


func (p *TLS5424) Run() {
	tls_config := &tls.Config{
		InsecureSkipVerify:	false,
		MinVersion:		tls.VersionTLS11,
		MaxVersion:		tls.VersionTLS12,
	}

	conn,err:= connect_tls_remote( p.Remote, tls_config )
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
					if (nOErr.Op != "write" ||(
					nOErr.Err.Error() != "connection refused" &&
					nOErr.Err.Error() != "broken pipe" )) {
						p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.Remote, err }
						return
					}
					time.Sleep(10 * time.Second)
					conn,err= connect_tls_remote( p.Remote, tls_config )
					if err != nil {
						p.errchan <- &OutputError { p.Driver, p.Id, "Open "+p.Remote, err }
						return
					}
					_, err	= conn.Write( []byte(text+"\n") )
				}
		}
	}
}
