package handlers

import (
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	functions.GetIpAddress(w, r)
	session, _ := configs.GetNamed(r, configs.SiteConfig.SessionName)

	session.Values["SessionBool"] = false
	//store.Options.MaxAge = -1
	err := sessions.Save(r, w)
	if err != nil {
		log.Fatal("Session save error:", err)
	}
	http.Redirect(w, r, "/", 302)
}
