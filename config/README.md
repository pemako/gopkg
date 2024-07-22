# 解析配置文件使用示例

- 针对 `cfg` 格式的配置文件 `server.cfg` 内容格式使用说明如下

```cfg
; 这里要求必须存在 sectionName 如下 app_mode = dev 没有指定 section 在解析的时候
; 获取不到 app_mode 的值。故这里在使用的时候要求必须指定 section
;
; app_mode = dev
;
[server]
name = golang_libs_test_data
http_port = 10000 # 这里为注释 需要显示指定解析到自定义结构上的名称
Debug = false     ; 这里是另一中注释形式
polaris-addr = polaris-sdk-in.domob-inc.com:7867 ; 支持中画线
trace.addr = http://0.0.0.0:7820/api/traces   ; 支持key中存在点号，这里在sdk中会默认设置分割符为 ::

[DB]
host = xxxxx
Port = 3306
USER = domob

[Note]
Content = This is a test data
UseLibs = viper,ini ; 支持直接映射为 golang 的 []string{} 类型 注意原来业务中使用该类值的调
```

> 问题: 使用原来 [`gcfg`](github.com/go-gcfg/gcfg) 库提供的 `ReadFileInto()` 方法解析配置文件到结构体上的时候有很多限制，现升级采用 `viper` 进行统一的处理配置文件。

对于上面的配置文件使用姿势如下

> 说明： `viper` 对设置的键大小写不敏感，如果存在中划线，下滑线等需要进行显示的设置 `tag`。Viper在后台使用 [mapstructure](github.com/mitchellh/mapstructure) 来解析值, 默认情况下使用 `mapstructure`

1. 把 `server.cfg` 解析到 对应的 `strcut` 上

```golang
package coning

import (
  tools "github.com/pemako/gopkg/config"
)

type Config struct {
 Server Server `mapstructure:"server"`
 DB     DB     `mapstructure:"DB"`
 Note
}

type Server struct {
 Name        string `mapstructure:"name"`
 HttPort     int    `mapstructure:"http_port"`
 Debug       bool   `mapstructure:"debug"`
 PolarisAddr string `mapstructure:"polaris-addr"`
 TraceAddr   string `mapstructure:"trace.addr"`
}

type DB struct {
 Host string `mapstructure:"host"`
 Port int    `mapstructure:"Port"`
 User string `mapstructure:"USER"`
}

type Note struct {
 Content string   `mapstructure:"Content"`
 UseLibs []string `mapstructure:"UseLibs"`
}

var c = new(Config)

v, err := tools.ReadConfig(c, "server.cfg")
if err != nil {
 panic(err)
}


// 下面有两种方式获取值
// 1. 直接诶采用解析到 struct 上的 c 进行获取结构体中的值, 推荐
fmt.Println(c.Server.Name)

// 2. 使用 v 进行获取值，判断key，设置默认值等操作; 注意这种方式取值时如果文件格式为cfg则默认采用双冒号 ::，其它格式默认为 .
v.GetString("server::name")
v.GetString("server.name") // 非 cfg 格式文件
v.GetString("server::http_port")
```


> 这时就可以使用返回的 viper 实例进行使用提供的各种方法操作。详见 [viper](https://pkg.go.dev/github.com/spf13/viper)

### 一些常用的方法

- `Get(key string) interface{}`
- `GetBool(key string) bool`
- `GetDuration(key string) time.Duration`
- `GetFloat64(key string) float64`
- `GetInt(key string) int`
- `GetInt32(key string) int32`
- `GetInt64(key string) int64`
- `GetIntSlice(key string) []int`
- `GetSizeInBytes(key string) uint`
- `GetString(key string) string`
- `GetStringMap(key string) map[string]interface{}`
- `GetStringMapString(key string) map[string]string`
- `GetStringMapStringSlice(key string) map[string][]string`
- `GetStringSlice(key string) []string`
- `GetTime(key string) time.Time`
- `GetUint(key string) uint`
- `GetUint32(key string) uint32`
- `GetUint64(key string) uint64`
- `InConfig(key string) bool`
- `IsSet(key string) bool`

### 已知问题

[issues](https://github.com/spf13/viper/issues/1402) 如果 `cfg` 配置文件中的 `value` 值包含 `#` 特殊字符的需要使用 反引号 `` 。

```cfg
materialPrefix = "/landing/#/material"
# 需要修改为
materialPrefix = `/landing/#/material`
```

### 其它使用参考

- [https://github.com/spf13/viper](https://github.com/spf13/viper)
- [https://www.liwenzhou.com/posts/Go/viper_tutorial/#autoid-1-5-0](https://www.liwenzhou.com/posts/Go/viper_tutorial/#autoid-1-5-0)
