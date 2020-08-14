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
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type structServerConfig struct {
	RootFolder string `yaml:"rootFolder"`
	ScriptFile string `yaml:"scriptFile"`
}
type structMCSHConfig struct {
	Servers map[string]structServerConfig `yaml:"servers"`
}

var mcshConfig = structMCSHConfig{
	Servers: map[string]structServerConfig{
		"serverName1": structServerConfig{
			RootFolder: "path/to/your/server/root/folder",
			ScriptFile: "your/script/file/name/in/the/server/root/folder",
		},
	},
}

func data2yaml(data structMCSHConfig) []byte {
	yaml, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(yaml)
	return yaml
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
	mcshConfig = structMCSHConfig{}
	err = yaml.Unmarshal(configYaml, &mcshConfig)
	// fmt.Printf("mcshConfig:\n%v\n", mcshConfig)
}

func runServer(name string, serverConfig structServerConfig, writeClosers map[string]io.WriteCloser) error {
	if err := os.Chdir(serverConfig.RootFolder); err != nil {
		if os.IsNotExist(err) {
			log.Printf("server<%s>: Cannot find server root folder, please check your \"config.yml\"", name)
			wg.Done()
			return err
		}
		log.Printf("server<%s>: Error when chdir to server root folder - %s", name, err.Error())
		wg.Done()
		return err
	}
	cmd := exec.Command(path.Join(wd, "Scripts", serverConfig.ScriptFile))

	writeClosers[name], _ = cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("server<%s>: Error startingup: %s......", name, err.Error())
		wg.Done()
		return err
	}
	go asyncLog(name, stdout)
	go asyncLog(name, stderr)

	if err := cmd.Wait(); err != nil {
		log.Printf("server<%s>: Error running: %s......", name, err.Error())
		wg.Done()
		return err
	}
	wg.Done()
	return nil
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
				_, valid := writeClosers[string(res[1])]
				if valid {
					_, errWrite := writeClosers[string(res[1])].Write(append(res[2], '\n'))
					if errWrite != nil {
						log.Println("MCSH[stdinForward/ERROR]: Server stdin write failed - ", errWrite)
					}
				} else {
					log.Printf("MCSH[stdinForward/ERROR]: Cannot find running server <%v>\n", string(res[1]))
				}
			}
		}
	}
}

var wg sync.WaitGroup

var writeClosers = make(map[string]io.WriteCloser)
var wd = ""

func init() {
	os.Mkdir("Scripts", 0666)
	wd, _ = os.Getwd()
	readConfig()
}

func main() {
	// readConfig()
	for name, serverConfig := range mcshConfig.Servers {
		// fmt.Printf("server<%v> \n\troot:%v\n\tscript:%v\n", name, serverConfig.RootFolder, serverConfig.ScriptFile)
		go runServer(name, serverConfig, writeClosers)
		wg.Add(1)
	}
	go asyncForwardStdin()
	wg.Wait()
}
