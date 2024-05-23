package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	shared "tunnel/src"
	"tunnel/src/types"
)

const (
	HELP_PROMPT = "\nUsage:\n\n" +
		"\ttunnel <command> [arguments]\n\n" +
		"Available Commands:\n\n" +
		"\thelp                            Shows help.\n" +
		"\tversion                         Shows tunnel version.\n" +
		"\tinit <project-name>             Initializes a new 'tunnel.json'\n" +
		"\tupdate                          Checks for updates.\n" +
		"\tupgrade <all/package-name>      Upgrades packages.\n" +
		"\tlist                            Lists all installed packages.\n" +
		"\t-i, install <package-name>      Install a package.\n" +
		"\t-u, uninstall <package-name>    Uninstalls a package.\n" +
		"\t-r, remove <package-name>       Removes a package from cache.\n" +
		"\t-b, build                       Builds current project and creates an executable.\n" +
		"\t-r, run                         Runs current project.\n\n"
)

type CmdArgs struct {
	Argc  int
	Index int
	Argv  *[]string
}

func cmdInit(cmdArgs *CmdArgs) (err error) {
	if shared.FileExists {
		return errors.New("'tunnel.json' already exists in '" + shared.WorkingDir + "'")
	}

	if cmdArgs.Index+1 >= cmdArgs.Argc {
		return errors.New("needs <project-name> for option 'init\ntry: tunnel help")
	}

	cmdArgs.Index++
	//'init.json'
	initFile := shared.UserHomeDir + "/" + shared.TUNNEL_DEF_PATH + "/" + shared.TUNNEL_INIT_FILE
	projectName := (*cmdArgs.Argv)[cmdArgs.Index]
	isLangC := false

	if cmdArgs.Index+1 < cmdArgs.Argc {
		cmdArgs.Index++

		if (*cmdArgs.Argv)[cmdArgs.Index] == "-c" {
			isLangC = true
		} else {
			return errors.New("invalid option '" + (*cmdArgs.Argv)[cmdArgs.Index] + "' for 'init'")
		}
	}

	initFileByte, err := os.ReadFile(initFile)
	if err != nil {
		return errors.New("failed to open file '" + initFile + "'")
	}

	var config types.Config_t
	err = json.Unmarshal(initFileByte, &config)
	if err != nil {
		return errors.New("failed to 'Unmarshal' json")
	}

	config.Project.Name = projectName
	if isLangC {
		config.Project.LanguageDesc.Language = "c"
		config.Project.LanguageDesc.Compiler = "gcc"
		config.Src.Main = "main.c"
	}

	fileByte, err := json.MarshalIndent(&config, "", "\t")
	if err != nil {
		return errors.New("failed to Marshal config data")
	}

	//current_working_dir/tunnel.json
	filePath := shared.WorkingDir + "/" + shared.TUNNEL_FILE
	file, err := os.Create(filePath)
	if err != nil {
		return errors.New("failed to create file '" + filePath + "'")
	}

	defer file.Close()

	_, err = file.Write(fileByte)
	if err != nil {
		return errors.New("failed to write to the file '" + filePath + "'")
	}

	return
}

func Start(cmdArgs *CmdArgs) (err error) {
	for cmdArgs.Index < cmdArgs.Argc {
		arg := (*cmdArgs.Argv)[cmdArgs.Index]

		switch arg {
		case "help":
			fmt.Fprintf(os.Stdout, "%s\n", HELP_PROMPT)
			return

		case "version":
			fmt.Fprintf(os.Stdout, "%s\n", shared.TUNNEL_VERSION)
			return

		case "init":
			return cmdInit(cmdArgs)
		}
	}

	return
}
