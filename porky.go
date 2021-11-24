package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	sshConfig := ".ssh/config"
	fmt.Println("go!")
	hosts := parse(sshConfig)
	raw, _ := json.MarshalIndent(hosts, "", "  ")
	fmt.Println(string(raw))
}
