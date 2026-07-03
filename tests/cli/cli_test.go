package cli_test

import (
	"testing"

	"github.com/marcuwynu23/git-share/internal/cli"
)

func TestParseOnCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on"})
	if opts.Command != cli.On {
		t.Errorf("Command = %v, want On", opts.Command)
	}
}

func TestParseOffCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "off"})
	if opts.Command != cli.Off {
		t.Errorf("Command = %v, want Off", opts.Command)
	}
}

func TestParseStatusCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "status"})
	if opts.Command != cli.Status {
		t.Errorf("Command = %v, want Status", opts.Command)
	}
}

func TestParseVersionCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "version"})
	if opts.Command != cli.VersionCmd {
		t.Errorf("Command = %v, want VersionCmd", opts.Command)
	}
}

func TestParseConfigCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "config"})
	if opts.Command != cli.ConfigCmd {
		t.Errorf("Command = %v, want ConfigCmd", opts.Command)
	}
}

func TestParseServeCommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "_serve"})
	if opts.Command != cli.ServeCmd {
		t.Errorf("Command = %v, want ServeCmd", opts.Command)
	}
}

func TestParsePortFlag(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--port", "9090"})
	if opts.Port != 9090 {
		t.Errorf("Port = %d, want 9090", opts.Port)
	}
}

func TestParseReadOnlyFlag(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--readonly"})
	if !opts.ReadOnly {
		t.Error("ReadOnly should be true")
	}
}

func TestParseReadWriteFlag(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--readwrite"})
	if !opts.ReadWrite {
		t.Error("ReadWrite should be true")
	}
}

func TestParseHostnameFlag(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--hostname", "0.0.0.0"})
	if opts.Hostname != "0.0.0.0" {
		t.Errorf("Hostname = %q, want 0.0.0.0", opts.Hostname)
	}
}

func TestParseDaemonFlag(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--daemon"})
	if !opts.Daemon {
		t.Error("Daemon should be true")
	}
}

func TestParseMultipleFlags(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "--port", "9000", "--readwrite", "--hostname", "192.168.1.1"})
	if opts.Port != 9000 {
		t.Errorf("Port = %d, want 9000", opts.Port)
	}
	if !opts.ReadWrite {
		t.Error("ReadWrite should be true")
	}
	if opts.Hostname != "192.168.1.1" {
		t.Errorf("Hostname = %q, want 192.168.1.1", opts.Hostname)
	}
}

func TestParseOnCommandOffSubcommand(t *testing.T) {
	opts := cli.Parse([]string{"git-share", "on", "off"})
	if opts.Command != cli.Off {
		t.Errorf("Command = %v, want Off", opts.Command)
	}
}

func TestAppVersion(t *testing.T) {
	if cli.AppVersion == "" {
		t.Error("AppVersion should not be empty")
	}
}

func TestRunVersion(t *testing.T) {
	opts := &cli.Options{Command: cli.VersionCmd}
	cli.Run(opts)
}

func TestRunStatus(t *testing.T) {
	opts := &cli.Options{Command: cli.Status}
	cli.Run(opts)
}

func TestRunOff(t *testing.T) {
	opts := &cli.Options{Command: cli.Off}
	cli.Run(opts)
}

func TestPrintUsage(t *testing.T) {
	cli.PrintUsage()
}
