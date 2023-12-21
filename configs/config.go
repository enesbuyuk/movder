// configs/config.go
package configs

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/enesbuyuk/movder/pkg/structures"
	"github.com/gorilla/sessions"
)

var (
	once           sync.Once
	SiteConfig     structures.ConfigStruct
	Language       structures.LanguageJSONStruct
	LanguageEN     structures.LanguageJSONStruct
	LanguageTR     structures.LanguageJSONStruct
	UserSession    structures.UserStruct
	ConfigFile     *os.File
	LanguageENjson *os.File
	LanguageTRjson *os.File
	Db             *sql.DB
	Stores         = sessions.NewCookieStore([]byte("sadasdasd"))
)
var Tmpl *template.Template

func InitConfig() {
	once.Do(func() {

		ConfigFileName := "config-linux.json"
		if runtime.GOOS == "windows" {
			ConfigFileName = "config-windows.json"
		}

		var cfgErr error
		ConfigFile, cfgErr = os.Open("configs/" + ConfigFileName)
		if cfgErr != nil {
			log.Println("Config Error:", cfgErr)
			return
		}
		ConfigByte, _ := io.ReadAll(ConfigFile)
		err := json.Unmarshal(ConfigByte, &SiteConfig)
		if err != nil {
			log.Println("Unmarshall Error:", err)
			return
		}

		LanguageENjson, cfgErr = os.Open("configs/languageen.json")
		if cfgErr != nil {
			log.Println("Config Error:", cfgErr)
			return
		}
		LanguageENbyte, _ := io.ReadAll(LanguageENjson)
		err = json.Unmarshal(LanguageENbyte, &LanguageEN)
		if err != nil {
			log.Println("Unmarshall Error:", err)
			return
		}

		LanguageTRjson, cfgErr = os.Open("configs/languagetr.json")
		if cfgErr != nil {
			log.Println("Config Error:", cfgErr)
			return
		}
		LanguageTRbyte, _ := io.ReadAll(LanguageTRjson)
		err = json.Unmarshal(LanguageTRbyte, &LanguageTR)
		if err != nil {
			log.Println("Unmarshall Error:", err)
			return
		}

		Db, err = sql.Open("mysql", SiteConfig.DBuser+":"+SiteConfig.DBpassword+"@tcp("+SiteConfig.DBhost+":"+SiteConfig.DBport+")/"+SiteConfig.DBname)
		if err != nil {
			log.Println("SQL Connection Err!")
			log.Print(err.Error())
		}
		Stores.Options = &sessions.Options{
			Domain:   SiteConfig.Domain,
			Path:     "/",
			MaxAge:   86400 * 30,
			HttpOnly: true,
		}
	})

	var templateErr error
	Tmpl, templateErr = template.ParseGlob("templates/*.gohtml")
	if templateErr != nil {
		log.Fatal(templateErr)
	}
}

func Get(req *http.Request) (*sessions.Session, error) {
	return Stores.Get(req, "default-session-name")
}

func GetNamed(req *http.Request, name string) (*sessions.Session, error) {
	return Stores.Get(req, name)
}
