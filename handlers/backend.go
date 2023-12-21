package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gorilla/sessions"
	"github.com/gosimple/slug"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func BackendHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.TrimPrefix(r.URL.Path, "/backend/")

	if getData == "language-tr" {
		session.Values["Language"] = "TR"
		err := sessions.Save(r, w)
		if err != nil {
			log.Println("Session Save Error:", err)
			return
		}
		_, err = fmt.Fprint(w, "TR")
		if err != nil {
			log.Fatal("Fprint error:", err)
			return
		}
	} else if getData == "language-en" {
		session.Values["Language"] = "EN"
		err := sessions.Save(r, w)
		if err != nil {
			println("Session Save Error:", err)
			return
		}
		_, err = fmt.Fprint(w, "EN")
		if err != nil {
			log.Fatal("Fprint error:", err)
			return
		}
	}

	if session.Values["SessionBool"] == true {
		if getData == "settings" {
			mail := functions.IsEmailValid(r.FormValue("email"))
			if mail == true {
				_, err := configs.Db.Exec("UPDATE kullanici SET kullanici_isim=?,kullanici_telephone=?,kullanici_posta=?,kullanici_hakkinda=? WHERE kullanici_id=?", r.FormValue("name"), r.FormValue("telephone"), r.FormValue("email"), r.FormValue("bio"), configs.UserSession.Id)
				if err != nil {
					log.Println("Query Error:", err)
				}
				_, err = fmt.Fprint(w, "ok")
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			} else {
				_, err := fmt.Fprint(w, "mail")
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}
		} else if getData == "changepassword" {
			if len(r.FormValue("oldpassword")) != 0 && len(r.FormValue("newpassword1")) != 0 && len(r.FormValue("newpassword2")) != 0 {
				if r.FormValue("newpassword1") == r.FormValue("newpassword2") {
					if r.FormValue("oldpassword") == configs.UserSession.Password {
						if passwordvalidator.Validate(r.FormValue("newpassword1"), 60) == nil {
							_, err := configs.Db.Exec("UPDATE kullanici SET kullanici_sifre=? WHERE kullanici_id=?", r.FormValue("newpassword1"), configs.UserSession.Id)
							if err != nil {
								log.Fatal("Fprint error:", err)
								return
							}
							_, err2 := fmt.Fprint(w, "ok")
							if err2 != nil {
								log.Fatal("Fprint error:", err)
								return
							}
						} else {
							_, err := fmt.Fprint(w, "insecure")
							if err != nil {
								log.Fatal("Fprint error:", err)
								return
							}
						}
					} else {
						_, err := fmt.Fprint(w, "oldpassword")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				} else {
					_, err := fmt.Fprint(w, "newpassword")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			} else {
				_, err := fmt.Fprint(w, "empty")
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}
		} else if getData == "messagesend" {
			if len(r.FormValue("target")) != 0 {
				if len(strings.TrimRight(strings.TrimLeft(r.FormValue("content"), " "), " ")) != 0 {
					_, err := configs.Db.Exec("INSERT INTO message(message_owner,message_target,message_content,message_view,message_delete,message_owner_delete,message_target_delete) VALUES (?,?,?,0,0,0,0);", session.Values["UserId"], r.FormValue("target"), strings.TrimRight(strings.TrimLeft(r.FormValue("content"), " "), " "))
					if err != nil {
						log.Println("Query Error:", err)
					}
					_, err = fmt.Fprint(w, "ok")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				} else {
					_, err := fmt.Fprint(w, "content")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			} else {
				fmt.Print(w, "target")
			}
		} else if getData == "commentfilm" {
			if len(r.FormValue("film")) != 0 && len(r.FormValue("filmtype")) != 0 {
				if len(strings.TrimRight(strings.TrimLeft(r.FormValue("content"), " "), " ")) != 0 {
					_, err := configs.Db.Exec("INSERT INTO comment(comment_film,comment_owner,comment_content,comment_rate,comment_like,comment_dislike,comment_filmtype) VALUES (?,?,?,?,0,0,?);", r.FormValue("film"), session.Values["UserId"], strings.TrimRight(strings.TrimLeft(r.FormValue("content"), " "), " "), 1, strings.Trim(r.FormValue("filmtype"), " "))
					if err != nil {
						log.Println("Query Error:", err)
					}
					_, err2 := fmt.Fprint(w, "ok")
					if err2 != nil {
						log.Fatal("Fprint error:", err)
					}
				} else {
					_, err := fmt.Fprint(w, "content")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			} else {
				_, err := fmt.Fprint(w, "film")
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}
		} else if getData == "commentblog" {
			if len(r.FormValue("blog")) != 0 {
				if len(r.FormValue("content")) != 0 {
					_, err := configs.Db.Exec("INSERT INTO comment_blog(comment_blog,comment_owner,comment_content,comment_rate,comment_like,comment_dislike) VALUES (?,?,?,?,0,0);", r.FormValue("blog"), session.Values["UserId"], r.FormValue("content"), 1)
					if err != nil {
						log.Println("Query Error:", err)
					}
					_, err = fmt.Fprint(w, "ok")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				} else {
					_, err := fmt.Fprint(w, "content")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			} else {
				_, err := fmt.Fprint(w, "blog")
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}
		} else if getData == "follow" {
			if r.Method == "POST" {
				datatoint, _ := strconv.Atoi(r.FormValue("follow_owner"))
				if (len(r.FormValue("follow_owner")) != 0) && session.Values["UserId"] != datatoint {
					query, _ := configs.Db.Query(`SELECT * FROM follow WHERE follow_user=? AND follow_owner=?`, session.Values["UserId"], r.FormValue("follow_owner"))
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						sil, err := configs.Db.Prepare("DELETE FROM follow WHERE follow_user=? AND follow_owner=?")
						if err != nil {
							panic(err.Error())
						}
						_, err = sil.Exec(session.Values["UserId"], r.FormValue("follow_owner"))
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "remove")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("INSERT INTO follow (follow_user,follow_owner) VALUES (?,?);", session.Values["UserId"], r.FormValue("follow_owner"))
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "createlist" {
			if r.Method == "POST" {
				if len(r.FormValue("list_title")) != 0 && len(r.FormValue("list_type")) != 0 && (r.FormValue("list_type") == "everyone") || r.FormValue("list_type") == "me" {
					reg, _ := regexp.Compile("[^A-Za-z0-9]+")
					listurl := strings.ToLower(strings.Trim(reg.ReplaceAllString(r.FormValue("list_title"), "-"), "-"))
					listtitle := strings.TrimRight(strings.TrimLeft(r.FormValue("list_title"), ""), "")
					_, err := configs.Db.Exec("INSERT INTO list (list_owner,list_title,list_url,list_count,list_favorite,list_type,list_image) VALUES (?,?,?,'0','0',?,?);", session.Values["UserId"], listtitle, listurl, r.FormValue("list_type"), "edb97c2b2f8ed10920bd4d3af014c360.jpg")
					if err != nil {
						log.Println("Query Error:", err)
					}
					_, err = fmt.Fprint(w, "add")
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			}
		} else if getData == "watchlist" {
			if r.Method == "POST" {
				if len(r.FormValue("film")) != 0 && (r.FormValue("filmtype") == "movie" || r.FormValue("filmtype") == "tv") {
					query, _ := configs.Db.Query(`SELECT * FROM watchlist WHERE watchlist_film=? AND watchlist_filmtype=? AND watchlist_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						_, err := configs.Db.Exec("DELETE FROM watchlist WHERE watchlist_film=? AND watchlist_filmtype=? AND watchlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "remove")
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("INSERT INTO watchlist (watchlist_film,watchlist_filmtype,watchlist_owner) VALUES (?,?,?);", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "watchedlist" {
			if r.Method == "POST" {
				if len(r.FormValue("film")) != 0 && (r.FormValue("filmtype") == "movie" || r.FormValue("filmtype") == "tv") {
					query, _ := configs.Db.Query(`SELECT * FROM watchedlist WHERE watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						_, err := configs.Db.Exec("DELETE FROM watchedlist WHERE watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "remove")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("DELETE FROM watchlist WHERE watchlist_film=? AND watchlist_filmtype=? AND watchlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = configs.Db.Exec("INSERT INTO watchedlist (watchedlist_film,watchedlist_filmtype,watchedlist_owner) VALUES (?,?,?);", r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
						if err != nil {
							log.Println("Query Error:", err)
							return
						}

						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "diary" {
			if r.Method == "POST" {
				if len(r.FormValue("film")) != 0 && (r.FormValue("filmtype") == "movie" || r.FormValue("filmtype") == "tv") {
					query, _ := configs.Db.Query(`SELECT * FROM watchedlist WHERE watchedlist_diary=0 AND watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						_, err := configs.Db.Exec("UPDATE watchedlist SET watchedlist_diary=1 WHERE watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("UPDATE watchedlist SET watchedlist_diary=0 WHERE watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "remove")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "favorite" {
			if r.Method == "POST" {
				if len(r.FormValue("film")) != 0 && (r.FormValue("filmtype") == "movie" || r.FormValue("filmtype") == "tv") {
					query, _ := configs.Db.Query(`SELECT * FROM favorite WHERE favorite_film=? AND favorite_filmtype=? AND favorite_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						_, err := configs.Db.Exec("DELETE FROM favorite WHERE favorite_film=? AND favorite_filmtype=? AND favorite_owner=?", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "remove")
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("INSERT INTO favorite (favorite_film,favorite_filmtype,favorite_owner) VALUES (?,?,?);", r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "rate" {
			if r.Method == "POST" {
				if len(r.FormValue("value")) != 0 && len(r.FormValue("film")) != 0 && (r.FormValue("filmtype") == "movie" || r.FormValue("filmtype") == "tv") {
					WacthedQuery, _ := configs.Db.Query(`SELECT * FROM watchedlist WHERE watchedlist_film=? AND watchedlist_filmtype=? AND watchedlist_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), session.Values["UserId"])
					defer func(WacthedQuery *sql.Rows) {
						err := WacthedQuery.Close()
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					}(WacthedQuery)
					if WacthedQuery.Next() {
					} else {
						_, err := configs.Db.Exec("DELETE FROM watchlist WHERE watchlist_film=? AND watchlist_filmtype=? AND watchlist_owner=?", r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
						_, err = configs.Db.Exec("INSERT INTO watchedlist (watchedlist_film,watchedlist_filmtype,watchedlist_owner) VALUES (?,?,?);", r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					}
					query, _ := configs.Db.Query(`SELECT * FROM rate WHERE rate_film=? AND rate_filmtype=? AND rate_owner=?`, r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
					defer func(query *sql.Rows) {
						err := query.Close()
						if err != nil {
							log.Println("Query Error:", err)
							return
						}
					}(query)
					if query.Next() {
						_, err := configs.Db.Exec("UPDATE rate SET rate_value=? WHERE rate_film=? AND rate_filmtype=? AND rate_owner=?", r.FormValue("value"), r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id)
						if err != nil {
							println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					} else {
						_, err := configs.Db.Exec("INSERT INTO rate (rate_film,rate_filmtype,rate_owner,rate_value) VALUES (?,?,?,?);", r.FormValue("film"), r.FormValue("filmtype"), configs.UserSession.Id, r.FormValue("value"))
						if err != nil {
							println("Query Error:", err)
							return
						}
						_, err = fmt.Fprint(w, "add")
						if err != nil {
							log.Fatal("Fprint error:", err)
							return
						}
					}
				}
			}
		} else if getData == "getnotification" {
			Query, _ := configs.Db.Query(`SELECT notification_id,notification_owner,notification_target,notification_content,notification_view FROM notification WHERE notification_owner=? && notification_view =0 ORDER BY notification_id DESC LIMIT 5`, configs.UserSession.Id)
			defer func(Query *sql.Rows) {
				err := Query.Close()
				if err != nil {
					log.Println("Query Error:", err)
				}
			}(Query)
			var Notification structures.NotificationStruct
			var Notifications []structures.NotificationStruct
			for Query.Next() {
				err := Query.Scan(&Notification.Id, &Notification.Owner.Id, &Notification.Target.Id, &Notification.Content, &Notification.View)
				if err != nil {
				}
				Notifications = append(Notifications, Notification)
			}
			err := configs.Tmpl.ExecuteTemplate(w, "getnotification", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.SiteConfig.TitleHome,

				"Notifications": Notifications,
			})
			if err != nil {
				return
			}
		} else if getData == "notification-count" {
			Query, _ := configs.Db.Query(`SELECT COUNT(notification_id) AS notification_count FROM notification WHERE notification_owner=? && notification_view= 0;`, configs.UserSession.Id)
			defer func(Query *sql.Rows) {
				err := Query.Close()
				if err != nil {

				}
			}(Query)
			var count int
			if Query.Next() {
				err := Query.Scan(&count)
				if err != nil {
					log.Fatal("Scan error:", err)
				}
				if count != 0 {
					_, err := fmt.Fprint(w, `<span id="notification">`+strconv.Itoa(count)+`</span>`)
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				} else {
					_, err := fmt.Fprint(w, `<span id="notification"></span>`)
					if err != nil {
						log.Fatal("Fprint error:", err)
					}
				}
			} else {
				_, err := fmt.Fprint(w, `<span id="notification"></span>`)
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}

			QueryDM, _ := configs.Db.Query(`SELECT COUNT(client_id) AS client_count FROM client_chat WHERE client_owner=? && client_didntviewmsg != 0;`, configs.UserSession.Id)
			defer func(QueryDM *sql.Rows) {
				err := QueryDM.Close()
				if err != nil {

				}
			}(QueryDM)
			var dmcount int
			if Query.Next() {
				err := Query.Scan(&dmcount)
				if err != nil {
					log.Fatal("Scan error:", err)
				}
				_, err = fmt.Fprint(w, `<span id="dm">`+strconv.Itoa(dmcount)+`</span>`)
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			} else {
				_, err := fmt.Fprint(w, `<span id="dm"></span>`)
				if err != nil {
					log.Fatal("Fprint error:", err)
				}
			}
		}

	} else if getData == "changepassword" || getData == "messagesend" || getData == "commentfilm" || getData == "commentblog" || getData == "follow" || getData == "createlist" || getData == "watchlist" || getData == "watchedlist" || getData == "diary" || getData == "favorite" ||
		getData == "rate" || getData == "getnotification" || getData == "notification-count" {
		_, err := fmt.Fprint(w, "permission")
		if err != nil {
			log.Fatal("Fprint error:", err)
		}
	}

	if getData == "films" {
		if r.FormValue("get") == "popular-movies" {
			var PopularMovies structures.AllFilmstruct
			err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/popular", "&page="+r.FormValue("page")), &PopularMovies)

			for i := 0; i < len(PopularMovies.Results); i++ {
				PopularMovies.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
				PopularMovies.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
				if PopularMovies.Results[i].Title == "" {
					PopularMovies.Results[i].Title = PopularMovies.Results[i].Name
					PopularMovies.Results[i].OTitle = PopularMovies.Results[i].OName
				}
				PopularMovies.Results[i].Url = slug.Make(PopularMovies.Results[i].OTitle)
			}

			err = configs.Tmpl.ExecuteTemplate(w, "getpopularmovies", map[string]interface{}{
				"Language":   configs.Language,
				"SiteConfig": configs.SiteConfig,
				"TitlePage":  configs.SiteConfig.TitleHome,

				"PopularMovies": PopularMovies,
			})
			if err != nil {
				log.Println("Template Error:", err)
				return
			}
		} else if r.FormValue("get") == "popular-series" {
			var PopularSeries structures.AllFilmstruct
			err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/tv/popular", "&page="+r.FormValue("page")), &PopularSeries)

			for i := 0; i < len(PopularSeries.Results); i++ {
				PopularSeries.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
				PopularSeries.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
				if PopularSeries.Results[i].Title == "" {
					PopularSeries.Results[i].Title = PopularSeries.Results[i].Name
					PopularSeries.Results[i].OTitle = PopularSeries.Results[i].OName
				}
				PopularSeries.Results[i].Url = slug.Make(PopularSeries.Results[i].OTitle)
			}

			err = configs.Tmpl.ExecuteTemplate(w, "getpopularseries", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.SiteConfig.TitleHome,

				"PopularSeries": PopularSeries,
			})
			if err != nil {
				log.Println("Template Error:", err)
				return
			}
		} else if r.FormValue("get") == "trend-movies" {
			var TrendMovies structures.AllFilmstruct
			err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/movie/day", "&page="+r.FormValue("page")), &TrendMovies)

			for i := 0; i < len(TrendMovies.Results); i++ {
				TrendMovies.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendMovies.Results[i].Id, configs.UserSession.Id)
				TrendMovies.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendMovies.Results[i].Id, configs.UserSession.Id)
				if TrendMovies.Results[i].Title == "" {
					TrendMovies.Results[i].Title = TrendMovies.Results[i].Name
					TrendMovies.Results[i].OTitle = TrendMovies.Results[i].OName
				}
				TrendMovies.Results[i].Url = slug.Make(TrendMovies.Results[i].OTitle)
			}

			err = configs.Tmpl.ExecuteTemplate(w, "gettrendmovies", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.SiteConfig.TitleHome,

				"TrendMovies": TrendMovies,
			})
			if err != nil {
				log.Println("Template Error:", err)
				return
			}
		} else if r.FormValue("get") == "trend-series" {
			var TrendSeries structures.AllFilmstruct
			err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/tv/day", "&page="+r.FormValue("page")), &TrendSeries)

			for i := 0; i < len(TrendSeries.Results); i++ {
				TrendSeries.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendSeries.Results[i].Id, configs.UserSession.Id)
				TrendSeries.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendSeries.Results[i].Id, configs.UserSession.Id)
				if TrendSeries.Results[i].Title == "" {
					TrendSeries.Results[i].Title = TrendSeries.Results[i].Name
					TrendSeries.Results[i].OTitle = TrendSeries.Results[i].OName
				}
				TrendSeries.Results[i].Url = slug.Make(TrendSeries.Results[i].OTitle)
			}

			err = configs.Tmpl.ExecuteTemplate(w, "gettrendseries", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.SiteConfig.TitleHome,

				"TrendSeries": TrendSeries,
			})
			if err != nil {
				log.Println("Template Error:", err)
				return
			}
		} else if r.FormValue("get") == "upcoming" {
			var Upcoming structures.AllFilmstruct
			err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/upcoming", "&page="+r.FormValue("page")), &Upcoming)

			for i := 0; i < len(Upcoming.Results); i++ {
				Upcoming.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, Upcoming.Results[i].Id, configs.UserSession.Id)
				Upcoming.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, Upcoming.Results[i].Id, configs.UserSession.Id)
				if Upcoming.Results[i].Title == "" {
					Upcoming.Results[i].Title = Upcoming.Results[i].Name
					Upcoming.Results[i].OTitle = Upcoming.Results[i].OName
				}
				Upcoming.Results[i].Url = slug.Make(Upcoming.Results[i].OTitle)
			}

			err = configs.Tmpl.ExecuteTemplate(w, "getupcoming", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   configs.SiteConfig.TitleHome,

				"Upcoming": Upcoming,
			})
			if err != nil {
				println("Template Error:", err)
				return
			}
		}

	}

}
