package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PaesslerAG/gval"
)

type Invocation struct {
	Values map[string]interface{} `json:"values"`
}

type Result struct {
	Value interface{} `json:"value"`
}

type Step struct {
	Kind       string    `json:"kind"`
	Expression string    `json:"expression"`
	Output     string    `json:"output"`
	Algorithm  Algorithm `json:"algorithm"`
}

type Algorithm struct {
	Steps []Step `json:"steps,omitempty"`
}

func execute(invocation Invocation, algorithm Algorithm) Result {
	var value interface{}
	var err error
	for _, step := range algorithm.Steps {
		value, err = gval.Evaluate(step.Expression, invocation.Values)
		if err != nil {
			fmt.Println("Error:", err)
			return Result{}
		}

		if step.Kind == "conditional" {
			fmt.Printf("conditional: %s\n", value == true)
			if value == true {
				value = execute(invocation, step.Algorithm).Value

				fmt.Printf("%s: %s = %d\n", step.Output, step.Expression, value)
				if step.Output != "" {
					invocation.Values[step.Output] = value
				} else {
					return Result{Value: value}
				}

			}
		} else {

			fmt.Printf("%s: %s = %d\n", step.Output, step.Expression, value)
			if step.Output != "" {
				invocation.Values[step.Output] = value
			} else {
				return Result{Value: value}
			}

		}

	}
	return Result{Value: value}
}

func main() {

	file, err := os.OpenFile("algo.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	var mainAlgo Algorithm
	err = decoder.Decode(&mainAlgo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	invocation := Invocation{
		Values: map[string]interface{}{
			"nodes": 100,
		},
	}

	result := execute(invocation, mainAlgo)
	fmt.Printf("result: %s\n\n", result.Value)

	// str, err := json.Marshal(mainAlgo)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }
	// fmt.Printf("%s", str)
}
