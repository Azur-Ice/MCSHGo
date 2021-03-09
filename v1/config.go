package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServerConfig ...
type ServerConfig struct {
	ExecOptions string `yaml:"execOptions"`
	ExecPath    string `yaml:"execPath"`
}

// Config ...
type Config struct {
	CommandPrefix string                  `yaml:"command_prefix"`
	BackupDir     string                  `yaml:"backup_dir"`
	Servers       map[string]ServerConfig `yaml:"servers"`
}

//MCSHConfig ...
var MCSHConfig Config

// ReadConfig ...
func ReadConfig(config *Config) {
	configYaml, err := ioutil.ReadFile("./config.yml")
	if err != nil { // 读取文件发生错误
		if os.IsNotExist(err) { // 文件不存在，创建并写入默认配置
			log.Println("MCSH: Cannot find config.yml, creating...")
			createDefaultConfigFile()
			log.Println("MCSH: Successful created config.yml, please complete the config.")
		}
		os.Exit(1)
	}
	err = yaml.Unmarshal(configYaml, config)
}

func init() {
	log.Println("MCSH[init/INFO]: Initializing config...")
	ReadConfig(&MCSHConfig)
	backupDir, _ = filepath.Abs(MCSHConfig.BackupDir)
	os.Mkdir(backupDir, 0666)
}

func createDefaultConfigFile() {
	// str := Data2yaml(mcshConfig)
	defaultCofig, _ := yaml.Marshal(Config{
		CommandPrefix: "#",
		BackupDir:     "./Backups",
		Servers: map[string]ServerConfig{
			"serverName1": {
				ExecOptions: "-Xms4G -Xmx4G",
				ExecPath:    "path/to/your/server/s/exec/jar/file",
			},
		},
	})
	ioutil.WriteFile("./config.yml", defaultCofig, 0666)
}
