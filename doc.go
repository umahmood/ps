/*
Package ps which runs 'ps' unix command and returns the process information.

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
*/
package ps
