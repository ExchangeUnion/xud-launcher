package service

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

type BaseConfig struct {
	Disable     bool
	ExposePorts []string
	Dir         string
	Image       string
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

	config BaseConfig
}

func newBase(name string) Base {
	return Base{
		Name:        name,
		Image:       "",
		Environment: make(map[string]string),
		Command:     []string{},
		Ports:       []string{},
		Volumes:     []string{},

		config: BaseConfig{},
	}
}

func (t *Base) ConfigureFlags(defaultValues *BaseConfig, cmd *cobra.Command) error {
	var key string

	key = fmt.Sprintf("%s.disabled", t.Name)
	cmd.PersistentFlags().BoolVar(
		&t.config.Disable,
		key,
		defaultValues.Disable,
		fmt.Sprintf("Enable/Disable %s service", t.Name),
	)
	if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
		return err
	}

	key = fmt.Sprintf("%s.expose-ports", t.Name)
	cmd.PersistentFlags().StringSliceVar(
		&t.config.ExposePorts,
		key,
		defaultValues.ExposePorts,
		fmt.Sprintf("Expose %s service ports to your host machine", t.Name),
	)
	if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
		return err
	}

	key = fmt.Sprintf("%s.dir", t.Name)
	cmd.PersistentFlags().StringVar(
		&t.config.Dir,
		key,
		defaultValues.Dir,
		fmt.Sprintf("Specify the main data directory of %s service", t.Name),
	)
	if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
		return err
	}

	key = fmt.Sprintf("%s.image", t.Name)
	cmd.PersistentFlags().StringVar(
		&t.config.Image,
		key,
		"",
		fmt.Sprintf("Specify the image of %s service", t.Name),
	)
	if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
		return err
	}

	return nil
}

func (t *Base) Apply(dir string, network string) error {
	for _, port := range t.config.ExposePorts {
		t.Ports = append(t.Ports, port)
	}

	t.Volumes = append(t.Volumes, fmt.Sprintf("%s:%s", t.config.Dir, dir))

	if t.config.Image == "" {
		t.Image = images[network][t.GetName()]
	} else {
		t.Image = t.config.Image
	}

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

func (t *Base) Disabled() bool {
	return t.config.Disable
}

type Service interface {
	ConfigureFlags(cmd *cobra.Command) error
	GetConfig() interface{}
	GetName() string
	Apply(config *SharedConfig, services map[string]Service) error

	GetImage() string
	GetCommand() []string
	GetEnvironment() map[string]string
	GetVolumes() []string
	GetPorts() []string
	Disabled() bool
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

func ReflectFlags(name string, config interface{}, cmd *cobra.Command) error {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fn := field.Name
		usage := field.Tag.Get("usage")
		ft := field.Type
		key := fmt.Sprintf("%s.%s", name, strcase.ToKebab(fn))
		p := v.FieldByName(fn).Addr().Interface()
		switch ft.Kind() {
		case reflect.String:
			cmd.PersistentFlags().StringVar(p.(*string), key, "", usage)
		case reflect.Bool:
			fmt.Printf("%s p=%p\n", name, p)
			cmd.PersistentFlags().BoolVar(p.(*bool), key, false, usage)
		case reflect.Uint16:
			cmd.PersistentFlags().Uint16Var(p.(*uint16), key, 0, usage)
		default:
			return errors.New("unsupported config struct field type: " + ft.Kind().String())
		}
		if err := viper.BindPFlag(key, cmd.PersistentFlags().Lookup(key)); err != nil {
			return err
		}
	}
	return nil
}
