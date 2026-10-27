package scriptrunner

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/open-function-computers-llc/server-run-inertia/account"
)

// StreamMessage is one line of script output on its way to the browser.
// Type tells the frontend how to render the line, so scripts don't have to
// fake separation with echoed hyphens or blank lines.
type StreamMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

const (
	MsgStdout  = "stdout"
	MsgMeta    = "meta"
	MsgSection = "section"
	MsgError   = "error"
)

// Scripts can prefix output lines with these to have them rendered as a
// section header or an error, instead of stuffing hyphens into the console.
const (
	sectionLinePrefix = "##section## "
	errorLinePrefix   = "##error## "
)

func message(text, msgType string) StreamMessage {
	return StreamMessage{Message: text, Type: msgType}
}

var redactedEnvKeyParts = []string{"PASSWORD", "SECRET", "TOKEN"}
const redactedValue = "********"

// redactEnvValue hides values for sensitive-looking env keys, so generated
// passwords aren't streamed to every browser watching the script console.
func redactEnvValue(key, value string) string {
	upper := strings.ToUpper(key)
	for _, part := range redactedEnvKeyParts {
		if strings.Contains(upper, part) {
			return redactedValue
		}
	}
	return value
}

// StreamScriptOutput runs a script, sending each line of output down the
// channel, and returns the script's error (if any) once it has exited.
func StreamScriptOutput(scriptName string, args, env map[string]string, communicationChannel *chan StreamMessage) error {
	defer close(*communicationChannel)

	scriptsRoot := os.Getenv("SCRIPTS_ROOT")

	t, err := scriptType(scriptName)
	if err != nil {
		*communicationChannel <- message(err.Error(), MsgError)
		return err
	}

	// validate what the browser sent us before touching the filesystem
	if t == ARGSCRIPT {
		if _, err := argumentScripts[scriptName].orderedArgs(args); err != nil {
			*communicationChannel <- message(err.Error(), MsgError)
			return err
		}
	}
	if t == ENVSCRIPT {
		if err := validateEnv(scriptName, env); err != nil {
			*communicationChannel <- message(err.Error(), MsgError)
			return err
		}
	}

	*communicationChannel <- message("running "+scriptName, MsgMeta)
	*communicationChannel <- message("args:", MsgMeta)
	for key, val := range args {
		*communicationChannel <- message(key+"="+val, MsgMeta)
	}
	*communicationChannel <- message("env:", MsgMeta)
	for key, val := range env {
		*communicationChannel <- message(key+"="+redactEnvValue(key, val), MsgMeta)
	}

	// time to run the script!

	// set up the command runner
	var cmd *exec.Cmd
	if t == ARGSCRIPT {
		argValues, _ := argumentScripts[scriptName].orderedArgs(args)

		// Guard the lock script - don't run if locked-down is disabled for this account
		if scriptName == "lock" {
			accountName := args["ACCOUNT_NAME"]
			if accountName != "" {
				a, err := account.Find(accountName)
				if err == nil && a.IsLockedDownDisabled {
					*communicationChannel <- message("Lock-down is disabled for this account (isLockedDownDisabled is set). The lock script will not run.", MsgMeta)
					return nil
				}
			}
		}

		cmd = exec.Command(filepath.Join(scriptsRoot, argumentScripts[scriptName].filename), argValues...)
	}

	if t == PLAINSCRIPT {
		cmd = exec.Command(filepath.Join(scriptsRoot, plainScripts[scriptName]))
	}

	if t == ENVSCRIPT {
		cmd = exec.Command(filepath.Join(scriptsRoot, envScripts[scriptName].filename))

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
			*communicationChannel <- classifyLine(scanner.Text())
		}
		wg.Done()
	}()

	err = cmd.Start()
	if err != nil {
		*communicationChannel <- message("Error running the command: "+err.Error(), MsgError)
		return err
	}

	wg.Wait()

	// reap the process and report how it exited
	err = cmd.Wait()
	if err != nil {
		*communicationChannel <- message("Script exited with an error: "+err.Error(), MsgError)
	}
	return err
}

// classifyLine turns a raw script output line into a typed message. Scripts
// opt into structure by prefixing lines with ##section## or ##error##.
func classifyLine(line string) StreamMessage {
	if text, ok := strings.CutPrefix(line, sectionLinePrefix); ok {
		return message(text, MsgSection)
	}
	if text, ok := strings.CutPrefix(line, errorLinePrefix); ok {
		return message(text, MsgError)
	}
	return message(line, MsgStdout)
}

func ReturnScriptOutput(scriptName string, args, env map[string]string) ([]byte, error) {
	t, err := scriptType(scriptName)
	if err != nil {
		return nil, err
	}
	scriptsRoot := os.Getenv("SCRIPTS_ROOT")

	// validate what we were given before touching the filesystem
	if t == ARGSCRIPT {
		if _, err := argumentScripts[scriptName].orderedArgs(args); err != nil {
			return nil, err
		}
	}
	if t == ENVSCRIPT {
		if err := validateEnv(scriptName, env); err != nil {
			return nil, err
		}
	}

	// set up the command runner
	var cmd *exec.Cmd
	if t == ARGSCRIPT {
		argValues, _ := argumentScripts[scriptName].orderedArgs(args)
		cmd = exec.Command(filepath.Join(scriptsRoot, argumentScripts[scriptName].filename), argValues...)
	}

	if t == PLAINSCRIPT {
		cmd = exec.Command(filepath.Join(scriptsRoot, plainScripts[scriptName]))
	}

	if t == ENVSCRIPT {
		cmd = exec.Command(filepath.Join(scriptsRoot, envScripts[scriptName].filename))

		for key, val := range env {
			cmd.Env = append(cmd.Env, key+"="+val)
		}
	}

	cmd.Env = append(cmd.Env, "NOCONFIRM=yes")
	cmd.Env = append(cmd.Env, "TERM=xterm-mono")

	return cmd.Output()
}
