package storage

var CreateShortURLTable = `
	CREATE TABLE IF NOT EXISTS urls (
		id integer,
		user_id text,
		short_url text,
		url text PRIMARY KEY
	);`

var GetURL = `
	SELECT url 
	FROM urls 
	WHERE short_url = $1
`

var GetUserURL = `
	SELECT short_url, url
	FROM urls
	WHERE user_id = $1
`

var SaveURL = `
	INSERT into urls (id, user_id, short_url, url)
	VALUES ($1, $2, $3, $4)
`
