package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
)

const binPath string = "../bin/vcoad"
const accountFile string = "../wallets/accounts.csv"
const nodeHomeDirectory string = "../nodehomedirectory"

func main() {

	chainHome := []string{"--home", nodeHomeDirectory}

	// Open the CSV file
	file, err := os.Open(accountFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Process the records
	for _, row := range records {
		if row[2] != "0" {
			cmdstr := []string{"add-genesis-account", row[1], row[2] + row[3]}
			cmdstr = append(cmdstr, chainHome...)
			fmt.Println(row[0], ": ", cmdstr)
			out, _ := execCmdOutputOnly(cmdstr)
			fmt.Println(string(out))
		}
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
