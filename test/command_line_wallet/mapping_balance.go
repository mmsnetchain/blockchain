package main

import (
	"mmschainnewaccount/config"

	"mmschainnewaccount/test/map_chain_tools/sqlite"
)

func main() {
	sqlite.Init()

	config

}

func getNewAddress(c *gcli.Command, _ []string) error {

	password := interact.ReadPassword("Enter Password:")
	var params map[string]interface{}
	params = make(map[string]interface{})
	params["password"] = password
	info := model.Info{
		Method: "getnewaddress",
		Params: params,
	}
	result := model.HttpPost(url, info)
	gcli.Println(result)
	return nil
}

type Info struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func HttpPost(url string, info Info) string {
	jsons, _ := json.Marshal(info)
	result := string(jsons)
	jsonInfo := strings.NewReader(result)
	req, _ := http.NewRequest("POST", url, jsonInfo)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("user", config.RPCUser)
	req.Header.Add("password", config.RPCPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error create client:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error getInfo:%v", err)
	}
	return string(body)
}
