package driver

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/wwqdrh/gokit/logger"
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
	SetIgnore(p []string)
	Auth(name, password string)
	IsAuth() bool
	Download(url string) error
	List(url string) ([]FileItem, error)
	Delete(url string) error
	Update(local, url string) error
	GetLastTimeline(name string) string
	GetLastTimelineMap() map[string]int64
}

type DriverConfigAll struct {
	cfg  string
	data map[string]*DriverConfig
}

type DriverConfig struct {
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Ignores   []string          `json:"ignores"`
	TimeLines map[string]string `json:"timelines"` // 存储各个文件的上次上传时间
}

type IDriverConfig interface {
	GetLastTimeline(mode string, pname string) string // 获取文件的上次上传时间
	SetLastTimeline(mode string, pname string)        // 设置文件上次上传时间
	GetLastTimelineMap(mode string) map[string]int64
	GetConfig(mode string) (*DriverConfig, bool)
}

var defaultDriverData = map[string]DriverConfig{
	"坚果云": {
		Username: "",
		Password: "",
	},
}

func NewDriverConfigAll(cfg string) (IDriverConfig, error) {
	data, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}

	var config map[string]*DriverConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &DriverConfigAll{
		cfg:  cfg,
		data: config,
	}, nil
}

func (c *DriverConfigAll) GetLastTimeline(mode string, pname string) string {
	cfg, ok := c.data[mode]
	if !ok {
		return ""
	}
	lastupdate, ok := cfg.TimeLines[pname]
	if !ok {
		return ""
	}
	return lastupdate
}

func (c *DriverConfigAll) GetLastTimelineMap(mode string) map[string]int64 {
	cfg, ok := c.data[mode]
	if !ok {
		return map[string]int64{}
	}

	res := map[string]int64{}
	for name, t := range cfg.TimeLines {
		tt, err := time.Parse(time.RFC3339, t)
		if err != nil {
			logger.DefaultLogger.Warn(err.Error())
			continue
		}
		res[name] = tt.UnixNano()
	}
	return res
}

func (c *DriverConfigAll) SetLastTimeline(mode string, pname string) {
	cfg, ok := c.data[mode]
	if !ok {
		return
	}
	if cfg.TimeLines == nil {
		cfg.TimeLines = map[string]string{}
	}
	cfg.TimeLines[pname] = time.Now().Format(time.RFC3339)
	c.dump()
}

func (c *DriverConfigAll) GetConfig(mode string) (*DriverConfig, bool) {
	config, ok := c.data[mode]
	return config, ok
}

func (c *DriverConfigAll) dump() {
	data, err := json.Marshal(c.data)
	if err != nil {
		logger.DefaultLogger.Warn(err.Error())
		return
	}
	err = os.WriteFile(c.cfg, data, 0644)
	if err != nil {
		logger.DefaultLogger.Warn(err.Error())
		return
	}
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

	da, err := NewDriverConfigAll(dataPath)
	if err != nil {
		return nil, err
	}

	driverConfig, ok := da.GetConfig(name)
	if !ok {
		return nil, errors.New("no auth, pls update config.json") // 如果name不存在，返回空DriverConfig
	}

	switch name {
	case "坚果云":
		d := NewJianguoDriver(da)
		d.Auth(driverConfig.Username, driverConfig.Password)
		d.SetIgnore(driverConfig.Ignores)
		return d, nil
	default:
		return nil, errors.New("no this driver")
	}
}
