package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	tps()

}

func tps() {
	url := "/rpc"
	method := "POST"
	params := map[string]interface{}{
		"address": "13gyTcKozkP5b1bSTxFWn8sYttcoWT7Bmt",
		"amount":  100000000,
		"gas":     0,
		"pwd":     "123456789",
		"comment": "test",
	}

	rpcParans := map[string]interface{}{
		"method": "sendtoaddress",
		"params": params,
	}

	header := http.Header{
		"user":         []string{"test"},
		"password":     []string{"testp"},
		"Content-Type": []string{"application/json"}}
	client := &http.Client{}

	bs, err := json.Marshal(rpcParans)
	req, err := http.NewRequest(method, "http:
	if err != nil {
		fmt.Println("request")
		return
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("")
		return
	}
	fmt.Println("response:", resp.StatusCode)
	if resp.StatusCode == 200 {
		robots, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("body")
			return
		}

		fmt.Println(len(robots), string(robots))
	}
}
