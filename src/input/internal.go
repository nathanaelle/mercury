package	input

import	(
	"os"
	"time"
	"runtime"
	"strconv"

	"../message"

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


func (intl *InternalReport)DriverName() string {
	return	"i_internal"
}


func stringify_statistics(sr statistic_report) string  {
	return	"Tasks: "	+ strconv.FormatInt(int64(sr.Tasks),10)+
		", Mem: "	+ human_scale(float64(sr.MemAllocated),1024,"o") +
		", Sys: "	+ human_scale(float64(sr.SysAlloc),1024,"o") +
		", Heap: "	+ human_scale(float64(sr.Heap),1024,"o") +
		", Stack: "	+ human_scale(float64(sr.Stack),1024,"o")
}


func (intl *InternalReport)Run(dest chan<- Message, errchan chan<- error) {
	memStats	:= new(runtime.MemStats)
	ticker		:= time.Tick(intl.Tick.Get().(time.Duration) )
	pid		:= strconv.Itoa(os.Getpid())
	intl.end	= make(chan bool, 1 )

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

				msg	:= message.CreateMessage(
						stringify_statistics( stat ),
						intl.AppName,
						message.LOG_SYSLOG|message.LOG_INFO ).ProcID(pid).MsgID("statistics")
				dest	<- packmsg(intl.Id, *msg)

			case <-intl.end:
				return
		}
	}
}
