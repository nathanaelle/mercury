# Input Plugins for Mercury


## `o_appendfile` driver

Append log to a file

  * type : output driver
  * arguments :
    * `file`
      * mandatory
      * meaning : name of the local file where the log are stored


## `o_stderr` driver

Log to stderr

  * type : output driver


## `o_tcp5424` driver

Send to a remote TCP host

* type : output driver
* arguments :
  * `host`
    * mandatory
    * meaning : `"hostname:port"` of the destination


## `o_tls5424` driver

Send to a remote TCP+TLS host

  * type : output driver
  * arguments :
    * `host`
      * mandatory
      * meaning : `"hostname:port"` of the destination
