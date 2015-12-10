package	input

import	(
	"os"
	"time"
	"runtime"
	"strconv"

	"github.com/nathanaelle/syslog5424"

	types	"github.com/nathanaelle/useful.types"
)

type	(
	InternalReport		struct {
		GenericInput
		Tick		types.Duration	`json:"tick"`
		AppName		string		`json:"appname"`
	}

	statistic_report	struct {
		Tasks		uint
		MemAllocated	uint
		SysAlloc	uint
		Heap		uint
		Stack		uint
	}
)


func (intl *InternalReport) DriverName() string {
	return	"i_internal"
}

func (intl *InternalReport) Configure(errchan chan<- error) {
	intl.end	= make(chan bool, 1 )
	intl.errchan	= errchan

	if intl.AppName == "" {
		panic("AppName mandatory")
	}

	if intl.Tick.Get().(time.Duration) <= time.Second {
		intl.Tick.Set("300s")
	}
}


func stringify_statistics(sr statistic_report) string  {
	return	"Tasks: "	+ strconv.FormatInt(int64(sr.Tasks),10)+
		", Mem: "	+ human_scale(float64(sr.MemAllocated),1024,"o") +
		", Sys: "	+ human_scale(float64(sr.SysAlloc),1024,"o") +
		", Heap: "	+ human_scale(float64(sr.Heap),1024,"o") +
		", Stack: "	+ human_scale(float64(sr.Stack),1024,"o")
}


func (intl *InternalReport) Run(dest chan<- Message) {
	memStats	:= new(runtime.MemStats)
	ticker		:= time.Tick(intl.Tick.Get().(time.Duration) )
	pid		:= strconv.Itoa(os.Getpid())

	for {
		select {
			case <-ticker:
				runtime.ReadMemStats(memStats)

				stat	:= statistic_report {
					Tasks:		uint(runtime.NumGoroutine()),
					MemAllocated:	uint(memStats.Alloc),
					SysAlloc:	uint(memStats.Sys),
					Heap:		uint(memStats.HeapAlloc),
					Stack:		uint(memStats.StackInuse),
				}

				dest	<- packmsg(intl.Id, syslog5424.CreateMessage(
						intl.AppName,
						syslog5424.LOG_SYSLOG|syslog5424.LOG_INFO,
						stringify_statistics( stat ) ).ProcID(pid).MsgID("statistics"))

			case <-intl.end:
				return
		}
	}
}
