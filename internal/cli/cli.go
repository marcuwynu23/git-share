package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/markwayne/git-share/internal/config"
	"github.com/markwayne/git-share/internal/git"
	"github.com/markwayne/git-share/internal/server"
	"github.com/markwayne/git-share/internal/util"
)

const AppVersion = "0.1.0"

type Command string

const (
	On          Command = "on"
	Off         Command = "off"
	Status      Command = "status"
	VersionCmd  Command = "version"
	ConfigCmd   Command = "config"
)

type Options struct {
	Command   Command
	Port      int
	ReadOnly  bool
	ReadWrite bool
	Hostname  string
	Help      bool
}

func Parse(args []string) *Options {
	opts := &Options{}

	if len(args) < 2 {
		PrintUsage()
		os.Exit(0)
	}

	subcommand := args[1]

	fs := flag.NewFlagSet("git-share", flag.ContinueOnError)
	fs.Usage = func() { PrintUsage() }

	fs.IntVar(&opts.Port, "port", 0, "port to listen on")
	fs.BoolVar(&opts.ReadOnly, "readonly", false, "enable read-only mode (clone only)")
	fs.BoolVar(&opts.ReadWrite, "readwrite", false, "enable read-write mode (allow push)")
	fs.StringVar(&opts.Hostname, "hostname", "", "hostname or IP to bind to")
	fs.BoolVar(&opts.Help, "help", false, "show help")

	remaining := []string{}
	for i, arg := range args[2:] {
		if !strings.HasPrefix(arg, "-") && i == 0 && Command(arg) == Off {
			subcommand = "off"
			continue
		}
		remaining = append(remaining, arg)
	}

	fs.Parse(remaining)

	switch Command(subcommand) {
	case On:
		opts.Command = On
	case Off:
		opts.Command = Off
	case Status:
		opts.Command = Status
	case VersionCmd:
		opts.Command = VersionCmd
	case ConfigCmd:
		opts.Command = ConfigCmd
	default:
		if subcommand == "--help" || subcommand == "-h" {
			PrintUsage()
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", subcommand)
		PrintUsage()
		os.Exit(1)
	}

	return opts
}

func PrintUsage() {
	fmt.Println(`Usage: git share <command> [options]

Commands:
  on         Start sharing the current repository
  off        Stop sharing
  status     Show sharing status
  version    Show version
  config     Configure git-share

Options:
  --port <port>       Port to listen on (default: 8080)
  --readonly          Enable read-only mode (clone only, no push)
  --readwrite         Enable read-write mode (allow push)
  --hostname <addr>   Hostname or IP to bind to
  --help              Show this help message

Examples:
  git share on
  git share on --port 9000
  git share on --readwrite
  git share on --readonly
  git share status
  git share version`)
}

func Run(opts *Options) {
	switch opts.Command {
	case On:
		runOn(opts)
	case Off:
		runOff()
	case Status:
		runStatus()
	case VersionCmd:
		runVersion()
	case ConfigCmd:
		runConfig()
	}
}

func runOn(opts *Options) {
	repo, err := git.FindRepository()
	if err != nil {
		util.Fatal("Error: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		util.Fatal("Error loading config: %v", err)
	}

	serverCfg := &server.ServerConfig{
		Port:     cfg.Port,
		ReadOnly: cfg.ReadOnly,
		Hostname: cfg.Hostname,
	}

	if opts.Port > 0 {
		serverCfg.Port = opts.Port
	}
	if opts.ReadWrite {
		serverCfg.ReadOnly = false
	} else if opts.ReadOnly {
		serverCfg.ReadOnly = true
	}
	if opts.Hostname != "" {
		serverCfg.Hostname = opts.Hostname
	}
	if cfg.Timeout > 0 {
		serverCfg.Timeout = time.Duration(cfg.Timeout) * time.Second
	}

	srv := server.New(repo, serverCfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		util.WaitForSignal()
		cancel()
	}()

	if err := srv.Start(ctx); err != nil {
		util.Fatal("Server error: %v", err)
	}
}

func runOff() {
	fmt.Println("Stop command received")
	fmt.Println("Note: git-share does not currently run as a background process.")
	fmt.Println("Press Ctrl+C in the terminal where git-share is running.")
}

func runStatus() {
	repo, err := git.FindRepository()
	if err != nil {
		fmt.Println("Not inside a Git repository")
		return
	}
	fmt.Printf("Repository: %s\n", repo.Root)
	fmt.Printf("Branch:     %s\n", repo.Branch)
}

func runVersion() {
	fmt.Printf("git-share version %s\n", AppVersion)
}

func runConfig() {
	cfg, err := config.Load()
	if err != nil {
		util.Fatal("Error loading config: %v", err)
	}

	if len(os.Args) < 4 {
		fmt.Printf("Current config:\n")
		fmt.Printf("  port:     %d\n", cfg.Port)
		fmt.Printf("  readonly: %v\n", cfg.ReadOnly)
		fmt.Printf("  hostname: %s\n", cfg.Hostname)
		fmt.Printf("  timeout:  %d\n", cfg.Timeout)
		fmt.Printf("  theme:    %s\n", cfg.Theme)
		fmt.Printf("\nUsage: git share config <key> <value>\n")
		return
	}

	key := os.Args[3]
	value := os.Args[4]

	switch key {
	case "port":
		p, err := strconv.Atoi(value)
		if err != nil {
			util.Fatal("Invalid port: %s", value)
		}
		cfg.Port = p
	case "readonly":
		cfg.ReadOnly = value == "true" || value == "yes" || value == "1"
	case "hostname":
		cfg.Hostname = value
	case "timeout":
		t, err := strconv.Atoi(value)
		if err != nil {
			util.Fatal("Invalid timeout: %s", value)
		}
		cfg.Timeout = t
	case "theme":
		cfg.Theme = value
	default:
		util.Fatal("Unknown config key: %s", key)
	}

	if err := config.Save(cfg); err != nil {
		util.Fatal("Error saving config: %v", err)
	}
	fmt.Printf("Config updated: %s = %s\n", key, value)
}
