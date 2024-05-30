package config

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signingKey" json:"signingKey"`
}

type AliSmsConfig struct {
	ApiKey    string `mapstructure:"apikey" json:"api_key"`
	ApiSecret string `mapstructure:"apisecret" json:"apisecret"`
	Expire    int    `mapstructure:"expire" json:"expire"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        uint64 `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
}

type AliPayConfig struct {
	AppID        string `mapstructure:"appId" json:"appId"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type ServerConfig struct {
	Name             string         `mapstructure:"name" json:"name"`
	Host             string         `mapstructure:"host" json:"host"`
	Port             int            `mapstructure:"port" json:"port"`
	Tags             []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo GoodsSrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	OrderSrvInfo     GoodsSrvConfig `mapstructure:"order_srv" json:"order_srv"`
	JWTInfo          JWTConfig      `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo       AliSmsConfig   `mapstructure:"sms" json:"sms"`
	RedisInfo        RedisConfig    `mapstructure:"redis" json:"redis"`
	ConsulInfo       ConsulConfig   `mapstructure:"consul" json:"consul"`
	AliPayInfo       AliPayConfig   `mapstructure:"ali_pay" json:"ali_pay"`
}
