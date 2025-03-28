package storage

var CreateShortURLTable = `
	CREATE TABLE IF NOT EXISTS urls (
		id integer,
		short_url text,
		url text
	);`

var GetURL = `
	SELECT short_url, url 
	FROM urls 
	WHERE short_url = $1"
`