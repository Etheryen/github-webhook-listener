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
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

func verifySignature(secret, signature string, body []byte) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func Handler(config *Config, githubSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if !verifySignature(githubSecret, signature, body) {
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}

		var payload pushPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		projectName := payload.Repository.Name

		details, err := config.get(projectName)
		if err != nil {
			http.Error(w, "Project doesn't exist", http.StatusBadRequest)
			return
		}

		if strings.HasSuffix(payload.Ref, "/"+details.Branch) {
			log.Printf(
				"%s -> push event to %s branch received, executing redeploy...\n",
				projectName,
				details.Branch,
			)

			cmd := exec.Command("sh", "-c", details.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				log.Printf("Error executing command: %v", err)
				http.Error(
					w,
					"Deployment failed",
					http.StatusInternalServerError,
				)
				return
			}

			log.Printf("%s -> redeployment successful\n", projectName)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Redeployment successful"))
			return
		}

		log.Printf(
			"%s -> push to branch '%s' ignored.",
			projectName,
			payload.Ref,
		)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Branch ignored"))
	}
}
