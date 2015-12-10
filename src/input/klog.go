// +build linux

package	input

import	(
)


type	KlogReader struct {
	GenericInput
}

func (klr *KlogReader) DriverName() string {
	return	"i_klog"
}

func (klr *KlogReader) Configure(errchan chan<- error) {
	klr.end		= make(chan bool,1)
	klr.errchan	= errchan
}


func (klr *KlogReader) Read(p []byte) (int,error) {
	n,err	:= klog_unread_size()
	if err!= nil {
		return 0,err
	}

	if n != 0 {
		if n > len(p) {
			n = len(p)
		}

		buf_n,err := klog_read( n )
		if err!= nil {
			return 0,err
		}
		copy( p, buf_n )
		return n,nil
	}

	b_1,err	:= klog_read( 3 )
	if err!= nil {
		return 0,err
	}

	n,err	= klog_unread_size()
	if err!= nil {
		return 0,err
	}

	n	+=len(b_1)

	if n > len(p) {
		n = len(p)
	}

	buffer,err	:= klog_read( n-len(b_1) )
	if err!= nil {
		return 0,err
	}
	copy( p, append( b_1, buffer... ) )

	return n,nil
}


func (klr *KlogReader)Run(dest chan<- Message) {
	boot_ts,err	:= boot_time()
	if err != nil {
		klr.errchan <- &InputError{ klr.Driver, klr.Id,"boot_time() ", err }
		return
	}

	klog_open()
	defer klog_close()

	raw_klog:= make(chan string,100)

	go reader_to_channel( &KlogReader {} , raw_klog )

	for {
		select {
			case line := <-raw_klog:
				l,err := ParseMessage_KLog( boot_ts, line )
				if err != nil {
					klr.errchan <- &InputError{ klr.Driver, klr.Id,"ParseMessage_KLog() ", err }
					continue
				}
				dest <- packmsg(klr.Id, l)

			case <-klr.end:
				return
		}
	}
}
