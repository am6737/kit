package config

import (
	"fmt"
	"github.com/cossim/kit/pkg/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type LogConfig struct {
	Stdout bool   `mapstructure:"stdout" yaml:"stdout"`
	V      int    `mapstructure:"v" yaml:"v"`
	Format string `mapstructure:"format" yaml:"format"`
}

type MySQLConfig struct {
	DSN     string `mapstructure:"dsn" yaml:"dsn"`
	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
}

func (c MySQLConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type RedisConfig struct {
	Proto    string `mapstructure:"proto" yaml:"proto"`
	Password string `mapstructure:"password" yaml:"password"`
	Address  string `mapstructure:"address" yaml:"address"`
	Port     int    `mapstructure:"port" yaml:"port"`
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type HTTPConfig struct {
	Address string        `mapstructure:"address" yaml:"address"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

func (c HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type GRPCConfig struct {
	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
}

func (c GRPCConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type RegisterConfig struct {
	Name    string   `mapstructure:"name" yaml:"name"`
	Tags    []string `mapstructure:"tags" yaml:"tags"`
	Address string   `mapstructure:"address" yaml:"address"`
	Port    int      `mapstructure:"port" yaml:"port"`
	Enabled bool     `mapstructure:"enabled" yaml:"enabled"`
}

func (c RegisterConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type DiscoversConfig map[string]ServiceConfig

type ServiceConfig struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
	Direct  bool   `mapstructure:"direct" yaml:"direct"`
}

func (c ServiceConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type MessageQueueConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Address  string `mapstructure:"address" yaml:"address"`
	Port     int    `mapstructure:"port" yaml:"port"`
}

func (c MessageQueueConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type DtmConfig struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Address string `mapstructure:"address" yaml:"address"`
	Port    int    `mapstructure:"port" yaml:"port"`
}

func (c DtmConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type OSSCommonConfig struct {
	Address   string `mapstructure:"address" yaml:"address"`
	Port      int    `mapstructure:"port" yaml:"port"`
	AccessKey string `mapstructure:"accessKey" yaml:"accessKey"`
	SecretKey string `mapstructure:"secretKey" yaml:"secretKey"`
	SSL       bool   `mapstructure:"ssl" yaml:"ssl"`
	//PresignedExpires int    `mapstructure:"presignedExpires"`
}

func (c OSSCommonConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

type OssConfig map[string]OSSCommonConfig

type Config struct {
	Log                 LogConfig                 `mapstructure:"log" yaml:"log"`
	MySQL               MySQLConfig               `mapstructure:"mysql" yaml:"mySQL"`
	Redis               RedisConfig               `mapstructure:"redis" yaml:"redis"`
	HTTP                HTTPConfig                `mapstructure:"http" yaml:"http"`
	GRPC                GRPCConfig                `mapstructure:"grpc" yaml:"grpc"`
	Register            RegisterConfig            `mapstructure:"register" yaml:"register"`
	Discovers           DiscoversConfig           `mapstructure:"discovers" yaml:"discovers"`
	Encryption          EncryptionConfig          `mapstructure:"encryption" yaml:"encryption"`
	MessageQueue        MessageQueueConfig        `mapstructure:"message_queue" yaml:"messageQueue"`
	MultipleDeviceLimit MultipleDeviceLimitConfig `mapstructure:"multiple_device_limit" yaml:"multiple_device_limit"`
	SystemConfig        SystemConfig              `mapstructure:"system" yaml:"system"`
	Dtm                 DtmConfig                 `mapstructure:"dtm" yaml:"dtm"`
	OSS                 OssConfig                 `mapstructure:"oss" yaml:"oss"`
	Email               EmailConfig               `mapstructure:"email" yaml:"email"`
	Livekit             LivekitConfig             `mapstructure:"livekit" yaml:"livekit"`
}

type EncryptionConfig struct {
	Enable     bool   `mapstructure:"enable" yaml:"enable"`
	Name       string `mapstructure:"name" yaml:"name"`
	Email      string `mapstructure:"email" yaml:"email"`
	RsaBits    int    `mapstructure:"rsaBits" yaml:"rsaBits"`
	Passphrase string `mapstructure:"passphrase" yaml:"passphrase"`
}

type MultipleDeviceLimitConfig struct {
	Enable bool `mapstructure:"enable" yaml:"enable"`
	Max    int  `mapstructure:"max" yaml:"max"`
}

type SystemConfig struct {
	Environment       string `mapstructure:"environment" yaml:"environment"`
	AvatarFilePath    string `mapstructure:"avatar_file_path" yaml:"avatar_file_path"`
	AvatarFilePathDev string `mapstructure:"avatar_file_path_dev" yaml:"avatar_file_path_dev"`
	GatewayAddress    string `mapstructure:"gateway_address" yaml:"gateway_address"`
	GatewayPort       string `mapstructure:"gateway_port" yaml:"gateway_port"`
	GatewayAddressDev string `mapstructure:"gateway_address_dev" yaml:"gateway_address_dev"`
	GatewayPortDev    string `mapstructure:"gateway_port_dev" yaml:"gateway_port_dev"`
}

type EmailConfig struct {
	Enable     bool   `mapstructure:"enable" yaml:"enable"`
	SmtpServer string `mapstructure:"smtp_server" yaml:"smtp_server"`
	Port       int    `mapstructure:"port" yaml:"port"`
	Username   string `mapstructure:"username" yaml:"username"`
	Password   string `mapstructure:"password" yaml:"password"`
}

type LivekitConfig struct {
	Address   string        `mapstructure:"address" yaml:"address"`
	Url       string        `mapstructure:"url" yaml:"livekit_server"`
	ApiKey    string        `mapstructure:"api_key" yaml:"api_key"`
	ApiSecret string        `mapstructure:"secret_key" yaml:"secret_key"`
	Timeout   time.Duration `mapstructure:"timeout" yaml:"timeout"`
	Port      int           `mapstructure:"port" yaml:"port"`
}

func (c LivekitConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

func (c *Config) HTTPClientFor(config *Config) (*http.Client, error) {
	//transport, err := TransportFor(config)
	//if err != nil {
	//	return nil, err
	//}
	//var httpClient *http.Client
	//if transport != http.DefaultTransport || config.HTTP.Timeout > 0 {
	//	httpClient = &http.Client{
	//		Transport: transport,
	//		Timeout:   config.HTTP.Timeout,
	//	}
	//} else {
	//	httpClient = http.DefaultClient
	//}
	//
	//return httpClient, nil
	return http.DefaultClient, nil
}

// TransportFor returns an http.RoundTripper that will provide the authentication
// or transport level security defined by the provided Config. Will return the
// default http.DefaultTransport if no special case behavior is needed.
func TransportFor(config *Config) (http.RoundTripper, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.New(cfg)
}

// TransportConfig converts a client config to an appropriate transport config.
func (c *Config) TransportConfig() (*transport.Config, error) {
	return nil, nil
}

func (c *Config) GRPCClientFor(config *Config) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(config.GRPC.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
