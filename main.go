package main

import (
	"fmt"
	fsnotify "gopkg.in/fsnotify.v1"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: gosup <program> [args...]")
		os.Exit(0)
	}
	chng := make(chan fsnotify.Op, 1)
	go WatchFiles(args, chng)
	recvSig := make(chan os.Signal, 10)
	signal.Notify(recvSig, os.Interrupt, os.Kill)

	kill := make(chan int)
	go waitAndKill(recvSig, kill)
	for {
		go StartProc(args, kill)
		for (<-chng & fsnotify.Write) <= 0 {
		}
		fmt.Println("Detected change... restarting")
		kill <- 1
		<-kill
	}
}

func waitAndKill(recvSig chan os.Signal, killChan chan int) {
	<-recvSig
	killChan <- 1
	<-killChan
	os.Exit(1)
}

func StartProc(args []string, kill chan int) {
	fmt.Println("Making a proc")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	check(err)
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	check(err)

	if 1 == <-kill {
		fmt.Println("Killing process")
		err = syscall.Kill(-pgid, syscall.SIGTERM)
		if err != nil {
			fmt.Printf("Error killing process: %v\n", err)
		}
		kill <- 2
	}
	cmd.Wait()
}

func WatchFiles(args []string, chng chan fsnotify.Op) {
	watcher, err := fsnotify.NewWatcher()
	check(err)
	watcher.Add(args[2])
	fmt.Printf("watching %v\n", args[2])
	defer watcher.Close()
	for {
		select {
		case event := <-watcher.Events:
			fmt.Printf("Watcher event: %v\n", event)
			chng <- event.Op
		case err := <-watcher.Errors:
			check(err)
		}
	}
}
