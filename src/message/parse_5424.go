package	message

import	(
	"strings"
	"time"
)




func ParseMessage_5424(data string) *Message  {
	parts := strings.SplitN(data, " ", 7)

	msg := EmptyMessage()

	msg.header	= parts[0]
	msg.timestamp,_	= time.Parse(RFC5424TimeStamp,parts[1])
	msg.hostname	= parts[2]
	msg.appname	= parts[3]
	msg.procid	= parts[4]
	msg.msgid	= parts[5]
	msg.message	= parts[6]

	return msg
}
