package assets

import (
	"github.com/go-floki/floki"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func checkPidAlive(pidFile string) bool {
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return false
	}

	pidStr, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Fatal(err)
	}

	pid, err := strconv.ParseInt(string(pidStr), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return false
		}
	}

	return true
}

func runWatchify() {
	javascriptDir := "../app/javascripts"

	pidFile := "watchify.pid"

	if checkPidAlive(pidFile) {
		return
	}

	log.Println("running npm start in", javascriptDir)

	cmd := exec.Command("npm", "start")
	cmd.Dir = javascriptDir

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(cmd.Process.Pid), 10)), 0644)
}

func runStylus() {
	cssDir := "../app/stylesheets"

	pidFile := "stylus.pid"

	if checkPidAlive(pidFile) {
		return
	}

	log.Println("running stylus in", cssDir)

	cmd := exec.Command("stylus", "-u", "./node_modules/nib", "-u", "./node_modules/bootstrap3-stylus", "-w", "-o", "../../static/")
	cmd.Dir = cssDir

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(cmd.Process.Pid), 10)), 0644)
}

func init() {
	if floki.Env == floki.Dev {
		runWatchify()
		runStylus()
	}
}
