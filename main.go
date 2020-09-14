package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fortis/go-steam-totp"
)

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func attemptFindSteamcmd() string {
	windowsPath := "C:\\steamcmd\\steamcmd.exe"
	unixPath := "/home/steam/steamcmd"
	dockerPath := "/home/steam/steamcmd.sh"

	if runtime.GOOS == "windows" {
		return windowsPath
	} else {
		if fileExist(dockerPath) {
			return dockerPath
		} else {
			return unixPath
		}
	}
}

func parseArgs() (string, string, string, string, string) {
	steamcmdPath := flag.String("path", attemptFindSteamcmd(), "Path to steamcmd executable")
	steamcmdArgs := flag.String("args", "", "Arguments to pass to steamcmd")
	steamcmdUser := flag.String("username", "", "Username to log in with")
	steamcmdPass := flag.String("password", "", "Password to log in with")
	tfaSeed := flag.String("seed", "", "The 2FA seed/secret")
	flag.Parse()

	if !fileExist(*steamcmdPath) {
		fmt.Println("Provided steamcmd path is invalid.")
		os.Exit(1)
	}

	if *steamcmdUser == "" || *steamcmdPass == "" {
		fmt.Println("Username or password not provided.")
		os.Exit(1)
	}

	if *tfaSeed == "" {
		fmt.Println("No 2FA seed provided.")
		os.Exit(1)
	}

	return *steamcmdPath, *steamcmdArgs, *steamcmdUser, *steamcmdPass, *tfaSeed
}

func main() {
	path, args, user, pass, seed := parseArgs()

	fmt.Println(path)

	code, err := steam_totp.GenerateAuthCode(seed, time.Now())
	if err != nil {
		fmt.Println("Error while generating code: ", err)
		os.Exit(1)
	}

	commandArgs := fmt.Sprintf("+login %s %s %s %s", user, pass, code, args)
	fmt.Println(fmt.Sprintf("steamcmd-2fa :: debug :: %s %s", path, strings.Replace(strings.Replace(strings.Replace(commandArgs, code, "*****", -1), pass, "*****", -1), user, "*****", -1)))
	cmd := exec.Command(path, commandArgs)

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error while creating StdoutPipe: ", err)
		os.Exit(1)
	}
	defer cmdOut.Close()

	scanner := bufio.NewScanner(cmdOut)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running command: ", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for command: ", err)
		return
	}
}
