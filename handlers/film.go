package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gosimple/slug"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FilmHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.Split(strings.TrimPrefix(r.URL.Path, "/film/"), "/")

	var MovieUrl structures.FilmStruct
	err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/"+getData[0], ""), &MovieUrl)

	var TvUrl structures.FilmStruct
	err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/tv/"+getData[0], ""), &TvUrl)

	if err != nil {
		return
	}

	if MovieUrl.Title == "" {
		MovieUrl.Title = MovieUrl.Name
		MovieUrl.OTitle = MovieUrl.OName
	}

	if MovieUrl.Date == "" {
		MovieUrl.Date = MovieUrl.FirstAir
	}

	if TvUrl.Title == "" {
		TvUrl.Title = TvUrl.Name
		TvUrl.OTitle = TvUrl.OName

	}

	if TvUrl.Date == "" {
		TvUrl.Date = TvUrl.FirstAir
	}

	MovieUrl.Url = slug.Make(MovieUrl.OTitle)
	TvUrl.Url = slug.Make(TvUrl.OTitle)

	var FilmType string
	var Film structures.FilmStruct

	if getData[1] == MovieUrl.Url {
		FilmType = "movie"
		Film = MovieUrl
	} else if getData[1] == TvUrl.Url {
		FilmType = "tv"
		Film = TvUrl
	}

	if strconv.Itoa(Film.Id) == "0" {
		http.Redirect(w, r, "/404", 302)
	} else {
		p := message.NewPrinter(language.English)
		Film.FilmBudget = p.Sprintf("%d", Film.Budget)

		if Film.Image == "" {
			Film.Image = Film.Backdrop
		}

		Film.Vote = math.Round(Film.Vote*10) / 10
		Film.FilmType = FilmType

		TotalHour := Film.Runtime / 60
		TotalMin := Film.Runtime - TotalHour*60
		Film.RuntimeString = strconv.Itoa(TotalHour) + " hr " + strconv.Itoa(TotalMin) + " mins"

		var rate int
		if session.Values["SessionBool"] == true {
			rateQuery, _ := configs.Db.Query(`SELECT rate_value FROM rate WHERE rate_film=? AND rate_filmtype=? AND rate_owner=?`, Film.Id, Film.FilmType, configs.UserSession.Id)
			defer func(rateQuery *sql.Rows) {
				err := rateQuery.Close()
				if err != nil {
					log.Println("Query Error:", err)
				}
			}(rateQuery)
			if rateQuery.Next() {
				err := rateQuery.Scan(&rate)
				if err != nil {
					log.Fatal("Scan error:", err)
				}
			} else {
				rate = 0
			}
		}

		type CastStruct struct {
			Id         int    `json:"id"`
			Name       string `json:"name"`
			OName      string `json:"original_name"`
			Image      string `json:"profile_path"`
			Department string `json:"known_for_department"`
			Job        string `json:"job"`
			Character  string `json:"character"`
			Url        string
		}
		type CreditStruct struct {
			Id   int `json:"id"`
			Cast []CastStruct
			Crew []CastStruct
		}
		type VideoStruct struct {
			Name     string `json:"name"`
			Site     string `json:"site"`
			Key      string `json:"key"`
			Official bool   `json:"official"`
			Date     string `json:"published_at"`
			Type     string `json:"type"`
		}
		type VideoResults struct {
			Results []VideoStruct `json:"results"`
		}

		type ProviderRentBuy struct {
			Logo string `json:"logo_path"`
			Name string `json:"provider_name"`
		}

		type ProviderStruct struct {
			Results struct {
				UnitedState struct {
					Link string            `json:"link"`
					Rent []ProviderRentBuy `json:"rent"`
					Buy  []ProviderRentBuy `json:"buy"`
				} `json:"US"`

				Turkiye struct {
					Link string            `json:"link"`
					Rent []ProviderRentBuy `json:"rent"`
					Buy  []ProviderRentBuy `json:"buy"`
				} `json:"TR"`
			} `json:"results"`
		}

		var Provider ProviderStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+FilmType+"/"+getData[0]+"/watch/providers", ""), &Provider)

		log.Println(len(Provider.Results.Turkiye.Link))

		//VIDEOS
		var Video VideoResults
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+FilmType+"/"+getData[0]+"/videos", ""), &Video)
		for i := 0; i < len(Video.Results); i++ {
			if Video.Results[i].Date != "" {
				t, _ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", Video.Results[i].Date)
				Video.Results[i].Date = t.Format("15:04 - 2 January 2006")
			}
		}

		var Credit CreditStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+FilmType+"/"+getData[0]+"/credits", ""), &Credit)

		for i := 0; i < len(Credit.Cast); i++ {
			Credit.Cast[i].Url = slug.Make(Credit.Cast[i].Name)
		}

		for i := 0; i < len(Credit.Crew); i++ {
			Credit.Crew[i].Url = slug.Make(Credit.Crew[i].Name)
		}

		var similarFilm structures.HomeFilmStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+FilmType+"/"+strconv.Itoa(Film.Id)+"/recommendations", ""), &similarFilm)
		for i := 0; i < len(similarFilm.Results); i++ {
			similarFilm.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, similarFilm.Results[i].Id, configs.UserSession.Id)
			similarFilm.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, similarFilm.Results[i].Id, configs.UserSession.Id)
			if similarFilm.Results[i].Title == "" {
				similarFilm.Results[i].Title = similarFilm.Results[i].OTitle
				similarFilm.Results[i].Title = similarFilm.Results[i].Name
				similarFilm.Results[i].OTitle = similarFilm.Results[i].OName
			}
			if similarFilm.Results[i].Date == "" {
				similarFilm.Results[i].Date = similarFilm.Results[i].FirstAir
			}
			similarFilm.Results[i].Url = slug.Make(similarFilm.Results[i].OTitle)
		}

		CommentsGet, _ := configs.Db.Query("SELECT comment_id,comment_film,comment_owner,comment_content,comment_date,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url,kullanici.kullanici_reviewer from (SELECT * FROM comment WHERE comment_film=? && comment_filmtype=? ORDER BY comment_id DESC) as comment INNER JOIN (select * from kullanici) as kullanici ON comment.comment_owner=kullanici.kullanici_id", Film.Id, Film.FilmType)
		defer func(CommentsGet *sql.Rows) {
			err := CommentsGet.Close()
			if err != nil {
				println("Query Error:", err)
				return
			}
		}(CommentsGet)
		var Comments []structures.CommentStruct
		var Comment structures.CommentStruct

		for CommentsGet.Next() {
			err := CommentsGet.Scan(&Comment.Id, &Comment.Film, &Comment.Owner.Id, &Comment.Content, &Comment.Date, &Comment.Owner.Nick, &Comment.Owner.Photo, &Comment.Owner.Link, &Comment.Owner.Reviewer)
			if err != nil {
				println("Scan Error:", err)
				return
			}

			rateQuery, _ := configs.Db.Query("SELECT rate_value from rate wHERE rate_owner=? && rate_film=? && rate_filmtype=? ", Comment.Owner.Id, Film.Id, Film.FilmType)
			if rateQuery.Next() {
				err := rateQuery.Scan(&Comment.Rate)
				if err != nil {
					log.Fatal("Scan error:", err)
				}
				Comment.RateRound = float32(Comment.Rate) / 2
			} else {
				Comment.Rate = 0
				Comment.RateRound = 0
			}
			Comments = append(Comments, Comment)
		}

		watchlist := 0
		watchlistsql, err := configs.Db.Query(`SELECT watchlist_film FROM watchlist WHERE watchlist_film=? AND watchlist_owner=?`, Film.Id, session.Values["UserId"])
		defer func(watchlistsql *sql.Rows) {
			err := watchlistsql.Close()
			if err != nil {
				log.Fatal("Query error:", err)
			}
		}(watchlistsql)
		if watchlistsql.Next() {
			watchlist = 1
		} else {
			watchlist = 0
		}

		watchedlist := 0
		diary := 0
		watchedlistsql, err := configs.Db.Query(`SELECT watchedlist_diary FROM watchedlist WHERE watchedlist_film=? AND watchedlist_owner=?`, Film.Id, session.Values["UserId"])
		defer func(watchedlistsql *sql.Rows) {
			err := watchedlistsql.Close()
			if err != nil {
				log.Fatal("Query error:", err)
			}
		}(watchedlistsql)
		if watchedlistsql.Next() {
			watchedlist = 1
			err := watchedlistsql.Scan(&diary)
			if err != nil {
				log.Fatal("Query error:", err)
			}

		} else {
			watchedlist = 0
		}

		Favoritelist := 0
		Favoritelistsql, err := configs.Db.Query(`SELECT favorite_film FROM favorite WHERE favorite_film=? AND favorite_owner=?`, Film.Id, session.Values["UserId"])
		defer func(Favoritelistsql *sql.Rows) {
			err := Favoritelistsql.Close()
			if err != nil {
				log.Println("Query Error:", err)
				return
			}
		}(Favoritelistsql)
		if Favoritelistsql.Next() {
			Favoritelist = 1
		} else {
			Favoritelist = 0
		}

		err = configs.Tmpl.ExecuteTemplate(w, "film", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   Film.Title + configs.SiteConfig.TitleLastPart,

			"Film":         Film,
			"Similar":      similarFilm,
			"Comments":     Comments,
			"Watchlist":    watchlist,
			"Watchedlist":  watchedlist,
			"Diary":        diary,
			"Favoritelist": Favoritelist,
			"Credit":       Credit,
			"Rate":         rate,
			"RateCount":    float32(rate) / 2,

			"Video":    Video.Results,
			"Provider": Provider.Results,
		})
		if err != nil {
			log.Fatal("Execute template error:", err)
		}
	}
}
