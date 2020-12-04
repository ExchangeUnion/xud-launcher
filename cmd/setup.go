package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	restore   bool
	backupDir string
)

const (
	DEFAULT_PASSWORD = "OpenDEX!Rocks"
)

func init() {
	setupCmd.PersistentFlags().String("wallet-password", "", "XUD master wallet password")
	setupCmd.PersistentFlags().StringVar(&backupDir, "backup-dir", "", "XUD backup location")
	setupCmd.PersistentFlags().BoolVar(&restore, "restore", true, "Restore wallets")

	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Bring up your XUD environment in one command",
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		logger.Debugf("Generate files")
		err = exec.Command(os.Args[0], "gen").Run()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Starting proxy")
		err = exec.Command(os.Args[0], "up", "--", "-d", "proxy").Run()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Starting lndbtc, lndltc and connext")
		err = exec.Command(os.Args[0], "up", "--", "-d", "lndbtc", "lndltc", "connext").Run()
		if err != nil {
			logger.Fatal(err)
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		logger.Debugf("Waiting for LNDs to be synced")
		waitLnds()

		logger.Debugf("Starting xud")
		err = exec.Command(os.Args[0], "up", "--", "-d", "xud").Run()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Debugf("Ensuring wallets are created and unlocked")
		ensureWallets()

		logger.Debugf("Starting boltz")
		err = exec.Command(os.Args[0], "up", "--", "-d", "boltz").Run()
		if err != nil {
			logger.Fatal(err)
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

func waitLnds() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		waitLnd("lndbtc")
		wg.Done()
	}()

	go func() {
		waitLnd("lndltc")
		wg.Done()
	}()

	wg.Wait()
}

func waitLnd(name string) {
	for {
		status := getServiceStatus(name)
		logger.Debugf("%s: %s", name, status)
		if strings.Contains(status, "100.00%") {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ensureWallets() {
	status := getServiceStatus("xud")
	if strings.Contains(status, "Wallet missing") {
		create(DEFAULT_PASSWORD)
	} else if strings.Contains(status, "Wallet locked") {
		unlock(DEFAULT_PASSWORD)
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
