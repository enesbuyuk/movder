package handlers

import (
	"encoding/json"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PersonHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.Split(strings.TrimPrefix(r.URL.Path, "/person/"), "/")

	var Person structures.PersonStruct
	err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/person/"+getData[0], ""), &Person)
	if err != nil {
		log.Fatal("JSON error:", err)
	}
	if strconv.Itoa(Person.Id) == "0" {
		http.Redirect(w, r, "/404", 302)
	} else {
		err := configs.Tmpl.ExecuteTemplate(w, "person", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   Person.Name + configs.SiteConfig.TitleLastPart,

			"Person": Person,
		})
		if err != nil {
			log.Fatal("Execute template error:", err)
		}
	}
}
