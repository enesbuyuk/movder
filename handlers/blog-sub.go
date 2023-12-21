package handlers

import (
	"database/sql"
	"github.com/enesbuyuk/movder/configs"
	"github.com/enesbuyuk/movder/pkg/functions"
	"github.com/enesbuyuk/movder/pkg/structures"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func BlogSubHandler(w http.ResponseWriter, r *http.Request) {
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

	getData := strings.TrimPrefix(r.URL.Path, "/blog/")

	Result, _ := configs.Db.Query("SELECT blog.blog_id, blog.blog_url, blog.blog_title, blog.blog_overview, blog.blog_content, blog.blog_owner, blog.blog_date, blog.blog_image,blog.blog_tags,kullanici.kullanici_nick,kullanici.kullanici_isim,kullanici.kullanici_url,kullanici.kullanici_avatar,kullanici.kullanici_hakkinda,category.category_id,category.category_url,category.category_name,category.category_color FROM (SELECT * FROM blog WHERE blog_url=?) as blog INNER JOIN kullanici ON kullanici.kullanici_id=blog.blog_owner INNER JOIN category_blog category ON category.category_id=blog.blog_category", getData)
	defer func(Result *sql.Rows) {
		err := Result.Close()
		if err != nil {
			println("Query Error:", err)
			return
		}
	}(Result)
	if Result.Next() {
		var Blog structures.BlogStruct
		var Owner structures.UserStruct
		err := Result.Scan(&Blog.BlogId, &Blog.BlogLink, &Blog.BlogTitle, &Blog.BlogOverview, &Blog.BlogContent, &Blog.BlogOwner, &Blog.BlogDate, &Blog.BlogImage, &Blog.BlogTags, &Owner.Nick, &Owner.Name, &Owner.Link, &Owner.Photo, &Owner.About, &Blog.CategoryId, &Blog.CategoryLink, &Blog.CategoryName, &Blog.CategoryColor)
		if err != nil {
			println("Scan Error:", err)
			return
		}
		BlogContent := template.HTML(functions.BBToHtml(Blog.BlogContent))
		BlogTags := strings.Split(Blog.BlogTags, ",")

		BlogsGet, _ := configs.Db.Query("SELECT blog.blog_id, blog.blog_url, blog.blog_title, blog.blog_overview, blog.blog_content, blog.blog_owner, blog.blog_date, blog.blog_image, kullanici.kullanici_nick,kullanici.kullanici_url FROM (SELECT * FROM blog ORDER BY blog_popular DESC LIMIT 10) AS blog INNER JOIN kullanici ON kullanici.kullanici_id=blog.blog_owner")
		defer func(BlogsGet *sql.Rows) {
			err := BlogsGet.Close()
			if err != nil {
				println("Query Error:", err)
				return
			}
		}(BlogsGet)
		var Blogs []structures.BlogStruct
		var LatestBlog structures.BlogStruct
		for BlogsGet.Next() {
			err := BlogsGet.Scan(&LatestBlog.BlogId, &LatestBlog.BlogLink, &LatestBlog.BlogTitle, &LatestBlog.BlogOverview, &LatestBlog.BlogContent, &LatestBlog.BlogOwner, &LatestBlog.BlogDate, &LatestBlog.BlogImage, &LatestBlog.BlogAuthorNick, &LatestBlog.BlogAuthorLink)
			if err != nil {
				log.Fatal("Scan error:", err)
			}
			Blogs = append(Blogs, LatestBlog)
		}

		commentsGet, _ := configs.Db.Query("SELECT comment_id,comment_blog,comment_owner,comment_content,comment_date,comment_rate,kullanici.kullanici_nick,kullanici.kullanici_avatar,kullanici.kullanici_url from (SELECT * FROM comment_blog WHERE comment_blog=? ORDER BY comment_id DESC) as comment INNER JOIN (select * from kullanici) as kullanici ON comment.comment_owner=kullanici.kullanici_id", Blog.BlogId)
		defer func(commentsGet *sql.Rows) {
			err := commentsGet.Close()
			if err != nil {
				log.Println("Query Error:", err)
				return
			}
		}(commentsGet)
		var comments []structures.CommentStruct
		var Comment structures.CommentStruct

		CommentCount := 0
		for commentsGet.Next() {
			CommentCount++
			err := commentsGet.Scan(&Comment.Id, &Comment.Blog, &Comment.Owner.Id, &Comment.Content, &Comment.Date, &Comment.Rate, &Comment.Owner.Nick, &Comment.Owner.Photo, &Comment.Owner.Link)
			if err != nil {
				log.Fatal("Scan error:", err)
			}
			comments = append(comments, Comment)
		}

		errT := configs.Tmpl.ExecuteTemplate(w, "blog_page", map[string]interface{}{
			"Language":    configs.Language,
			"SessionBool": session.Values["SessionBool"],
			"Session":     configs.UserSession,
			"SiteConfig":  configs.SiteConfig,
			"TitlePage":   Blog.BlogTitle + configs.SiteConfig.TitleLastPart,

			"MetaAciklama": Blog.BlogOverview,
			"MetaKelime":   Blog.BlogTags,
			"Blog":         Blog,
			"BlogContent":  BlogContent,
			"Owner":        Owner,
			"Blogs":        Blogs,
			"Comment":      comments,
			"CommentCount": CommentCount,
			"BlogTags":     BlogTags,
		})
		if errT != nil {
			log.Fatal("Execute template error:", err)
		}
	} else {
		http.Redirect(w, r, "/404", 302)
	}
}
