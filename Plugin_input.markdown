# Input Plugins for Mercury


## `i_devlog` driver

  * type : input driver
  * arguments :
    * `devlog`
      * optional
      * default : `"/dev/log"`
      * meaning : the name of the datagram unix socket to listen to


## `i_krnlog` driver

  * type : input driver
  * arguments :
    * none


## `i_journald` driver

  * type : input driver
  * arguments :
    * `journald`
      * optional
      * default : `"/run/systemd/journal/syslog"`
      * meaning : the name of the datagram unix socket to listen to


## `i_internal` driver

* type : input driver
* arguments :
  * `appname`
    * mandatory
    * meaning : application name of any new message
  * `tick`
    * optional
    * default : `"300"`
    * meaning : the name of the datagram unix socket to listen to


## `i_fifo` driver

  * type : input driver
    * arguments :
    * `source`
      * mandatory
      * meaning : name of the local file where the log are stored
    * `appname`
      * mandatory
      * meaning : application name of any new message
    * `priority`
      * mandatory
      * meaning : `"facility.level"` of any new message


## `i_tailfile` driver

  * type : input driver
  * arguments :
    * `source`
      * mandatory
      * meaning : name of the local file where the log are stored
    * `appname`
      * mandatory
      * meaning : application name of any new message
    * `priority`
      * mandatory
      * meaning : `"facility.level"` of any new message
