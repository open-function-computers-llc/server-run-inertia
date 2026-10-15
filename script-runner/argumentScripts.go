package scriptrunner

import "fmt"

// argumentScripts take positional arguments, so each one lists the arg names it
// expects in order. fixedArgs are passed first, before any named params.
var argumentScripts = map[string]argumentScript{
	"unlock": {
		filename: "ofco-unlock-site.sh",
		params:   []string{"ACCOUNT_NAME"},
	},
	"lock": {
		filename: "ofco-lock-site-production.sh",
		params:   []string{"ACCOUNT_NAME"},
	},
	// ofco-stream-file.sh is `tail $1 $2 $3 $4`. LOG_FILE must be built by the
	// server, never taken from the browser, or it could read any file.
	"log-viewer": {
		filename:  "ofco-stream-file.sh",
		fixedArgs: []string{"-n", "200"},
		params:    []string{"LOG_FILE"},
	},
}

type argumentScript struct {
	filename  string
	fixedArgs []string
	params    []string
}

// orderedArgs returns the script's arguments in the order the script expects.
func (s argumentScript) orderedArgs(args map[string]string) ([]string, error) {
	values := append([]string{}, s.fixedArgs...)
	for _, p := range s.params {
		v, ok := args[p]
		if !ok || v == "" {
			return nil, fmt.Errorf("%w: %s", ErrorMissingArgument, p)
		}
		values = append(values, v)
	}
	return values, nil
}
