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

	"gopkg.in/yaml.v3"
)

var (
	workDir, _    = os.Getwd()
	scriptsDir, _ = filepath.Abs("./Scripts")
	scriptSuffix  = ""
)
var (
	commandRegexString string
	commandRegex       *regexp.Regexp
)

var cmds = make(map[string]interface{})

var wg sync.WaitGroup

///// command part /////

// Command describe args of a command
type Command struct {
	cmd  string
	args []string
}

func backup(server *Server, args []string) {
	log.Println("MCSH[command-backup/INFO]: Backup!")
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
	go asyncLog(server.name, server.stdout)
	go asyncLog(server.name, server.stderr)
	if err := cmd.Wait(); err != nil {
		log.Panicf("server<%s>: Error when running:\n%s", server.name, err.Error())
	}
}
func asyncLog(name string, readCloser io.ReadCloser) error {
	var outputReplaceRegString = `(\[\d\d:\d\d:\d\d\]) *\[.+?\/(.+?)\]`
	outputReplaceReg, err := regexp.Compile(outputReplaceRegString)
	if err != nil {
		log.Println("MCSH[outputForward/ERROR]: Regex compile failed - ", err)
	}
	cache := ""
	buf := make([]byte, 1024)
	for {
		num, err := readCloser.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			// b := buf[:num]
			s := outputReplaceReg.ReplaceAllString(string(buf[:num]), "["+name+"/$2]")
			lines := strings.Split(s, "\n")
			lines[0] = cache + lines[0]
			for i := 0; i < len(lines)-1; i++ {
				log.Println(lines[i])
			}
			cache = lines[len(lines)-1]
		}
	}
}
func asyncForwardStdin() {
	var forwardRegString = `(.+?) *\| *(.+)`
	forwardReg, errCompile := regexp.Compile(forwardRegString)
	if errCompile != nil {
		log.Println("MCSH[stdinForward/ERROR]: Regex compile failed - ", errCompile)
	}

	stdinReader := bufio.NewReader(os.Stdin)
	for {
		line, errRead := stdinReader.ReadBytes('\n')
		if errRead != nil {
			log.Println("MCSH[stdinForward/ERROR]: ", errRead)
		} else {
			line = line[:len(line)-1]
			if line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			// log.Println(line)
			res := forwardReg.FindSubmatch(line)
			if res != nil {
				// log.Println(res)
				server, valid := servers[string(res[1])]
				if valid {
					if command, ok := getCommand(res[2]); ok { // is #Command, execute
						fmt.Println(command)
						if command.cmd == "run" && !server.online {
							server.run()
						}

						cmdFun, exist := cmds[command.cmd]
						if !exist {
							log.Println("MCSH[stdinForward/ERROR]: Command \"" + command.cmd + "\" not found.")
						} else {
							cmdFun.(func(server *Server, args []string))(server, command.args)
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
			} else {
				for _, server := range servers {
					server.stdin.Write(append(line, '\n'))
				}
			}
		}
	}
}

///// util part /////

func data2yaml(data Config) []byte {
	yaml, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Println(err)
	}
	return yaml
}

func getCommand(str []byte) (Command, bool) {
	command := Command{}
	commandStr := commandRegex.FindSubmatch(str)[1]
	fmt.Println(commandStr)
	if string(commandStr) == "" { // 命令为空
		return command, false
	}

	cmd := strings.Split(string(commandStr), " ")
	command.cmd = cmd[0]
	if len(cmd) > 1 {
		command.args = cmd[1:]
	}
	fmt.Println(command)

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
func init() {
	os.Mkdir(scriptsDir, 0666)
	readConfig()
	commandRegexString = "^" + mcshConfig.CommandPrefix + "(.*)"
	fmt.Println(commandRegexString)
	commandRegex = regexp.MustCompile(commandRegexString)
	log.Println("MCSH[init/INFO]: Running on", runtime.GOOS)
	if runtime.GOOS == "windows" {
		scriptSuffix = ".bat"
	} else {
		scriptSuffix = ".sh"
	}
}

var servers = make(map[string]*Server)

func main() {
	// readConfig()
	for name, serverConfig := range mcshConfig.Servers {
		servers[name] = &Server{name: name, config: serverConfig}
		go servers[name].run()
		wg.Add(1)
	}
	go asyncForwardStdin()
	wg.Wait()
}
