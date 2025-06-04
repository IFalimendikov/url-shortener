package storage

func (s *Storage) PingDB() bool {
	if s.DB != nil {
		err := s.DB.Ping()
		return err == nil
	}
	return false
}
