package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/disgoorg/json"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/disgo/rest"
)

var (
	letters      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	clientID     = snowflake.GetEnv("client_id")
	clientSecret = os.Getenv("client_secret")
	baseURL      = os.Getenv("base_url")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo %s", disgo.Version)

	s := &server{
		client: oauth2.New(clientID, clientSecret,
			oauth2.WithRestClientConfigOpts(
				rest.WithHTTPClient(&http.Client{
					Timeout: 5 * time.Second,
				}),
			),
		),
		sessions: map[string]oauth2.Session{},
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/trylogin", s.handleTryLogin)
	_ = http.ListenAndServe(":6969", mux)
}

type server struct {
	client   *oauth2.Client
	sessions map[string]oauth2.Session
	rand     *rand.Rand
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		writeHTML(w, `<button><a href="/login">login</a></button>`)
	}
	session, ok := s.sessions[cookie.Value]
	if !ok {
		writeHTML(w, `<button><a href="/login">login</a></button>`)
	}

	session, ok, err = s.client.VerifySession(session)
	if err != nil {
		writeError(w, "error while verifying or refresh session", err)
		return
	}
	if ok {
		s.sessions[cookie.Value] = session
	}

	user, err := s.client.GetUser(session)
	if err != nil {
		writeError(w, "error while getting user data", err)
		return
	}
	connections, err := s.client.GetConnections(session)
	if err != nil {
		writeError(w, "error while getting connections data", err)
		return
	}

	userJSON, err := json.MarshalIndent(user, "<br />", "&ensp;")
	if err != nil {
		writeError(w, "error while formatting user data", err)
		return
	}
	connectionsJSON, err := json.MarshalIndent(connections, "<br />", "&ensp;")
	if err != nil {
		writeError(w, "error while formatting connections data", err)
		return
	}

	writeHTML(w, fmt.Sprintf("user:<br />%s<br />connections: <br />%s", userJSON, connectionsJSON))
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.client.GenerateAuthorizationURL(baseURL+"/trylogin", discord.PermissionsNone, 0, false, discord.OAuth2ScopeIdentify, discord.OAuth2ScopeGuilds, discord.OAuth2ScopeEmail, discord.OAuth2ScopeConnections, discord.OAuth2ScopeWebhookIncoming), http.StatusSeeOther)
}

func (s *server) handleTryLogin(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL.Query()
		code  = query.Get("code")
		state = query.Get("state")
	)
	if code != "" && state != "" {
		identifier := s.randStr(32)
		session, _, err := s.client.StartSession(code, state)
		if err != nil {
			writeError(w, "error while starting session", err)
			return
		}
		s.sessions[identifier] = session
		http.SetCookie(w, &http.Cookie{Name: "token", Value: identifier})
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *server) randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[s.rand.Intn(len(letters))]
	}
	return string(b)
}

func writeError(w http.ResponseWriter, text string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(text + ": " + err.Error()))
}

func writeHTML(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(text))
}
