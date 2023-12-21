package handlers

import (
	"database/sql"
	"fmt"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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
			if len(r.FormValue("nick")) != 0 && len(r.FormValue("name")) != 0 && len(r.FormValue("email")) != 0 && len(r.FormValue("password1")) != 0 && len(r.FormValue("password2")) != 0 {
				if r.FormValue("password1") == r.FormValue("password2") {
					if passwordvalidator.Validate(r.FormValue("password1"), 60) == nil {
						mail := functions.IsEmailValid(r.FormValue("email"))
						query, _ := configs.Db.Query(`SELECT kullanici_id FROM kullanici WHERE kullanici_nick=?`, r.FormValue("nick"))
						defer func(query *sql.Rows) {
							err := query.Close()
							if err != nil {
								log.Fatal("Query error:", err)
								return
							}
						}(query)
						if query.Next() {
							_, err := fmt.Fprint(w, "nick")
							if err != nil {
								log.Fatal("Fprint error:", err)
								return
							}
						} else {
							queryt, _ := configs.Db.Query(`SELECT kullanici_id FROM kullanici WHERE kullanici_posta=?`, r.FormValue("email"))
							defer func(queryt *sql.Rows) {
								err := queryt.Close()
								if err != nil {
									log.Fatal("Query error:", err)
									return
								}
							}(queryt)
							if queryt.Next() || mail == false {
								_, err := fmt.Fprint(w, "email")
								if err != nil {
									log.Fatal("Fprint error:", err)
									return
								}
							} else {
								ekle, err := configs.Db.Prepare("INSERT INTO kullanici (kullanici_nick, kullanici_isim, kullanici_posta, kullanici_sifre, kullanici_yetki, kullanici_url, kullanici_avatar, kullanici_hakkinda, kullanici_yazi) VALUES (?,?,?,?,?,?,?,?,?);")
								if err != nil {
									panic(err.Error())
								} else {
									reg, err := regexp.Compile("[^A-Za-z0-9]+")
									if err != nil {
										log.Fatal(err)
									}
									kullaniciurl := strings.ToLower(strings.Trim(reg.ReplaceAllString(r.FormValue("nick"), "-"), "-"))
									useremail := html.EscapeString(strings.ToLower(strings.Trim(r.FormValue("email"), "")))
									_, err = ekle.Exec(r.FormValue("nick"), strings.TrimRight(strings.TrimLeft(r.FormValue("name"), ""), ""), useremail, r.FormValue("password1"), "0", kullaniciurl, "default.jpg", "movder.com", "0")
									if err != nil {
										log.Fatal("Query error:", err)
										return
									}
									_, err = fmt.Fprint(w, "ok")
									if err != nil {
										log.Fatal("Fprint error:", err)
									}
									http.Redirect(w, r, "/signin", 302)
								}
							}
						}
					} else {
						_, err := fmt.Fprint(w, "insecure")
						if err != nil {
							log.Fatal("Fprint error:", err)
						}
					}
				} else {
					_, err := fmt.Fprint(w, "password")
					if err != nil {
						log.Fatal("Fprint error:", err)
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
			err := configs.Tmpl.ExecuteTemplate(w, "signup", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.Language.Signup + configs.SiteConfig.TitleLastPart,
			})
			if err != nil {
				log.Fatal("Execute template error:", err)
			}
		}
	}
}
