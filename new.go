package cloudinitgo

import (
	"net/http"
)

type Config struct {
	// Example: "admin"
	UserName string
	// Example: "secretpassword123"
	Password string
	// Example: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."
	SshPublicKey string
	// Example: "cloud-init/data"
	DataDir string
}

type cloudInit struct {
	Port   int
	Config Config
	server *http.Server
}

// New creates a new cloud-init server.
// example dataDir: "cloud-init/data" default: "cloud-init/data"
func New(c Config) *cloudInit {

	return &cloudInit{
		Port:   8181,
		Config: c,
	}
}
