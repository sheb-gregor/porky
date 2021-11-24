package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"strings"
)

func parse(path string) map[string]HostParams {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("[ERROR]:", err.Error())
	}
	buf := bytes.NewBuffer(rawFile)

	scanner := bufio.NewScanner(buf)

	collecting := false
	hosts := map[string]HostParams{}
	host := new(HostParams)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) != 2 {
			continue
		}

		switch parts[0] {
		case "Host":
			if collecting {
				hosts[host.Host] = *host
				host = new(HostParams)
			}
			host.Host = parts[1]
			host.Port = "22"
			host.KeyFile = "~/.ssh/.id_rsa"
			collecting = true
		case "Hostname":
			host.Origin = parts[1]
		case "User":
			host.User = parts[1]
		case "Port":
			host.Port = parts[1]
		case "IdentityFile":
			host.KeyFile = parts[1]
		}
	}
	return hosts
}
