package	main

import	(
	"flag"
	"runtime"

	"./input"
	"./output"

	types	"github.com/nathanaelle/useful.types"
)


const	MERCURYNAME	string		= "mercury"
const	DEFAULT_CONF	types.Path	= "/etc/mercury/"


func main()  {
	conf_path	:= new(types.Path)
	outputs		:= new(stringList)
	*conf_path	= DEFAULT_CONF
	*outputs	= stringList([]string{"o_stderr"})

	var	numcpu	= flag.Int("cpu", 1, "maximum number of logical CPU that can be executed simultaneously")
	flag.Var(conf_path, "conf", "path to the director" )
	flag.Var(outputs, "outputs", "coma ',' separated list of outputs for the internal logs" )

	flag.Parse()

	switch {
		case *numcpu >runtime.NumCPU():	runtime.GOMAXPROCS(runtime.NumCPU())
		case *numcpu <1:		runtime.GOMAXPROCS(1)
		default:			runtime.GOMAXPROCS(*numcpu)
	}

	mercury	:= NewMercury()

	mercury.Register(
		new(input.InternalReport),
		new(input.DevLogReader),
		new(input.KlogReader),
		new(input.FileReader),
		new(input.FIFOReader),
		new(input.JournalReader),
		new(output.AppendFile),
		new(output.StdErr),
		new(output.TCP5424),
		new(output.TLS5424)	)

	mercury.Load(string(*conf_path),[]string(*outputs))
	mercury.SetSignals()

	go mercury.MainLoop()

	mercury.End()
}
