package service

import "testing"

func Test1(t *testing.T) {
	err := checkProvider("http://eth.kilrau.com:52041")
	if err != nil {
		println(err)
	}
}
