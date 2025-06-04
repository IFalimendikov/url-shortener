package storage

// PingDB checks if the database connection is alive
func (s *Storage) PingDB() bool {
	if s.DB != nil {
		err := s.DB.Ping()
		return err == nil
	}
	return false
}
