package main

import (
	"fmt"
	"os/exec"
)

const binPath string = "../bin/vcoad"
const nodeHomeDirectory string = "../nodehomedirectory"
const moniker string = "base"
const chainID string = "vmt_mainnet-1"

func main() {

	chainHome := []string{"--home", nodeHomeDirectory}
	chainIdFlag := []string{"--chain-id", chainID}
	initWorkingChain := append([]string{"init", moniker}, chainHome...)
	initWorkingChain = append(initWorkingChain, chainIdFlag...)

	// Start a new local chain for working
	fmt.Println("Start a new working chain...")
	out, err := execCmdOutputOnly(initWorkingChain)
	if err != nil {
		fmt.Println("Working chain initialization error:", err)
		return
	}
	fmt.Println(string(out), "...done")
}

func execCmdOutputOnly(cmdstr []string) (out []byte, err error) {
	cmd := exec.Command(binPath, cmdstr...)
	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return
}
