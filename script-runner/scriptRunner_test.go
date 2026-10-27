package scriptrunner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func writeScript(t *testing.T, root, filename, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, filename), []byte(contents), 0755); err != nil {
		t.Fatal(err)
	}
}

func runStream(t *testing.T, script string, args, env map[string]string) ([]StreamMessage, error) {
	t.Helper()

	ch := make(chan StreamMessage, 64)
	var msgs []StreamMessage
	done := make(chan struct{})
	go func() {
		for m := range ch {
			msgs = append(msgs, m)
		}
		close(done)
	}()

	err := StreamScriptOutput(script, args, env, &ch)
	<-done
	return msgs, err
}

func registerPlainScript(t *testing.T, name, filename string) {
	t.Helper()
	original, existed := plainScripts[name]
	plainScripts[name] = filename
	t.Cleanup(func() {
		if existed {
			plainScripts[name] = original
		} else {
			delete(plainScripts, name)
		}
	})
}

func registerEnvScript(t *testing.T, name string, spec envScript) {
	t.Helper()
	original, existed := envScripts[name]
	envScripts[name] = spec
	t.Cleanup(func() {
		if existed {
			envScripts[name] = original
		} else {
			delete(envScripts, name)
		}
	})
}

func TestStreamScriptOutputTypedLines(t *testing.T) {
	root := t.TempDir()
	t.Setenv("SCRIPTS_ROOT", root)

	writeScript(t, root, "showcase.sh", `#!/bin/bash
echo "##section## Doing the thing"
echo "step one complete"
echo "##error## uh oh"
`)
	registerPlainScript(t, "test-showcase", "showcase.sh")

	msgs, err := runStream(t, "test-showcase", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := []StreamMessage{
		{"running test-showcase", MsgMeta},
		{"args:", MsgMeta},
		{"env:", MsgMeta},
		{"Doing the thing", MsgSection},
		{"step one complete", MsgStdout},
		{"uh oh", MsgError},
	}
	if !slices.Equal(msgs, want) {
		t.Fatalf("got %+v, want %+v", msgs, want)
	}
}

func TestStreamScriptOutputExitError(t *testing.T) {
	root := t.TempDir()
	t.Setenv("SCRIPTS_ROOT", root)

	writeScript(t, root, "fail.sh", "#!/bin/bash\necho \"working\"\nexit 3\n")
	registerPlainScript(t, "test-fail", "fail.sh")

	msgs, err := runStream(t, "test-fail", nil, nil)
	if err == nil {
		t.Fatal("expected an error for a script that exits 3")
	}

	last := msgs[len(msgs)-1]
	if last.Type != MsgError || !strings.Contains(last.Message, "exit status 3") {
		t.Fatalf("expected exit error line, got %+v", last)
	}
}

func TestStreamScriptOutputRedactsSecretEnvValues(t *testing.T) {
	root := t.TempDir()
	t.Setenv("SCRIPTS_ROOT", root)

	writeScript(t, root, "secrecy.sh", "#!/bin/bash\ntrue\n")
	registerEnvScript(t, "test-secret", envScript{
		filename:     "secrecy.sh",
		requiredEnvs: []string{"IPADDRESS"},
		optionalEnvs: []string{"ADMIN_PASSWORD", "GITEA_TOKEN"},
	})

	msgs, err := runStream(t, "test-secret", nil, map[string]string{
		"IPADDRESS":      "203.0.113.7",
		"ADMIN_PASSWORD": "hunter2",
		"GITEA_TOKEN":    "abc123",
	})
	if err != nil {
		t.Fatal(err)
	}

	joined := ""
	for _, m := range msgs {
		joined += m.Message + "\n"
	}

	for _, want := range []string{"ADMIN_PASSWORD=********", "GITEA_TOKEN=********", "IPADDRESS=203.0.113.7"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected %q in stream output, got:\n%s", want, joined)
		}
	}
	for _, leak := range []string{"hunter2", "abc123"} {
		if strings.Contains(joined, leak) {
			t.Fatalf("secret %q leaked into stream output:\n%s", leak, joined)
		}
	}
}

func TestStreamScriptOutputRejectsBadInputBeforeRunning(t *testing.T) {
	root := t.TempDir()
	t.Setenv("SCRIPTS_ROOT", root)

	// override the real ban-lookup file so we can tell whether it ran
	writeScript(t, root, "ban-lookup-never.sh", "#!/bin/bash\necho \"should never run\"\n")
	registerEnvScript(t, "ban-lookup", envScript{filename: "ban-lookup-never.sh", requiredEnvs: []string{"IPADDRESS"}})

	msgs, err := runStream(t, "ban-lookup", nil, map[string]string{"IPADDRESS": "not-an-ip"})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if len(msgs) == 0 || msgs[len(msgs)-1].Type != MsgError {
		t.Fatalf("expected error message on channel, got %+v", msgs)
	}
	for _, m := range msgs {
		if strings.Contains(m.Message, "should never run") {
			t.Fatal("script ran despite invalid env")
		}
	}
}
