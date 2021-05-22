package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Command ...
type Command struct {
	Cmd  string
	Args []string
}

// Cmds ...
var Cmds = make(map[string]interface{})

// func clone(server *Server, args []string) error {
// 	server2, exist := servers[args[0]]
// 	if exist {
// 		server2.Write("stop")
// 		backup(server2, []string{"make", "[before clone]"})
// 	} else {
// 		log.Printf("MCSH[stdinForward/ERROR]: Cannot find running server <%v>\n", string(args[0]))
// 	}
// 	return nil
// }

func backup(server *Server, args []string) error {
	if args[0] == "make" {
		comment := GetTimeStamp()
		if len(args) > 1 {
			comment = comment + " " + strings.Join(args[1:], " ")
		}
		dst := path.Join(backupDir, fmt.Sprintf("%s - %s", server.ServerName, comment))
		src := path.Join(filepath.Dir(server.ServerConfig.ExecPath), "world")
		log.Printf("[%s/INFO]: Making backup to %s...\n", server.ServerName, dst)
		server.Write(fmt.Sprintf("say Making backup to %s...", dst))
		err := CopyDir(src, dst)
		if err != nil {
			log.Printf("[%s/ERROR]: Backup making failed.\n", server.ServerName)
			server.Write("say Backup making failed.")
			return err
		}
		log.Printf("[%s/INFO]: Backup making successed.\n", server.ServerName)
		server.Write("say Backup making successed.")
	} else if args[0] == "" || args[0] == "list" {
		// log.Printf("[%s/INFO]: Listing backup.\n", server.ServerName)
		res, _ := ioutil.ReadDir(backupDir)
		for i, f := range res {
			fmt.Printf("[%v] %s\n", i, f.Name())
			server.Write(fmt.Sprintf("say [%v] %s", i, f.Name()))
		}
	} else if args[0] == "load" {
		i, err := strconv.Atoi(strings.Join(args[1:], ""))
		if err == nil {
			load(server, i)
		}
	}
	return nil
}

func load(server *Server, i int) error {
	res, _ := ioutil.ReadDir(backupDir)
	backup(server, []string{"make", fmt.Sprintf("Before loading %s", res[i].Name())})

	wg.Add(1)
	server.Write("stop")
	for server.online {
		time.Sleep(time.Second)
	}
	backupSavePath := path.Join(backupDir, res[i].Name())
	serverSavePath := path.Join(filepath.Dir(server.ServerConfig.ExecPath), "world")
	os.RemoveAll(serverSavePath)
	log.Printf("[%s/INFO]: Loading backup %s...\n", server.ServerName, res[i].Name())
	err := CopyDir(backupSavePath, serverSavePath)
	if err != nil {
		log.Printf("[%s/ERROR]: Backup loading failed.\n", server.ServerName)
		wg.Done()
		return err
	}
	log.Printf("[%s/INFO]: Backup loading successed.\n", server.ServerName)

	go server.Run(&wg)
	return nil
}

func start(server *Server, args []string) error {
	if server.online {
		return nil
	} else {
		server.Run(&wg)
	}

	return nil
}

func restart(server *Server, args []string) error {
	wg.Add(1)
	server.Write("stop")
	for server.online {
		time.Sleep(time.Second)
	}
	go server.Run(&wg)
	return nil
}

func init() {
	log.Println("MCSH[init/INFO]: Initializing commands...")
	Cmds["backup"] = backup
	Cmds["start"] = start
	Cmds["restart"] = restart
	// Cmds["clone"] = clone
}
