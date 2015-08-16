package	output


import	(
	"sync/atomic"
	"time"
	"os"
	"bufio"
	//"encoding/base64"
//	types	"github.com/nathanaelle/useful.types"
)



type	queue	struct {
	q_next	chan<- next
	q_peek	chan<- read
	q_push	chan<- string
	end	chan<- bool
	r_pos	int64
	w_pos	int64
}


type	next	struct {
	answer	chan<- string
}


type	read	struct {
	answer	chan<- string
}


const	(
	Q_SYNC	int32 = iota
	Q_DESYNC
)





func	CreateQueue(path string) queue {
	q_n	:= make(chan next)
	q_pe	:= make(chan read)
	q_pu	:= make(chan string)
	end	:= make(chan bool,1)

	q	:= queue { q_n, q_pe, q_pu, end, 0, 0 }

	go q.EvLoop( path, q_n, q_pe, q_pu, end )

	return	q
}


func	(q queue)Push(v string){
	q.q_push <- v
}


func	(q queue)End(){
	q.end <- true
}


func	(q queue)Peek() string {
	r_peek	:= make(chan string, 1)
	query	:= read { r_peek }
	q.q_peek<-query
	r	:=<-r_peek

	return r
}


func	(q queue)Next() string {
	r_next	:= make(chan string, 1)
	query	:= next { r_next }
	q.q_next<-query
	r	:=<-r_next

	return r
}


func	(q queue)EvLoop( q_path string, next <-chan next, peek <-chan read, push <-chan string, end <-chan bool )  {
	q_file,err	:= os.OpenFile(q_path, os.O_RDWR | os.O_CREATE, 0600)

	if err != nil {
		panic(err)
	}
	defer		q_file.Close()

	q_synced	:= Q_SYNC
	r_pos		:= int64(0)
	w_pos,_		:= q_file.Seek(0, os.SEEK_END)
	ticker		:= time.Tick(500 * time.Second)

	for {
		select {
			case	<-ticker:
				if atomic.CompareAndSwapInt32( &q_synced, Q_DESYNC, Q_SYNC ) {
					err := q_file.Sync()
					if err != nil {
						panic(err)
					}
				}

			case	query	:=<-next:
				p	:= atomic.LoadInt64( &r_pos )
				q_file.Seek(p, os.SEEK_SET)
				v,_	:=bufio.NewReader(q_file).ReadBytes('\n')
				atomic.AddInt64(&r_pos, int64(len(v)) )
				query.answer<-string(v)

			case	query	:=<-peek:
				p	:= atomic.LoadInt64( &r_pos )
				q_file.Seek(p, os.SEEK_SET)
				v,_	:=bufio.NewReader(q_file).ReadBytes('\n')
				query.answer<-string(v)

			case	v	:=<-push:
				raw_v	:= []byte(v)
				//s_raw_v	:= len(raw_v)
				p	:= atomic.LoadInt64( &w_pos )
				w_s,err	:= q_file.WriteAt(raw_v, p )
				if err != nil {
					panic(err)
				}
				atomic.AddInt64(&w_pos		, int64(w_s)	)
				atomic.StoreInt32(&q_synced	, Q_DESYNC	)

			case	<-end:
				break
		}
	}
	q_file.Sync()
}
