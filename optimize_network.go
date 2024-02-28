package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main()  {
	res := json_to_map("capacity_graph.json")

	fmt.Println(res["source1"])

	cap_to_aug_flow(res)
}

type edge struct {
	dest string
	cap int
	flow int
}

func json_to_map(file_name string) map[string]interface{} {
	cap_graph_json, err := os.Open(file_name)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File opened successfully")

	defer cap_graph_json.Close()

	byte_value, _ := ioutil.ReadAll(cap_graph_json)

	var res map[string]interface{}
	json.Unmarshal([]byte(byte_value), &res)

	return res
}

func cap_to_aug_flow(cap_graph map[string]interface{}) map[string][]edge {
	sources := make(map[string]bool)
	sinks := make(map[string]bool)

	for node := range cap_graph {
		fmt.Println(node)
		if len(node) >=6 && node[0:6] == "source" {
			sources[node] = true
		}
		if len(node) >=4 &&node[0:4] == "sink"{
			sinks[node] = true
		}
	}

	fmt.Println(sources)
	fmt.Println(sinks)

	var aug_flow map[string][]edge


	return aug_flow
}