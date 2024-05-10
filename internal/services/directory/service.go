package directory

import (
	"os"
	"path"
)

const (
	Local       = `.queeck` // path related to home directory, ex: /home/username/.queeck
	Permissions = os.ModePerm
)

type Service struct {
	working string
	home    string
}

func New() (service *Service, err error) {
	s := &Service{}
	if s.working, err = os.Getwd(); err != nil {
		return nil, err
	}
	if s.home, err = os.UserHomeDir(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) Working() string {
	return s.working
}

func (s *Service) Home() string {
	return s.home
}

func (s *Service) Local() string {
	return path.Join(s.home, Local)
}

func (s *Service) IsLocalCreated() bool {
	_, err := os.Stat(s.Local())
	return err == nil
}

func (s *Service) CreateLocal() error {
	return os.MkdirAll(s.Local(), Permissions)
}
