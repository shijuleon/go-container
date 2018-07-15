package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main(){
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("What?")
	}
}

func run(){
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...) // trick to execute the same program as contained child process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS | // hostname
			syscall.CLONE_NEWIPC | // ipc
			syscall.CLONE_NEWPID | // process id
			syscall.CLONE_NEWNET | // networking
			syscall.CLONE_NEWUSER, // user
		UidMappings: []syscall.SysProcIDMap{ // uid
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{ // gid
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	// for completion chroot into a filesystem
	// mount proc so that the process shows up in the ps list

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
func child(){
	fmt.Printf("running %v as PID: %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

