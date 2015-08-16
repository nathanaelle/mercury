package	message

import	(
	"os"
	"strconv"
	"time"
)

type	Message		struct {
	header		string
	timestamp	time.Time
	hostname	string
	appname		string
	procid		string
	msgid		string
	message		string
}


const	RFC5424TimeStamp string = "2006-01-02T15:04:05.999999Z07:00"




var	hostname,_	= os.Hostname()



func EmptyMessage() *Message {
	return &Message { "", time.Unix(0,0), "-", "-", "-", "-", "-" }
}


func (msg *Message)Now() *Message {
	return &Message { msg.header, time.Now(), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}


func stamp_to_ts(stamp string) time.Time  {
	now	:= time.Now()
	ts,_	:= time.Parse(time.Stamp, stamp)
	year	:= now.Year()

	if (now.Month() ==1 && ts.Month()==12){
		year--
	}

	return time.Date( year, ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond(), ts.Location() )
}

func (msg *Message)Stamp(stamp	string) *Message  {
	return &Message { msg.header, stamp_to_ts(stamp), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}

func delta_boot_to_ts(boot_ts time.Time, s_sec string, s_nsec string) time.Time {
	sec,_	:= strconv.ParseInt(s_sec, 10, 64)
	nsec,_	:= strconv.ParseInt(s_nsec, 10, 64)

	return	boot_ts.Add( time.Duration(nsec)*time.Nanosecond + time.Duration(sec)*time.Second )
}


func (msg *Message)Delta(boot_ts time.Time, s_sec string, s_nsec string) *Message  {
	return &Message { msg.header, delta_boot_to_ts(boot_ts, s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}

func epoc_to_ts(s_sec string, s_nsec string) time.Time {
	sec,_	:= strconv.ParseInt(s_sec, 10, 64)
	nsec,_	:= strconv.ParseInt(s_nsec, 10, 64)

	return time.Unix(sec, nsec)
}

func (msg *Message)Epoch(s_sec string, s_nsec string) *Message {
	return &Message { msg.header, epoc_to_ts(s_sec, s_nsec), msg.hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}


func (msg *Message)App(appname string) *Message {
	return &Message { msg.header, msg.timestamp, msg.hostname, appname, msg.procid, msg.msgid, msg.message }
}

func (msg *Message)ProcID(procid string) *Message {
	return &Message { msg.header, msg.timestamp, msg.hostname, msg.appname, procid, msg.msgid, msg.message }
}

func (msg *Message)MsgID(msgid string) *Message {
	return &Message { msg.header, msg.timestamp, msg.hostname, msg.appname, msg.procid, msgid, msg.message }
}




func (msg *Message)Priority(prio string) *Message {
	return &Message { "<"+prio+">1", msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}


func (msg *Message)LocalHost() *Message {
	return &Message { msg.header, msg.timestamp, hostname, msg.appname, msg.procid, msg.msgid, msg.message }
}


func (msg *Message)Data(data string) *Message {
	switch data {
		case "":
			return &Message { msg.header, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, "-" }

		default:
			return &Message { msg.header, msg.timestamp, msg.hostname, msg.appname, msg.procid, msg.msgid, "- "+data }
	}
}


func (msg *Message)Stringify() string  {
	return msg.header + " " + msg.timestamp.Format(RFC5424TimeStamp) + " " + msg.hostname + " " + msg.appname + " " + msg.procid + " " + msg.msgid + " " + msg.message
}

func CreateMessage(data string, appname string, prio int) *Message  {
	return	EmptyMessage().App(appname).Priority(strconv.Itoa(prio)).LocalHost().Now().Data(data)
}
