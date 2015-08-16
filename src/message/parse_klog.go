package	message

import	(
	"unicode"
	"log"
	"time"
)



//	statefull tokenizer for linux printk() message buffer
//
//	BUG(nath): may need some generic API
func get_klog_tokenizer() func(rune)bool {
	started	:= false
	state	:= "priority"

	return func (c rune) bool {
		switch state {

			case	"dispatch":
			switch {
				case	c == '<':
				state	= "priority"
				started	= true
				return	true

				case	c == '[':
				state	= "date"
				started	= true
				return	true

				case	started:
				state	= "message"
				return	unicode.IsSpace(c)

				default:
				started	= true
				return true
			}

			case	"priority":
			switch {
				case	c == '<':
				return	true

				case	c == '>':
				state	= "dispatch"
				return	true

				default:
				return	!unicode.IsDigit(c)
			}

			case	"date":
			switch {
				case	c == '[' || c == '.':
					return	true

				case	c == ']':
				state	= "dispatch"
				return	true

				default:
				return	!unicode.IsDigit(c)
			}

			default:
			return	false
		}
	}
}


//	statefull parser for linux printk() message buffer
//
//	BUG(nath): may need some generic API
func ParseMessage_KLog(boot_ts time.Time, data string) *Message  {
	log.SetFlags(log.Ltime | log.Lshortfile)

	part	:= FieldsFuncN(data, 4, get_klog_tokenizer())

	if (len(part) < 4){
		log.Println(data)
		for pi := range part {
			log.Println(part[pi])
		}
	}

	switch len(part) {

		case 2:
			return	EmptyMessage().LocalHost().Priority(part[0]).Now().Data(part[1])

		case 3:
			// (kern) 0 * 8 + 6 (info)
			return	EmptyMessage().LocalHost().Priority("6").Delta(boot_ts, part[0], part[1]).Data(part[2])

		case 4:
			return	EmptyMessage().LocalHost().Priority(part[0]).Delta(boot_ts, part[1], part[2]).Data(part[3])

		default:
			// (kern) 0 * 8 + 6 (info)
			return	EmptyMessage().LocalHost().Now().Priority("6").Data(data)
	}
}
