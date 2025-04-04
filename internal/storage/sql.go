package storage

var CreateShortURLTable = `
	CREATE TABLE IF NOT EXISTS urls (
		id integer,
		short_url text,
		url text PRIMARY KEY
	);`

var GetURL = `
	SELECT url 
	FROM urls 
	WHERE short_url = $1
`

var SaveURL = `
	INSERT into urls (id, short_url, url)
	VALUES ($1, $2, $3)
`
