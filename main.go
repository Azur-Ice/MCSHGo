package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Dirs
var (
	workDir, _    = os.Getwd()
	scriptsDir, _ = filepath.Abs("./Scripts")
	backupsDir, _ = filepath.Abs("./Backups")
)
var scriptSuffix = ""
var (
	commandRegex    *regexp.Regexp
	forwardReg      *regexp.Regexp
	outputFormatReg *regexp.Regexp
)

var cmds = make(map[string]interface{})

var wg sync.WaitGroup

///// command part /////

// Command describe args of a command
type Command struct {
	cmd  string
	args []string
}

func backup(server *Server, args []string) error {
	if args[0] == "make" {
		comment := ""
		if len(args) > 1 {
			comment = strings.Join(args[1:], "")
		}
		dst := path.Join(backupsDir, fmt.Sprintf("%s - %s %s", server.name, getTimeStamp(), comment))
		src := path.Join(server.config.RootFolder, "world")
		log.Printf("[%s/INFO]: Making backup to %s...\n", server.name, dst)
		err := copyDir(src, dst)
		if err != nil {
			log.Printf("[%s/ERROR]: Backup making failed.\n", server.name)
			return err
		}
		log.Printf("[%s/INFO]: Backup making successed.\n", server.name)
	}
	return nil
}

///// config part /////

// ServerConfig holds all fields of every server in "config.yml/servers"
type ServerConfig struct {
	RootFolder string `yaml:"rootFolder"`
}

// Config holds all fields in "config.yml"
type Config struct {
	CommandPrefix string                  `yaml:"command_prefix"`
	Servers       map[string]ServerConfig `yaml:"servers"`
}

var mcshConfig = Config{
	CommandPrefix: "#",
	Servers: map[string]ServerConfig{
		"serverName1": ServerConfig{
			RootFolder: "path/to/your/server/root/folder",
		},
	},
}

///// server part /////

// Server contains info of a server
type Server struct {
	name           string
	config         ServerConfig
	online         bool
	stdin          io.WriteCloser
	stdout, stderr io.ReadCloser
}

func (server *Server) run() {
	defer func() {
		recover()
		server.online = false
		wg.Done()
		return
	}()
	// fmt.Println(exec.LookPath(path.Join(scriptsDir, server.name+scriptSuffix)))
	cmd := exec.Command(path.Join(scriptsDir, server.name+scriptSuffix))
	server.stdin, _ = cmd.StdinPipe()
	server.stdout, _ = cmd.StdoutPipe()
	server.stderr, _ = cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Panicf("server<%s>: Error when starting:\n%s", server.name, err.Error())
	}
	server.online = true
	go server.forwardStdout()
	if err := cmd.Wait(); err != nil {
		log.Panicf("server<%s>: Error when running:\n%s", server.name, err.Error())
	}
}
func (server *Server) forwardStdout() {
	defer func() {
		recover()
		return
	}()
	cache := ""
	buf := make([]byte, 1024)
	for {
		num, err := server.stdout.Read(buf)
		if err != nil && err != io.EOF { //非EOF错误
			log.Panicln(err)
		}
		if num > 0 {
			str := outputFormatReg.ReplaceAllString(cache+string(buf[:num]), "["+server.name+"/$2]") // 格式化读入的字符串
			lines := strings.SplitAfter(str, "\n")                                                   // 按行分割开
			for i := 0; i < len(lines)-1; i++ {
				log.Print(lines[i])
			}
			cache = lines[len(lines)-1] //最后一行下次循环处理
		}
	}
}
func forwardStdin() {
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
			if res != nil { // 特定服务器
				server, valid := servers[string(res[1])]
				if valid {
					if command, ok := getCommand(res[2]); ok { // is #Command, execute

						if command.cmd == "run" && !server.online {
							server.run()
						}

						cmdFun, exist := cmds[command.cmd]
						if !exist {
							log.Println("MCSH[stdinForward/ERROR]: Command \"" + command.cmd + "\" not found.")
						} else {
							cmdFun.(func(server *Server, args []string) error)(server, command.args)
						}

					} else {
						_, errWrite := server.stdin.Write(append(res[2], '\n')) // is not #Command, forward
						if errWrite != nil {
							log.Println("MCSH[stdinForward/ERROR]: Server stdin write failed - ", errWrite)
						}
					}
				} else {
					log.Printf("MCSH[stdinForward/ERROR]: Cannot find running server <%v>\n", string(res[1]))
				}
			} else { // 全部服务器
				for _, server := range servers {
					server.stdin.Write(append(line, '\n'))
				}
			}
		}
	}
}

///// util part /////
func getTimeStamp() string {
	return time.Now().Format("2006-01-02 15-04-05")
}
func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
func copyDir(srcDir string, dstDir string) error {
	err := os.Mkdir(dstDir, 0666)
	if err != nil {
		log.Println(err)
	}
	fileInfoList, _ := ioutil.ReadDir(srcDir)
	for i := 0; i < len(fileInfoList); i++ {
		// fmt.Println("Copying: ", fileInfoList[i].Name(), fileInfoList[i].IsDir(), "...")
		if fileInfoList[i].IsDir() {
			copyDir(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		} else {
			copyFile(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		}
	}
	return nil
}
func data2yaml(data Config) []byte {
	yaml, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Println(err)
	}
	return yaml
}
func getCommand(str []byte) (Command, bool) {
	command := Command{}
	commandStr := ""
	if res := commandRegex.FindSubmatch(str); len(res) > 1 {
		commandStr = string(res[1])
	}
	if commandStr == "" { // 命令为空
		return command, false
	}

	cmd := strings.Split(string(commandStr), " ")
	command.cmd = cmd[0]
	if len(cmd) > 1 {
		command.args = cmd[1:]
	}
	// fmt.Println(command)

	return command, true
}
func readConfig() {
	configYaml, err := ioutil.ReadFile("./config.yml")
	if err != nil { // 读取文件发生错误
		if os.IsNotExist(err) { // 文件不存在，创建并写入默认配置
			log.Println("MCSH: Cannot find config.yml, creating...")
			ioutil.WriteFile("./config.yml", data2yaml(mcshConfig), 0666)
			log.Println("MCSH: Successful created config.yml, please complete the config.")
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	mcshConfig = Config{}
	err = yaml.Unmarshal(configYaml, &mcshConfig)
}

///// others /////
func initCommands() {
	cmds["backup"] = backup
}
func initRegexs() {
	commandRegex = regexp.MustCompile("^" + mcshConfig.CommandPrefix + "(.*)")
	forwardReg = regexp.MustCompile(`(.+?) *\| *(.+)`)
	outputFormatReg = regexp.MustCompile(`(\[\d\d:\d\d:\d\d\]) *\[.+?\/(.+?)\]`)
}
func initDirs() {
	os.Mkdir(scriptsDir, 0666)
	os.Mkdir(backupsDir, 0666)
}
func init() {
	initRegexs()
	initCommands()
	initDirs()
	readConfig()
	log.Println("[MCSH/INFO]: Running on", runtime.GOOS)
	if runtime.GOOS == "windows" {
		scriptSuffix = ".bat"
	} else {
		scriptSuffix = ".sh"
	}
}

var servers = make(map[string]*Server)

func main() {
	// goto test
	for name, serverConfig := range mcshConfig.Servers {
		servers[name] = &Server{name: name, config: serverConfig}
		go servers[name].run()
		wg.Add(1)
	}
	go forwardStdin()
	wg.Wait()
	// test:
	// 	os.Mkdir(path.Join(workDir, "test"), 0666)
	// 	fmt.Println(path.Join(workDir, "test"))
	// 	fmt.Println(getTimeStamp())
}
