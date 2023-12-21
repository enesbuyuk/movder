package functions

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/enesbuyuk/movder/pkg/structures"
)

func BBToHtml(str string) string {
	r := strings.NewReplacer("[P]", "<p>",
		"[/P]", "</p>",
		"[P SMALL]", "<p class=\"section__text section__text--small\">",
		"[BR]", "<br>",
		"[B]", "<b>",
		"[/B]", "</b>",
		"[H2]", "<h2>",
		"[/H2]", "</h2>",
		"[H3]", "<h3>",
		"[/H3]", "</h3>",
		"[ALINTI]", "<blockquote>",
		"[/ALINTI]", "</blockquote>",
		"[KOD]", "<pre><code>",
		"[/KOD]", "</code></pre>",
		"[UL]", "<ul>",
		"[/UL]", "</ul>",
		"[LI]", "<li>",
		"[/LI]", "</li>",
		"[A][LINK]", "<a href=\"",
		"[/LINK]", "\" target=\"_blank\">",
		"[/A]", "</a>",
		"[FOTO]", `<figure class="kg-card kg-image-card kg-card-hascaption"><img src="`,
		"[/FOTO]", `" class="kg-image medium-zoom-image" alt="" loading="lazy"><figcaption>Photo by <a href="https://unsplash.com/@freestocks?utm_source=ghost&amp;utm_medium=referral&amp;utm_campaign=api-credit">freestocks.org</a> / <a href="https://unsplash.com/?utm_source=ghost&amp;utm_medium=referral&amp;utm_campaign=api-credit">Unsplash</a></figcaption></figure>`,
	)

	return r.Replace(str)
}

func SessionFunc(db *sql.DB, kullanici_id interface{}) structures.UserStruct {
	var UserSession structures.UserStruct
	User, _ := db.Query("SELECT kullanici_id,kullanici_nick, kullanici_isim, kullanici_url, kullanici_avatar, kullanici_yetki,kullanici_posta, kullanici_sifre,kullanici_hakkinda,kullanici_telephone FROM kullanici WHERE kullanici_id=?", kullanici_id)
	defer User.Close()
	if User.Next() {
		User.Scan(&UserSession.Id, &UserSession.Nick, &UserSession.Name, &UserSession.Link, &UserSession.Photo, &UserSession.Authority, &UserSession.Mail, &UserSession.Password, &UserSession.About, &UserSession.Phone)
	} else {
		UserSession.Id = 0
	}
	return UserSession
}

func TMDBApi(LastValue string, GetUrl string, Extra string) []byte {
	Url, _ := http.Get("https://api.themoviedb.org/" + GetUrl + "?api_key=" + LastValue + Extra)
	body, _ := ioutil.ReadAll(Url.Body)
	return body
}

func GetIpAddress(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	xforward := r.Header.Get("X-Forwarded-For")
	fmt.Println("Url:" + r.URL.Path + "<->IP : " + ip + " <-> X-Forwarded-For : " + xforward)
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CheckMyList(listType string, db *sql.DB, movieId int, userId int) bool {
	if listType == "watch" {
		query, _ := db.Query(`SELECT watchlist_film FROM watchlist WHERE watchlist_film=? AND watchlist_owner=? LIMIT 1`, movieId, userId)
		defer func(query *sql.Rows) {
			err := query.Close()
			if err != nil {
				log.Fatal("Query error:", err)
				return
			}
		}(query)
		if query.Next() {
			return true
		} else {
			return false
		}
	} else if listType == "watched" {
		query, _ := db.Query(`SELECT watchedlist_film FROM watchedlist WHERE watchedlist_film=? AND watchedlist_owner=? LIMIT 1`, movieId, userId)
		defer func(query *sql.Rows) {
			err := query.Close()
			if err != nil {
				log.Fatal("Query error:", err)
				return
			}
		}(query)
		if query.Next() {
			return true
		} else {
			return false
		}
	} else {
		return false
	}

}
