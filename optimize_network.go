package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"sort"
)

func main()  {
	res := json_to_map("capacity_graph.json")

	aug_flow := cap_to_aug_flow(res)

	node_keys := []string{}
	for key := range aug_flow {
		node_keys = append(node_keys, key)
	}
	sort.Slice(node_keys, func(i, j int) bool {
		return node_keys[i] > node_keys[j]
	})
	
	fmt.Println("Augmented Flow Graph Init")
	fmt.Println("--------------------------")
	for _ , node := range node_keys {
		fmt.Println(node)
		fmt.Println(aug_flow[node])
	}
	
	
	fmt.Println("\nFinal Flow Graph")
	fmt.Println("--------------------------")
	final_flow := solve(aug_flow)
	for _ , node := range node_keys {
		fmt.Println(node)
		fmt.Println(final_flow[node])
	}
	
}

type edge struct {
	dest string
	cap int
	flow int
	residual bool
}

func dfs(graph map[string][]edge, cur string, end string, path []string) []string{
	path = append(path, cur)
	for i := range graph[cur] {
		edg := graph[cur][i]
		visited := false
		for _ , node := range path {
			if edg.dest == node {
				visited = true
				break
			}
		}
		if visited {
			continue
		}
		
		if edg.flow >= edg.cap {
			continue
		}
		if edg.dest == "main_source" {
			return nil
		}
		if edg.dest == end {
			path = append(path, end)
			return path
		}
		out := dfs(graph, edg.dest, end, path)
		if out != nil {
			return out
		}
	}
	return nil
}

func solve(init_graph map[string][]edge) map[string][]edge {
	path := []string{}
	aug_flow := init_graph
	
	for {
		path = dfs(aug_flow, "main_source", "main_sink", path)
		if path == nil {
			break
		}
		min_dif := 1000000000
		for i := 0; i < len(path) - 1; i++ {
			edg := edge{dest: "temp", cap: 0, flow: 0}
			for e := range aug_flow[path[i]] {
				if aug_flow[path[i]][e].dest == path[i + 1] {
					edg = aug_flow[path[i]][e]
					break
				}
			}
			dif := edg.cap - edg.flow
			if min_dif > dif {
				min_dif = dif
			}
		}
		for i := 0; i < len(path) - 1; i++ {
			for e := range aug_flow[path[i]] {
				edg := &aug_flow[path[i]][e]
				if edg.dest == path[i + 1] {
					if edg.residual {
						edg.cap -= min_dif
						break
					}
					edg.flow += min_dif
					for r := range aug_flow[edg.dest] {
						if path[i] == aug_flow[edg.dest][r].dest {
							aug_flow[edg.dest][r].cap += min_dif
						}
					}
					break
				}
			}
		}
		path = nil
	}

	return aug_flow
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
			edg = edge{dest: dest, cap: cap, flow: flow, residual: false}
			resid := edge{dest: source, cap: 0, flow: 0, residual: true}
			aug_flow[source] = append(aug_flow[source], edg)
			aug_flow[dest] = append(aug_flow[dest], resid)
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
			edg := edge{dest: dest, cap: cap, flow: flow, residual: false}
			resid := edge{dest: substation, cap: 0, flow: 0, residual: true}
			aug_flow[substation] = append(aug_flow[substation], edg)
			aug_flow[dest] = append(aug_flow[dest], resid)
		}
	}


	// Connect sinks to main sink
	for sink := range sinks {
		cap := int(cap_graph[sink].(float64))
		edg := edge{dest: "main_sink", cap: cap, flow: 0, residual: false}
		resid := edge{dest: sink, cap: 0, flow: 0, residual: true}
		aug_flow[sink] = append(aug_flow[sink], edg)
		aug_flow["main_sink"] = append(aug_flow["main_sink"], resid)
	}

	return aug_flow
}