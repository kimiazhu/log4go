This project forked from [alecthomas/log4go](http://github.com/alecthomas/log4go) and adapted with golang v1.5.3.

I hava change some features and done with some bugfixes. Include:

1. Reuse the exist log file (only if the size and maxline not reached) rather than create a new file every boot time.

2. Rotate log filename bug fixes with filelog

3. log4go_test.go now can run on go1.5.3

4. More accurate log time

5. Exclude specified log by adding the exclude pattern, which will match the *source* field. See &lt;exclude&gt; attribute in [example.xml](https://github.com/kimiazhu/log4go/blob/master/example.xml)

6. Auto load configuration in init() method, when the `log4go.xml` placed in exec dir or `{exec_dir}/conf` dir.

### Installation:
- Run `go get github.com/kimiazhu/log4go`

### Usage:
- Add the following import:
import log "github.com/kimiazhu/log4go"

### Acknowledgements:
- pomack
  For providing awesome patches to bring log4go up to the latest Go spec