package utils

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Version information.
var (
	Version   = "None"
	GitHash   = "None"
	GitBranch = "None"
	BuildTS   = "None"
)

// PrintRawInfo prints the version information without long info.
func PrintRawInfo(app string) {
	fmt.Printf("Welcome to %s\n", app)
	fmt.Printf("Release Version: %s\n", Version)
	fmt.Printf("Git Commit Hash: %s\n", GitHash)
	fmt.Printf("Git Branch: %s\n", GitBranch)
	fmt.Printf("UTC Build Time: %s\n", BuildTS)
}

// LogRawInfo prints the version information.
func LogRawInfo(app string) {
	log.Infof("Welcome to %s", app)
	log.Infof("Release Version: %s", Version)
	log.Printf("Git Commit Hash: %s", GitHash)
	log.Printf("Git Branch: %s", GitBranch)
	log.Printf("UTC Build Time: %s", BuildTS)
}
