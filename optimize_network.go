package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main()  {
	res := json_to_map("capacity_graph.json")

	aug_flow := cap_to_aug_flow(res)

	fmt.Println("Augmented Flow Graph Init")
	fmt.Println("--------------------------")
	for node := range aug_flow {
		fmt.Println(node)
		fmt.Println(aug_flow[node])
	}

	fmt.Println("DFS Path")
	fmt.Println("--------------------------")
	path := []string{}
	fmt.Println(dfs(aug_flow, "main_source", "main_sink", path))
}

type edge struct {
	dest string
	cap int
	flow int
}

func dfs(graph map[string][]edge, cur string, end string, path []string) []string{
	if len(graph[cur]) == 0 {
		path = nil
		return path
	}
	path = append(path, cur)
	for i := range graph[cur] {
		edg := graph[cur][i]
		if edg.dest == end {
			path = append(path, edg.dest)
			return path
		}
		out := dfs(graph, edg.dest, end, path)
		if out != nil {
			return out
		}
	}
	return nil
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
	substations := make(map[string]bool)
	sinks := make(map[string]bool)

	for node := range cap_graph {
		if len(node) >=6 && node[0:6] == "source" {
			sources[node] = true
		}
		if len(node) >= 10 && node[0:10] == "substation"{
			substations[node] = true
		}
		if len(node) >=4 &&node[0:4] == "sink"{
			sinks[node] = true
		}
	}


	aug_flow := make(map[string][]edge)
	
	// Add sources and main source to graph
	for source := range sources {
		// Connect sources to main source
		src := cap_graph[source].(map[string]interface{})
		cap := int(src["produced"].(float64))
		edg := edge{dest: source, cap: cap, flow: 0}
		aug_flow["main_source"] = append(aug_flow["main_source"], edg)

		// Connect sources to substations
		edges := src["edges"].([]interface{})

		for i := range edges {
			eg := edges[i].(map[string]interface{})
			dest := string(eg["dest"].(string))
			cap = int(eg["cap"].(float64))
			flow := int(eg["flow"].(float64))
			edg = edge{dest: dest, cap: cap, flow: flow}
			aug_flow[source] = append(aug_flow[source], edg)
		}
	}
	
	// Add substations to graph
	for substation := range substations {
		edges := cap_graph[substation].([]interface{})
		for i := range edges {
			eg := edges[i].(map[string]interface{})
			dest := string(eg["dest"].(string))
			cap := int(eg["cap"].(float64))
			flow := int(eg["flow"].(float64))
			edg := edge{dest: dest, cap: cap, flow: flow}
			aug_flow[substation] = append(aug_flow[substation], edg)
		}
	}

	// Connect sinks to main sink
	for sink := range sinks {
		cap := int(cap_graph[sink].(float64))
		edge := edge{dest: "main_sink", cap: cap, flow: 0}
		aug_flow[sink] = append(aug_flow["main_sink"], edge)
	}

	return aug_flow
}