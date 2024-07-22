package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	ErrFileTypeNotAllow = errors.New("file type not allow")
	ErrReadInConfig     = errors.New("fatal error config file")
	ErrUnmarshal        = errors.New("unmarshal config failed")
)

func allowType(t string) bool {
	l := []string{
		"json",
		"toml",
		"yaml",
		"yml",
		"properties",
		"props",
		"prop",
		"hcl",
		"tfvars",
		"dotenv",
		"env",
		"ini",
		"cfg", // viper 不直接支持 cfg 格式的配置文件解析，需要转为 ini 进行解析
	}

	for _, v := range l {
		if v == t {
			return true
		}
	}

	return false
}

// ReadConfig 使用 viper 读取配置文件 支持文件类型 JSON, TOML, YAML, HCL, INI, envfile or Java properties
// viper的配置的key值目前是不区分大小写, 如果文件后缀为 cfg 格式则这里采用默认的分割符为 ::
func ReadConfig(config any, filename string) (v *viper.Viper, err error) {
	if _, err = os.Stat(filename); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	fileType := ""
	if ext != ".atlantis" { // 这里说名文件名称包含后缀
		fileType = ext[1:] // 注意这里获取到的文件后缀是包含 . 号的，再进行和 allowType 进行比较的时候需要去除开头的
	} else {
		// 如果文件名称不包含后缀这里为了兼容 atlantis-agent
		// 写到配置文件目录下的 .atlantis 文件(内容格式为 cfg)
		// 这里默认设置为 cfg 格式
		fileType = "cfg"
	}

	fileType = strings.ToLower(fileType)
	if !allowType(fileType) {
		return nil, ErrFileTypeNotAllow
	}

	v = viper.NewWithOptions()

	if fileType == "cfg" {
		fileType = "ini"
		// 如果你想要解析那些键本身就包含.(默认的键分隔符）的配置，需要修改分隔符, 这里默认设置为 ::
		v = viper.NewWithOptions(viper.KeyDelimiter("::"))
	}

	v.SetConfigFile(filename)
	v.SetConfigType(fileType)

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Join(ErrReadInConfig, err)
	}

	if err := v.Unmarshal(config); err != nil {
		return nil, errors.Join(ErrUnmarshal, err)
	}

	return v, nil
}
