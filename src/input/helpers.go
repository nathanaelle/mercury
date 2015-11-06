package	input

import	(
	"os"
	"io"
	"math"
	"bytes"
	"bufio"
	"errors"
	"strconv"
	"io/ioutil"

	"github.com/nathanaelle/syslog5424"
)


type	Message	struct	{
	Source	string
	Message	syslog5424.Message
}


func packmsg(id string, msg syslog5424.Message) Message {
	return Message{ id, msg }
}


type	GenericInput struct {
	Id	string		`json:"id"`
	Driver	string		`json:"driver"`
	Output	[]string	`json:"output"`

	end	chan bool
	errchan	chan<-error
}


func (genin *GenericInput)DriverType() string {
	return	"INPUT"
}


func (genin *GenericInput)End() {
	genin.end <- true
	close(genin.end)
}


func (genin *GenericInput)SendTo() []string {
	return genin.Output
}


func human_scale(value float64, base float64, unit string) string {
	exp	:= []string { "y","z","a","f","p","n","µ","m","","k","M","G","T","P","E","Z","Y" }
	s	:= math.Floor(math.Log2(value)/math.Log2(base))
	h_v	:= value / math.Pow(base, s)

	if s > -9 && s < 9 {
		return	strconv.FormatFloat(h_v,'f',2,64)+" "+exp[int(s)+8]+unit
	}

	return strconv.FormatFloat(value,'E',6,64)+" "+unit
}


func reader_to_channel(r io.Reader, dest chan<- string)  {
	scanner := bufio.NewScanner( r )

	for scanner.Scan() {
		dest <- scanner.Text()
	}
}


func file_size(filename string) (n int64) {
	fst,err	:= os.Stat(filename)
	if err != nil {
		return 0
	}

	return fst.Size()
}


func file_exists(filename string) (bool) {
	_,err	:= os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}


func file_read(filename string) (string) {
	ret,err	:= ioutil.ReadFile(filename)
	if err == nil {
		return string(rtrim_blank(ret))
	}

	return ""
}


func file_write(filename string, data string) {
	ioutil.WriteFile(filename, []byte(data+"\n"), 0644 )
}


func rtrim_blank(d []byte) []byte {
	return bytes.TrimRight( d, "\n\r \t\000" )
}



func FieldsFuncN(s string, hope int, f func(rune) bool) []string {
	p_is_sep:= true
	is_sep	:= true
	begin	:= -1
	end	:= -1
	res	:= make( []string, 0, hope )

	for i,rune := range s {
		p_is_sep = is_sep
		is_sep = f(rune)
		switch {
			case is_sep && !p_is_sep:
				end = i
				res = append( res, s[begin:end] )

			case !is_sep && p_is_sep:
				begin	= i
		}
	}

	if(begin>-1 && begin>end ) {
		res = append( res, s[begin:len(s)] )
	}

	return res
}



func get_3164tokenizer() func(rune)bool {
	colon	:= 2
	state	:= "priority"

	return func (c rune) bool {
		switch	state {
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
			case	colon >  0 && c != ':':
				return	false

			case	colon == 0 && c != ' ':
				return	false

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
func parse_3164(data []byte) (syslog5424.Message,error)  {
	part	:= FieldsFuncN(string(data), 10, get_3164tokenizer() )

	switch len(part) {
	case 0, 1, 2, 3:
		return syslog5424.EmptyMessage(), errors.New("wrong format 3164 : "+ string(data))

	case 4:
		prio, err := strconv.Atoi(part[0])
		if err != nil {
			return syslog5424.EmptyMessage(),errors.New("Wrong Priority :"+string(part[0]))
		}
		return syslog5424.CreateMessage(part[2], syslog5424.Priority(prio), part[3]).Stamp(part[1]),nil

	default:
		prio, err := strconv.Atoi(part[0])
		if err != nil {
			return syslog5424.EmptyMessage(),errors.New("Wrong Priority :"+string(part[0]))
		}
		return syslog5424.CreateMessage(part[2], syslog5424.Priority(prio), part[4]).Stamp(part[1]).ProcID(part[3]),nil
	}
}


func parse_3164_or_5424(source string, data []byte) (Message,error){
	pos	:= bytes.Index(data, []byte{'>'})

	if  pos > 0 && data[pos+1] >= '0' && data[pos+1] <= '9' {
		l, err := syslog5424.Parse( data )
		if err != nil {
			return Message{}, err
		}

		return packmsg(source, l ),nil
	}

	l, err := parse_3164( data )
	if err != nil {
		return Message{}, err
	}

	return packmsg(source, l ),nil
}
