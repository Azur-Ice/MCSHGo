package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	workDir   string
	backupDir string
)

var servers = make(map[string]*Server)
var wg sync.WaitGroup

func processInput() {
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		line, errRead := stdinReader.ReadBytes('\n')
		if errRead != nil {
			log.Println("MCSH[stdinForward/ERROR]: ", errRead)
		} else {
			// 去掉换行符
			if i := strings.LastIndex(string(line), "\r"); i > 0 {
				line = line[:i]
			} else {
				line = line[:len(line)-1]
			}
			// 转发正则
			res := forwardReg.FindSubmatch(line)
			if res != nil { // 转发到特定服务器
				server, exist := servers[string(res[1])]
				if exist {
					server.InChan <- string(res[2])
				} else {
					log.Printf("MCSH[stdinForward/ERROR]: Cannot find running server <%v>\n", string(res[1]))
				}
			} else { // 转发到所有服务器
				for _, server := range servers {
					server.InChan <- string(line)
				}
			}
		}
	}
}

func main() {
	// goto test

	for name, serverConfig := range MCSHConfig.Servers {
		servers[name] = NewServer(name, serverConfig)
		wg.Add(1)
		go servers[name].Run(&wg) // TODO: complete runFunc
	}
	go processInput()
	wg.Wait()

	// test:
	// 	dst := path.Join(backupDir,
	// 		fmt.Sprintf("%s - %s %s", "serverName", GetTimeStamp(), "comment"))
	// 	src := path.Join(filepath.Dir(MCSHConfig.Servers["serverName1"].ExecPath), "world")
	// 	CopyDir(src, dst)
	// 	a := [...]int{1}
	// 	fmt.Println(a[1:])
	// 	res := i_regex.PlayerOutputReg.FindStringSubmatch(`[22:25:31] [Server thread/INFO]: <_AzurIce_> sbbssbs`)
	// 	fmt.Println(res)
}

func init() {
	log.Println("MCSH[init/INFO]: Initializing...")
	workDir, _ = os.Getwd()
}
