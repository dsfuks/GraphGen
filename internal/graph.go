package GraphGen

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"time"
)

// Количество числовых атрибутов в узлах и ребрах.
const numAttr int = 5

type Node struct {
	Attributes [numAttr]int `json:"attributes"`
	Id         int          `json:"id"`
}

//	Nodes - массив узлов, которые соединяет ребро
type Edge struct {
	Attributes [numAttr]int `json:"attributes"`
	Nodes      [2]*Node     `json:"nodes"`
}

type Graph struct {
	Nodes []Node `json:"Nodes"`
	Edges []Edge `json:"Edges"`
}

func (g *Graph) NewGraph() {
	rand.Seed(time.Now().UnixNano())
	numOfNodes := 5 + rand.Intn(6)
	setNodes := map[*Node]struct{}{}
	// Множество для проверки на кратные ребра.
	setConnections := map[string]struct{}{}
	// Число ребер.
	numOfEdges := numOfNodes - 1
	// Добавляем узлы в граф.
	for i := 0; i < numOfNodes; i++ {
		var tempNode Node
		tempNode.Id = i + 1
		// Заполнение узла атрибутами.
		for k := range tempNode.Attributes {
			tempNode.Attributes[k] = rand.Intn(1000)
		}
		g.Nodes = append(g.Nodes, tempNode)
		setNodes[&g.Nodes[len(g.Nodes)-1]] = struct{}{}
	}
	// Добавляем ребра.
	for i := 0; i < numOfEdges; i++ {
		var tempEdge Edge
		count := 0
		// Соединяем узлы, избегая кратных ребер.
		for key := range setNodes {
			switch count {
			case 0:
				tempEdge.Nodes[count] = key
			case 1:
				if _, ok := setConnections[fmt.Sprint(tempEdge.Nodes[0], key)]; !ok {
					tempEdge.Nodes[count] = key
					setConnections[fmt.Sprint(tempEdge.Nodes[0], key)] = struct{}{}
				} else {
					continue
				}
			default:
				break
			}
			count++
		}
		// Добавляем атрибуты в ребро.
		for k := range tempEdge.Attributes {
			tempEdge.Attributes[k] = rand.Intn(1000)
		}
		g.Edges = append(g.Edges, tempEdge)
	}
}
func (g *Graph) Handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		g.getGraph(w)
	} else {
		http.Error(w, "", http.StatusNotFound)
	}
}
func (g *Graph) getGraph(w http.ResponseWriter) {
	js, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err2 := w.Write(js); err2 != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

}

// Принимает слайс Node или Edge. Возвращает max, min, среднее значение атрибутов, медиану.
func attrCalculation(a interface{}) (max int, min int, average float64, median float64) {
	allAttributes := make([]int, 0)
	max = math.MinInt
	min = math.MaxInt
	switch a.(type) {
	case []Node:
		for _, v := range a.([]Node) {
			for _, k := range v.Attributes {
				allAttributes = append(allAttributes, k)
				if min > k {
					min = k
				}
				if max < k {
					max = k
				}
				average += float64(k)
			}
		}
	case []Edge:
		for _, v := range a.([]Edge) {
			for _, k := range v.Attributes {
				allAttributes = append(allAttributes, k)
				if min > k {
					min = k
				}
				if max < k {
					max = k
				}
				average += float64(k)
			}
		}

	}
	average /= float64(len(allAttributes))
	sort.Ints(allAttributes)
	if len(allAttributes)%2 == 0 {
		median = float64(allAttributes[len(allAttributes)/2]+allAttributes[len(allAttributes)/2-1]) / 2
	} else {
		median = float64(allAttributes[len(allAttributes)/2])
	}
	return max, min, average, median
}

// NodeAttributesInfo Информация по атрибутам узлов.
func (g *Graph) NodeAttributesInfo() string {
	s := "Сводная информация по атрибутам узлов:\n"
	max, min, average, median := attrCalculation(g.Nodes)
	s += fmt.Sprintf("Максимум: %v, минимум: %v, медиана: %v, среднее значение: %v", max, min, median, average)
	return s
}

// EdgeAttributesInfo Информация по атрибутам ребер.
func (g *Graph) EdgeAttributesInfo() string {
	s := "Сводная информация по атрибутам ребер:\n"
	max, min, average, median := attrCalculation(g.Edges)
	s += fmt.Sprintf("Максимум: %v, минимум: %v, медиана: %v, среднее значение: %v", max, min, median, average)
	return s
}
