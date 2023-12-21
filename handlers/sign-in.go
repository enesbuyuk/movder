package handlers

import (
	"database/sql"
	"fmt"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
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
	if session.Values["SessionBool"] == true {
		http.Redirect(w, r, "/", 302)
	} else {
		if r.Method == "POST" {
			if len(r.FormValue("kullaniciadi")) != 0 && len(r.FormValue("sifre")) != 0 {
				var kullanici structures.UserStruct
				signInQuery, hata := configs.Db.Query(`SELECT kullanici_id, kullanici_nick, kullanici_isim, kullanici_posta, kullanici_url, kullanici_yetki, kullanici_avatar, kullanici_hakkinda FROM kullanici WHERE kullanici_nick=? && kullanici_sifre=? LIMIT 1`, r.FormValue("kullaniciadi"), r.FormValue("sifre"))
				defer func(query *sql.Rows) {
					err := query.Close()
					if err != nil {
						log.Fatal("Query error:", err)
						return
					}
				}(signInQuery)
				if hata != nil {
					panic(hata.Error())
				}
				if signInQuery.Next() {
					err := signInQuery.Scan(&kullanici.Id, &kullanici.Nick, &kullanici.Name, &kullanici.Mail, &kullanici.Link, &kullanici.Authority, &kullanici.Photo, &kullanici.About)
					if err != nil {
						log.Fatal("Scan error:", err)
						return
					}

					session.Values["SessionBool"] = true
					session.Values["UserId"] = kullanici.Id

					err = sessions.Save(r, w)
					if err != nil {
						log.Println("Session Save Error:", err)
						return
					}
					_, err = fmt.Fprint(w, "ok")
					if err != nil {
						log.Fatal("Fprint error:", err)
						return
					}
					http.Redirect(w, r, "/", 302)
				} else {
					_, err := fmt.Fprint(w, "error")
					if err != nil {
						log.Fatal("Fprint error:", err)
						return
					}
				}

			} else {
				_, err := fmt.Fprint(w, "error")
				if err != nil {
					log.Fatal("Fprint error:", err)
					return
				}
			}

		} else {
			err := configs.Tmpl.ExecuteTemplate(w, "signin", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				//		"Session":UserSession,
				"SiteConfig": configs.SiteConfig,
				"TitlePage":  configs.Language.Signin + configs.SiteConfig.TitleLastPart,
			})
			if err != nil {
				log.Fatal("Execute template error:", err)
			}
		}
	}
}
