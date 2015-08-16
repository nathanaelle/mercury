package	message

import	(
)


const severityMask = 0x07
const facilityMask = 0xf8


const (
	LOG_EMERG int = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)


const (
	LOG_KERN int = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)


func PriorityDecode(s_prio string) (int,error) {
	facility:= map[string]int{
		"kern":		LOG_KERN,
		"user":		LOG_USER,
		"mail":		LOG_MAIL,
		"daemon":	LOG_DAEMON,
		"auth":		LOG_AUTH,
		"syslog":	LOG_SYSLOG,
		"lpr":		LOG_LPR,
		"news":		LOG_NEWS,
		"uucp":		LOG_UUCP,
		"cron":		LOG_CRON,
		"authpriv":	LOG_AUTHPRIV,
		"ftp":		LOG_FTP,
		"local0":	LOG_LOCAL0,
		"local1":	LOG_LOCAL1,
		"local2":	LOG_LOCAL2,
		"local3":	LOG_LOCAL3,
		"local4":	LOG_LOCAL4,
		"local5":	LOG_LOCAL5,
		"local6":	LOG_LOCAL6,
		"local7":	LOG_LOCAL7,
	}

	level:= map[string]int{
		"emerg":	LOG_EMERG,
		"alert":	LOG_ALERT,
		"crit":		LOG_CRIT,
		"err":		LOG_ERR,
		"warning":	LOG_WARNING,
		"notice":	LOG_NOTICE,
		"info":		LOG_INFO,
		"debug":	LOG_DEBUG,
	}

	prios	:= FieldsFuncN(s_prio, 2, func(r rune)bool{
		switch {
			case r=='.':
				return true
			case (r>64 && r<91) || (r>96 && r<123):
				return false
			default:
				return true
		}
	})

	if len(prios)>2 {
		return -1,error_invalid_prio(s_prio)
	}

	f,ok	:= facility[prios[0]]
	if !ok {
		return -1,error_invalid_prio(s_prio)
	}

	l,ok	:= level[prios[1]]
	if !ok {
		return -1,error_invalid_prio(s_prio)
	}

	return f|l , nil
}
