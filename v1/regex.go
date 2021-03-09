package main

import (
	"log"
	"regexp"
)

var forwardReg *regexp.Regexp
var commandReg *regexp.Regexp
var playerOutputReg *regexp.Regexp
var outputFormatReg *regexp.Regexp

func init() {
	log.Println("MCSH[init/INFO]: Initializing regexps...")
	forwardReg = regexp.MustCompile(`(.+?) *\| *(.+)`)
	commandReg = regexp.MustCompile("^" + string(MCSHConfig.CommandPrefix) + "(.*)")
	playerOutputReg = regexp.MustCompile(`\]: <(.*?)> (.*)`)
	outputFormatReg = regexp.MustCompile(`(\[\d\d:\d\d:\d\d\]) *\[.+?\/(.+?)\]`)
}
