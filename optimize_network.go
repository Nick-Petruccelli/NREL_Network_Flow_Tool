package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main()  {
	cap_graph_json, err := os.Open("capacity_graph.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File opened successfully")

	defer cap_graph_json.Close()

	byte_value, _ := ioutil.ReadAll(cap_graph_json)

	var res map[string]interface{}
	json.Unmarshal([]byte(byte_value), &res)

	fmt.Println(res["source1"])
}