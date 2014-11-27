package assets

import (
	"github.com/go-floki/floki"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
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
		} else {
			// if process is alive - wait for it to finish
			process.Wait()
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

	var cmd *exec.Cmd

	if _, err := os.Stat(javascriptDir + "/Gulpfile.js"); os.IsNotExist(err) {
		log.Println("running `npm start`")
		cmd = exec.Command("npm", "start")
	} else {
		log.Println("running `gulp watchify`")
		cmd = exec.Command("gulp", "watchify")
	}

	cmd.Dir = javascriptDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	//var out bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(cmd.Process.Pid), 10)), 0644)
}

func init() {
	floki.RegisterAppEventHandler("ConfigureAppEnd", func(f *floki.Floki) {
		compileAssets := f.Config.Bool("compileAssets", true)
		if compileAssets && floki.Env == floki.Dev {

			go func() {
				for {
					runWatchify()
					log.Println("Watchify crashed. Sleeping 5sec and restarting..")
					time.Sleep(5 * time.Second)
				}
			}()

			go func() {
				for {
					runStylus()
					log.Println("Stylus crashed. Sleeping 5sec and restarting..")
					time.Sleep(5 * time.Second)
				}
			}()

		}
	})
}
