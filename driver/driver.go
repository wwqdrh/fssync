package driver

import "errors"

var (
	ErrNotAuth = errors.New("no auth")
	ErrInvUrl  = errors.New("invalid url")
)

type FileItem struct {
	Name          string
	Href          string
	Owner         string
	Status        string
	ResourceType  interface{}
	ContentType   string
	ContentLength int64
	LastModify    string
	Privileges    []string
}

type IDriver interface {
	Auth(name, password string)
	IsAuth() bool
	Download(url string) error
	List(url string) ([]FileItem, error)
	Delete(url string) error
	Update(local, url string) error
}
