package service

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

var (
	logger = logrus.New()
)

const (
	PROXY_DATA_DIR = "/root/network/data"
)

func init() {
	logger.SetLevel(logrus.DebugLevel)
}

type BaseConfig struct {
	Disabled    bool     `usage:"Enable/Disable service"`
	ExposePorts []string `usage:"Expose service ports to your host machine"`
	Dir         string   `usage:"Specify the main data directory of service"`
	Image       string   `usage:"Specify the image of service"`
}

type SharedConfig struct {
	Network        string
	HomeDir        string
	NetworkDir     string
	ExternalIp     string
	Dev            bool
	UseLocalImages string
}

type Base struct {
	Name        string
	Image       string
	Environment map[string]string
	Command     []string
	Ports       []string
	Volumes     []string
	Hostname    string

	Disabled bool
	Network  string

	Config BaseConfig
}

func newBase(name string) Base {
	return Base{
		Name:        name,
		Image:       "",
		Environment: make(map[string]string),
		Command:     []string{},
		Ports:       []string{},
		Volumes:     []string{},
		Hostname:    "",
	}
}

func (t *Base) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := ReflectFlags(t.Name, &t.Config, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}, cmd); err != nil {
		return err
	}
	return nil
}

func (t *Base) Apply(dir string, network string) error {

	ReflectFillConfig(t.Name, &t.Config)

	for _, port := range t.Config.ExposePorts {
		t.Ports = append(t.Ports, port)
	}

	t.Volumes = append(t.Volumes, fmt.Sprintf("%s:%s", t.Config.Dir, dir))

	t.Image = t.Config.Image

	t.Disabled = t.Config.Disabled

	t.Network = network

	return nil
}

func (t *Base) GetName() string {
	return t.Name
}

func (t *Base) GetImage() string {
	return t.Image
}

func (t *Base) GetCommand() []string {
	return t.Command
}

func (t *Base) GetEnvironment() map[string]string {
	return t.Environment
}

func (t *Base) GetVolumes() []string {
	return t.Volumes
}

func (t *Base) GetPorts() []string {
	return t.Ports
}

func (t *Base) GetHostname() string {
	return t.Hostname
}

func (t *Base) IsDisabled() bool {
	return t.Disabled
}

func (t *Base) ToJson() map[string]interface{} {
	var result = make(map[string]interface{})
	result["name"] = t.Name
	result["disabled"] = t.Disabled
	result["rpc"] = make(map[string]interface{})
	return result
}

func NewService(name string) Service {
	switch name {
	case "bitcoind":
		s := newBitcoind("bitcoind")
		return &s
	case "litecoind":
		s := newLitecoind("litecoind")
		return &s
	case "geth":
		s := newGeth("geth")
		return &s
	case "lndbtc":
		s := newLnd("lndbtc", "bitcoin")
		return &s
	case "lndltc":
		s := newLnd("lndltc", "litecoin")
		return &s
	case "connext":
		s := newConnext("connext")
		return &s
	case "xud":
		s := newXud("xud")
		return &s
	case "arby":
		s := newArby("arby")
		return &s
	case "boltz":
		s := newBoltz("boltz")
		return &s
	case "webui":
		s := newWebui("webui")
		return &s
	case "proxy":
		s := newProxy("proxy")
		return &s
	}

	return nil
}

func getDefaultValue(dv reflect.Value, fieldName string) interface{} {
	f := dv.FieldByName(fieldName)
	return f.Interface()
}

func ReflectFlags(name string, config interface{}, defaultValues interface{}, cmd *cobra.Command) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	dv := reflect.ValueOf(defaultValues).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		fn := field.Name
		usage := field.Tag.Get("usage")
		ft := field.Type

		key := fmt.Sprintf("%s.%s", name, strcase.ToKebab(fn))

		//p := v.FieldByName(fn).Addr().Interface()

		value := getDefaultValue(dv, fn)

		switch ft.Kind() {
		case reflect.String:
			//cmd.PersistentFlags().StringVar(p.(*string), key, value.(string), usage)
			cmd.PersistentFlags().String(key, value.(string), usage)
		case reflect.Bool:
			//cmd.PersistentFlags().BoolVar(p.(*bool), key, value.(bool), usage)
			cmd.PersistentFlags().Bool(key, value.(bool), usage)
		case reflect.Uint16:
			//cmd.PersistentFlags().Uint16Var(p.(*uint16), key, value.(uint16), usage)
			cmd.PersistentFlags().Uint16(key, value.(uint16), usage)
		case reflect.Slice:
			// FIXME differentiate slice item type
			//cmd.PersistentFlags().StringSliceVar(p.(*[]string), key, value.([]string), usage)
			cmd.PersistentFlags().StringSlice(key, value.([]string), usage)
		default:
			return errors.New("unsupported config struct field type: " + ft.Kind().String())
		}
		if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
			return err
		}
	}
	return nil
}

func ReflectFillConfig(name string, config interface{}) {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fn := field.Name
		ft := field.Type
		key := fmt.Sprintf("%s.%s", name, strcase.ToKebab(fn))
		//flag := cmd.PersistentFlags().Lookup(key)
		p := v.FieldByName(fn).Addr().Interface()

		switch ft.Kind() {
		case reflect.String:
			//if ! flag.Changed {
			//	*p.(*string) = viper.GetString(key)
			//}
			*p.(*string) = viper.GetString(key)
		case reflect.Bool:
			//if ! flag.Changed {
			//	*p.(*bool) = viper.GetBool(key)
			//}
			*p.(*bool) = viper.GetBool(key)
		case reflect.Uint16:
			//if ! flag.Changed {
			//	*p.(*uint16) = uint16(viper.GetUint(key))
			//}
			*p.(*uint16) = uint16(viper.GetUint(key))
		case reflect.Slice:
			//if ! flag.Changed {
			//	*p.(*[]string) = viper.GetStringSlice(key)
			//}
			*p.(*[]string) = viper.GetStringSlice(key)
		}
	}
}
