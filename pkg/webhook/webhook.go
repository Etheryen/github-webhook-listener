package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type pushPayload struct {
	Ref string `json:"ref"`
}

func verifySignature(secret, signature string, body []byte) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "Missing signature", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg := "Failed to read request body"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if !verifySignature(os.Getenv("GITHUB_SECRET"), signature, body) {
		http.Error(w, "Invalid signature", http.StatusForbidden)
		return
	}

	var payload pushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	branch := os.Getenv("BRANCH")

	if strings.HasSuffix(payload.Ref, "/"+branch) {
		log.Printf(
			"Push event to %s branch received, executing redeploy...\n",
			branch,
		)

		command := os.Getenv("COMMAND")
		program := strings.Split(command, " ")[0]
		args := strings.Split(command, " ")[1:]

		cmd := exec.Command(program, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Error executing command: %v", err)
			http.Error(w, "Deployment failed", http.StatusInternalServerError)
			return
		}

		log.Println("Redeployment successful")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Redeployment successful"))
		return
	}

	log.Printf("Push to branch '%s' ignored.", payload.Ref)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Branch ignored"))
}
