package handlers

import (
	"database/sql"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"log"
	"net/http"
)

func BlogHandler(w http.ResponseWriter, r *http.Request) {
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
	BlogsGet, _ := configs.Db.Query("SELECT blog.blog_id, blog.blog_url, blog.blog_title, blog.blog_overview, blog.blog_content, blog.blog_owner, blog.blog_date, blog.blog_image, blog.blog_comment, kullanici.kullanici_nick,kullanici.kullanici_url,category.category_id,category.category_url,category.category_name,category.category_color FROM (SELECT * FROM blog ORDER BY blog_id DESC) AS blog INNER JOIN kullanici ON kullanici.kullanici_id=blog.blog_owner INNER JOIN category_blog category ON category.category_id=blog.blog_category")
	defer func(BlogsGet *sql.Rows) {
		err := BlogsGet.Close()
		if err != nil {
			println("Query Error:", err)
			return
		}
	}(BlogsGet)
	var Blogs []structures.BlogStruct
	var Blog structures.BlogStruct
	for BlogsGet.Next() {
		err := BlogsGet.Scan(&Blog.BlogId, &Blog.BlogLink, &Blog.BlogTitle, &Blog.BlogOverview, &Blog.BlogContent, &Blog.BlogOwner, &Blog.BlogDate, &Blog.BlogImage, &Blog.BlogComment, &Blog.BlogAuthorNick, &Blog.BlogAuthorLink, &Blog.CategoryId, &Blog.CategoryLink, &Blog.CategoryName, &Blog.CategoryColor)
		if err != nil {
			println("Scan Error:", err)
			return
		}
		Blogs = append(Blogs, Blog)
	}

	err := configs.Tmpl.ExecuteTemplate(w, "blog", map[string]interface{}{
		"Language":    configs.Language,
		"SessionBool": session.Values["SessionBool"],
		"Session":     configs.UserSession,
		"SiteConfig":  configs.SiteConfig,
		"TitlePage":   "Blog" + configs.SiteConfig.TitleLastPart,

		"MetaAciklama": Blog.BlogOverview,
		"MetaKelime":   Blog.BlogTags,
		"Blogs":        Blogs,
	})
	if err != nil {
		log.Fatal("Execute template error:", err)
	}
}
