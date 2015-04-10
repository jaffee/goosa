package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	cmd := exec.Command("go", "run", "mt.go")
	// cmd := exec.Command("./run.sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	fmt.Println("Starting....")

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting process")
		os.Exit(1)
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		fmt.Printf("Error getting pgid %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Sleeping.... ")
	time.Sleep(3 * time.Second)

	fmt.Println("Awake.... Killing....")
	err = syscall.Kill(-pgid, syscall.SIGTERM)
	fmt.Printf("Result of pgid Kill: %v\n", err)

	err = cmd.Wait()
	fmt.Printf("Result of Wait: %v\n", err)

	fmt.Println(cmd)
	// fmt.Println("Sleeping again...")
	// time.Sleep(3 * time.Second)
	// fmt.Println("And we're done.")
}
