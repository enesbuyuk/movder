package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gosimple/slug"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.Split(strings.TrimPrefix(r.URL.Path, "/@"), "/")
	if strings.HasPrefix(r.URL.Path, "/@") == true {

		query, err := configs.Db.Query("SELECT kullanici_id, kullanici_nick, kullanici_isim, kullanici_posta, kullanici_yetki, kullanici_url, kullanici_avatar, kullanici_hakkinda, kullanici_yazi, kullanici_instagram, kullanici_twitter, kullanici_facebook, kullanici_website, kullanici_location FROM kullanici WHERE kullanici_url=? LIMIT 1", getData[0])
		if err != nil {
			log.Fatal(err)
		}
		defer func(query *sql.Rows) {
			err := query.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(query)
		var MeProfile bool
		if query.Next() {
			var kullanici structures.UserStruct
			err := query.Scan(&kullanici.Id, &kullanici.Nick, &kullanici.Name, &kullanici.Mail, &kullanici.Authority, &kullanici.Link, &kullanici.Photo, &kullanici.About, &kullanici.ArticleCount, &kullanici.Instagram, &kullanici.Twitter, &kullanici.Facebook, &kullanici.Website, &kullanici.Location)
			if err != nil {
				return
			}
			kullanicihakkinda := kullanici.About

			LenWatchlistQuery, _ := configs.Db.Query(`SELECT COUNT(watchlist_id) AS count FROM watchlist WHERE watchlist_owner=?`, kullanici.Id)
			defer func(LenWatchlistQuery *sql.Rows) {
				err := LenWatchlistQuery.Close()
				if err != nil {

				}
			}(LenWatchlistQuery)
			LenWatchlist := 0
			if LenWatchlistQuery.Next() {
				err := LenWatchlistQuery.Scan(&LenWatchlist)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
			}

			LenWatchedlistQuery, _ := configs.Db.Query(`SELECT COUNT(watchedlist_id) AS count FROM watchedlist WHERE watchedlist_owner=?`, kullanici.Id)
			defer func(LenWatchedlistQuery *sql.Rows) {
				err := LenWatchedlistQuery.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(LenWatchedlistQuery)
			LenWatchedlist := 0
			if LenWatchedlistQuery.Next() {
				err := LenWatchedlistQuery.Scan(&LenWatchedlist)
				if err != nil {
					return
				}
			}

			LenDiaryQuery, _ := configs.Db.Query(`SELECT COUNT(watchedlist_id) AS count FROM watchedlist WHERE watchedlist_owner=? && watchedlist_diary=1`, kullanici.Id)
			defer func(LenDiaryQuery *sql.Rows) {
				err := LenDiaryQuery.Close()
				if err != nil {

				}
			}(LenDiaryQuery)
			LenDiary := 0
			if LenDiaryQuery.Next() {
				err := LenDiaryQuery.Scan(&LenDiary)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
			}

			LenFavoriteQuery, _ := configs.Db.Query(`SELECT COUNT(favorite_id) AS count FROM favorite WHERE favorite_owner=?`, kullanici.Id)
			defer func(LenFavoriteQuery *sql.Rows) {
				err := LenFavoriteQuery.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(LenFavoriteQuery)
			LenFavorite := 0
			if LenFavoriteQuery.Next() {
				err := LenFavoriteQuery.Scan(&LenFavorite)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
			}

			FollowersGet, _ := configs.Db.Query("SELECT kullanici.kullanici_nick,kullanici.kullanici_isim, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM follow WHERE follow_owner=?) as follow INNER JOIN kullanici ON follow.follow_user=kullanici.kullanici_id", kullanici.Id)
			defer func(FollowersGet *sql.Rows) {
				err := FollowersGet.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(FollowersGet)

			var Followers []structures.UserStruct
			var Follower structures.UserStruct
			for FollowersGet.Next() {
				err := FollowersGet.Scan(&Follower.Nick, &Follower.Name, &Follower.Link, &Follower.Photo)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				Followers = append(Followers, Follower)
			}

			FollowingsGet, _ := configs.Db.Query("SELECT kullanici.kullanici_nick,kullanici.kullanici_isim, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM follow WHERE follow_user=?) as follow INNER JOIN kullanici ON follow.follow_owner=kullanici.kullanici_id", kullanici.Id)
			defer func(FollowingsGet *sql.Rows) {
				err := FollowingsGet.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(FollowingsGet)

			var Followings []structures.UserStruct
			var Following structures.UserStruct
			for FollowingsGet.Next() {
				err := FollowingsGet.Scan(&Following.Nick, &Following.Name, &Following.Link, &Following.Photo)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				Followings = append(Followings, Following)
			}

			commentsQuery, _ := configs.Db.Query("SELECT Comment_id,Comment_Film,Comment_owner,Comment_content,Comment_date,Comment_rate,Comment_like,Comment_dislike,comment_filmtype FROM comment WHERE Comment_owner=? ORDER BY Comment_id DESC", kullanici.Id)
			defer func(commentsQuery *sql.Rows) {
				err := commentsQuery.Close()
				if err != nil {
					log.Fatal("Query error:", err)
				}
			}(commentsQuery)

			var comments []structures.CommentStruct
			var Comment structures.CommentStruct
			for commentsQuery.Next() {
				err := commentsQuery.Scan(&Comment.Id, &Comment.Film, &Comment.Owner.Id, &Comment.Content, &Comment.Date, &Comment.Rate, &Comment.Like, &Comment.Dislike, &Comment.FilmType)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				comments = append(comments, Comment)
			}
			for i := 0; i < len(comments); i++ {
				var Single structures.FilmStruct
				err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+comments[i].FilmType+"/"+strconv.Itoa(comments[i].Film), ""), &Single)

				if comments[i].FilmType == "movie" {
					if Single.Title == "" {
						comments[i].Title = Single.OTitle
					} else {
						comments[i].Title = Single.Title
					}
					comments[i].Link = slug.Make(Single.OTitle)
				} else if comments[i].FilmType == "tv" {
					if Single.Name == "" {
						comments[i].Title = Single.OName
					} else {
						comments[i].Title = Single.Name
					}
					comments[i].Link = slug.Make(Single.OName)
				}
			}

			MyWatchedget, _ := configs.Db.Query("SELECT watchedlist_film FROM watchedlist WHERE watchedlist_owner=?", configs.UserSession.Id)
			defer func(MyWatchedget *sql.Rows) {
				err := MyWatchedget.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(MyWatchedget)
			var MyWatcheds []structures.FilmStruct
			var MyWatched structures.FilmStruct
			for MyWatchedget.Next() {
				err := MyWatchedget.Scan(&MyWatched.Id)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				MyWatcheds = append(MyWatcheds, MyWatched)
			}

			if session.Values["SessionBool"] == true {
				if kullanici.Id == session.Values["UserId"] {
					MeProfile = true
				} else {
					MeProfile = false
				}
			} else {
				MeProfile = false
			}

			var Watcheds []structures.FilmStruct
			var WatchedsDiary []structures.FilmStruct
			var Watchlists []structures.FilmStruct
			var Favorites []structures.FilmStruct
			var LatestLists []structures.ListStruct

			if len(getData) > 1 && (getData[1] == "films" || getData[1] == "diary") {
				getWatchedQuery, _ := configs.Db.Query("SELECT watchedlist_film,watchedlist_id,watchedlist_date, watchedlist_owner,watchedlist_diary,watchedlist_filmtype FROM watchedlist WHERE watchedlist_owner=? ORDER BY watchedlist_id DESC", kullanici.Id)
				var Watched structures.FilmStruct
				defer func(getWatchedQuery *sql.Rows) {
					err := getWatchedQuery.Close()
					if err != nil {
						log.Fatal("Query error:", err)
					}
				}(getWatchedQuery)
				for getWatchedQuery.Next() {
					err := getWatchedQuery.Scan(&Watched.Id, &Watched.Url, &Watched.Url, &Watched.Url, &Watched.Diary, &Watched.FilmType)
					if err != nil {
						log.Fatal("Scan error:", err)
					}
					if Watched.Diary == 1 {
						WatchedsDiary = append(WatchedsDiary, Watched)
					}
					Watcheds = append(Watcheds, Watched)
				}
				for i := 0; i < len(Watcheds); i++ {
					Watcheds[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, Watcheds[i].Id, configs.UserSession.Id)
					Watcheds[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, Watcheds[i].Id, configs.UserSession.Id)
					var Single structures.FilmStruct
					err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+Watcheds[i].FilmType+"/"+strconv.Itoa(Watcheds[i].Id), ""), &Single)

					if Watcheds[i].FilmType == "movie" && Single.Title == "" {
						Single.Title = Single.OTitle
					}
					if Watcheds[i].FilmType == "tv" && Single.Name == "" {
						if Single.Name == "" {
							Single.Title = Single.OName
						} else {
							Single.Title = Single.Name
						}
					}
					Watcheds[i].Title = Single.Title
					Watcheds[i].OTitle = Single.OTitle
					Watcheds[i].Image = Single.Image
					Watcheds[i].Vote = Single.Vote
					Watcheds[i].Date = Single.Date
					Watcheds[i].Url = slug.Make(Single.OTitle)
				}
				for i := 0; i < len(WatchedsDiary); i++ {
					WatchedsDiary[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, WatchedsDiary[i].Id, configs.UserSession.Id)
					WatchedsDiary[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, WatchedsDiary[i].Id, configs.UserSession.Id)

					var WatchedFilm structures.FilmStruct
					err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/"+strconv.Itoa(WatchedsDiary[i].Id), ""), &WatchedFilm)

					if WatchedsDiary[i].FilmType == "movie" && WatchedFilm.Title == "" {
						WatchedFilm.Title = WatchedFilm.OTitle
					}
					if WatchedsDiary[i].FilmType == "tv" && WatchedFilm.Name == "" {
						if WatchedFilm.Name == "" {
							WatchedFilm.Title = WatchedFilm.OName
						} else {
							WatchedFilm.Title = WatchedFilm.Name
						}
					}
					WatchedsDiary[i].Title = WatchedFilm.Title
					WatchedsDiary[i].OTitle = WatchedFilm.OTitle
					WatchedsDiary[i].Image = WatchedFilm.Image
					WatchedsDiary[i].Vote = WatchedFilm.Vote
					WatchedsDiary[i].Date = WatchedFilm.Date
					WatchedsDiary[i].Url = slug.Make(WatchedFilm.OTitle)
				}
			} else if len(getData) > 1 && getData[1] == "watchlist" {
				WatchlistQuery, _ := configs.Db.Query("SELECT watchlist_film,watchlist_date, watchlist_filmtype FROM watchlist WHERE watchlist_owner=? ORDER BY watchlist_id DESC", kullanici.Id)
				defer func(WatchlistQuery *sql.Rows) {
					err := WatchlistQuery.Close()
					if err != nil {
						log.Fatal("Query error:", err)
					}
				}(WatchlistQuery)
				var Watchlist structures.FilmStruct
				for WatchlistQuery.Next() {
					err := WatchlistQuery.Scan(&Watchlist.Id, &Watchlist.Url, &Watchlist.FilmType)
					if err != nil {
						log.Fatal("Scan error:", err)
					}
					Watchlists = append(Watchlists, Watchlist)
				}
				for i := 0; i < len(Watchlists); i++ {
					Watchlists[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, WatchedsDiary[i].Id, configs.UserSession.Id)
					Watchlists[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, WatchedsDiary[i].Id, configs.UserSession.Id)
					var Single structures.FilmStruct
					err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+Watchlists[i].FilmType+"/"+strconv.Itoa(Watchlists[i].Id), ""), &Single)

					if Watchlists[i].FilmType == "movie" && Single.Title == "" {
						Single.Title = Single.OTitle
					}
					if Watchlists[i].FilmType == "tv" && Single.Name == "" {
						if Single.Name == "" {
							Single.Title = Single.OName
						} else {
							Single.Title = Single.Name
						}
					}
					Watchlists[i].Title = Single.Title
					Watchlists[i].OTitle = Single.OTitle
					Watchlists[i].Image = Single.Image
					Watchlists[i].Vote = Single.Vote
					Watchlists[i].Date = Single.Date
					Watchlists[i].Url = slug.Make(Single.OTitle)
				}
			} else if len(getData) > 1 && getData[1] == "favorite" {
				FavoriteQuery, _ := configs.Db.Query("SELECT favorite_film,favorite_date, favorite_filmtype  FROM favorite WHERE favorite_owner=? ORDER BY favorite_id DESC", kullanici.Id)
				defer func(FavoriteQuery *sql.Rows) {
					err := FavoriteQuery.Close()
					if err != nil {
						log.Fatal("Query error:", err)
					}
				}(FavoriteQuery)
				var Favorite structures.FilmStruct
				for FavoriteQuery.Next() {
					err := FavoriteQuery.Scan(&Favorite.Id, &Favorite.Url, &Favorite.FilmType)
					if err != nil {
						log.Fatal("Scan error:", err)
					}
					Favorites = append(Favorites, Favorite)
				}
				for i := 0; i < len(Favorites); i++ {
					Favorites[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, Favorites[i].Id, configs.UserSession.Id)
					Favorites[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, Favorites[i].Id, configs.UserSession.Id)
					var Single structures.FilmStruct
					err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+Favorites[i].FilmType+"/"+strconv.Itoa(Favorites[i].Id), ""), &Single)

					if Favorites[i].FilmType == "movie" && Single.Title == "" {
						Single.Title = Single.OTitle
					}
					if Favorites[i].FilmType == "tv" && Single.Name == "" {
						if Single.Name == "" {
							Single.Title = Single.OName
						} else {
							Single.Title = Single.Name
						}
					}
					Favorites[i].Title = Single.Title
					Favorites[i].OTitle = Single.OTitle
					Favorites[i].Image = Single.Image
					Favorites[i].Vote = Single.Vote
					Favorites[i].Date = Single.Date
					Favorites[i].Url = slug.Make(Single.Title)
				}
			} else if len(getData) > 1 && getData[1] == "lists" {
				ListType := "everyone"
				var LatestListsGet *sql.Rows
				if MeProfile == true {
					LatestListsGet, _ = configs.Db.Query("SELECT list.list_id, list.list_owner, list.list_title, list.list_url, list.list_count, list.list_favorite, list.list_date, list.list_image, list.list_type, kullanici.kullanici_nick, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM list WHERE list_owner=?) AS list INNER JOIN kullanici ON kullanici.kullanici_id=list.list_owner ORDER BY list.list_id DESC", kullanici.Id)
					defer func(LatestListsGet *sql.Rows) {
						err := LatestListsGet.Close()
						if err != nil {
							println("Query error:", err)
						}
					}(LatestListsGet)
				} else {
					LatestListsGet, _ = configs.Db.Query("SELECT list.list_id, list.list_owner, list.list_title, list.list_url, list.list_count, list.list_favorite, list.list_date, list.list_image, list.list_type, kullanici.kullanici_nick, kullanici.kullanici_url, kullanici.kullanici_avatar FROM (SELECT * FROM list WHERE list_owner=? AND list_type=?) AS list INNER JOIN kullanici ON kullanici.kullanici_id=list.list_owner ORDER BY list.list_id DESC", kullanici.Id, ListType)
					defer func(LatestListsGet *sql.Rows) {
						err := LatestListsGet.Close()
						if err != nil {
							println("Query error:", err)
						}
					}(LatestListsGet)
				}
				var LatestList structures.ListStruct
				for LatestListsGet.Next() {
					err := LatestListsGet.Scan(&LatestList.Id, &LatestList.Owner, &LatestList.Title, &LatestList.Url, &LatestList.Count, &LatestList.Favorite, &LatestList.Date, &LatestList.Image, &LatestList.Type, &LatestList.OwnerNick, &LatestList.OwnerUrl, &LatestList.OwnerImage)
					if err != nil {
						log.Fatal("Scan error:", err)
					}
					LatestLists = append(LatestLists, LatestList)
				}
			} else {
				//http.Redirect(w, r, "/@"+kullanici.Link, 302)
			}
			followget, _ := configs.Db.Query(`SELECT * FROM follow WHERE follow_user=? AND follow_owner=?`, configs.UserSession.Id, kullanici.Id)
			defer func(followget *sql.Rows) {
				err := followget.Close()
				if err != nil {
					println("Query error:", err)
					return
				}
			}(followget)
			var FollowState bool
			if followget.Next() {
				FollowState = true
			} else {
				FollowState = false
			}

			//SIMILAR

			mat := 0
			var i, j int
			for i = 0; i < len(MyWatcheds); i++ {
				for j = 0; j < len(Watcheds); j++ {
					if MyWatcheds[i].Id == Watcheds[j].Id {
						mat++
					}
				}
			}
			var SimilarPercent float64
			if len(MyWatcheds) == 0 {
				SimilarPercent = 0
			} else {
				SimilarPercent = math.Floor(float64(mat) / float64(len(MyWatcheds)) * 100)
			}

			errT := configs.Tmpl.ExecuteTemplate(w, "profile", map[string]interface{}{
				"Language":    configs.Language,
				"SessionBool": session.Values["SessionBool"],
				"Session":     configs.UserSession,
				"SiteConfig":  configs.SiteConfig,
				"TitlePage":   kullanici.Nick + " (" + kullanici.Name + ")" + configs.SiteConfig.TitleLastPart,

				"MeProfile":         MeProfile,
				"FollowState":       FollowState,
				"kullanici":         kullanici,
				"KullaniciHakkinda": kullanicihakkinda,
				"Watched":           Watcheds,
				"Diary":             WatchedsDiary,
				"SimilarPercent":    SimilarPercent,
				"Watchlist":         Watchlists,
				"Favorite":          Favorites,
				"Comments":          comments,
				"Followers":         Followers,
				"Followings":        Followings,
				"LenWatchlist":      LenWatchlist,
				"LenWatchedlist":    LenWatchedlist,
				"LenDiary":          LenDiary,
				"LenFavorite":       LenFavorite,
				"LatestLists":       LatestLists,
			})
			if errT != nil {
				log.Fatal("Execute template error:", err)
			}
		} else {
			http.Redirect(w, r, "/404", 302)
		}
	} else {
		//HOME

		LatestrateQuery, _ := configs.Db.Query("SELECT rate_film,rate_owner,rate_date,rate_filmtype,rate_value,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url from rate INNER JOIN kullanici ON rate.rate_owner=kullanici.kullanici_id ORDER BY rate.rate_id DESC LIMIT 6")
		defer func(LatestrateQuery *sql.Rows) {
			err := LatestrateQuery.Close()
			if err != nil {
				println("Query error:", err)
				return
			}
		}(LatestrateQuery)
		var Rates []structures.CommentStruct
		var Rate structures.CommentStruct
		for LatestrateQuery.Next() {
			err := LatestrateQuery.Scan(&Rate.Film, &Rate.Owner.Id, &Rate.Date, &Rate.FilmType, &Rate.Rate, &Rate.Owner.Nick, &Rate.Owner.Photo, &Rate.Owner.Link)
			if err != nil {
				log.Fatal("Scan error:", err)
				return
			}
			Rate.RateRound = float32(Rate.Rate) / 2

			var Single structures.FilmStruct
			err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+Rate.FilmType+"/"+strconv.Itoa(Rate.Film), ""), &Single)

			if Rate.FilmType == "movie" {
				if Single.Title == "" {
					Rate.Title = Single.OTitle
				} else {
					Rate.Title = Single.Title
				}
				Rate.Link = slug.Make(Single.OTitle)
			} else if Rate.FilmType == "tv" {
				if Single.Name == "" {
					Rate.Title = Single.OName
				} else {
					Rate.Title = Single.Name
				}
				Rate.Link = slug.Make(Single.OName)
			}

			Rates = append(Rates, Rate)
		}

		CommentsGet, _ := configs.Db.Query("SELECT comment_film,comment_owner,comment_content,comment_date,comment_filmtype,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url from (SELECT * FROM comment WHERE comment_date BETWEEN DATE_SUB( CURDATE( ) ,INTERVAL 10 DAY ) AND CURDATE( )) as comment INNER JOIN kullanici ON comment.comment_owner=kullanici.kullanici_id ORDER BY comment.comment_like DESC LIMIT 6")
		defer func(CommentsGet *sql.Rows) {
			err := CommentsGet.Close()
			if err != nil {
				println("Query error:", err)
				return
			}
		}(CommentsGet)
		var comments []structures.CommentStruct
		var Comment structures.CommentStruct
		for CommentsGet.Next() {
			err := CommentsGet.Scan(&Comment.Film, &Comment.Owner.Id, &Comment.Content, &Comment.Date, &Comment.FilmType, &Comment.Owner.Nick, &Comment.Owner.Photo, &Comment.Owner.Link)
			if err != nil {
				log.Fatal("Scan error:", err)
				return
			}

			rateQuery, _ := configs.Db.Query("SELECT rate_value from rate wHERE rate_owner=? && rate_film=? && rate_filmtype=?", Comment.Owner.Id, Comment.Film, Comment.FilmType)
			func(rateQuery *sql.Rows) {
				err := rateQuery.Close()
				if err != nil {
					println("Query error:", err)
					return
				}
			}(rateQuery)
			if rateQuery.Next() {
				err := rateQuery.Scan(&Comment.Rate)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				Comment.RateRound = float32(Comment.Rate) / 2
			} else {
				Comment.Rate = 0
				Comment.RateRound = 0
			}

			var Single structures.FilmStruct
			err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+Comment.FilmType+"/"+strconv.Itoa(Comment.Film), ""), &Single)

			if Comment.FilmType == "movie" {
				if Single.Title == "" {
					Comment.Title = Single.OTitle
				} else {
					Comment.Title = Single.Title
				}
				Comment.Link = slug.Make(Single.OTitle)
			} else if Comment.FilmType == "tv" {
				if Single.Name == "" {
					Comment.Title = Single.OName
				} else {
					Comment.Title = Single.Name
				}
				Comment.Link = slug.Make(Single.OName)
			}

			comments = append(comments, Comment)
		}

		LatestCommentsQuery, _ := configs.Db.Query("SELECT comment_film,comment_owner,comment_content,comment_date,comment_filmtype,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url from comment INNER JOIN kullanici ON comment.comment_owner=kullanici.kullanici_id ORDER BY comment.comment_id DESC LIMIT 6")
		defer func(LatestCommentsQuery *sql.Rows) {
			err := LatestCommentsQuery.Close()
			if err != nil {
				println("Query error:", err)
				return
			}
		}(LatestCommentsQuery)
		var LatestComments []structures.CommentStruct
		var LatestComment structures.CommentStruct
		for LatestCommentsQuery.Next() {
			err := LatestCommentsQuery.Scan(&LatestComment.Film, &LatestComment.Owner.Id, &LatestComment.Content, &LatestComment.Date, &LatestComment.FilmType, &LatestComment.Owner.Nick, &LatestComment.Owner.Photo, &LatestComment.Owner.Link)
			if err != nil {
				println("Query error:", err)
				return
			}

			rateQuery, _ := configs.Db.Query("SELECT rate_value from rate wHERE rate_owner=? && rate_film=? && rate_filmtype=?", LatestComment.Owner.Id, LatestComment.Film, LatestComment.FilmType)
			func(rateQuery *sql.Rows) {
				err := rateQuery.Close()
				if err != nil {
					println("Query error:", err)
					return
				}
			}(rateQuery)
			if rateQuery.Next() {
				err := rateQuery.Scan(&LatestComment.Rate)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				LatestComment.RateRound = float32(LatestComment.Rate) / 2
			} else {
				LatestComment.Rate = 0
				LatestComment.RateRound = 0
			}

			var Single structures.FilmStruct
			err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+LatestComment.FilmType+"/"+strconv.Itoa(LatestComment.Film), ""), &Single)

			if LatestComment.FilmType == "movie" {
				if Single.Title == "" {
					LatestComment.Title = Single.OTitle
				} else {
					LatestComment.Title = Single.Title
				}
				LatestComment.Link = slug.Make(Single.OTitle)
			} else if LatestComment.FilmType == "tv" {
				if Single.Name == "" {
					LatestComment.Title = Single.OName
				} else {
					LatestComment.Title = Single.Name
				}
				LatestComment.Link = slug.Make(Single.OName)
			}

			LatestComments = append(LatestComments, LatestComment)
		}

		LatestReviewersQuery, _ := configs.Db.Query("SELECT comment_film,comment_owner,comment_content,comment_date,comment_filmtype,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url from comment INNER JOIN (SELECT * FROM kullanici WHERE kullanici_reviewer=1) as kullanici ON comment.comment_owner=kullanici.kullanici_id ORDER BY comment.comment_id DESC LIMIT 6")
		defer func(LatestReviewersQuery *sql.Rows) {
			err := LatestReviewersQuery.Close()
			if err != nil {
				log.Fatal("Query error:", err)
				return
			}
		}(LatestReviewersQuery)
		var LatestReviewers []structures.CommentStruct
		var LatestReviewer structures.CommentStruct
		for LatestReviewersQuery.Next() {
			err := LatestReviewersQuery.Scan(&LatestReviewer.Film, &LatestReviewer.Owner.Id, &LatestReviewer.Content, &LatestReviewer.Date, &LatestReviewer.FilmType, &LatestReviewer.Owner.Nick, &LatestReviewer.Owner.Photo, &LatestReviewer.Owner.Link)
			if err != nil {
				println("Scan Error:", err)
				return
			}

			rateQuery, _ := configs.Db.Query("SELECT rate_value from rate wHERE rate_owner=? && rate_film=? && rate_filmtype=?", LatestReviewer.Owner.Id, LatestReviewer.Film, LatestReviewer.FilmType)
			func(rateQuery *sql.Rows) {
				err := rateQuery.Close()
				if err != nil {
					log.Fatal("Query error:", err)
					return
				}
			}(rateQuery)
			if rateQuery.Next() {
				err := rateQuery.Scan(&LatestReviewer.Rate)
				if err != nil {
					log.Fatal("Scan error:", err)
					return
				}
				LatestReviewer.RateRound = float32(LatestReviewer.Rate) / 2
			} else {
				LatestReviewer.Rate = 0
				LatestReviewer.RateRound = 0
			}

			var Single structures.FilmStruct
			err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/"+LatestReviewer.FilmType+"/"+strconv.Itoa(LatestReviewer.Film), ""), &Single)

			if LatestReviewer.FilmType == "movie" {
				if Single.Title == "" {
					LatestReviewer.Title = Single.OTitle
				} else {
					LatestReviewer.Title = Single.Title
				}
				LatestReviewer.Link = slug.Make(Single.OTitle)
			} else if LatestReviewer.FilmType == "tv" {
				if Single.Name == "" {
					LatestReviewer.Title = Single.OName
				} else {
					LatestReviewer.Title = Single.Name
				}
				LatestReviewer.Link = slug.Make(Single.OName)
			}

			LatestReviewers = append(LatestReviewers, LatestReviewer)
		}

		BlogsGet, _ := configs.Db.Query("SELECT blog.blog_id, blog.blog_url, blog.blog_title, blog.blog_overview, blog.blog_content, blog.blog_owner, blog.blog_date, blog.blog_image,  kullanici.kullanici_nick,kullanici.kullanici_url,category.category_id,category.category_url,category.category_name,category.category_color FROM (SELECT * FROM blog ORDER BY blog_popular DESC LIMIT 2) AS blog INNER JOIN kullanici ON kullanici.kullanici_id=blog.blog_owner INNER JOIN category_blog category ON category.category_id=blog.blog_category")
		defer func(BlogsGet *sql.Rows) {
			err := BlogsGet.Close()
			if err != nil {
				log.Fatal("Query error:", err)
				return
			}
		}(BlogsGet)
		var Blogs []structures.BlogStruct
		var Blog structures.BlogStruct
		for BlogsGet.Next() {
			err := BlogsGet.Scan(&Blog.BlogId, &Blog.BlogLink, &Blog.BlogTitle, &Blog.BlogOverview, &Blog.BlogContent, &Blog.BlogOwner, &Blog.BlogDate, &Blog.BlogImage, &Blog.BlogAuthorNick, &Blog.BlogAuthorLink, &Blog.CategoryId, &Blog.CategoryLink, &Blog.CategoryName, &Blog.CategoryColor)
			if err != nil {
				println("Query error:", err)
				return
			}
			Blogs = append(Blogs, Blog)
		}

		Blogs2Get, _ := configs.Db.Query(`SELECT blog.blog_id, blog.blog_url, blog.blog_title, blog.blog_overview, blog.blog_content, blog.blog_owner, blog.blog_date, blog.blog_image, kullanici.kullanici_nick,kullanici.kullanici_url,category.category_id,category.category_url,category.category_name,category.category_color FROM (SELECT * FROM blog ORDER BY blog_popular DESC LIMIT 2,2) AS blog INNER JOIN kullanici ON kullanici.kullanici_id=blog.blog_owner INNER JOIN category_blog category ON category.category_id=blog.blog_category`)
		defer func(Blogs2Get *sql.Rows) {
			err := Blogs2Get.Close()
			if err != nil {
				log.Fatal("Query error:", err)
				return
			}
		}(Blogs2Get)
		var Blogs2 []structures.BlogStruct
		var Blog2 structures.BlogStruct
		for Blogs2Get.Next() {
			err := Blogs2Get.Scan(&Blog2.BlogId, &Blog2.BlogLink, &Blog2.BlogTitle, &Blog2.BlogOverview, &Blog2.BlogContent, &Blog2.BlogOwner, &Blog2.BlogDate, &Blog2.BlogImage, &Blog2.BlogAuthorNick, &Blog2.BlogAuthorLink, &Blog2.CategoryId, &Blog2.CategoryLink, &Blog2.CategoryName, &Blog2.CategoryColor)
			if err != nil {
				println("Query error:", err)
				return
			}
			Blogs2 = append(Blogs2, Blog2)
		}

		var PopularMovies structures.HomeFilmStruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/popular", ""), &PopularMovies)

		var PopularSeries structures.HomeFilmStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/tv/popular", ""), &PopularSeries)

		var TrendMovie structures.HomeFilmStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/movie/day", ""), &TrendMovie)

		var TrendTV structures.HomeFilmStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/tv/day", ""), &TrendTV)

		var Upcomings structures.HomeFilmStruct
		err = json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/upcoming", ""), &Upcomings)

		for i := 0; i < len(PopularMovies.Results); i++ {
			PopularMovies.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
			PopularMovies.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
			if PopularMovies.Results[i].Title == "" {
				PopularMovies.Results[i].Title = PopularMovies.Results[i].Name
				PopularMovies.Results[i].OTitle = PopularMovies.Results[i].OName
			}
			PopularMovies.Results[i].Url = slug.Make(PopularMovies.Results[i].OTitle)
		}

		for i := 0; i < len(PopularSeries.Results); i++ {
			PopularSeries.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
			PopularSeries.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
			if PopularSeries.Results[i].Title == "" {
				PopularSeries.Results[i].Title = PopularSeries.Results[i].Name
				PopularSeries.Results[i].OTitle = PopularSeries.Results[i].OName
			}
			PopularSeries.Results[i].Url = slug.Make(PopularSeries.Results[i].OTitle)
			if PopularSeries.Results[i].Date == "" {
				PopularSeries.Results[i].Date = PopularSeries.Results[i].FirstAir
			}
		}

		for i := 0; i < len(TrendMovie.Results); i++ {
			TrendMovie.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendMovie.Results[i].Id, configs.UserSession.Id)
			TrendMovie.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendMovie.Results[i].Id, configs.UserSession.Id)
			if TrendMovie.Results[i].Title == "" {
				TrendMovie.Results[i].Title = TrendMovie.Results[i].Name
				TrendMovie.Results[i].OTitle = TrendMovie.Results[i].OName
			}
			TrendMovie.Results[i].Url = slug.Make(TrendMovie.Results[i].OTitle)
			if TrendMovie.Results[i].Date == "" {
				TrendMovie.Results[i].Date = TrendMovie.Results[i].FirstAir
			}
		}

		for i := 0; i < len(TrendTV.Results); i++ {
			TrendTV.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendTV.Results[i].Id, configs.UserSession.Id)
			TrendTV.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendTV.Results[i].Id, configs.UserSession.Id)
			if TrendTV.Results[i].Title == "" {
				TrendTV.Results[i].Title = TrendTV.Results[i].Name
				TrendTV.Results[i].OTitle = TrendTV.Results[i].OName
			}
			TrendTV.Results[i].Url = slug.Make(TrendTV.Results[i].OTitle)
			if TrendTV.Results[i].Date == "" {
				TrendTV.Results[i].Date = TrendTV.Results[i].FirstAir
			}
		}

		for i := 0; i < len(Upcomings.Results); i++ {
			Upcomings.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, Upcomings.Results[i].Id, configs.UserSession.Id)
			Upcomings.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, Upcomings.Results[i].Id, configs.UserSession.Id)
			if Upcomings.Results[i].Title == "" {
				Upcomings.Results[i].Title = Upcomings.Results[i].Name
				Upcomings.Results[i].OTitle = Upcomings.Results[i].OName
			}
			Upcomings.Results[i].Url = slug.Make(Upcomings.Results[i].OTitle)
			if Upcomings.Results[i].Date == "" {
				Upcomings.Results[i].Date = Upcomings.Results[i].FirstAir
			}
		}

		err = configs.Tmpl.ExecuteTemplate(w, "home", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.SiteConfig.TitleHome,

			"LatestReviewers": LatestReviewers,
			"comments":        comments,
			"LatestComments":  LatestComments,
			"Rates":           Rates,
			"Blogs":           Blogs,
			"Blogs2":          Blogs2,
			"PopularMovies":   PopularMovies,
			"PopularSeries":   PopularSeries,
			"TrendMovie":      TrendMovie,
			"TrendTV":         TrendTV,
			"Upcomings":       Upcomings,
		})
		if err != nil {
			log.Fatal("Execute template error:", err)
		}
	}
}
