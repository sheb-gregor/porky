package main

type HostParams struct {
	Host    string
	Origin  string
	User    string
	Port    string
	KeyFile string
}

const sshTemplate = `
Host {{ .Host  }}
  Hostname {{ .Origin }}
  User {{ .User  }}
  Port {{ .Port }}
  PasswordAuthentication no
  IdentityFile {{ .KeyFile }}
`
