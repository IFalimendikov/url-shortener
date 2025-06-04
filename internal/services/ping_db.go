package services

func (s *URLs) PingDB() bool {
	return s.Storage.PingDB()
}
