package service

import (
	"os"
	"syscall"
)

var Signals []os.Signal = []os.Signal{os.Interrupt,os.Kill,syscall.SIGKILL,syscall.SIGSTOP,syscall.SIGHUP,syscall.SIGINT,syscall.SIGQUIT,syscall.SIGILL,syscall.SIGABRT,syscall.SIGTERM,syscall.SIGTERM}


