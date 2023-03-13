package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://randomuser.me/api/"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	err := ioutil.WriteFile("test.json", body, 0644)
	if err != nil {
		panic(err)
	}

}
