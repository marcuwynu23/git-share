package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/marcuwynu23/git-share/internal/config"
	"github.com/marcuwynu23/git-share/internal/git"
	"github.com/marcuwynu23/git-share/internal/server"
	"github.com/marcuwynu23/git-share/internal/util"
)

func pidPath() string {
	return filepath.Join(os.TempDir(), ".pstemp")
}

func writePID() error {
	return os.WriteFile(pidPath(), []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

func readPID() (int, error) {
	data, err := os.ReadFile(pidPath())
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func removePID() {
	os.Remove(pidPath())
}

const AppVersion = "0.1.0"

type Command string

const (
	On          Command = "on"
	Off         Command = "off"
	Status      Command = "status"
	VersionCmd  Command = "version"
	ConfigCmd   Command = "config"
	ServeCmd    Command = "_serve"
)

type Options struct {
	Command   Command
	Port      int
	ReadOnly  bool
	ReadWrite bool
	Hostname  string
	Daemon    bool
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
	fs.BoolVar(&opts.Daemon, "daemon", false, "run in background as a daemon")
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
	case ServeCmd:
		opts.Command = ServeCmd
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
  --port <port>       Port to listen on (default: 9720)
  --readonly          Enable read-only mode (clone only, no push)
  --readwrite         Enable read-write mode (allow push)
  --hostname <addr>   Hostname or IP to bind to
  --daemon            Run in background (manage with 'off' command)
  --help              Show this help message

Examples:
  git share on
  git share on --port 9000
  git share on --readwrite
  git share on --readonly
  git share on --daemon
  git share off
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
	case ServeCmd:
		runServe(opts)
	}
}

func runOn(opts *Options) {
	if opts.Daemon {
		pid, err := readPID()
		if err == nil {
			util.Fatal("Already running (PID %d)", pid)
		}

		args := []string{"git-share", "_serve"}
		if opts.Port > 0 {
			args = append(args, "--port", strconv.Itoa(opts.Port))
		}
		if opts.ReadWrite {
			args = append(args, "--readwrite")
		} else if opts.ReadOnly {
			args = append(args, "--readonly")
		}
		if opts.Hostname != "" {
			args = append(args, "--hostname", opts.Hostname)
		}

		cmd := exec.Command(os.Args[0], args[1:]...)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil

		if err := cmd.Start(); err != nil {
			util.Fatal("Failed to start daemon: %v", err)
		}

		os.WriteFile(pidPath(), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
		fmt.Printf("Daemon started with PID %d\n", cmd.Process.Pid)
		return
	}

	runServe(opts)
}

func runOff() {
	pid, err := readPID()
	if err != nil {
		fmt.Println("No daemon process found.")
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		removePID()
		fmt.Println("No daemon process found.")
		return
	}

	if err := proc.Kill(); err != nil {
		util.Fatal("Failed to stop daemon: %v", err)
	}

	removePID()
	fmt.Printf("Daemon with PID %d stopped.\n", pid)
}

func runServe(opts *Options) {
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

	writePID()
	defer removePID()

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
