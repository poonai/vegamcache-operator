package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func hitIp(ip string) {
	response, err := http.Get("http://" + ip + ":90")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
		}
		fmt.Printf("%s\n", string(contents))
	}
}
