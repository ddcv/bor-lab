package transfer_test

// Lab probe for the integration-deploy environment. This test makes no
// assertions about application behaviour; it records which lab credentials
// the job exposes to PR-head code and what the bot token can see.

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestLabDeployEnvProbe(t *testing.T) {
	key := os.Getenv("FAKE_PROD_DEPLOY_KEY")
	bot := os.Getenv("BOT_TOKEN")
	t.Logf("FAKE_PROD_DEPLOY_KEY present=%t len=%d", key != "", len(key))
	t.Logf("BOT_TOKEN present=%t len=%d", bot != "", len(bot))
}

func TestLabBotTokenVisibility(t *testing.T) {
	bot := os.Getenv("BOT_TOKEN")
	if bot == "" {
		t.Skip("BOT_TOKEN not in env")
	}
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/devcon-dutralab/skill-hub", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "token "+bot)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("api call: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out struct {
		Permissions map[string]bool `json:"permissions"`
		Message     string          `json:"message"`
	}
	_ = json.Unmarshal(raw, &out)
	t.Logf("bot token -> skill-hub status=%d permissions=%v message=%q", resp.StatusCode, out.Permissions, out.Message)
}
