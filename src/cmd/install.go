package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	shared "tunnel/src"
	"tunnel/src/logerr"
	"tunnel/src/types"

	"github.com/schollz/progressbar/v3"
)

func download(url string, headers *[2][]string, descp string, buffer *bytes.Buffer) (err error) {
	resp, err := fetch(url, headers)

	if err != nil {
		return
	}

	bar := progressbar.DefaultBytes(resp.ContentLength, descp)
	_, err = io.Copy(io.MultiWriter(buffer, bar), resp.Body)
	resp.Body.Close()
	fmt.Println()
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

func containsTunnelConf(zf *[]*zip.File) (bool, uint64) {
	for i, f := range *zf {
		if strings.Contains(f.Name, shared.TUNNEL_FILE) {
			return true, uint64(i)
		}
	}

	return false, 0
}

func buildExecutable() (err error) {
	return
}

func buildProject(zf *[]*zip.File, pkgName string, i uint64) (err error) {
	confZipFile := (*zf)[i] //tunnel.json file index
	confFile, err := confZipFile.Open()
	if err != nil {
		return errors.New(err.Error() + ": failed to open config file")
	}

	defer confFile.Close()

	confByte, err := io.ReadAll(confFile)
	if err != nil {
		return
	}

	var config types.Config_t
	err = json.Unmarshal(confByte, &config)

	if err != nil {
		return errors.New(err.Error() + ": failed to unmarshall json")
	}

	pkgVer := config.Project.Version
	if checkIfProjectExists(pkgName, pkgVer) {
		return errors.New("package '" + pkgName + "' already installed")
	}

	for _, pkg := range config.Dependencies {
		err = cmdInstall(&CmdArgs{
			Index: 0,
			Argc:  1,
			Argv:  &[]string{pkg.PackageUrl + "@" + pkg.Version},
		})

		if err != nil {
			if err.Error() != ("package '" + pkg.PackageUrl + "' already installed") {
				return err
			}

			logerr.Warn("package '" + pkg.PackageUrl + "' already installed")
			err = nil
			continue
		}
	}

	if config.Project.BuildExecutable {
		err = buildExecutable()
		return
	}

	if len(config.Src.Includes) == 0 {
		return errors.New("package '" + pkgName + "' has no 'include' files")
	}

	rootDir := (*zf)[0].Name

	libsPathList := []string{}
	includePathList := []string{}
	for _, file := range *zf {
		if file.Name == rootDir {
			continue
		}

		fileName := strings.Replace(file.Name, rootDir, "", 1)

		if fileName[0] == '.' {
			continue
		}

		fullPath := ""
		inInclude := false

		for _, includePaths := range config.Src.Includes {
			if strings.Contains(fileName, includePaths) {
				fullPath = filepath.Join(
					shared.UserHomeDir,
					shared.TUNNEL_DEF_PATH,
					shared.TUNNEL_PKG_PATH,
					pkgName,
					pkgVer,
					fileName,
				)

				inInclude = true
			}
		}

		if !inInclude {
			for _, libPath := range config.Src.Libs {
				if strings.Contains(fileName, libPath) {
					fullPath = filepath.Join(
						shared.UserHomeDir,
						shared.TUNNEL_DEF_PATH,
						shared.TUNNEL_CACHE_PATH,
						pkgName,
						pkgVer,
						fileName,
					)
				}
			}
		}

		if len(fullPath) == 0 {
			continue
		}

		if file.FileInfo().IsDir() {
			if inInclude {
				includePathList = append(includePathList, filepath.Join(
					shared.UserHomeDir,
					shared.TUNNEL_DEF_PATH,
					shared.TUNNEL_PKG_PATH,
					pkgName,
					pkgVer,
					fileName,
				))
			}

			os.MkdirAll(fullPath, file.Mode())
			continue
		}

		if !inInclude {
			libsPathList = append(libsPathList, fullPath)
		}

		rf, err := file.Open()
		if err != nil {
			return errors.New(err.Error() + ": failed to read file")
		}

		defer rf.Close()

		wf, err := os.Create(fullPath)
		if err != nil {
			return errors.New(err.Error() + ": failed to create file")
		}

		defer wf.Close()

		_, err = io.Copy(wf, rf)
		if err != nil {
			return errors.New(err.Error() + ": failed to copy file")
		}
	}

	if len(libsPathList) > 0 {
		libPath := filepath.Join(
			shared.UserHomeDir,
			shared.TUNNEL_DEF_PATH,
			shared.TUNNEL_LIB_PATH,
		)

		err = os.MkdirAll(libPath, os.ModePerm)
		if err != nil {
			return errors.New(err.Error() + ": failed to create 'lib' path")
		}

		cmd := exec.Command(
			config.Project.LanguageDesc.Compiler,
			strings.Join(config.Project.LanguageDesc.Flags, " "),
			"-c",
			"-o",
			filepath.Join(libPath, config.Project.Name+".o"),
			strings.Join(libsPathList, " "),
			strings.Join(config.Src.LinkerFlags, " "),
			"-I",
			strings.Join(includePathList, " -I "),
		)

		fmt.Println(cmd.String())
		cmd.Start()
	}

	return
}

func checkIfProjectExists(pkgName, pkgVer string) bool {
	path := filepath.Join(
		shared.UserHomeDir,
		shared.TUNNEL_DEF_PATH,
		shared.TUNNEL_PKG_PATH,
		pkgName,
		pkgVer,
	)

	_, err := os.Stat(path)
	return err == nil
}

func cmdInstall(cmdArgs *CmdArgs) (err error) {
	if cmdArgs.Index >= cmdArgs.Argc {
		return errors.New("needs <package-name> for option 'install'")
	}

	for cmdArgs.Index < cmdArgs.Argc {
		pkgName := (*cmdArgs.Argv)[cmdArgs.Index]
		atIndex := strings.LastIndexByte(pkgName, '@')
		headers := [2][]string{
			{"User-Agent", "tunnel"},
			{"Accept", "application/json"},
		}

		url := ""
		pkgFileName := pkgName

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
		} else {
			strArr := strings.Split(pkgName, "@")
			pkgFileName = strArr[0]
			pkgVer := strArr[1]
			url = GIT_VER_LINK_1 + pkgFileName + GIT_VER_LINK_2 + pkgVer + ".zip"

			if checkIfProjectExists(pkgFileName, pkgVer) {
				return errors.New("package '" + pkgFileName + "' already installed")
			}
		}

		var resp bytes.Buffer
		err := download(url, &headers, "downloading '"+pkgFileName+"' ", &resp)
		if err != nil {
			return errors.New(err.Error() + ": failed to fetch 'zip' file")
		}

		body := resp.Bytes()
		if err != nil {
			return errors.New(err.Error() + ": failed to read body")
		}

		zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
		if err != nil {
			return errors.New(err.Error() + ": failed to open zip")
		}

		var i uint64
		contains := false
		if contains, i = containsTunnelConf(&zipReader.File); !contains {
			return errors.New("package '" + pkgName + "' do not contains any 'tunnel.json'")
		}

		err = buildProject(&zipReader.File, pkgFileName, i)
		if err != nil {
			return err
		}

		cmdArgs.Index++
	}

	return
}
