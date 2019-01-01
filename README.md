# PS

A Go library which runs the Unix 'ps' command and returns structured information 
about the process.

This library essentially runs:

```
bash -c ps -aux | grep [p]rocess_name
```
And parses the output of the command.

# Installation

> $ go get github.com/umahmood/ps

# Usage

```
package main

import (
    "fmt"

    "github.com/umahmood/ps"
)

func main() {
    proc, err := ps.Snapshot("firefox")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(proc.PID)
    fmt.Println("")
    fmt.Println(proc)
}
```
Output:
```
19352

User   : usman
PID    : 19352
CPU    : 0.00
MEM    : 8.10
VSZ    : 8353516
RSS    : 1350828
TTY    : ??
STAT   : S
START  : 22Dec18
TIME   : 334:44.58
COMMAND: /Applications/Firefox.app/Contents/MacOS/firefox
```

[What do the output fields of the ps command mean?](https://kb.iu.edu/d/afnv)

# Limitations

The library only matches unique running processes. If there are multiple processes 
with the same name, then an error is returned. i.e. if we have the following 
running processes:

```
- proc_foo --arg=1
- proc_foo --arg=2
- proc_zap
```

As there are multiple `proc_foo` processes, calling `ps.Snapshot("proc_foo")` will 
throw an error. I may update the library to change this, but currently it fits 
my needs. 

This library has been tested Linux and MacOS, it may not work on windows. 

# Documentation

> http://godoc.org/github.com/umahmood/ps

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
