package handlers

import (
	"encoding/json"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gosimple/slug"
	"io"
	"log"
	"net/http"
	"net/url"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := url.QueryEscape(r.FormValue("q"))

	var getData2 string
	if r.FormValue("in") == "movies" {
		getData2 = "movie"
	} else if r.FormValue("in") == "people" {
		getData2 = "person"
	} else if r.FormValue("in") == "tvshow" {
		getData2 = "tv"
	} else {
		getData2 = "movie"
	}

	type Genres []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	type Genress struct {
		Genres Genres
	}

	var Genre Genress
	err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/genre/movie/list", ""), &Genre)

	type PersonNested struct {
		Page    int `json:"page"`
		Results []structures.PersonStruct
	}

	UpcomingUrl, _ := http.Get("https://api.themoviedb.org/3/search/" + getData2 + "?api_key=" + configs.SiteConfig.TMDBApiLink + "&query=" + getData)
	BodyUpcoming, _ := io.ReadAll(UpcomingUrl.Body)
	var UpcomingFilm structures.AllFilmstruct
	var People PersonNested

	err = json.Unmarshal(BodyUpcoming, &UpcomingFilm)
	if err != nil {
		return
	}
	err = json.Unmarshal(BodyUpcoming, &People)
	if err != nil {
		return
	}

	for i := 0; i < len(UpcomingFilm.Results); i++ {
		UpcomingFilm.Results[i].FilmType = getData2
		UpcomingFilm.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, UpcomingFilm.Results[i].Id, configs.UserSession.Id)
		UpcomingFilm.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, UpcomingFilm.Results[i].Id, configs.UserSession.Id)
		if UpcomingFilm.Results[i].Title == "" {
			UpcomingFilm.Results[i].Title = UpcomingFilm.Results[i].Name
			UpcomingFilm.Results[i].OTitle = UpcomingFilm.Results[i].OName
		}
		UpcomingFilm.Results[i].Url = slug.Make(UpcomingFilm.Results[i].OTitle)
		if UpcomingFilm.Results[i].Date == "" {
			UpcomingFilm.Results[i].Date = UpcomingFilm.Results[i].FirstAir
		}
	}

	for i := 0; i < len(People.Results); i++ {
		People.Results[i].Url = slug.Make(People.Results[i].Name)
	}

	err = configs.Tmpl.ExecuteTemplate(w, "search", map[string]interface{}{
		"Language":    configs.Language,
		"SessionBool": session.Values["SessionBool"],
		"Session":     configs.UserSession,
		"SiteConfig":  configs.SiteConfig,
		"TitlePage":   configs.Language.Search + configs.SiteConfig.TitleLastPart,

		"Genre":    Genre,
		"Upcoming": UpcomingFilm,
		"People":   People,
		"getdata":  getData,
		"getdata2": getData2,
	})
	if err != nil {
		log.Fatal("Execute template error:", err)
	}
}
