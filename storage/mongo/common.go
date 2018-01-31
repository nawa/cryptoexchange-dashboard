package mongo

import "github.com/globalsign/mgo"

type baseStorage struct {
	baseSession    *mgo.Session
	refreshSession bool
}

func (s *baseStorage) getDB() (db *mgo.Database, closeSession func()) {
	session := s.baseSession
	if s.refreshSession {
		session = session.Copy()
		closeSession = func() {
			session.Close()
		}
	} else {
		closeSession = func() {
			//do nothing with one base session
		}
	}
	return session.DB(""), closeSession
}
