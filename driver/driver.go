package driver

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
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
	UploadExtras(p map[string]string) // local, remove
	Auth(name, password string)
	IsAuth() bool
	Download(url string) error
	List(url string) ([]FileItem, error)
	Delete(url string) error
	Update(local, url string) error
	GetLastTimeline(name string) int64
	GetLastTimelineMap() map[string]int64
}

type DriverConfigAll struct {
	cfg  string
	data *sync.Map // map[string]*DriverConfig
}

type DriverConfig struct {
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Extras    map[string]string `json:"extras"`
	Ignores   []string          `json:"ignores"`
	TimeLines map[string]int64  `json:"timelines"` // 存储各个文件的上次上传时间
}

type IDriverConfig interface {
	GetLastTimeline(mode string, pname string) int64 // 获取文件的上次上传时间
	SetLastTimeline(mode string, pname string)       // 设置文件上次上传时间
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
	m := &sync.Map{}
	for key, v := range config {
		m.Store(key, v)
	}
	return &DriverConfigAll{
		cfg:  cfg,
		data: m,
	}, nil
}

func (c *DriverConfigAll) GetLastTimeline(mode string, pname string) int64 {
	cfg, ok := c.data.Load(mode)
	if !ok {
		return 0
	}
	lastupdate, ok := cfg.(*DriverConfig).TimeLines[pname]
	if !ok {
		return 0
	}
	return lastupdate
}

func (c *DriverConfigAll) GetLastTimelineMap(mode string) map[string]int64 {
	cfg, ok := c.data.Load(mode)
	if !ok {
		return map[string]int64{}
	}

	res := map[string]int64{}
	for name, t := range cfg.(*DriverConfig).TimeLines {
		res[name] = t
	}
	return res
}

func (c *DriverConfigAll) SetLastTimeline(mode string, pname string) {
	cfgR, ok := c.data.Load(mode)
	if !ok {
		return
	}
	cfg := cfgR.(*DriverConfig)
	if cfg.TimeLines == nil {
		cfg.TimeLines = map[string]int64{}
	}
	cfg.TimeLines[pname] = time.Now().UnixNano()
	c.dump()
}

func (c *DriverConfigAll) GetConfig(mode string) (*DriverConfig, bool) {
	config, ok := c.data.Load(mode)
	return config.(*DriverConfig), ok
}

func (c *DriverConfigAll) dump() {
	configs := map[string]*DriverConfig{}
	c.data.Range(func(key, value any) bool {
		configs[key.(string)] = value.(*DriverConfig)
		return true
	})

	data, err := json.Marshal(configs)
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
		d.UploadExtras(driverConfig.Extras)
		return d, nil
	default:
		return nil, errors.New("no this driver")
	}
}
