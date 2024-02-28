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
	s1 := res["source1"].(map[string]interface{})
	fmt.Println(s1)
	prod1 := s1["produced"]
	fmt.Println(prod1)

	aug_flow := cap_to_aug_flow(res)
	fmt.Println(aug_flow)
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
		if len(node) >=6 && node[0:6] == "source" {
			sources[node] = true
		}
		if len(node) >=4 &&node[0:4] == "sink"{
			sinks[node] = true
		}
	}

	aug_flow := make(map[string][]edge)
	
	for source := range sources {
		// Connect sources to main source
		src := cap_graph[source].(map[string]interface{})
		cap := int(src["produced"].(float64))
		edg := edge{dest: source, cap: cap, flow: 0}
		aug_flow["main_source"] = append(aug_flow["main_source"], edg)

		// Connect sources to substations
		edges := src["edges"].([]interface{})
		fmt.Println(edges)
		for i := range edges {
			eg := edges[i].(map[string]interface{})
			dest := string(eg["dest"].(string))
			cap = int(eg["cap"].(float64))
			flow := int(eg["flow"].(float64))
			edg = edge{dest: dest, cap: cap, flow: flow}
			aug_flow[source] = append(aug_flow[source], edg)
		}
	}
	

	return aug_flow
}