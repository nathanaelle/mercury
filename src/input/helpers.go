package	input

import	(
	"os"
	"io"
	"math"
	"bufio"
	"strconv"
	"strings"
	"io/ioutil"

	"../message"
)


type	Message	struct	{
	Source	string
	Message	message.Message
}


func packmsg(id string, msg message.Message) Message {
	return Message{ id, msg }
}


type	GenericInput struct {
	Id	string		`json:"id"`
	Driver	string		`json:"driver"`
	Output	[]string	`json:"output"`

	end	chan bool
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
		return strings.TrimRight(string(ret),"\n\r \t")
	}

	return ""
}


func file_write(filename string, data string) {
	ioutil.WriteFile(filename, []byte(data+"\n"), 0644 )
}
