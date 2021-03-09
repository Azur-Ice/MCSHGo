package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
)

// Command ...
type Command struct {
	Cmd  string
	Args []string
}

// Cmds ...
var Cmds = make(map[string]interface{})

func backup(server *Server, args []string) error {
	if args[0] == "make" {
		comment := GetTimeStamp()
		if len(args) > 1 {
			comment = comment + " " + strings.Join(args[1:], " ")
		}
		dst := path.Join(backupDir, fmt.Sprintf("%s - %s", server.ServerName, comment))
		src := path.Join(filepath.Dir(server.ServerConfig.ExecPath), "world")
		log.Printf("[%s/INFO]: Making backup to %s...\n", server.ServerName, dst)
		err := CopyDir(src, dst)
		if err != nil {
			log.Printf("[%s/ERROR]: Backup making failed.\n", server.ServerName)
			return err
		}
		log.Printf("[%s/INFO]: Backup making successed.\n", server.ServerName)
	} else if args[0] == "" {

	} else if args[0] == "restore" {

	}
	return nil
}

func init() {
	log.Println("MCSH[init/INFO]: Initializing commands...")
	Cmds["backup"] = backup
}
