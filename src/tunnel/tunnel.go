package tunnel

import (
	"os"
	shared "tunnel/src"
	"tunnel/src/logerr"
)

func setHomeDir() {
	var err error
	shared.UserHomeDir, err = os.UserHomeDir()

	if err != nil {
		logerr.Error("failed to get '$HOME' environment variable")
		os.Exit(1)
	}
}

func setWorkingDir() {
	var err error
	shared.WorkingDir, err = os.Getwd()

	if err != nil {
		logerr.Error("failed to get working directory")
		os.Exit(1)
	}
}

func getFileExists() {
	_, err := os.Stat(shared.WorkingDir + "/" + shared.TUNNEL_FILE)
	shared.FileExists = (err == nil)
}

func Init() {
	setHomeDir()
	setWorkingDir()
	getFileExists()
}
