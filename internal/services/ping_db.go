package services

import ()

func (s *URLs) PingDB() bool {
	return s.Storage.PingDB()
}
