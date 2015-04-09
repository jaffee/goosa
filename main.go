package main

import (
	"fmt"
	fsnotify "gopkg.in/fsnotify.v1"
	"os"
	"os/exec"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	chng := make(chan fsnotify.Op, 1)
	go WatchFiles(args, chng)
	for {
		kill := make(chan int, 1)
		go StartProc(args, kill)
		for (<-chng & fsnotify.Write) > 0 {
		}
		fmt.Println("Detected change... restarting")
		kill <- 1
		time.Sleep(time.Second * 2)
	}
}

func StartProc(args []string, kill chan int) {
	fmt.Println("Making a proc")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	check(err)
	if 1 == <-kill {
		fmt.Println("Killing process")
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Printf("Error killing process: %v\n", err)
		}
		if !cmd.ProcessState.Exited() {
			fmt.Printf("Not exited? kill harder")
			cmd.Process.Kill()
		}
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
