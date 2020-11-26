package cmd

import (
	"testing"
)

func Test1(t *testing.T) {
	genCmd.Parent().SetArgs([]string{"-n", "simnet", "gen"})
	//genCmd.SetArgs([]string{})
	err := genCmd.Execute()
	if err != nil {
		t.Error(err)
	}
}
