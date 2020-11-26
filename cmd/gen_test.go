package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"testing"
)

func Test1(t *testing.T) {
	rootCmd.SetArgs([]string{"-n", "simnet", "gen"})
	//genCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	if err != nil {
		t.Error(err)
	}
}

func Test2(t *testing.T) {
	genCmd.Parent().SetArgs([]string{"-n", "simnet", "gen", "--help"})
	//genCmd.SetArgs([]string{})
	err := genCmd.Execute()
	if err != nil {
		t.Error(err)
	}
}

type Base struct {
	A1 string
	A2 string
}

type Config struct {
	Base
	inner InnerConfig
}

type InnerConfig struct {
	Foo string
	Bar string
}

func (t Config) f1() {
	println("===> VALUE RECEIVER")
	fmt.Printf("%p\n", &t)
	fmt.Printf("cfg=%p, cfg.inner=%p, cfg.inner.Foo=%p cfg.innter.Bar=%p\n", &t, &t.inner, &t.inner.Foo, &t.inner.Bar)
	fmt.Printf("cfg.Base=%p, cfg.Base.A1=%p, cfg.Base.A2=%p\n", &t.Base, &t.Base.A1, &t.Base.A2)
}

func (t *Config) f2() {
	println("===> POINTER RECEIVER")
	fmt.Printf("%p\n", t)
	fmt.Printf("cfg=%p, cfg.inner=%p, cfg.inner.Foo=%p cfg.innter.Bar=%p\n", t, &t.inner, &t.inner.Foo, &t.inner.Bar)
	fmt.Printf("cfg.Base=%p, cfg.Base.A1=%p, cfg.Base.A2=%p\n", &t.Base, &t.Base.A1, &t.Base.A2)
}

func Test3(t *testing.T) {
	cfg := Config{}
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("test")
		},
	}

	fmt.Printf("cfg=%p, cfg.inner=%p, cfg.inner.Foo=%p cfg.innter.Bar=%p\n", &cfg, &cfg.inner, &cfg.inner.Foo, &cfg.inner.Bar)
	fmt.Printf("cfg.Base=%p, cfg.Base.A1=%p, cfg.Base.A2=%p\n", &cfg.Base, &cfg.Base.A1, &cfg.Base.A2)
	cmd.PersistentFlags().StringVar(&cfg.inner.Foo, "foo", "bar", "")

	cmd.SetArgs([]string{})
	cmd.Execute()
	println(cfg.inner.Foo)

	cfg.f1()
	cfg.f2()
}
