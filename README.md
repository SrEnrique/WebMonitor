# Web Server Monitor

Is a service no configurable simple monitor, run on port 2221

## Get Information 

* [x] Hostname
* [x] Memory
  * [x] Total
  * [x] Used
  * [x] Free
* [x] CPU
  * [x] Percent one second
* [x] Disk
  * [x] External
  * [x] Partitions

## Windows server service

Create service in CMD

```
sc.exe create MiniMonSys binPath=PATH_TO_WebMonitor.exe
```

to run open windows service manager and search MiniMonSys and click in start

