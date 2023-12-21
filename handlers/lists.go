package handlers

import (
	"database/sql"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"log"
	"net/http"
)

func ListsHandler(w http.ResponseWriter, r *http.Request) {
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

	PopularListsGet, _ := configs.Db.Query("SELECT list.list_id, list.list_owner, list.list_title, list.list_url, list.list_count, list.list_favorite, list.list_date, list.list_image, list.list_type, kullanici.kullanici_nick, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM list WHERE list_type='everyone') AS list INNER JOIN kullanici ON kullanici.kullanici_id=list.list_owner ORDER BY list.list_favorite DESC LIMIT 12")
	defer func(PopularListsGet *sql.Rows) {
		err := PopularListsGet.Close()
		if err != nil {
			println("Query Error:", err)
		}
	}(PopularListsGet)
	var PopularLists []structures.ListStruct
	var PopularList structures.ListStruct
	for PopularListsGet.Next() {
		err := PopularListsGet.Scan(&PopularList.Id, &PopularList.Owner, &PopularList.Title, &PopularList.Url, &PopularList.Count, &PopularList.Favorite, &PopularList.Date, &PopularList.Image, &PopularList.Type, &PopularList.OwnerNick, &PopularList.OwnerUrl, &PopularList.OwnerImage)
		if err != nil {
			log.Fatal("Scan error:", err)
		}
		PopularLists = append(PopularLists, PopularList)
	}

	LatestListsGet, _ := configs.Db.Query("SELECT list.list_id, list.list_owner, list.list_title, list.list_url, list.list_count, list.list_favorite, list.list_date, list.list_image, list.list_type, kullanici.kullanici_nick, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM list WHERE list_type='everyone') AS list INNER JOIN kullanici ON kullanici.kullanici_id=list.list_owner ORDER BY list.list_id DESC LIMIT 12")
	defer func(LatestListsGet *sql.Rows) {
		err := LatestListsGet.Close()
		if err != nil {
			log.Fatal("Query error:", err)
		}
	}(LatestListsGet)
	var LatestLists []structures.ListStruct
	var LatestList structures.ListStruct
	for LatestListsGet.Next() {
		err := LatestListsGet.Scan(&LatestList.Id, &LatestList.Owner, &LatestList.Title, &LatestList.Url, &LatestList.Count, &LatestList.Favorite, &LatestList.Date, &LatestList.Image, &LatestList.Type, &LatestList.OwnerNick, &LatestList.OwnerUrl, &LatestList.OwnerImage)
		if err != nil {
			log.Fatal("Scan error:", err)
		}
		LatestLists = append(LatestLists, LatestList)
	}

	err := configs.Tmpl.ExecuteTemplate(w, "lists", map[string]interface{}{
		"Language":    configs.Language,
		"SessionBool": session.Values["SessionBool"],
		"Session":     configs.UserSession,
		"SiteConfig":  configs.SiteConfig,
		"TitlePage":   configs.Language.Lists + configs.SiteConfig.TitleLastPart,

		"PopularLists": PopularLists,
		"LatestLists":  LatestLists,
	})
	if err != nil {
		log.Fatal("Execute template error:", err)
	}
}
