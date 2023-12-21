package structures

type M map[string]interface{}

type SettingStruct struct {
	SettingId    int
	SettingName  string
	SettingValue string
}

type UserStruct struct {
	Id           int
	Nick         string
	Name         string
	Mail         string
	Password     string
	Authority    string
	Url          string
	Link         string
	Photo        string
	About        string
	ArticleCount int
	Instagram    string
	Twitter      string
	Facebook     string
	Website      string
	Location     string
	Online       string
	Phone        string
	Reviewer     int
}

type FilmStruct struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Overview string `json:"overview"`

	Image    string `json:"poster_path"`
	Backdrop string `json:"backdrop_path"`

	Vote       float64 `json:"vote_average"`
	VoteCount  int     `json:"vote_count"`
	Date       string  `json:"release_date"`
	FirstAir   string  `json:"first_air_date"`
	Adult      bool    `json:"adult"`
	Genres     []int   `json:"genre_ids"`
	OTitle     string  `json:"original_title"`
	Url        string
	Name       string `json:"name"`
	OName      string `json:"original_name"`
	Budget     int    `json:"budget"`
	FilmBudget string
	Tagline    string `json:"tagline"`

	CheckWatchlist   bool
	CheckWatchedlist bool
	Diary            int
	FilmType         string

	Runtime       int `json:"runtime"`
	RuntimeString string
}

type PersonStruct struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ProfilePath  string `json:"profile_path"`
	Url          string
	OName        string `json:"original_name"`
	Department   string `json:"known_for_department"`
	Biography    string `json:"biography"`
	PlaceofBirth string `json:"place_of_birth"`
	Birthday     string `json:"birthday"`
	Deathday     string `json:"deathday"`
	Gender       int    `json:"gender"`
	Overview     string `json:"overview"`
}

type BlogStruct struct {
	BlogId       int
	BlogLink     string
	BlogTitle    string
	BlogOverview string
	BlogContent  string
	BlogOwner    int
	BlogDate     string
	BlogImage    string
	BlogTags     string
	BlogComment  string

	BlogAuthorAvatar string
	BlogAuthorNick   string
	BlogAuthorLink   string

	CategoryId    int
	CategoryLink  string
	CategoryName  string
	CategoryColor string
}

type CommentStruct struct {
	Id        int
	Date      string
	Blog      int
	Film      int
	FilmType  string
	Content   string
	Owner     UserStruct
	Like      int
	Dislike   int
	Link      string
	Title     string
	Rate      int
	RateRound float32
}

type ListStruct struct {
	Id       int
	Owner    string
	Title    string
	Url      string
	Count    int
	Favorite int
	Date     string
	Image    string
	Type     string

	OwnerNick  string
	OwnerUrl   string
	OwnerImage string
}

type MessageStruct struct {
	Message_id     int
	Owner          int
	Message_target UserStruct
	Content        string
	Date           string
}

type BoxStruct struct {
	Id     int
	Owner  int
	Count  int
	Target UserStruct
}

type NotificationStruct struct {
	Id      int
	Owner   UserStruct
	Target  UserStruct
	Content string
	View    int
}

type HomeFilmStruct struct {
	Page         int `json:"page"`
	Results      [10]FilmStruct
	TotalResults int `json:"total_results"`
}

type AllFilmstruct struct {
	Page    int `json:"page"`
	Results []FilmStruct
}

type LanguageJSONStruct struct {
	Language           string `json:"Language"`
	Home               string `json:"Home"`
	About              string `json:"About"`
	Contact            string `json:"Contact"`
	Films              string `json:"Films"`
	Social             string `json:"Social"`
	News               string `json:"News"`
	Signin             string `json:"Signin"`
	Signup             string `json:"Signup"`
	Forgotpassword     string `json:"Forgotpassword"`
	Privacy            string `json:"Privacy"`
	Extra              string `json:"Extra"`
	Popular            string `json:"Popular"`
	Dailytrendsmovie   string `json:"Dailytrendsmovie"`
	Dailytrendstv      string `json:"Dailytrendstv"`
	Upcoming           string `json:"Upcoming"`
	Readmore           string `json:"Readmore"`
	Latestfromblog     string `json:"Latestfromblog"`
	Lists              string `json:"Lists"`
	Watchlist          string `json:"Watchlist"`
	Watchedlist        string `json:"Watchedlist"`
	Advices            string `json:"Advices"`
	PopularLists       string `json:"PopularLists"`
	CreateListsMessage string `json:"CreateListsMessage"`
	LatestLists        string `json:"LatestLists"`
	Search             string `json:"Search"`
	Chat               string `json:"Chat"`
	StartMessagingNow  string `json:"StartMessagingNow"`
	Typeyourmessage    string `json:"Typeyourmessage"`
	CreateOwnList      string `json:"CreateOwnList"`
	Member             string `json:"Member"`
	PopularMovies      string `json:"PopularMovies"`
	PopularSeries      string `json:"PopularSeries"`
	PrivacyPolicy      string `json:"PrivacyPolicy"`
	PopularReviews     string `json:"PopularReviews"`
	LatestReviews      string `json:"LatestReviews"`
	LatestRates        string `json:"LatestRates"`
	LatestReviewers    string `json:"LatestReviewers"`
	PopularReviewers   string `json:"PopularReviewers"`
}

type ConfigStruct struct {
	DBhost     string `json:"dbhost"`
	DBport     string `json:"dbport"`
	DBname     string `json:"dbname"`
	DBuser     string `json:"dbuser"`
	DBpassword string `json:"dbpassword"`

	TMDBApiKey  string `json:"tmdbapikey"`
	TMDBApiLink string `json:"tmdbApiLink"`

	Port            string `json:"port"`
	Domain          string `json:"domain"`
	SiteUrl         string `json:"siteurl"`
	SessionKey      string `json:"sessionkey"`
	SessionName     string `json:"sessinname"`
	TitleLastPart   string `json:"titlelastpart"`
	TitleHome       string `json:"titlehome"`
	MetaDescription string `json:"metadescription"`
	MetaKeywords    string `json:"metakeywords"`
}
