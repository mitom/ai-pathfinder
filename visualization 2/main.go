package main;

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "math"
    "flag"
    "os"
    "bufio"
    "encoding/json"
    "log"
)

//Represent a cavern linking to other caverns it can reach
type Cavern struct {
    Id    int `json:"i"`
    X     int `json:"x"`
    Y     int `json:"y"`
    paths []*Cavern
}

//Represent a visited node in the graph with the underlying cavern and the heuristic cost
type Node struct {
    parent *Node
    cavern *Cavern
    Cost   float64 `json:"c"`
    Id     int `json:"i"`
}

//Only required for visualizing the cave
type Edge struct {
    From int `json:"f"`
    To int `json:"t"`
}

//Only required for visualizing the cave
type VisualCave struct {
    Nodes []*Cavern `json:"nodes"`
    Edges []Edge `json:"edges"`
}


//Check whether the given list contains a certain node by evaluating if it relies on the same cave
//Return the node if found, or nil if not
func contains_node(s []*Node, e Node) *Node {
    for _, a := range s {
        if a.cavern.equals(*e.cavern) {
            return a
        }
    }
    return nil
}

//check whether 2 caverns are the same
func (c Cavern) equals(t Cavern) bool {
    return c.Id == t.Id
}

//Calculate the distance between 2 caverns
func (c Cavern) distance_to(t Cavern) float64 {
    return math.Sqrt(math.Pow(float64(c.X - t.X), 2) + math.Pow(float64(c.Y - t.Y), 2))
}

//panic helper for errors
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// read in a file and return it as a string
func read_file(file string) string {
    dat, err := ioutil.ReadFile(file)
    check(err)

    return string(dat[:])
}

//build the structs for the cave based on the input
func build_cave(input []string, visualize bool) (*Cavern, *Cavern) {
    caverns_num, _ := strconv.Atoi(input[0])
    caverns := make([]*Cavern, caverns_num)

    visualization := VisualCave{nil, make([]Edge, 0)}

    //the loop must step by 2 as there are 2 coordiantes
    for i := 1; i <= caverns_num * 2; i = i + 2 {
        x, _ := strconv.Atoi(input[i])
        y, _ := strconv.Atoi(input[i + 1])

        //the cavern id starts from 1, the list is indexed from 0
        caverns[(i - 1) / 2] = &Cavern{(i - 1) / 2 + 1, x, y, make([]*Cavern, 0)}
    }

    visualization.Nodes = caverns
    //get rid of everything in the input we have processed so far
    input = input[caverns_num * 2 + 1:]

    for i := 0; i < len(input); i++ {
        //only need to process connections
        if input[i] == "1" {
            //due to the matrix the remainder is the index of the source and the result int is the index of the target
            from := caverns[i % caverns_num]
            to := caverns[i / caverns_num]
            from.paths = append(from.paths, to)

            if visualize {
                visualization.Edges = append(visualization.Edges, Edge{from.Id, to.Id})
            }
        }
    }
    //dump the visualization data
    if visualize {
        strB, _ := json.Marshal(visualization)
        err := ioutil.WriteFile("visualization/instructions.js", []byte(fmt.Sprintf("setCave(%v)", string(strB))), 0644)
        check(err)
    }

    //return the first and last cavern for start ang goal
    return caverns[0], caverns[len(caverns) - 1]

}

//recursively calculate the distance for a node through its parents
func calculate_path_distance(node Node, distance float64) float64 {
    if node.parent == nil {
        return distance
    }

    return calculate_path_distance(*node.parent, distance + node.cavern.distance_to(*node.parent.cavern))
}

//recursively build the path for this node from the start
func build_path(node Node, path string) string {
    if node.parent == nil {
        return strconv.Itoa(node.cavern.Id) + path
    }

    return build_path(*node.parent, "->" + strconv.Itoa(node.cavern.Id) + path)
}

//represent a node when printed to the console
func display_node(node *Node) string {
    return fmt.Sprintf("Cavern %v(heur. cost: %v)", node.cavern.Id, node.Cost)
}

//dump debug information to the console
func dump(iteration int, open_list []*Node, closed_list []*Node) {
    fmt.Printf("Iteration: %v\n", iteration)
    fmt.Println("Open list:")
    open := make([]string, len(open_list))
    for i := 0; i < len(open_list); i++ {
        open[i] = display_node(open_list[i])
    }
    fmt.Println(open)
    fmt.Println("Closed list:")
    closed := make([]string, len(closed_list))
    for i := 0; i < len(closed_list); i++ {
        closed[i] = display_node(closed_list[i])
    }
    fmt.Println(closed)
    fmt.Println("========================")
}



//A* algorithm
func search(start Cavern, goal Cavern, verbose bool, visualize bool) {
    fmt.Println("Starting search..")
    open_list := make([]*Node, 0)
    closed_list := make([]*Node, 0)
    //add the starting node to the open list
    open_list = append(open_list, &Node{nil, &start, start.distance_to(goal)*2, start.Id})
    //count the iterations for debugging purposes
    iteration := 1

    file, err := os.OpenFile("visualization/instructions.js", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Cannot create file", err)
    }
    defer file.Close()

    //loop while we have any options left
    for len(open_list) > 0 {
        //find cheapest node based on the heuristic cost
        var current *Node
        var current_i int
        for i := 0; i < len(open_list); i++ {
            if nil == current || current.Cost > open_list[i].Cost {
                current = open_list[i]
                current_i = i
            }
        }
        //remove the node from the open list
        open_list = append(open_list[:current_i], open_list[current_i + 1:]...)

        //evaluate all children whether they are the goal or if they are a potential candidate
        for i := 0; i < len(current.cavern.paths); i++ {
            child := current.cavern.paths[i]
            child_node := Node{current, child, calculate_path_distance(*current, 0.0) + child.distance_to(*current.cavern) + child.distance_to(goal)*2, child.Id}

            //check if the target is the goal
            if child.equals(goal) {
                if verbose {
                    dump(iteration, open_list, closed_list);
                }

                if visualize {
                    curr_json, _ := json.Marshal(child_node)
                    opn_json, _ := json.Marshal(open_list)
                    clsd_json, _ := json.Marshal(closed_list)
                    file.WriteString(fmt.Sprintf("\naddState(%v, \"%v\", %v, %v)", string(curr_json), build_path(child_node, ""), string(opn_json), string(clsd_json)))
                    file.WriteString("\ninit()")
                }

                // if its the goal, we are done
                fmt.Printf("Found path: %v\n", build_path(child_node, ""))
                fmt.Printf("Total distance: %v\n", calculate_path_distance(child_node, 0.0))
                fmt.Printf("Iterations: %v\n", iteration)
                return
            }

            //If we already know a shorter path to this node, skip it
            in_open := contains_node(open_list, child_node)
            if in_open != nil && in_open.Cost < child_node.Cost {
                continue
            }

            //If we already evaluated this node via a shorter path and dismissed it, skip it
            in_closed := contains_node(closed_list, child_node)
            if in_closed != nil && in_closed.Cost < child_node.Cost {
                continue
            }

            //add the child as a candidate to the open list
            open_list = append(open_list, &child_node)
        }

        if verbose {
            fmt.Println("Considering:")
            fmt.Println(display_node(current))
            fmt.Printf("Current distance on the path: %v\n", calculate_path_distance(*current, 0.0))
            dump(iteration, open_list, closed_list)
            fmt.Println("Press enter to continue..")
            bufio.NewReader(os.Stdin).ReadBytes('\n')
        }

        if visualize {
            curr_json, _ := json.Marshal(current)
            opn_json, _ := json.Marshal(open_list)
            clsd_json, _ := json.Marshal(closed_list)
            file.WriteString(fmt.Sprintf("\naddState(%v, \"%v\", %v, %v)", string(curr_json), build_path(*current, ""), string(opn_json), string(clsd_json)))
        }

        //add the current node to the closed list
        closed_list = append(closed_list, current)
        iteration++
    }

    if visualize {
        file.WriteString("\ninit()")

    }

    fmt.Println("Could not find a path")
    fmt.Printf("Iterations: %v\n", iteration)
    return
}

func main() {
    verbose := flag.Bool("verbose", false, "Wether or not to log verbosely")
    visualize := flag.Bool("visualize", false, "Wether or not to generate visualization")
    input := flag.String("input", "input.cav", "Input file")
    flag.Parse()
    result := strings.Split(read_file(*input), ",")
    start, goal := build_cave(result, *visualize)

    search(*start, *goal, *verbose, *visualize)
}