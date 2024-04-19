package driver

import (
	"encoding/json"
	"errors"
	"os"
)

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

type DriverConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var defaultDriverData = map[string]DriverConfig{
	"坚果云": {
		Username: "",
		Password: "",
	},
}

// 检查homePath下是否存在config.json, 不存在就新建
// {"driver1": {"username": "", "password": ""}}
func initDriver(configPath string) error {
	// 检查文件是否存在
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		data, err := json.MarshalIndent(defaultDriverData, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(configPath, data, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

// 读取homePath下config.json中name部分的数据
func LoadDriver(dataPath string, name string) (IDriver, error) {
	if err := initDriver(dataPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}

	var config map[string]DriverConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	driverConfig, ok := config[name]
	if !ok {
		return nil, errors.New("no auth, pls update config.json") // 如果name不存在，返回空DriverConfig
	}

	switch name {
	case "坚果云":
		d := NewJianguoDriver()
		d.Auth(driverConfig.Username, driverConfig.Password)
		return d, nil
	default:
		return nil, errors.New("no this driver")
	}
}
