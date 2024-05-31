package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	shared "tunnel/src"
	"tunnel/src/types"

	"github.com/schollz/progressbar/v3"
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

	GIT_LATEST_LINK_1 = "https://api.github.com/repos/"
	GIT_LATEST_LINK_2 = "/releases/latest"
	GIT_VER_LINK_1    = "https://github.com/"
	GIT_VER_LINK_2    = "/archive/"
)

type CmdArgs struct {
	Argc  int
	Index int
	Argv  *[]string
}

func download(url string, headers *[2][]string, descp string, file *os.File) (err error) {
	resp, err := fetch(url, headers)
	if err != nil {
		return
	}

	bar := progressbar.DefaultBytes(resp.ContentLength, descp)
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	resp.Body.Close()
	return
}

func fetch(url string, headers *[2][]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	for _, item := range headers {
		req.Header.Set(item[0], item[1])
	}

	return http.DefaultClient.Do(req)
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

func cmdInstall(cmdArgs *CmdArgs) (err error) {
	if cmdArgs.Index+1 >= cmdArgs.Argc {
		return errors.New("needs <package-name> for option 'install'")
	}

	cmdArgs.Index++

	for cmdArgs.Index < cmdArgs.Argc {
		pkgName := (*cmdArgs.Argv)[cmdArgs.Index]
		atIndex := strings.LastIndexByte(pkgName, '@')
		headers := [2][]string{
			{"User-Agent", "tunnel"},
			{"Accept", "application/json"},
		}
		url := ""
		pkgFileName := ""

		if atIndex == -1 {
			jsonResp, err := fetch(
				GIT_LATEST_LINK_1+pkgName+GIT_LATEST_LINK_2,
				&headers,
			)

			if err != nil {
				return errors.New(err.Error() + ": failed to fetch")
			}

			defer jsonResp.Body.Close()

			if jsonResp.StatusCode != 200 {
				return errors.New("no package '" + pkgName + "' found")
			}

			respBytes, err := io.ReadAll(jsonResp.Body)
			if err != nil {
				return errors.New(err.Error() + ": failed to read 'json' body")
			}

			jsonBody := make(map[string]interface{})
			err = json.Unmarshal(respBytes, &jsonBody)
			if err != nil {
				return errors.New(err.Error() + ": failed to unmarshal 'json'")
			}

			url = jsonBody["zipball_url"].(string)
			pkgFileName = strings.Replace(pkgName, "/", "_", 1)
		} else {
			strArr := strings.Split(pkgName, "@")
			pkgFileName = strArr[0]
			pkgVer := strArr[1]

			url = GIT_VER_LINK_1 + pkgFileName + GIT_VER_LINK_2 + pkgVer + ".zip"
			pkgFileName = strings.Replace(pkgFileName, "/", "_", 1)
		}

		file, err := os.Create(
			shared.UserHomeDir + "/" + shared.TUNNEL_DEF_PATH + "/" +
				shared.TUNNEL_CACHE_PATH + "/" + pkgFileName + ".zip")

		if err != nil {
			return errors.New(err.Error() + ": failed to create a file")
		}

		defer file.Close()

		err = download(url, &headers, "downloading: "+pkgName, file)
		if err != nil {
			return errors.New(err.Error() + ": failed to fetch 'zip' file")
		}

		cmdArgs.Index++
	}

	return
}

func Start(cmdArgs *CmdArgs) (err error) {
	for cmdArgs.Index < cmdArgs.Argc {
		arg := &(*cmdArgs.Argv)[cmdArgs.Index]

		switch *arg {
		case "help":
			fmt.Fprintf(os.Stdout, "%s\n", HELP_PROMPT)
			return

		case "version":
			fmt.Fprintf(os.Stdout, "%s\n", shared.TUNNEL_VERSION)
			return

		case "init":
			return cmdInit(cmdArgs)

		case "install":
			return cmdInstall(cmdArgs)
		}

		cmdArgs.Index++
	}

	return
}
