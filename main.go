package main

import (
	"fmt"
	"os"
	"os/exec"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	for {
		// TODO
	}
	fmt.Println(args)
}

func StartProc(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()
	check(err)
	err = cmd.Start()
	check(err)
	a := make([]byte, 3)
	n, err := stdout.Read(a)
	fmt.Println(n)
	fmt.Println(err)
	fmt.Println(string(a))
	err = cmd.Wait()
	fmt.Println(err)

}

func WatchFiles(args []string) {
	// TODO
}
