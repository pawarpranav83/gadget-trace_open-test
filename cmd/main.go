package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// for testing commands (temporary)
	// Build
	cmdBuild := exec.Command("make", "build")
	cmdBuild.Stdin = os.Stdin
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err := cmdBuild.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// buildOut, err := cmdBuild.Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit code =", e.ExitCode())
		default:
			panic(err)
		}
	}

	// Run
	cmdRun := exec.Command("make", "run")
	cmdRun.Stdin = os.Stdin
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	err = cmdRun.Run()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit code =", e.ExitCode())
		default:
			panic(err)
		}
	}
}
