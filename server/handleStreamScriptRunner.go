package server

import (
	"math/rand"
	"net/http"
	"time"
)

func (s *server) handleStreamScriptRunner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := s.upgrader.Upgrade(w, r, nil)

		if err != nil {
			sendJSONError(w, http.StatusInternalServerError, map[string]string{
				"error": "Problem setting up the web socket connection",
			})
		}
		defer c.Close()

		maxLines := rand.Intn(55) + 40
		i := 0
		for {
			message := "yo"
			for j := 0; j <= i; j++ {
				message += "!"
			}
			c.WriteJSON(map[string]string{"message": message})

			seconds := rand.Intn(500) + 50
			time.Sleep(time.Duration(seconds) * time.Millisecond)

			i++
			if i > maxLines {
				break
			}
		}
	}
}
