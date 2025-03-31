package storage

var CreateShortURLTable = `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        short_url TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL
    );`

var GetURL = `
	SELECT url 
	FROM urls 
	WHERE short_url = $1
`