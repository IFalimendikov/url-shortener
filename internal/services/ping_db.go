package services

// PingDB checks if the database connection is alive and returns true if successful
func (s *URLs) PingDB() bool {
	return s.Storage.PingDB()
}
