package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	restore   bool
	backupDir string
)

const (
	DefaultPassword = "OpenDEX!Rocks"
)

func init() {
	setupCmd.PersistentFlags().String("wallet-password", "", "XUD master wallet password")
	setupCmd.PersistentFlags().StringVar(&backupDir, "backup-dir", "", "XUD backup location")
	setupCmd.PersistentFlags().BoolVar(&restore, "restore", true, "Restore wallets")

	rootCmd.AddCommand(setupCmd)
}

func runCommand(name string, args ...string) {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	color.Blue(c.String())
	err := c.Run()
	if err != nil {
		fmt.Printf("🐞 %s\n", err)
		os.Exit(1)
	}
	fmt.Println()
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Bring up your XUD environment in one command",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		logfile := filepath.Join(networkDir, "logs", fmt.Sprintf("%s.log", network))
		f, err := os.Create(logfile)
		if err != nil {
			logger.Fatalf("Failed to create %s: %s", logfile, err)
		}
		defer f.Close()

		launcher := os.Args[0]

		logger.Debugf("Generate files")
		runCommand(launcher, "gen")

		logger.Debugf("Pulling images")
		runCommand(launcher, "pull")

		logger.Debugf("Starting proxy")
		runCommand(launcher, "up", "--", "-d", "proxy")

		logger.Debugf("Starting lndbtc, lndltc and connext")
		runCommand(launcher, "up", "--", "-d", "lndbtc", "lndltc", "connext")

		// FIXME enable tls verification
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		logger.Debugf("Waiting for LNDs to be synced")
		_, err = f.WriteString("Waiting for XUD dependencies to be ready\n")
		if err != nil {
			logger.Fatalf("Failed to write to %s: %s", logfile, err)
		}

		var lnds = map[string]string{
			"lndbtc": "0.00% (0/0)",
			"lndltc": "0.00% (0/0)",
		}
		var lndsMutex = &sync.Mutex{}
		waitLnds(func(service string, status string) {
			lndsMutex.Lock()
			defer lndsMutex.Unlock()
			lnds[service] = status
			_, err := f.WriteString(fmt.Sprintf(" [LightSync] lndbtc: %s | lndltc: %s\n", lnds["lndbtc"], lnds["lndltc"]))
			if err != nil {
				logger.Fatalf("Failed to write to %s: %s", logfile, err)
			}
		})

		logger.Debugf("Starting xud")
		runCommand(launcher, "up", "--", "-d", "xud")

		logger.Debugf("Ensuring wallets are created and unlocked")
		_, err = f.WriteString("Setup wallets\n")
		if err != nil {
			logger.Fatalf("Failed to write to %s: %s", logfile, err)
		}
		ensureWallets()

		logger.Debugf("Starting boltz")
		runCommand(launcher, "up", "--", "-d", "boltz")

		_, err = f.WriteString("Start shell\n")
		if err != nil {
			logger.Fatalf("Failed to write to %s: %s", logfile, err)
		}
	},
}

type StatusResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func getServiceStatus(name string) string {
	resp, err := http.Get(fmt.Sprintf("https://localhost:8889/api/v1/status/%s", name))
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var r StatusResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		logger.Fatal(err)
	}
	return r.Status
}

func waitLnds(onChange func(service string, status string)) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		waitLnd("lndbtc", func(status string) {
			onChange("lndbtc", status)
		})
		wg.Done()
	}()

	go func() {
		waitLnd("lndltc", func(status string) {
			onChange("lndltc", status)
		})
		wg.Done()
	}()

	wg.Wait()
}

func waitLnd(name string, onChange func(status string)) {
	for {
		status := getServiceStatus(name)
		logger.Debugf("%s: %s", name, status)
		onChange(status)
		if strings.Contains(status, "100.00%") {
			break
		}
		if strings.Contains(status, "99.99%") {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ensureWallets() {
	status := getServiceStatus("xud")
	logger.Debugf("xud: %s", status)
	if strings.Contains(status, "Wallet missing") {
		logger.Debug("Creating wallets")
		create(DefaultPassword)
	} else if strings.Contains(status, "Wallet locked") {
		logger.Debug("Unlocking wallets")
		unlock(DefaultPassword)
	}
}

func create(password string) {
	var payload = []byte(`{"password":"` + password + `"}`)
	resp, err := http.Post("https://localhost:8889/api/v1/xud/create", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(body))
}

func unlock(password string) {
	var payload = []byte(`{"password":"` + password + `"}`)
	resp, err := http.Post("https://localhost:8889/api/v1/xud/unlock", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(body))
}
