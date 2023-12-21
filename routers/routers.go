package routers

import (
	"fmt"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/handlers"
	"log"
	"net/http"
	"path/filepath"
)

func StartServer() {
	configs.InitConfig()
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(NeuteredFileSystem{http.Dir("static/")})))

	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/all/", handlers.AllHandler)
	mux.HandleFunc("/lists", handlers.ListsHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/blog", handlers.BlogHandler)
	mux.HandleFunc("/blog/", handlers.BlogSubHandler)
	mux.HandleFunc("/film/", handlers.FilmHandler)
	mux.HandleFunc("/person/", handlers.PersonHandler)
	mux.HandleFunc("/signin", handlers.SignInHandler)
	mux.HandleFunc("/signup", handlers.SignUpHandler)
	mux.HandleFunc("/signout", handlers.SignOutHandler)
	mux.HandleFunc("/backend/", handlers.BackendHandler)
	mux.HandleFunc("/about", handlers.BackendHandler)
	mux.HandleFunc("/privacy", handlers.BackendHandler)
	mux.HandleFunc("/faq", handlers.BackendHandler)
	mux.HandleFunc("/chat/", handlers.ChatHandler)
	mux.HandleFunc("/404", handlers.Error404Handler)
	mux.HandleFunc("/panel", handlers.AdminHomeHandler)

	//WebSocket
	mux.HandleFunc("/ws", handlers.WebsocketHandler)

	fmt.Println("Web server is started on ", configs.SiteConfig.Port)
	err := http.ListenAndServe(configs.SiteConfig.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

type NeuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
