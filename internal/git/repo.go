package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrNotARepo     = errors.New("not inside a Git repository")
	ErrGitNotFound  = errors.New("git executable not found")
	ErrNotOnBranch  = errors.New("not on a branch (detached HEAD)")
)

type RepoInfo struct {
	Root       string
	Worktree   string
	Branch     string
	IsBare     bool
	GitDir     string
}

func CheckGit() error {
	_, err := exec.LookPath("git")
	return err
}

func FindRepository() (*RepoInfo, error) {
	if err := CheckGit(); err != nil {
		return nil, ErrGitNotFound
	}

	root, err := execGit("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, ErrNotARepo
	}
	root = strings.TrimSpace(root)

	gitDir, err := execGit("rev-parse", "--git-dir")
	if err != nil {
		return nil, ErrNotARepo
	}
	gitDir = strings.TrimSpace(gitDir)

	bare, _ := execGit("rev-parse", "--is-bare")
	isBare := strings.TrimSpace(bare) == "true"

	branch, err := execGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, ErrNotOnBranch
	}
	branch = strings.TrimSpace(branch)

	worktree := filepath.Dir(gitDir)
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	if !filepath.IsAbs(worktree) {
		worktree = root
	}

	return &RepoInfo{
		Root:     root,
		Worktree: worktree,
		Branch:   branch,
		IsBare:   isBare,
		GitDir:   gitDir,
	}, nil
}

func Validate(repo *RepoInfo) error {
	info, err := os.Stat(repo.GitDir)
	if err != nil {
		return ErrNotARepo
	}
	if !info.IsDir() {
		return ErrNotARepo
	}
	return nil
}

func execGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
