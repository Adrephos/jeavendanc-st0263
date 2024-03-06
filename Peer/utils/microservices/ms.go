package microservices

import (
	"fmt"
	"os/exec"

	utils_cmd "github.com/Adrephos/jeavendanc-st0263/Peer/utils/cmd"
)

func StartMicroservices(binaries ...string) {
	for _, command := range binaries {
		cmd := exec.Command(command)
		cmd.Start()
		fmt.Printf("%smicroservice %s with PID %d started%s\n",
			utils_cmd.YELLOW, cmd.String(), cmd.Process.Pid, utils_cmd.WHITE)
	}
	fmt.Println()
}
