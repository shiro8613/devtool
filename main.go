package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "./devtool.yml"
const (
	TYPE_SYNC = "sync"
	TYPE_ASYNC = "async"
)

type Config struct {
	Scripts ScriptsType `yaml:"scripts"`
}

type ScriptsType map[string]ScriptType

type ScriptType struct {
	Type string `yaml:"type"`
	Command map[string]string `yaml:"command"`
}

var nowProcess map[string]*exec.Cmd = make(map[string]*exec.Cmd)

func ReadConfig() Config {
	b, err := os.ReadFile(CONFIG_PATH);
	if err != nil {
		log.Fatalln("devtool.yml does not exist.")
	}

	var c Config

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatalln("Failed to parse devtool.yml.")
	}

	if !checkConfig(c) {
		log.Fatalln("Failed to parse devtool.yml ")
	}

	return c

}

func checkConfig(config Config) bool {
	scripts := config.Scripts
	for k, v := range scripts {
		if k == "" {
			return false
		}

		if len(v.Command) == 0 {
			return false
		}

		if (!slices.Contains([]string{TYPE_SYNC, TYPE_ASYNC}, v.Type)) {
			return false
		}		

		if (len(v.Command) == 1 && v.Type == TYPE_ASYNC) {
			return false
		}

		for k, cs := range v.Command {
			if k == "" {
				return false
			}
			
			if cs == "" {
				return false
			}
		}
	}

	return true

}

func WriteConfig() {
	c := Config{
		Scripts: ScriptsType{
			"test": ScriptType{
				Type: "sync",
				Command: map[string]string{
					"echo_aaa": "echo aaa",
				},
			},
		}, 
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalln("Internal error")
	}

	err = os.WriteFile(CONFIG_PATH, b, os.ModeAppend)
	if err != nil {
		log.Fatalln("Failed to create file.")
	}
}

func signalCatcher() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL)
	
	<- ch
	
	for k,v := range nowProcess {
		v.Process.Signal(syscall.SIGTERM)
		fmt.Printf("%s is sttoped\n", k)
	}
}

func cmdRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("devtool [init, run]")
		return nil
	}

	if !slices.Contains([]string{"init", "run"}, args[0]) {
		fmt.Printf("'%s' is a non-existent command.\n", args[0])
		return nil
	}

	go signalCatcher()

	switch (args[0]) {
	case "init":
		WriteConfig()
		fmt.Println("File created.")
		return nil
	
	case "run":
		c := ReadConfig()
		a := args[1]

		if a == "" {
			fmt.Println("More args please")
			return nil
		}

		if _, ok := c.Scripts[a]; !ok {
			fmt.Printf("'%s' is a non-existent command in devtool.yml.\n", a)
		}

		cmd := c.Scripts[a]
		commandRunner(cmd)
		return nil
	}


	return nil
}

func commandRunner(cmd ScriptType) {
	switch cmd.Type {
	case TYPE_SYNC:
		for k, c := range cmd.Command {
			processOne(k, c)
		}
	case TYPE_ASYNC:
		var wg sync.WaitGroup
		for k, c := range cmd.Command {
			k := k
			c := c
			wg.Add(1)
			go func () {
				processOne(k, c)
				wg.Done()
			}()
		}
		wg.Wait()
		fmt.Println("Done")
	}
}

func processOne(name string, cmdString string) {
    osname := runtime.GOOS
    var cmd *exec.Cmd
    if osname == "windows" {
        shell := os.Getenv("COMSPEC")
        cmd = exec.Command(shell, "/c", cmdString)
    } else {
        shell := "/bin/sh"
        cmd = exec.Command(shell, "-c", cmdString)
    }

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
        log.Fatalln(err)
    }

	nowProcess[name] = cmd

	go func () {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("[%s] %s\n", name, scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
        log.Fatalln(err)
    }

	delete(nowProcess, name)
}

func main() {
	rootCmd := &cobra.Command{
		Use: "devtool [init, run]",
		RunE: cmdRun,
	}
	
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln("Internal error")
	}
}
