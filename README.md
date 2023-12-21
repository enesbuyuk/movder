# MOVDER
<center><img src="https://github.com/enesbuyuk/movder/assets/82279640/97a96f01-10ba-4e0a-830b-1a8040fbaca8"></center>

## Overview
A social media movie website written with Go programming language using TMDB API.

## Features
- **Movie List:** View the most popular, latest, and category-wise movies.
- **Movie Details:** Discover detailed information about selected movies.
- **User Profiles:** Users can save the movies they watch and customize their profile pages.
- **Real-time Messaging:** Connect with other users through real-time chat functionality.
- **Multilingual Support:** Users can choose and use the language they want. (For now, TR and ENG)
- **Logging System:** Basic logging to track important events.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/enesbuyuk/movder.git
   cd movder
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```
   
3. Enter your data for you in `configs/config-windows.json` & `configs/config-linux.json`. (Note: `tmdbApiLink` value must be empty.)
   ```json
   {
	"dbhost" 		: "DB_HOST",
	"dbport" 		: "DB_PORT",
	"dbname" 		: "DB_NAME",
	"dbuser" 		: "DB_USER" ,
	"dbpassword" 		: "DB_PASSWORD", 

	"tmdbapikey"		: "TMDB_APIKEY",
	"tmdbApiLink"		: "", 

	"port"			: ":80",
	"domain"		: " YOUR_DOMAIN",
   	"siteurl"		: "https://YOUR_DOMAIN/",
   	"sessionkey"		: "YOUR_SESSION_KEY",
   	"sessinname"		: "YOUR_SESSION_NAME",
   	"titlelastpart"		: " | MOVDER",
   	"titlehome"		: "MOVDER",
   	"metadescription"	: "This is description.",
   	"metakeywords"		: "movie, movies, series, film, imdb, tmdb, movie social media, social media, tv, tv show, movies social media"
   }
   ```
   
5. Create database like in `database/movder.sql`:

6. Start or compile the application:
   ```bash
   go build main.go; ./main
   ```

   ```bash
   go run main.go
   ```

## Screenshots

### Homepage
![image](https://github.com/enesbuyuk/movder/assets/82279640/6d9cb226-fb92-45c9-a282-527ece8ecb11)

### Profile Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/d159be20-f4e9-4d03-8490-30902c70d19d)

### Chat Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/4d42a396-4ebc-4925-957a-dd73ac7d27fc)

### Films Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/bd6877c5-fde6-4f15-a517-96d2a3261a83)

### Films - Popular Movies Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/13f9e83c-49c6-40dd-a803-d12d3bef0e7c)
 
### Films - Popular Series Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/d207bed6-8da0-48cc-b626-66182846af23)

### Films - Trending in Movies Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/e4178932-3c91-4172-b3b8-38608a3c9a71)

### Films - Trending in Series Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/808d7059-14c7-4f08-a008-562f5cabf778)

### Films - Upcoming Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/8664552f-c4b4-4812-90fa-1e612b5d1e24)

### Movie Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/88bd3575-29ca-4676-8f60-a6fe5919d4b7)

### Blog Homepage
![image](https://github.com/enesbuyuk/movder/assets/82279640/85427794-b4bd-4169-9f2f-99a0a0d3fb7c)

### Lists Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/f32b17a9-e470-459d-9801-4dc2f714b7aa)


## Contributing
- Use GitHub Issues to report bugs or make suggestions.
- Feel free to submit pull requests for code improvements.
