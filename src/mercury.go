package	main

import	(
	"os"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"reflect"
	"strconv"
	"syscall"
	"io/ioutil"
	"os/signal"
	"encoding/json"
	"path/filepath"

	"./input"
	"./output"
	"./message"

	"gopkg.in/fsnotify.v1"
)


type	(
	Input	interface {
		Driver
		SendTo()	[]string
		Run(chan<- input.Message, chan<- error)
	}

	Output	interface {
		Driver
		Send(string)
		Run(chan<-error)
	}

	Driver	interface {
		DriverName()	string
		DriverType()	string
		End()
	}

	mercury_conf	struct {
		drivers		map[string]Driver
		instances	map[string]Driver
		file2id		map[string]string

		errchan		chan error
		inputchan	chan input.Message
		internalchan	chan internal_msg
		end		chan bool
		watchend	chan bool
		mainend		chan bool
		outputs		[]string

		glock		*sync.Mutex
	}

	internal_msg		struct {
		Severity	int
		Category	string
		Message		string
	}
)


func NewMercury() *mercury_conf {
	mc		:=new(mercury_conf)
	mc.glock	= new(sync.Mutex)

	mc.drivers	= make( map[string]Driver		)
	mc.instances	= make( map[string]Driver		)
	mc.file2id	= make( map[string]string		)
	mc.inputchan	= make( chan input.Message	, 1000	)
	mc.internalchan	= make( chan internal_msg	, 20	)
	mc.errchan	= make( chan error		, 20	)
	mc.watchend	= make( chan bool		, 1	)
	mc.end		= make( chan bool		, 1	)
	mc.mainend	= make( chan bool		, 1	)

	return mc
}


func (mc *mercury_conf)Register(drivers ...Driver) {
	for _, in := range drivers {
		mc.drivers[in.DriverName()]	= in
	}
}


func sd_notify(state,message string) bool {
	if os.Getenv("NOTIFY_SOCKET") == "" {
		return false
	}

	conn, err := net.Dial("unixgram", os.Getenv("NOTIFY_SOCKET"))
	if err != nil {
		return false
	}
	defer	conn.Close()

	_, err = conn.Write([]byte(state+"="+message))
	if err != nil {
		return false
	}

	return true
}


func (mc *mercury_conf)Log(sev int, cat string, msg string) {
	log.Printf("trace: %v, %s, %s\n", sev, cat, msg)

	sd_notify("STATUS",cat+": "+msg+"\n")
	mc.internalchan <- internal_msg { sev, cat, msg }
}


func (mc *mercury_conf)MainLoop(){
	sd_notify("READY","1")
	mc.Log(message.LOG_NOTICE, "start", "drivers available=" + strconv.Itoa(len(mc.drivers)))
	log.SetFlags(log.Ltime | log.Lshortfile)
	pid	:= strconv.Itoa(os.Getpid())

	for {
		select {
			case log := <-mc.internalchan:
				msg		:= message.CreateMessage( log.Message, MERCURYNAME, message.LOG_SYSLOG|log.Severity ).ProcID(pid)
				if log.Category != "" && log.Category != "-" {
					msg	= msg.MsgID( log.Category )
				}
				s_msg		:= msg.Stringify()

				for _,v := range mc.outputs {
					if out,ok := mc.instances[v]; ok {
						if _,ok := out.(Output); ok {
							out.(Output).Send(s_msg)
						}
					}
				}

			case	datum := <-mc.inputchan:
				msg	:= datum.Message.Stringify()
				// TODO defensive cast to avoid panic
				if _,ok := mc.instances[datum.Source].(Input); !ok {
					mc.Log(	message.LOG_ERR, "dispatcher",
						fmt.Sprintf("[%s] is not INPUT",datum.Source))
					continue
				}

				outputs	:= mc.instances[datum.Source].(Input).SendTo()
				for _,v := range outputs {
					if out,ok := mc.instances[v]; ok {
						if _,ok := out.(Output); !ok {
							mc.Log(	message.LOG_ERR, "dispatcher",
								fmt.Sprintf("[%s] from [%s] is not OUTPUT", v, datum.Source))
							continue
						}
						out.(Output).Send(msg)
					}
				}

			case	err := <-mc.errchan:
				mc.Log(message.LOG_ERR, "plugin", err.Error())

			case	<-mc.mainend:
				return
		}
	}
}


func (mc *mercury_conf)End(){
	<-mc.end
	mc.Log(message.LOG_NOTICE, "end", "stopping")

	for _,v := range mc.instances {
		if v.DriverType() == "INPUT" {
			v.End()
		}
	}

	for _,v := range mc.instances {
		if v.DriverType() == "OUTPUT" {
			v.End()
		}
	}

	sd_notify("STOPPING","1")
	time.Sleep(500*time.Millisecond)
	mc.mainend<-true
}


func (mc *mercury_conf)SetSignals() {
	signalChannel	:= make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChannel
		switch sig {
			case syscall.SIGTERM,os.Interrupt:
				mc.watchend <- true
				mc.end <- true
				return
		}
	}()
}


func (mc *mercury_conf)watch(watcher *fsnotify.Watcher) {
	defer watcher.Close()

	for {
		select {
			case event := <-watcher.Events:
				fullname	:= event.Name
				filename	:= filepath.Base(fullname)
				if filename[len(filename)-5:len(filename)] != ".conf" || filename[0] == '.' {
					continue
				}

				if	event.Op&fsnotify.Remove == fsnotify.Remove	||
					event.Op&fsnotify.Rename == fsnotify.Rename	{
					mc.stop_instance(fullname)
				}

				if	event.Op&fsnotify.Write == fsnotify.Write	{
					mc.stop_instance(fullname)
					mc.start_instance(fullname)
				}

				if	event.Op&fsnotify.Create == fsnotify.Create	{
					mc.start_instance(fullname)
				}


			case err := <-watcher.Errors:
				mc.errchan <- err

			case <-mc.watchend:
				return
		}
	}
}


func (mc *mercury_conf)Load(dirpath string, outputs []string) {
	log.Printf("start_instance %s\n", "stderr" )
	stderr_inst	:= new(output.StdErr)
	mc.instances["stderr"]	= Driver(stderr_inst)
	go stderr_inst.Run(mc.errchan)
	mc.outputs = outputs

	log.Printf("load path %s\n", dirpath )

	dir,err	:= ioutil.ReadDir(dirpath)
	exterminate(err)

	for _,file := range dir {
		filename	:= file.Name()
		if filename[len(filename)-5:len(filename)] != ".conf" || filename[0] == '.' {
			continue
		}
		mc.start_instance( filepath.Join(dirpath,filename) )
	}

	watcher,err	:= fsnotify.NewWatcher()
	exterminate(err)

	go mc.watch(watcher)

	log.Printf("watch path %s\n", dirpath )

	exterminate(watcher.Add(dirpath))
}






func (mc *mercury_conf)start_instance(fullname string) {
	mc.glock.Lock()
	defer	mc.glock.Unlock()

	//mc.Log(	message.LOG_DEBUG, "start_instance", fullname )
	log.Printf("start_instance %s\n", fullname )

	if _,ok	:= mc.file2id[fullname]; ok {
		return
	}

	t_conf		:= make(map[string]interface{})

	raw_conf,err	:= ioutil.ReadFile( fullname )
	exterminate(err)

	err	= json.Unmarshal(raw_conf, &t_conf)
	exterminate(err)

	id,ok	:= t_conf["id"].(string)
	if !ok	{
		return
	}

	driver,ok:= t_conf["driver"].(string)
	if !ok	{
		return
	}

	drv,ok	:= mc.drivers[driver]
	if !ok	{
		return
	}

	//	this line create a pointer to a new struct with the type of drv
	conf	:= reflect.New(reflect.ValueOf(drv).Elem().Type()).Interface().(Driver)
	err	= json.Unmarshal(raw_conf, conf)
	exterminate(err)

	switch drv.DriverType() {

		case	"INPUT":
			go conf.(Input).Run(mc.inputchan, mc.errchan)

		case	"OUTPUT":
			go conf.(Output).Run(mc.errchan)

		default:
			return
	}

	mc.instances[id]	= conf
	mc.file2id[fullname]	= id
	mc.Log(	message.LOG_NOTICE, "instance", fmt.Sprintf("start [%s] as %s[%s]",fullname,  drv.DriverType(),  drv.DriverName() ))
	//log.Printf("start [%s] as %s[%s]\n",fullname,  drv.DriverType(),  drv.DriverName() )
}



func (mc *mercury_conf)stop_instance(fullname string) {
	mc.glock.Lock()
	defer	mc.glock.Unlock()

	if _,ok	:= mc.file2id[fullname]; !ok {
		return
	}

	if _,ok	:= mc.instances[mc.file2id[fullname]]; !ok {
		return
	}

	mc.instances[mc.file2id[fullname]].End()

	delete(mc.file2id, fullname)
	mc.Log(	message.LOG_NOTICE, "instance", fmt.Sprintf("stop [%s]",fullname ))
}
