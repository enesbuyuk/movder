package handlers

import (
	"database/sql"
	"fmt"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/websocket"
	"log"
	"net/http"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	functions.GetIpAddress(w, r)
	session, _ := configs.GetNamed(r, configs.SiteConfig.SessionName)
	configs.UserSession = functions.SessionFunc(configs.Db, session.Values["UserId"])

	if session.Values["Language"] == "TR" {
		configs.SiteConfig.TMDBApiLink = configs.SiteConfig.TMDBApiKey + "&language=tr"
		configs.Language = configs.LanguageTR
	} else {
		configs.SiteConfig.TMDBApiLink = configs.SiteConfig.TMDBApiKey + "&language=en"
		configs.Language = configs.LanguageEN
	}

	getKey := r.URL.Query().Get("key")
	ClientIdKey := ""
	if len(getKey) != 0 {
		ClientGet, _ := configs.Db.Query("SELECT client_key FROM client_chat WHERE client_key=? AND (client_owner=? || client_target=?) LIMIT 1", getKey, configs.UserSession.Id, configs.UserSession.Id)
		defer func(ClientGet *sql.Rows) {
			err := ClientGet.Close()
			if err != nil {

			}
		}(ClientGet)
		if ClientGet.Next() {
			err := ClientGet.Scan(&ClientIdKey)
			if err != nil {
				log.Fatal("Fprint error:", err)
			}
			websocketmovder.ServeWs(w, r, ClientIdKey)
		} else {
			_, err := fmt.Fprint(w, "error")
			if err != nil {
				log.Fatal("Fprint error:", err)
			}
		}
	} else {
		_, err := fmt.Fprint(w, "error")
		if err != nil {
			log.Fatal("Fprint error:", err)
		}
	}
}
