package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"flag"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Task struct {
	User        string `yaml:"user"`
	CommandDesc string `yaml:"command_description"`
	Period      string `yaml:"period"`
	Logging     string `yaml: log`
	Command     string `yaml:"command"`
}

var (
	showHelp bool
	yamlFile string
	logFile  string
	// Define other flags here
)

func init() {
	// Define command-line flags
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.StringVar(&yamlFile, "config", "/etc/bettercron/config.yaml", "Define yaml file location")
	flag.StringVar(&logFile, "log", "/var/log/bettercron.log", "Define log file location")
	// Define other flags here
}

func executeCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command: %s\n", err)
		return
	}

	fmt.Printf("Command output: %s\n", output)
}

func main() {
	// Parse command-line arguments
	flag.Parse()

	// Show help information if requested
	if showHelp {
		flag.Usage()
		return
	}

	fmt.Println(yamlFile)

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile("file.yml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Parse YAML into tasks slice
	var tasks []Task
	err = yaml.Unmarshal(yamlFile, &tasks)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Create a new cron scheduler
	c := cron.New()

	// Schedule tasks
	for _, task := range tasks {
		command := task.Command
		cronExpression := task.Period
		_, err := c.AddFunc(cronExpression, func() {
			executeCommand(command)
		})
		if err != nil {
			log.Printf("Failed to schedule task: %s", err)
		}
	}

	// Start the scheduler
	c.Start()

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Stop the scheduler
	c.Stop()
}
