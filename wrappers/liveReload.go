package main

/*
	liveReload.go
	Author: Escher Wright-Dykhouse

	The liveReload wrapper is a program that runs above server-run and is responsible for:
	- watching for file changes in backend (.go) folders
	- running 'npx mix watch' to watch for frontend changes (/frontend)
	- rebuilding and running the program (./build.sh) on file changes
	- broadcasting a websocket message to the frontend to reload the page
*/

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// file watching will ignore these directories
var autoRebuildIgnoreDirs = []string{
	"node_modules",
	"frontend",
	"wrappers",
	"cmd",
	"dist",
	".git",
}

// file watching will watch these file extensions
var autoRebuildFileExtensions = []string{
	".go",
	".env",
}

// output prefixes
const outputPrefix = "[Live Reload] "
const outputPrefixLaravelMix = "[Laravel Mix] "
const outputPrefixMain = "[Main] "

// messages sent to this channel will trigger a rebuild
const frontendRebuild = 0
const backendRebuild = 1

var rebuildChannel = make(chan int)

// messages sent to this channel will trigger a frontend page reload
var frontendReloadChannel = make(chan bool)

// map of active websocket connections for live reload feature
var frontendReloadClients = make(map[*websocket.Conn]bool)

// define the websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// let all connections through; we're generous like that
		return true
	},
}

// define the main function
func main() {

	// check global env file
	_, err := os.Stat("/etc/server-run.env")
	if err == nil {
		// global env file exists! let's load it
		godotenv.Load("/etc/server-run.env")
	}

	// check for local env file (will take precedence over global env file)
	_, err = os.Stat(".env")
	if err == nil {
		// local env file exists! let's load it
		godotenv.Load(".env")
	}

	// get the port from the environment variable
	liveReloadPort := os.Getenv("APP_LIVERELOAD_PORT")
	if liveReloadPort == "" {
		liveReloadPort = "8981" // default port
	}

	fmt.Println(outputPrefix + "Starting live reload websocket server at :" + liveReloadPort)

	// create a new http server
	server := &http.Server{
		Addr:    ":" + liveReloadPort,
		Handler: nil,
	}

	// setup the frontend reload handler
	http.HandleFunc("/frontend-reload", handleFrontendReload)

	// compile the directories with watched files
	go watchDirectories(compileDirectoriesWithWatchedFiles())

	// start the laravel mix watcher
	// this will block until the initial build has completed
	laravelMixWatcher()

	// start the frontend reload broadcaster
	go startFrontendReloadBroadcaster()

	// start the program runner
	go programRunner(liveReloadPort)

	// start the http server
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(outputPrefix+"Error starting live reload websocket server:", err)
		os.Exit(1)
	}
}

// compileDirectoriesWithWatchedFiles is a function that will compile a list of
// all directorys that contain go code and are children of the current project
// root folder
func compileDirectoriesWithWatchedFiles() []string {
	discoveredDirs := []string{"."}
	watchedDirs := []string{}

	// recursively find all folders that contain watched files
	for len(discoveredDirs) > 0 {
		dir := discoveredDirs[0]
		discoveredDirs = discoveredDirs[1:]
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(outputPrefix+"Error reading directory:", err)
			continue
		}

		for _, file := range files {
			if file.IsDir() && !slices.Contains(autoRebuildIgnoreDirs, file.Name()) {
				discoveredDirs = append(discoveredDirs, dir+"/"+file.Name())
			}
		}

		for _, file := range files {
			if !file.IsDir() && slices.Contains(autoRebuildFileExtensions, filepath.Ext(file.Name())) {
				watchedDirs = append(watchedDirs, dir)
				break
			}
		}
	}

	// return the directories
	return watchedDirs
}

func watchDirectories(dirs []string) {
	// Create a new watcher using fsnotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(outputPrefix+"Error creating watcher:", err)
		os.Exit(1)
	}
	defer watcher.Close()

	// Listen for events on the watcher
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if slices.Contains(autoRebuildFileExtensions, filepath.Ext(event.Name)) {
						fmt.Println(outputPrefix+"Watched file modified:", event.Name)
						rebuildChannel <- backendRebuild
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println(outputPrefix+"Error:", err)
			}
		}
	}()

	for _, dir := range dirs {

		// Add the frontend directory to the watcher
		err = watcher.Add(dir)
		if err != nil {
			fmt.Println(outputPrefix+"Error adding directory to watcher:", err)
			os.Exit(1)
		}
	}

	fmt.Println(outputPrefix+"Watching directories:", dirs, "with extensions:", autoRebuildFileExtensions)

	// Wait for the watcher to finish
	<-make(chan struct{})
}

func pingServerUntilResponse(appPort string) bool {
	pingURL := "http://localhost:" + appPort

	for i := 0; i < 40; i++ {
		_, err := http.Get(pingURL)
		if err != nil {
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return true
	}

	return false
}

// instantiates a new server-run process
func createProcess(reloadPort string) *exec.Cmd {
	program := exec.Command("./build.sh", os.Args[1:]...)
	program.Env = os.Environ()
	program.Env = append(program.Env, "APP_ENV=development", "APP_LIVERELOAD_PORT="+reloadPort)

	// stdErrChannel := make(chan string)
	go pipeOutput(program.StdoutPipe, outputPrefixMain, nil)
	go pipeOutput(program.StderrPipe, outputPrefixMain, nil)

	go func() {
		// find the app port
		appPort := os.Getenv("APP_PORT")
		if appPort == "" {
			return
		}

		// ping the server until it responds
		shouldReload := pingServerUntilResponse(appPort)

		// send a reload message to the frontend
		if shouldReload {
			fmt.Println(outputPrefix + "Webserver started. Reloading frontend...")
			frontendReloadChannel <- true
		} else {
			panic("Couldn't ping server, gotta go!")
		}
	}()

	return program
}

// runs the program and restarts on reloadChannel events
func programRunner(port string) {
	// run the program
	program := createProcess(port)
	err := program.Start()
	if err != nil {
		fmt.Println(outputPrefix+"Error starting program:", err)
		os.Exit(1)
	}

	// listen for changes
	for range rebuildChannel {
		fmt.Println(outputPrefix + "Rebuilding program...")

		// send an interrupt signal to the program
		err := program.Process.Kill()
		if err != nil {
			fmt.Println(outputPrefix+"Error sending interrupt signal to program:", err)
		}

		// wait for the program to finish
		err = program.Wait()
		if err != nil {
			fmt.Println(outputPrefix+"Error waiting for program to finish:", err)
		}

		// create a new program
		program = createProcess(port)
		err = program.Start()
		if err != nil {
			fmt.Println(outputPrefix+"Error starting program:", err)
			os.Exit(1)
		}
	}
}

// watches for frontend code changes and sends a message to the reload channel
func laravelMixWatcher() {

	fmt.Println(outputPrefixLaravelMix + "Building...")

	// run the command
	command := exec.Command("npm", "exec", "mix", "watch")
	go pipeOutput(command.StderrPipe, outputPrefixLaravelMix, nil)

	// listen to the command output
	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println(outputPrefix+"Error getting laravel mix watcher stdout pipe:", err)
		os.Exit(1)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "webpack compiled successfully") {
				fmt.Println(outputPrefixLaravelMix + "Build successful")
				// trigger a rebuild
				rebuildChannel <- frontendRebuild
			}
		}
	}()

	err = command.Start()
	if err != nil {
		fmt.Println(outputPrefix+"Error starting laravel mix watcher:", err)
		os.Exit(1)
	}

	// wait for initial mix build to complete
	<-rebuildChannel
}

// handles the frontend reload websocket connection
func handleFrontendReload(w http.ResponseWriter, r *http.Request) {

	// Create a new websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(outputPrefix+"Error upgrading websocket connection:", err)
		return
	}

	// add the client to the list of clients
	frontendReloadClients[conn] = true

	conn.WriteMessage(websocket.TextMessage, []byte("connected"))

	// listen for disconnect
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(frontendReloadClients, conn)
				break
			}
		}
	}()
}

// listens for messages on the frontend reload channel and broadcasts them to all connected clients
func startFrontendReloadBroadcaster() {

	// listen for messages on the frontend reload channel
	for range frontendReloadChannel {
		for conn := range frontendReloadClients {
			conn.WriteMessage(websocket.TextMessage, []byte("reload"))
		}
	}
}

// pipes the output of a command to the console and an optional channel
func pipeOutput(pipeGetter func() (io.ReadCloser, error), prefix string, msgChannel chan string) {
	pipe, err := pipeGetter()
	if err != nil {
		fmt.Println(outputPrefix+"Error getting pipe:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(prefix + line)
		if msgChannel != nil {
			msgChannel <- line
		}
	}
}
