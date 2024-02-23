package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Account struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	PubKey  string `json:"pubkey"`
	Address string `json:"address"`
}

const binPath string = "../bin/vcoad"
const walletPath string = "../wallets/"
const nodeHomeDirectory string = "../nodehomedirectory"

const mnemonicFile string = "../wallets/mnemfile.dat"
const addressFile string = "../wallets/addrfile.dat"

const moniker string = "base"
const chainID string = "vmt_mainnet-1"

func main() {

	args := os.Args[1:] // Exclude the first argument, which is the program name

	if len(args) == 0 {
		fmt.Println("Usage: genWallet numberOfAccounts")
		return
	}

	numOfaccounts, _ := strconv.Atoi(args[0])
	fmt.Println("Number of Accounts to create: ", args[0])

	chainHome := []string{"--home", nodeHomeDirectory}
	keyringBackend := []string{"--keyring-backend", "memory"}

	keysMnemonic := []string{"keys", "mnemonic"}
	keysMnemonic = append(keysMnemonic, chainHome...)

	keysAdd := []string{"keys", "add"}
	// keysDelete := []string{"keys", "delete"}
	recoverFlag := "--recover"
	outputJson := []string{"--output", "json"}

	// chainHome := []string{"--home", nodeHomeDirectory}
	// chainIdFlag := []string{"--chain-id", chainID}
	// initWorkingChain := append([]string{"init", moniker}, chainHome...)
	// initWorkingChain = append(initWorkingChain, chainIdFlag...)

	fmt.Println("Binary Path: ", binPath)
	fmt.Println("Wallet Path: ", walletPath)
	// fmt.Println("Home Directory: ", nodeHomeDirectory)

	// Start a new local chain for working
	// fmt.Println("Start a new working chain...")
	// out, err := execCmdOutputOnly(initWorkingChain)
	// if err != nil {
	// 	fmt.Println("Working chain initialization error:", err)
	// 	return
	// }
	// fmt.Println(string(out), "...done")

	mnemfile, err := os.Create(mnemonicFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer mnemfile.Close()

	addrfile, err := os.Create(addressFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer addrfile.Close()

	for i := 0; i < numOfaccounts; i++ {

		mnemonic, _ := execCmdOutputOnly(keysMnemonic)
		// fmt.Println(string(mnemonic))
		_, err = mnemfile.WriteString(strconv.Itoa(i) + ":\"" + strings.ReplaceAll(string(mnemonic), "\n", "\"\n"))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		accName := "acc" + strconv.Itoa(i)
		recoverKey := append(keysAdd, accName, recoverFlag)
		recoverKey = append(recoverKey, chainHome...)
		recoverKey = append(recoverKey, outputJson...)
		recoverKey = append(recoverKey, keyringBackend...)
		fmt.Println(recoverKey)

		output, _ := execCmdInputOutput(recoverKey, mnemonic)
		var acc Account
		err = json.Unmarshal(output, &acc)
		if err != nil {
			fmt.Println("Error unmarshaling to JSON:", err)
			return
		}

		fmt.Println(acc.Name, acc.Address, acc.PubKey, acc.Type)

		_, err = addrfile.WriteString(strconv.Itoa(i) + "," + acc.Address + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		// delete the key after
		// deleteKey := append(keysDelete, accName)
		// deleteKey = append(deleteKey, "--yes")
		// out, _ := execCmdOutputOnly(deleteKey)
		// fmt.Println(string(out))
	}
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

func execCmdInputOutput(cmdstr []string, in []byte) (out []byte, err error) {
	cmd := exec.Command(binPath, cmdstr...)
	cmdInput, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error setting up StdinPipe:", err)
		return
	}

	cmdOutput, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error setting up StdoutPipe:", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error to start cmd:", err)
		return
	}

	// input := []byte(string(mnemonic))
	_, err = cmdInput.Write(in)
	if err != nil {
		fmt.Println("Error writing mnemonic:", err)
		return
	}
	cmdInput.Close()

	// Read the output
	out = make([]byte, 0)
	buf := make([]byte, 1024) // Read in chunks of 1024 bytes
	for {
		n, err := cmdOutput.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading cmdOutput:", err)
			}
			break
		}
		out = append(out, buf[:n]...)
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error during waiting:", err)
		return
	}
	return
}
