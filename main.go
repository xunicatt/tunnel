package main

import (
	"os"
	"tunnel/src/cmd"
	"tunnel/src/logerr"
	"tunnel/src/tunnel"
)

func main() {
	tunnel.Init()
	
	err := cmd.Start(
		&cmd.CmdArgs{
			Argv:  &os.Args,
			Argc:  len(os.Args),
			Index: 1,
		},
	)

	if err != nil {
		logerr.Error("%s", err.Error())
		os.Exit(1)
	}
}
