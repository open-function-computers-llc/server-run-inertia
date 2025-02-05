package server

import (
	"encoding/json"
	"net/http"

	scriptrunner "github.com/open-function-computers-llc/server-run-inertia/script-runner"
)

func (s *server) handleStreamScriptRunner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		script := r.FormValue("script")

		// parse request `args`
		args := map[string]string{}
		requestArgs := r.FormValue("args")
		requestArgsErr := json.Unmarshal([]byte(requestArgs), &args)
		if requestArgsErr != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Couldn't parse required param `args`",
			})
			return
		}

		// parse request `env`
		env := map[string]string{}
		requestEnv := r.FormValue("env")
		requestEnvErr := json.Unmarshal([]byte(requestEnv), &env)
		if requestEnvErr != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Couldn't parse required param `env`",
			})
			return
		}

		// set up webhook
		c, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Problem setting up the web socket connection",
			})
		}
		defer c.Close()

		commChannel := make(chan string)

		go scriptrunner.StreamScriptOutput(script, args, env, &commChannel)

		for incomingMessage := range commChannel {
			message := map[string]interface{}{
				"message": incomingMessage,
			}
			c.WriteJSON(message)
		}

		c.Close()
	}
}
