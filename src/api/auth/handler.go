package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"supportflow/core"
	"supportflow/db/postgre"
)

var (
	stateStore   = make(map[string]time.Time)
	stateStoreMu sync.Mutex
)

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func getOAuthConfig() *oauth2.Config {
	redirectURL := core.GetString("google.oauth.redirect.url", "")
	if redirectURL == "" {
		port := core.GetString("service.port", "8080")
		redirectURL = fmt.Sprintf("http://localhost:%s/api/auth/google/callback", port)
	}
	return &oauth2.Config{
		ClientID:     core.GetString("google.oauth.client.id", ""),
		ClientSecret: core.GetString("google.oauth.client.secret", ""),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := getOAuthConfig()
	if cfg.ClientID == "" {
		log.Println("[Google Auth] client_id not configured")
		http.Error(w, `{"error":"google oauth not configured"}`, http.StatusServiceUnavailable)
		return
	}

	state, err := generateState()
	if err != nil {
		log.Printf("[Google Auth] failed to generate state: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	stateStoreMu.Lock()
	stateStore[state] = time.Now().Add(10 * time.Minute)
	stateStoreMu.Unlock()
	log.Printf("[Google Auth] State saved, total states: %d", len(stateStore))

	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Printf("[Google Auth] Generated auth URL, state=%s", state[:8])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errParam := r.URL.Query().Get("error")

	if errParam != "" {
		log.Printf("[Google Auth] Error from Google: %s", errParam)
		writeCallbackHTML(w, "", errParam)
		return
	}

	if code == "" || state == "" {
		log.Println("[Google Auth] Missing code or state")
		writeCallbackHTML(w, "", "Missing code or state parameter")
		return
	}

	stateStoreMu.Lock()
	exp, valid := stateStore[state]
	if valid && time.Now().Before(exp) {
		delete(stateStore, state)
	} else {
		valid = false
	}
	stateStoreMu.Unlock()

	log.Printf("[Google Auth] Callback state=%s valid=%v", state[:8], valid)

	if !valid {
		writeCallbackHTML(w, "", "Invalid or expired state")
		return
	}

	cfg := getOAuthConfig()
	ctx := context.Background()

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Printf("[Google Auth] Token exchange failed: %v", err)
		writeCallbackHTML(w, "", "Token exchange failed")
		return
	}

	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[Google Auth] Failed to get user info: %v", err)
		writeCallbackHTML(w, "", "Failed to get user info")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var profile googleUserInfo
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Printf("[Google Auth] Failed to parse user info: %v", err)
		writeCallbackHTML(w, "", "Failed to parse user info")
		return
	}

	log.Printf("[Google Auth] Profile received: email=%s name=%s", profile.Email, profile.Name)

	user, err := postgre.FindUserByGoogleSub(r.Context(), profile.ID)
	if err != nil {
		user, err = postgre.FindUserByEmail(r.Context(), profile.Email)
		if err != nil {
			user, err = postgre.CreateUser(r.Context(), profile.Email, profile.Name, profile.ID)
			if err != nil {
				log.Printf("[Google Auth] Failed to create user: %v", err)
				writeCallbackHTML(w, "", "Failed to create user")
				return
			}
			log.Printf("[Google Auth] Created new user: %s", user.ID)
		} else {
			_ = postgre.UpdateUserGoogleSub(r.Context(), user.ID, profile.ID)
			log.Printf("[Google Auth] Linked Google to existing user: %s", user.ID)
		}
	}

	userJSON, _ := json.Marshal(user)
	log.Printf("[Google Auth] Login successful for user: %s (%s)", user.Email, user.ID)
	writeCallbackHTML(w, string(userJSON), "")
}

func writeCallbackHTML(w http.ResponseWriter, userJSON string, errMsg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var script string
	if errMsg != "" {
		script = fmt.Sprintf(`
			if (window.opener) {
				window.opener.postMessage({ type: 'google-auth-error', error: %q }, '*');
				window.close();
			} else {
				document.body.innerText = %q;
				setTimeout(function() { window.location.href = '/login'; }, 3000);
			}
		`, errMsg, errMsg)
	} else {
		script = fmt.Sprintf(`
			var user = %s;
			if (window.opener) {
				window.opener.postMessage({ type: 'google-auth-success', user: user }, '*');
				window.close();
			} else {
				try { localStorage.setItem('sf-auth-user', JSON.stringify(user)); } catch(e) {}
				window.location.href = '/dashboard';
			}
		`, userJSON)
	}

	fmt.Fprintf(w, `<!DOCTYPE html>
<html><head><title>Kairon Auth</title></head>
<body><script>%s</script></body></html>`, script)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user, err := postgre.FindUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Printf("[Auth] Login failed for %s: user not found", req.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}

	if user.Password == nil || *user.Password != req.Password {
		log.Printf("[Auth] Login failed for %s: wrong password", req.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid credentials"}`))
		return
	}

	log.Printf("[Auth] Login successful for: %s (level %d)", user.Email, user.Level)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func HandleGenerateInvite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreatedBy string `json:"created_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	req.CreatedBy = strings.TrimSpace(req.CreatedBy)
	if req.CreatedBy == "" {
		http.Error(w, `{"error":"created_by is required"}`, http.StatusBadRequest)
		return
	}

	token, err := postgre.CreateInviteToken(r.Context(), req.CreatedBy)
	if err != nil {
		log.Printf("[Auth] Failed to create invite token: %v", err)
		http.Error(w, `{"error":"failed to create invite"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("[Auth] Invite token created by user %s", req.CreatedBy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"link":  "/register/" + token,
	})
}

func HandleValidateInvite(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	if token == "" {
		http.Error(w, `{"valid":false,"error":"token is required"}`, http.StatusBadRequest)
		return
	}

	err := postgre.ValidateInviteToken(r.Context(), token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"valid":false,"error":"invalid or expired token"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"valid": true,
	})
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Name     string `json:"name"`
		Company  string `json:"company"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	req.Name = strings.TrimSpace(req.Name)
	req.Company = strings.TrimSpace(req.Company)
	req.Email = strings.TrimSpace(req.Email)

	if req.Token == "" || req.Name == "" || req.Company == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"all fields are required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"error":"password must be at least 6 characters"}`, http.StatusBadRequest)
		return
	}

	if err := postgre.ValidateInviteToken(r.Context(), req.Token); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"invalid or expired token"}`))
		return
	}

	existing, _ := postgre.FindUserByEmail(r.Context(), req.Email)
	if existing != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":"email already taken"}`))
		return
	}

	user, err := postgre.CreateUserWithPassword(r.Context(), req.Email, req.Name, req.Password, req.Company, 4, "Admin")
	if err != nil {
		log.Printf("[Auth] Failed to create user via invite: %v", err)
		http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
		return
	}

	if err := postgre.ConsumeInviteToken(r.Context(), req.Token, user.ID); err != nil {
		log.Printf("[Auth] Failed to consume invite token: %v", err)
	}

	log.Printf("[Auth] User registered via invite: %s (company: %s, level: 4)", user.Email, req.Company)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
