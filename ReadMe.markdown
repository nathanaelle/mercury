# Mercury


## What is Mercury ?

Mercury is a syslog forwarder.

## Why another one ?

the main goals are :

  * [x] providing isolation between sources, and between destinations
  * [x] providing loading/unloading conf without restart
  * [x] providing FIFO support for software with file-only logging
  * [x] outputing only rfc5424 message
  * if you use a log aggregator, all the filters and complex-dispatching will run on the aggregator not the source host
  * if you have lots of vms, containers, jails, chroot, you may need a lightweight forwarder


## configuration

### rules

  1.  Each configuration block is a specific file.
  2.  Each configuration block is a json serialization in a file.
  3.  There is two kinds of drivers : `input driver` and `output driver`
  4.  The only globally mandatories keys are : `"id"` and `"driver"`
  5.  The key `"id"` is the name of a configuration block
  6.  The key `"driver"` is the used driver for a configuration block
  7.  An `instance` is a `configuration block` with an `id` and a `driver`
  8.  An `output instance` is a `configuration block` with an `id` and a `driver` using an `output driver`
  9.  The key `"output"` is mandatory for `input drivers`
  10. The key `"output"` is a list of `"id"` of `output instance`
