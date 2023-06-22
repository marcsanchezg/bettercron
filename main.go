package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"syscall"
	"time"

	"flag"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Task struct {
	User        string `yaml:"user"`
	CommandDesc string `yaml:"command_description"`
	Period      string `yaml:"period"`
	Command     string `yaml:"command"`
}

var (
	// Define flags here
	showHelp bool
	yamlFile string
	logFile  string

	// Define prometheus metrics here
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func init() {
	// Define command-line flags
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.StringVar(&yamlFile, "config", "/etc/bettercron/config.yaml", "Define yaml file location")
	flag.StringVar(&logFile, "log", "/var/log/bettercron.log", "Define log file location")
	// Define other flags here
}

func executeCommand(command string, logger *log.Logger, task Task) {
	var cmd *exec.Cmd
	if task.User == "" {
		cmd = exec.Command("bash", "-c", command)
	} else {
		logger.Println()
		cmd = exec.Command("sudo", "-u", task.User, "bash", "-c", command)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error executing command: %s\n", err)
		return
	}

	logger.Printf("\nCommand description "+task.CommandDesc+"\nRunning as "+task.User+"\nCommand output: \n%s\n", output)
}

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
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

	logger.Println("------Starting bettercron------")
	logger.Println("Config file location: " + yamlFile)
	logger.Println("Log file location: " + logFile + "\n")

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

	/* 	recordMetrics()
	   	http.Handle("/metrics", promhttp.Handler())
	   	http.ListenAndServe(":2112", nil) */

	// Create a new cron scheduler
	c := cron.New()

	// Schedule tasks
	for _, task := range tasks {

		// Retrieve the user information
		_, err := user.Lookup(task.User)
		if err != nil {
			log.Print(err)
		} else {
			command := task.Command

			switch cronExpression := task.Period; cronExpression {
			case "@reboot":
				go executeCommand(command, logger, task)
			default:
				_, err := c.AddFunc(cronExpression, func() {
					go executeCommand(command, logger, task)
				})
				if err != nil {
					log.Printf("Failed to schedule task: %s", err)
				}
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
