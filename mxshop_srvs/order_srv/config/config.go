package config

type ServerConfig struct {
	Name        string       `mapstructure:"name" json:"name"`
	Host        string       `mapstructure:"host" json:"host"`
	Tags        []string     `mapstructure:"tags" json:"tags"`
	MysqlInfo   MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
	JaegerInfo  JaegerConfig `mapstructure:"jaeger" json:"jaeger"`
	RedisConfig RedisConfig  `mapstructure:"redis" json:"redis"`

	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`         //商品微服务配置
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"` //库存微服务配置
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
