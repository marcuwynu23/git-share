package git

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleSmartHTTP(info *RepoInfo, w http.ResponseWriter, r *http.Request, repoName string, readOnly bool) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		http.Error(w, "git executable not found", http.StatusInternalServerError)
		return
	}

	gitDir := info.GitDir

	args := []string{"http-backend"}
	cmd := exec.CommandContext(r.Context(), gitPath, args...)
	cmd.Dir = gitDir

	cmd.Env = buildCGIEnv(r, gitDir, repoName, info.IsBare, readOnly)

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/x-git-") && r.Body != nil {
		cmd.Stdin = r.Body
	}
	if r.Body != nil {
		defer r.Body.Close()
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "failed to create stdout pipe", http.StatusInternalServerError)
		return
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, "failed to start git backend", http.StatusInternalServerError)
		return
	}

	reader := bufio.NewReader(stdout)
	statusCode := 200
	headersDone := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			headersDone = true
			break
		}

		if strings.HasPrefix(line, "Status: ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				code, err := strconv.Atoi(parts[1])
				if err == nil {
					statusCode = code
				}
			}
		} else if strings.Contains(line, ": ") {
			parts := strings.SplitN(line, ": ", 2)
			if len(parts) == 2 {
				w.Header().Add(parts[0], parts[1])
			}
		}
	}

	if !headersDone {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("git backend returned no headers"))
		cmd.Wait()
		return
	}

	w.WriteHeader(statusCode)
	io.Copy(w, reader)
	cmd.Wait()
}

func buildCGIEnv(r *http.Request, gitDir, repoName string, isBare bool, readOnly bool) []string {
	env := os.Environ()

	path := r.URL.Path
	path = strings.TrimPrefix(path, "/git")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	suffix := path
	repoPrefix := "/" + repoName + ".git"
	if strings.HasPrefix(path, repoPrefix) {
		suffix = path[len(repoPrefix):]
		if suffix == "" {
			suffix = "/"
		}
	}

	projectRoot := filepath.Dir(gitDir)

	pathInfo := suffix
	if !strings.HasPrefix(pathInfo, "/") {
		pathInfo = "/" + pathInfo
	}

	if isBare {
		pathInfo = "/" + filepath.Base(gitDir) + pathInfo
	}

	env = append(env,
		"GIT_PROJECT_ROOT="+projectRoot,
		"GIT_HTTP_EXPORT_ALL=1",
		"PATH_INFO="+pathInfo,
		"REQUEST_METHOD="+r.Method,
		"QUERY_STRING="+r.URL.RawQuery,
		"CONTENT_TYPE="+r.Header.Get("Content-Type"),
		"CONTENT_LENGTH="+r.Header.Get("Content-Length"),
		"SCRIPT_NAME=/" + repoName + ".git",
		"REQUEST_URI="+r.URL.RequestURI(),
		"SERVER_PROTOCOL="+r.Proto,
		"SERVER_SOFTWARE=git-share/0.1.0",
		"SERVER_NAME=git-share",
		"GIT_TERMINAL_PROMPT=0",
	)

	if !readOnly {
		env = append(env,
			"GIT_CONFIG_COUNT=1",
			"GIT_CONFIG_KEY_0=http.receivepack",
			"GIT_CONFIG_VALUE_0=true",
		)
	}

	return env
}
