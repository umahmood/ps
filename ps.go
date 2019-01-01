package ps

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	// ErrEmptyProcName process name is empty
	ErrEmptyProcName = errors.New("ps process name is empty")
	// ErrProcNotFound process not found
	ErrProcNotFound = errors.New("ps process not found")
	// ErrMultipleProcs matched multiple processes
	ErrMultipleProcs = errors.New("ps matched multiple processes")
	// ErrParsingPSOutput error parsing ps output
	ErrParsingPSOutput = errors.New("ps error parsing ps output")
)

// processStatusCodes current status of the process.
var processStatusCodes = map[string]string{
	"D": "uninterruptible sleep (usually IO)",
	"R": "running or runnable (on run queue)",
	"S": "interruptible sleep (waiting for an event to complete)",
	"T": "stopped, either by a job control signal or because it is being traced",
	"W": "paging (not valid since the 2.6.xx kernel)",
	"X": "dead (should never be seen)",
	"Z": "defunct ('zombie') process, terminated but not reaped by its parent",
	"<": "high-priority (not nice to other users)",
	"N": "low-priority (nice to other users)",
	"L": "has pages locked into memory (for real-time and custom IO)",
	"s": "is a session leader",
	"l": "is multi-threaded (using CLONE_THREAD, like NPTL pthreads do)",
	"+": "is in the foreground process group",
}

// PS process information
type PS struct {
	USER    string  // username of the process's owner
	PID     int64   // process ID number
	CPU     float64 // how much of the CPU the process is using
	MEM     float64 // how much memory the process is using
	VSZ     int64   // virtual memory usage
	RSS     int64   // real memory usage
	TTY     string  // terminal associated with the process
	STAT    string  // process status code
	START   string  // time when the process started
	TIME    string  // total CPU usage
	COMMAND string  // name of the process, including arguments, if any
}

// String format process information.
func (p PS) String() string {
	const tmpl = `USER   : %s
PID    : %d
CPU    : %.2f
MEM    : %.2f
VSZ    : %d
RSS    : %d
TTY    : %s
STAT   : %s
START  : %s
TIME   : %s
COMMAND: %s
`
	return fmt.Sprintf(tmpl, p.USER,
		p.PID,
		p.CPU,
		p.MEM,
		p.VSZ,
		p.RSS,
		p.TTY,
		p.STAT,
		p.START,
		p.TIME,
		p.COMMAND)
}

// Snapshot take a snapshot of the process.
func Snapshot(name string) (*PS, error) {
	output, err := runCmd(name)
	if err != nil {
		return nil, err
	}
	return parsePSOutput(string(output))
}

func runCmd(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyProcName
	}
	var psArg string
	switch runtime.GOOS {
	case "darwin":
		psArg = "aux"
	case "linux":
		psArg = "-aux"
	}
	arg := fmt.Sprintf("ps %s | grep [%s]%s", psArg, name[0:1], name[1:])
	cmd := exec.Command("bash", "-c", arg)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	slurp, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, ErrProcNotFound
		}
		return nil, err
	}
	return slurp, nil
}

// parsePSOutput parses the output from the ps command.
func parsePSOutput(output string) (*PS, error) {
	procs := strings.Split(output, "\n")
	psCount := 0
	for _, proc := range procs {
		if proc == "" {
			continue
		}
		psCount++
	}
	if psCount > 1 {
		return nil, ErrMultipleProcs
	}
	fields := strings.Fields(output)
	if len(fields) < 11 {
		return nil, ErrParsingPSOutput
	}
	if len(fields) > 9 {
		a := strings.Join(fields[10:], " ")
		fields = fields[:10]
		fields = append(fields, a)
	}
	ps := &PS{
		USER:    fields[0],
		PID:     toInt64(fields[1]),
		CPU:     toFloat64(fields[2]),
		MEM:     toFloat64(fields[3]),
		VSZ:     toInt64(fields[4]),
		RSS:     toInt64(fields[5]),
		TTY:     fields[6],
		STAT:    fields[7],
		START:   fields[8],
		TIME:    fields[9],
		COMMAND: fields[10],
	}
	return ps, nil
}

// toInt64 convert a string number to int64
func toInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return n
}

// toFloat64 convert a string number to float64
func toFloat64(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1
	}
	return n
}
