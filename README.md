# MOVDER
<center><img src="https://github.com/enesbuyuk/movder/assets/82279640/d3f308c8-4c3e-40c6-b5f7-a0afff6580ff"></center>

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
![image](https://github.com/enesbuyuk/movder/assets/82279640/59ad36d8-0d6b-4e51-9051-995a52183882)

### Profile Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/22ba8187-6b8d-4030-833b-80d3a1737e6d)

### Chat Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/5f7b99d8-fe4d-426b-8c94-e600c7c91022)

### Films Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/8d9f877a-da08-4446-af3f-9c1b36b19bde)

### Films - Popular Films Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/da834e62-1c72-42ec-9f8f-e1f499231fbf)
 
### Films - Popular Series Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/be5de66e-d2ae-4b09-8f22-b2b634fa0143)

### Films - Trending in Movies Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/5b0d0e87-01f3-45fc-9fe9-d00d63c68016)

### Films - Trending in Series Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/e4a42328-1e14-4a72-91a3-612b63dab9c8)

### Films - Upcoming Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/09214084-0c12-48d8-8d4d-d384496381f6)

## Movie Page
![image](https://github.com/enesbuyuk/movder/assets/82279640/c3f01456-0e3b-4d78-9b38-c148c133c99f)


## Blog Homepage
![image](https://github.com/enesbuyuk/movder/assets/82279640/5a502029-7122-421f-833d-206286eed2d7)

## Contributing
- Use GitHub Issues to report bugs or make suggestions.
- Feel free to submit pull requests for code improvements.
