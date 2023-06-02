package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"

	"flag"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"
)

type Task struct {
	User        string `yaml:"user"`
	CommandDesc string `yaml:"command_description"`
	Period      string `yaml:"period"`
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

func executeCommand(command string, logger *log.Logger) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing command: %s\n", err)
		return
	}

	logger.Printf("Command output: %s\n", output)
}

func main() {
	// Parse command-line arguments
	flag.Parse()

	// Show help information if requested
	if showHelp {
		flag.Usage()
		return
	}

	// Open the log file
	logOutput := os.Stdout
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logOutput = file
		defer file.Close()
	}

	// Set up the logger
	logger := log.New(logOutput, "", log.LstdFlags)

	logger.Println("Config file location: " + yamlFile)
	logger.Println("Log file location: " + logFile)

	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(yamlFile)
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

		// Retrieve the user information
		u, err := user.Lookup("username")
		if err != nil {
			log.Fatal(err)
		}

		// Convert the UID and GID to integers
		uid := u.Uid
		gid := u.Gid

		uidInt, err := strconv.Atoi(uid)
		if err != nil {
			log.Fatal(err)
		}

		gidInt, err := strconv.Atoi(gid)
		if err != nil {
			log.Fatal(err)
		}

		command := task.Command
		cronExpression := task.Period

		if cronExpression == "@reboot" {
			executeCommand(command, logger)
		} else {

			_, err := c.AddFunc(cronExpression, func() {
				executeCommand(command, logger)
			})
			if err != nil {
				log.Printf("Failed to schedule task: %s", err)
			}
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
