// Code generated by go2struct. DO NOT EDIT.
package config

type Settings struct {
	Application struct {
		Mode string `mapstructure:"mode"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"application"`

	MySQL struct {
		Master            string `mapstructure:"master"`
		Slave             string `mapstructure:"slave"`
		Amazon            string `mapstructure:"amazon"`
		MasterTablePrefix string `mapstructure:"masterTablePrefix"`
	} `mapstructure:"mysql"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Aliyun struct {
		AccessKeyId     string `mapstructure:"accessKeyId"`
		AccessKeySecret string `mapstructure:"accessKeySecret"`
		RegionId        string `mapstructure:"regionId"`
	} `mapstructure:"aliyun"`

	Qiniu struct {
		AccessKey string `mapstructure:"accessKey"`
		SecretKey string `mapstructure:"secretKey"`
		Bucket    string `mapstructure:"bucket"`
		Domain    string `mapstructure:"domain"`
	} `mapstructure:"qiniu"`
}
