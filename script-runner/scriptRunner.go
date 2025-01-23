package scriptrunner

import (
	"bufio"
	"os"
	"os/exec"
	"sync"
)

func StreamScriptOutput(scriptName string, args, env map[string]string, communicationChannel *chan string) {
	defer close(*communicationChannel)

	scriptsRoot := os.Getenv("SCRIPTS_ROOT")

	*communicationChannel <- "running " + scriptName
	*communicationChannel <- "args:"
	for key, val := range args {
		*communicationChannel <- key + "=" + val
	}
	*communicationChannel <- "-"
	*communicationChannel <- "env:"
	for key, val := range env {
		*communicationChannel <- key + "=" + val
	}
	*communicationChannel <- "-"
	*communicationChannel <- "-"

	t, err := scriptType(scriptName)

	if err != nil {
		*communicationChannel <- err.Error()
		close(*communicationChannel)
		return
	}

	// time to run the script!

	// set up the command runner
	var cmd *exec.Cmd
	if t == ARGSCRIPT {
		// TODO: add the arguments here correctly
		cmd = exec.Command(scriptsRoot + argumentScripts[scriptName])
	}

	if t == PLAINSCRIPT {
		cmd = exec.Command(scriptsRoot + plainScripts[scriptName])
	}

	if t == ENVSCRIPT {
		cmd = exec.Command(scriptsRoot + envScripts[scriptName])

		for key, val := range env {
			cmd.Env = append(cmd.Env, key+"="+val)
		}
	}

	// set up our communication pipes
	outPipe, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	cmd.Env = append(cmd.Env, "NOCONFIRM=yes")
	cmd.Env = append(cmd.Env, "TERM=xterm-mono")
	scanner := bufio.NewScanner(outPipe)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			*communicationChannel <- line
		}
		wg.Done()
	}()

	err = cmd.Start()
	if err != nil {
		*communicationChannel <- "Error running the command: " + err.Error()
		return
	}

	wg.Wait()
}
