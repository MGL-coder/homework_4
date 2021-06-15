package tetris

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unsafe"
)

// wordSize equals to the size of one word on the computer
// types contains type names and corresponding sizes in byte
const wordSize = int(unsafe.Sizeof(int(0)))

var types = map[string]int{"int8": 1, "int16": 2, "int32": 4, "int64": 8, "int": wordSize,
	"uint8": 1, "uint16": 2, "uint32": 4, "uint64": 8, "uint": wordSize,
	"float32": 4, "float64": 8, "complex64": 8, "complex128": 16, "byte": 1, "rune": 4,
	"uintptr": wordSize, "string": 2 * wordSize,
}

// checkFields checks fields for errors
func checkFields(fields []string) error {
	for _, field := range fields {
		if len(field) == 0 {
			return fmt.Errorf("please delete empty line in the struct body\n")
		}

		if len(strings.Split(field, " ")) < 2 {
			return fmt.Errorf("invalid field in the struct body %s\n", field)
		}
		typeName := strings.Split(field, " ")[1]

		if typeName[0] == '*' && len(typeName) <= 1 {
			return fmt.Errorf("invalid type '*'\n")
		}

		if _, ok := types[typeName]; !ok && typeName[0] != '*' {
			return fmt.Errorf("could not recognize type %s\n", typeName)
		}
	}
	return nil
}

// getSize return size of the given field in bytes
func getSize(field string) int {
	typeName := strings.Split(field, " ")[1]

	if typeName[0] == '*' && len(typeName) > 1 {
		return wordSize
	}

	size, _ := types[typeName]
	return size
}

// solveTetris solves tetris problem by greedy algorithm (for any number of fields)
func solveTetris(fields []string) []string {
	sort.Slice(fields, func(i, j int) bool {
		return getSize(fields[i]) > getSize(fields[j])
	})

	return fields
}

// solveTetrisBruteForce solves tetris problem by brute force (for any number of fields)
// returns 3 best results
func solveTetrisBruteForce(fields []string) [][]string {
	best := make([][]string, 3)
	min := make([]int, 3)
	for i := range best {
		best[i] = make([]string, len(fields))
		min[i] = 100000
	}

	// permutations by Heap's Algorithm
	var generate func(int, []string)
	k := len(fields)
	generate = func(k int, s []string) {
		if k == 1 {
			// updating min and best structs
			memory := calculateMemory(s)
			if memory < min[0] {
				copy(best[2], best[1])
				copy(best[1], best[0])
				min[2], min[1] = min[1], min[0]
				min[0] = memory
				copy(best[0], s)
			} else if memory < min[1] {
				copy(best[2], best[1])
				min[2] = min[1]
				min[1] = memory
				copy(best[1], s)
			} else if memory < min[2] {
				min[2] = memory
				copy(best[2], s)
			}
		} else {
			generate(k-1, s)

			for i := 0; i < k-1; i++ {
				if k%2 == 0 {
					s[i], s[k-1] = s[k-1], s[i]
				} else {
					s[0], s[k-1] = s[k-1], s[0]
				}
				generate(k-1, s)
			}
		}
	}

	f := make([]string, len(fields))
	copy(f, fields)
	generate(k, f)

	return best
}

// calculateMemory calculates and returns the memory used by struct with given fields
func calculateMemory(fields []string) int {
	sizes := make([]int, len(fields), 6)
	for i, field := range fields {
		sizes[i] = getSize(field)
	}

	capMul := 1
	cap := 0
	mem := 0
	for _, size := range sizes {
		switch size {
		case 1:
			mem++
		case 2:
			if capMul < 2 {
				capMul = 2
			}
			mem += mem%2 + 2
		case 4:
			if capMul < 4 {
				capMul = 4
			}
			if mem%4 == 0 {
				mem += 4
			} else {
				mem += 8 - mem%4
			}
		case 8:
			if capMul < 8 {
				capMul = 8
			}
			if mem%8 == 0 {
				mem += 8
			} else {
				mem += 16 - mem%8
			}
		case 16:
			if capMul < 8 {
				capMul = 8
			}
			if mem%8 == 0 {
				mem += 16
			} else {
				mem += 24 - mem%8
			}
		}

		if mem > cap {
			if mem%capMul == 0 {
				cap = mem
			} else {
				cap = (mem/capMul + 1) * capMul
			}
		}
	}

	return cap
}

// printStruct prints struct fields and the amount of memory it uses
func printStruct(fields []string) {
	memory := calculateMemory(fields)

	for _, field := range fields {
		fmt.Println("\t" + field)
	}

	fmt.Printf("Memory used = %v bytes\n\n", memory)
}

// Tetris prints struct fields of file at the given path,
// Optimizes struct by greedy algorithm and prints result
// Optimizes struct by brute force and prints 3 best solutions
// Writes optimized struct by greedy algorithm into the initial file
func Tetris(path string) (funcErr error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file: %s", err)
	}

	fields := make([]string, 0, 6) // fields of the struct

	// saving lines of the file and fields of the struct
	lines := strings.Split(string(input), "\n")
	inStructBlock := false
	var structLine int
	for i, line := range lines {
		if strings.Contains(line, " struct {") {
			structLine = i + 1
			inStructBlock = true
		} else if strings.Contains(line, "}") && inStructBlock {
			break
		} else if inStructBlock {
			fields = append(fields, strings.TrimSpace(line))
		}
	}

	// checking field
	if err := checkFields(fields); err != nil {
		return fmt.Errorf("struct field error: %s", err)
	}

	// printing initial struct
	fmt.Println("TETRIS:")
	fmt.Println("Initial struct:")
	printStruct(fields)

	// solving the tetris problem by greedy algorithm
	solveTetris(fields)

	// printing best solution
	fmt.Println("Best solution by greedy algorithm:")
	printStruct(fields)

	// solving the tetris problem by brute force
	fmt.Println("Top 3 solutions by brute force:")
	structs := solveTetrisBruteForce(fields)
	printStruct(structs[0])
	printStruct(structs[1])
	printStruct(structs[2])

	// writing the optimized struct into file
	for i, field := range fields {
		lines[i+structLine] = "\t" + field
	}
	output := strings.Join(lines, "\n")

	err = os.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("could write to the file: %s", err)
	}

	return nil
}
