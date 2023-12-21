package handlers

import (
	"database/sql"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"log"
	"net/http"
)

func FaqHandler(w http.ResponseWriter, r *http.Request) {
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

	query, _ := configs.Db.Query("SELECT ayar_id, ayar_ad, ayar_deger FROM ayar WHERE ayar_ad = 'about' LIMIT 1")
	defer func(query *sql.Rows) {
		err := query.Close()
		if err != nil {
			println("Query Error:", err)
			return
		}
	}(query)
	if query.Next() {
		var settings structures.SettingStruct
		err := query.Scan(&settings.SettingId, &settings.SettingName, &settings.SettingValue)
		if err != nil {
			log.Println("Scan Error:", err)
			return
		}
		about := functions.BBToHtml(settings.SettingValue)

		err = configs.Tmpl.ExecuteTemplate(w, "faq", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.About + configs.SiteConfig.TitleLastPart,

			"About": about,
		})
		if err != nil {
			log.Fatal("Execute template error:", err)
			log.Println("Template Error:", err)
			return
		}
	}
}
