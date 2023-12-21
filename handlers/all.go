package handlers

import (
	"encoding/json"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gosimple/slug"
	"net/http"
	"strings"
)

func AllHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.TrimPrefix(r.URL.Path, "/all/")
	if getData == "popular-movies" {
		var PopularMovies structures.AllFilmstruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/popular", ""), &PopularMovies)

		for i := 0; i < len(PopularMovies.Results); i++ {
			PopularMovies.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
			PopularMovies.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularMovies.Results[i].Id, configs.UserSession.Id)
			if PopularMovies.Results[i].Title == "" {
				PopularMovies.Results[i].Title = PopularMovies.Results[i].Name
				PopularMovies.Results[i].OTitle = PopularMovies.Results[i].OName
			}
			PopularMovies.Results[i].Url = slug.Make(PopularMovies.Results[i].OTitle)
		}

		err = configs.Tmpl.ExecuteTemplate(w, "popular-movies", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.PopularMovies + configs.SiteConfig.TitleLastPart,

			"PopularMovies": PopularMovies,
		})
		if err != nil {
			println("Template Error:", err)
			return
		}
	} else if getData == "popular-series" {
		var PopularSeries structures.AllFilmstruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/tv/popular", ""), &PopularSeries)

		for i := 0; i < len(PopularSeries.Results); i++ {
			PopularSeries.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
			PopularSeries.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, PopularSeries.Results[i].Id, configs.UserSession.Id)
			if PopularSeries.Results[i].Title == "" {
				PopularSeries.Results[i].Title = PopularSeries.Results[i].Name
				PopularSeries.Results[i].OTitle = PopularSeries.Results[i].OName
			}
			PopularSeries.Results[i].Url = slug.Make(PopularSeries.Results[i].OTitle)
		}

		err = configs.Tmpl.ExecuteTemplate(w, "popular-series", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.PopularSeries + configs.SiteConfig.TitleLastPart,

			"PopularSeries": PopularSeries,
		})
		if err != nil {
			println("Template Error:", err)
			return
		}
	} else if getData == "trend-movies" {
		var TrendMovies structures.AllFilmstruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/movie/day", ""), &TrendMovies)

		for i := 0; i < len(TrendMovies.Results); i++ {
			TrendMovies.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendMovies.Results[i].Id, configs.UserSession.Id)
			TrendMovies.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendMovies.Results[i].Id, configs.UserSession.Id)
			if TrendMovies.Results[i].Title == "" {
				TrendMovies.Results[i].Title = TrendMovies.Results[i].Name
				TrendMovies.Results[i].OTitle = TrendMovies.Results[i].OName
			}
			TrendMovies.Results[i].Url = slug.Make(TrendMovies.Results[i].OTitle)
		}

		err = configs.Tmpl.ExecuteTemplate(w, "trend-movies", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.Dailytrendsmovie + configs.SiteConfig.TitleLastPart,

			"TrendMovies": TrendMovies,
		})
		if err != nil {
			println("Template Error:", err)
			return
		}
	} else if getData == "trend-series" {
		var TrendSeries structures.AllFilmstruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/trending/tv/day", ""), &TrendSeries)

		for i := 0; i < len(TrendSeries.Results); i++ {
			TrendSeries.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, TrendSeries.Results[i].Id, configs.UserSession.Id)
			TrendSeries.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, TrendSeries.Results[i].Id, configs.UserSession.Id)
			if TrendSeries.Results[i].Title == "" {
				TrendSeries.Results[i].Title = TrendSeries.Results[i].Name
				TrendSeries.Results[i].OTitle = TrendSeries.Results[i].OName
			}
			TrendSeries.Results[i].Url = slug.Make(TrendSeries.Results[i].OTitle)
		}

		err = configs.Tmpl.ExecuteTemplate(w, "trend-series", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.Dailytrendstv + configs.SiteConfig.TitleLastPart,

			"TrendSeries": TrendSeries,
		})
		if err != nil {
			println("Template Error:", err)
			return
		}
	} else if getData == "upcoming" {
		var Upcoming structures.AllFilmstruct
		err := json.Unmarshal(functions.TMDBApi(configs.SiteConfig.TMDBApiLink, "3/movie/upcoming", ""), &Upcoming)

		for i := 0; i < len(Upcoming.Results); i++ {
			Upcoming.Results[i].CheckWatchlist = functions.CheckMyList("watch", configs.Db, Upcoming.Results[i].Id, configs.UserSession.Id)
			Upcoming.Results[i].CheckWatchedlist = functions.CheckMyList("watched", configs.Db, Upcoming.Results[i].Id, configs.UserSession.Id)
			if Upcoming.Results[i].Title == "" {
				Upcoming.Results[i].Title = Upcoming.Results[i].Name
				Upcoming.Results[i].OTitle = Upcoming.Results[i].OName
			}
			Upcoming.Results[i].Url = slug.Make(Upcoming.Results[i].OTitle)
		}

		err = configs.Tmpl.ExecuteTemplate(w, "upcoming", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.Upcoming + configs.SiteConfig.TitleLastPart,

			"Upcoming": Upcoming,
		})
		if err != nil {
			return
		}
	} else {
		err := configs.Tmpl.ExecuteTemplate(w, "all", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   configs.Language.Films + configs.SiteConfig.TitleLastPart,
		})
		if err != nil {
			return
		}
	}
}
