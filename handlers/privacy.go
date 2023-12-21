package handlers

import (
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"log"
	"net/http"
)

func PrivacyHandler(w http.ResponseWriter, r *http.Request) {
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

	err := configs.Tmpl.ExecuteTemplate(w, "privacy", map[string]interface{}{
		"Language":    configs.Language,
		"SessionBool": session.Values["SessionBool"],
		"Session":     configs.UserSession,
		"SiteConfig":  configs.SiteConfig,
		"TitlePage":   configs.Language.Privacy + configs.SiteConfig.TitleLastPart,
	})
	if err != nil {
		log.Fatal("Execute template error:", err)
	}
}
