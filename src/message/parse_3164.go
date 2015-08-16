package	message

import	(
	"log"
)



func get_3164tokenizer() func(rune)bool {
	colon	:= 2
	state	:= "priority"

	return func (c rune) bool {
		switch state {
			case	"priority":
				switch {
					case	c == '<':
						return	true

					case	c == '>':
						state	= "date"
						return	true

					default:
						return	false
				}

			case	"date":
				switch {
					case	colon >  0 && c != ':':	return	false
					case	colon == 0 && c != ' ':	return	false
					case	colon >  0 && c == ':':
						colon--
						return	false

					case	colon == 0 && c == ' ':
						state	= "appname"
						return	true
				}
				panic("state "+state+" for ["+string(c)+"] WTF")

			case	"appname":
				switch {
					case	c == ':' || c == ' ':
						state	= "premessage"
						return	true

					case	c == '[':
						state	= "pid"
						return	true

					default:
						return false
				}

			case	"pid":
				switch {
					case	c == ':' || c == ' ' || c == ']':
						state	= "premessage"
						return	true

					default:
						return false
				}

			case	"premessage":
				switch c {
					case	':', ' ', ']':
						return	true

					default:
						state	= "message"
						return	false
				}

			case	"message":
				return	false

			default:
				panic("state "+state+" for ["+string(c)+"] WTF")
		}
	}
}


// <29>Dec  4 11:02:35 pdns[2030]:
func ParseMessage_3164(data string) *Message  {
	part	:= FieldsFuncN(data, 10, get_3164tokenizer() )

	switch len(part) {
		case 0, 1, 2, 3:
			log.Println("P_3164\t|",data,"|")
			panic(data)

		case 4:
			return EmptyMessage().LocalHost().Priority(part[0]).Stamp(part[1]).App(part[2]).Data(part[3])

		default:
			return EmptyMessage().LocalHost().Priority(part[0]).Stamp(part[1]).App(part[2]).ProcID(part[3]).Data(part[4])
	}
}
