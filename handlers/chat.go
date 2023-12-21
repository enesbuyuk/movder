package handlers

import (
	"database/sql"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"log"
	"net/http"
	"strings"
	"time"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
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
		getData := strings.TrimPrefix(r.URL.Path, "/chat/")

		var Messages []structures.MessageStruct
		var Message structures.MessageStruct

		ClientIdKey := ""
		var ChosedBox structures.BoxStruct
		MessageCount := 0
		if len(getData) != 0 {
			TargetIDGet, _ := configs.Db.Query(`SELECT kullanici_id FROM kullanici WHERE kullanici_url=? LIMIT 1`, getData)
			defer func(TargetIDGet *sql.Rows) {
				err := TargetIDGet.Close()
				if err != nil {
					log.Fatal("Query error:", err)
				}
			}(TargetIDGet)

			if TargetIDGet.Next() {
				var targetId int
				err := TargetIDGet.Scan(&targetId)
				if err != nil {
					log.Fatal("Scan error:", err)
				}

				if targetId != configs.UserSession.Id {
					ChosedBoxGET, _ := configs.Db.Query(`SELECT kullanici_id,kullanici_nick,kullanici_url,kullanici_avatar,kullanici_online kullanici_id FROM kullanici WHERE kullanici_id=? LIMIT 1`, targetId)
					defer func(ChosedBoxGET *sql.Rows) {
						err := ChosedBoxGET.Close()
						if err != nil {
							log.Fatal("Query error:", err)
						}
					}(ChosedBoxGET)
					//		ChosedBoxGET, _ := configs.Db.Query("SELECT kullanici_id,kullanici.kullanici_nick,kullanici.kullanici_url,kullanici.kullanici_avatar,kullanici.kullanici_online FROM (select * FROM message WHERE (message_owner=? && message_target=?) || (message_owner=? && message_target=?) ORDER BY message_id DESC LIMIT 1) as message INNER JOIN kullanici ON kullanici.kullanici_id=message.message_target",configs.UserSession.Id, targetId,targetId,configs.UserSession.Id)
					if ChosedBoxGET.Next() {
						err := ChosedBoxGET.Scan(&ChosedBox.Target.Id, &ChosedBox.Target.Nick, &ChosedBox.Target.Link, &ChosedBox.Target.Photo, &ChosedBox.Target.Online)
						if err != nil {
							log.Fatal("Scan error:", err)
						}
					}

					MessageGet, _ := configs.Db.Query("SELECT message_id,message_owner,message_content,message_date FROM message WHERE ( (message_owner=? && message_target=?) || (message_owner=? && message_target=?)) AND message_delete=0 ORDER BY message_id ASC", ChosedBox.Target.Id, configs.UserSession.Id, configs.UserSession.Id, ChosedBox.Target.Id)
					defer func(MessageGet *sql.Rows) {
						err := MessageGet.Close()
						if err != nil {
							log.Fatal("Query error:", err)
						}
					}(MessageGet)
					for MessageGet.Next() {
						MessageCount++
						err := MessageGet.Scan(&Message.Message_id, &Message.Owner, &Message.Content, &Message.Date)
						if err != nil {
							log.Fatal("Scan error:", err)
						}
						Messages = append(Messages, Message)
					}

					for i := 0; i < len(Messages); i++ {
						t, _ := time.Parse("2006-01-02 15:04:05", Messages[i].Date)
						Messages[i].Date = t.Format("15:04 - 2 January 2006")
					}

					ClientGet, _ := configs.Db.Query("SELECT client_key FROM client_chat WHERE ((client_owner=? && client_target=?) || (client_owner=? && client_target=?)) LIMIT 1", ChosedBox.Target.Id, configs.UserSession.Id, configs.UserSession.Id, ChosedBox.Target.Id)
					defer func(ClientGet *sql.Rows) {
						err := ClientGet.Close()
						if err != nil {
							log.Fatal("Fprint error:", err)
						}
					}(ClientGet)
					if ClientGet.Next() {
						err := ClientGet.Scan(&ClientIdKey)
						if err != nil {
							log.Fatal("Scan error:", err)
						}
					} else {
						const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
						ClientIdKey = functions.StringWithCharset(40, charset)
						_, err := configs.Db.Exec("INSERT INTO client_chat(client_key,client_owner,client_target) VALUES (?,?,?);", ClientIdKey, configs.UserSession.Id, targetId)
						if err != nil {
							log.Fatal("Query error:", err)
						}
						_, err = configs.Db.Exec("INSERT INTO client_chat(client_key,client_owner,client_target) VALUES (?,?,?);", ClientIdKey, targetId, configs.UserSession.Id)
						if err != nil {
							log.Fatal("Query error:", err)
						}
					}

				} else {
					http.Redirect(w, r, "/chat", 302)
				}
			} else {
				http.Redirect(w, r, "/chat", 302)
			}
		}

		BoxGet, _ := configs.Db.Query("SELECT kullanici.kullanici_Id,kullanici.kullanici_nick,kullanici.kullanici_isim,kullanici.kullanici_url,kullanici.kullanici_avatar,kullanici.kullanici_online FROM (select * FROM client_chat WHERE client_owner=? ORDER BY client_id DESC) as message INNER JOIN kullanici ON kullanici.kullanici_id=message.client_target", configs.UserSession.Id)
		defer func(BoxGet *sql.Rows) {
			err := BoxGet.Close()
			if err != nil {
				log.Fatal("Query error:", err)
			}
		}(BoxGet)

		//	BoxGet, _ := configs.Db.Query("SELECT box_target,kullanici.kullanici_nick,kullanici.kullanici_isim,kullanici.kullanici_url,kullanici.kullanici_avatar,kullanici.kullanici_online FROM (SELECT * FROM message_box WHERE box_owner=? ORDER BY box_id DESC) as box INNER JOIN kullanici ON box.box_target=kullanici.kullanici_id",configs.UserSession.Id)
		var Boxes []structures.UserStruct
		var Box structures.UserStruct

		for BoxGet.Next() {
			err := BoxGet.Scan(&Box.Id, &Box.Nick, &Box.Name, &Box.Link, &Box.Photo, &Box.Online)
			if err != nil {
				log.Fatal("Scan error:", err)
			}
			if Box.Id != configs.UserSession.Id {
				Boxes = append(Boxes, Box)
			}
		}

		err := configs.Tmpl.ExecuteTemplate(w, "chat", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.Chat + configs.SiteConfig.TitleLastPart,

			"Boxes":        Boxes,
			"Messages":     Messages,
			"ChosedBox":    ChosedBox,
			"ClientIdKey":  ClientIdKey,
			"MessageCount": MessageCount,
		})
		if err != nil {
			println("Template Error:", err)
			return
		}
	} else {
		http.Redirect(w, r, "/signin", 302)
	}
}
