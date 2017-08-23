package utils

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "step 1\n")
	time.Sleep(1000. * time.Millisecond)
	fmt.Fprintf(os.Stdout, "step 2\n")
	os.Exit(0)
}

//
//func TestMain(t *testing.T) {
//	run(fakeExecCommand("sleep"))
//}
//
//func TestMainGgn(t *testing.T) {
//	run(exec.Command("ggn", "-L", "debug", "pp-dc3", "list-units"))
//}
