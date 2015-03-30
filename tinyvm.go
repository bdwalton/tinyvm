package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MEM_SIZE = 1024 // How much data memory? (Used to size both data and program memory.)
const NUM_REGS = 8    // The total number of registers available.
const PC_REG = 7      // The registered used as the program counter.

// Instructions are composed of one operation and up to three
// arguments.
type TinyInstruction struct {
	iop   string
	iargs []int
}

type TinyCPUState int

const (
	cpuOK TinyCPUState = iota
	cpuHALTED
	cpuDIV_ZERO
	cpuIMEM_ERR
	cpuDMEM_ERR
)

/* A structure representing a tiny machine */
type TinyMachine struct {
	stdin              *bufio.Reader             // To handle data input
	registers          [NUM_REGS]int             // 8 registers
	data_memory        [MEM_SIZE]int             // Data memory
	instruction_memory [MEM_SIZE]TinyInstruction // Instruction memory
	trace              bool                      // Output instructions as they're executed
	cpustate           TinyCPUState              // See cpu* constants above
}

// Operands are of the form r,s,t where r, s and t are all integers
func parseROop(args string) ([]int, error) {
	string_args := strings.Split(args, ",")
	converted_args := make([]int, 3)

	if len(string_args) != 3 {
		return nil, errors.New("Invalid arguments: " + args)
	} else {
		for i := 0; i < 3; i++ {
			num, err := strconv.Atoi(string_args[i])
			if err != nil {
				return nil, errors.New("Invalid arguments: " + args)
			} else {
				converted_args[i] = num
			}
		}
	}

	return converted_args, nil
}

// Operands are of the form r,s(t) where r, s and t are all integers
func parseRMop(args string) ([]int, error) {
	converted_args := make([]int, 3)

	x := strings.Index(args, ",")
	y := strings.Index(args, "(")
	z := strings.Index(args, ")")

	if x < 1 || y < x || z < y {
		return nil, errors.New("Invalid arguments: " + args)
	} else {
		indexes := [][]int{[]int{0, x}, []int{x + 1, y}, []int{y + 1, z}}

		for i, bounds := range indexes {
			num, err := strconv.Atoi(args[bounds[0]:bounds[1]])
			if err != nil {
				return nil, errors.New("Invalid arguments: " + args)
			} else {
				converted_args[i] = num
			}
		}
	}

	return converted_args, nil
}

func parseInstruction(line string) (TinyInstruction, error) {
	var args []int
	var err error
	var ti TinyInstruction

	// Chop the newline off and then split on spaces
	r := regexp.MustCompile(" +")
	stripped_line := strings.TrimSpace(r.ReplaceAllString(line, " "))
	line_parts := strings.Split(stripped_line, " ")

	if len(line_parts) != 2 {
		return ti, errors.New("Invalid instruction: '" + stripped_line + "'")
	} else {
		switch line_parts[0] {
		case "HALT", "IN", "OUT", "ADD", "SUB", "MUL", "DIV":
			args, err = parseROop(line_parts[1])
		case "LD", "ST", "LDA", "LDC", "JLT", "JLE", "JGT", "JGE", "JEQ", "JNE":
			args, err = parseRMop(line_parts[1])
		default:
			return ti, errors.New("Invalid opcode: '" + line_parts[0] + "'")
		}

		if err != nil {
			m := "Invalid arguments for opcode " + line_parts[0] + ": '" + line_parts[1] + "'"
			return ti, errors.New(m)
		} else {
			ti.iop = line_parts[0]
			ti.iargs = args
		}
	}

	return ti, nil
}

func (tm *TinyMachine) initializeMachine(clearprogram bool) {
	for i := 0; i < NUM_REGS; i++ {
		tm.registers[i] = 0
	}

	for i := 0; i < MEM_SIZE; i++ {
		tm.data_memory[i] = 0
	}

	if clearprogram {
		for i := 0; i < MEM_SIZE; i++ {
			tm.instruction_memory[i] = TinyInstruction{"HALT", []int{0, 0, 0}}
		}
	}

	// Store the size of the memory in the first memory element.
	tm.data_memory[0] = MEM_SIZE - 1
	tm.cpustate = cpuOK
	tm.registers[PC_REG] = 0
	tm.stdin = bufio.NewReader(os.Stdin) // An io helper.
}

// Leave the loaded program intact, but re-initialize the machine to a
// clean state otherwise.
func (tm *TinyMachine) resetState() {
	// Reset memory and registers, but leave program intact.
	tm.initializeMachine(false)
}

func (tm *TinyMachine) stepProgram() {
	if tm.cpustate != cpuOK {
		tm.handleCpuState()
		return
	}

	pc := tm.registers[PC_REG]
	if pc < 0 || pc > MEM_SIZE-1 {
		tm.cpustate = cpuIMEM_ERR
	} else {
		// Step the program counter
		tm.registers[PC_REG] = pc + 1

		instruction := tm.instruction_memory[pc]
		if tm.trace {
			fmt.Println("Executing:", instruction)
		}

		r := instruction.iargs[0]
		s := instruction.iargs[1]
		t := instruction.iargs[2]
		a := s + tm.registers[t]

		switch instruction.iop {
		case "HALT":
			tm.cpustate = cpuHALTED
		case "IN":
			m := fmt.Sprintf("Enter number to store in register %d", r)
			n := tm.readNumber(m, 0)
			tm.registers[r] = n
		case "OUT":
			fmt.Println(tm.registers[r])
		case "ADD":
			tm.registers[r] = tm.registers[s] + tm.registers[t]
		case "SUB":
			tm.registers[r] = tm.registers[s] - tm.registers[t]
		case "MUL":
			tm.registers[r] = tm.registers[s] * tm.registers[t]
		case "DIV":
			if tm.registers[t] == 0 {
				tm.cpustate = cpuDIV_ZERO
			} else {
				tm.registers[r] = tm.registers[s] / tm.registers[t]
			}
		case "LDA":
			if a < 0 || a > MEM_SIZE {
				tm.cpustate = cpuDMEM_ERR
			} else {
				tm.registers[r] = a
			}
		case "LDC":
			tm.registers[r] = s
		case "LD":
			if a < 0 || a > MEM_SIZE {
				tm.cpustate = cpuDMEM_ERR
			} else {
				tm.registers[r] = tm.data_memory[a]
			}
		case "ST":
			if a < 0 || a > MEM_SIZE {
				tm.cpustate = cpuDMEM_ERR
			} else {
				tm.data_memory[a] = tm.registers[r]
			}
		case "JLT":
			if tm.registers[r] < 0 {
				tm.registers[PC_REG] = a
			}
		case "JLE":
			if tm.registers[r] <= 0 {
				tm.registers[PC_REG] = a
			}
		case "JGE":
			if tm.registers[r] >= 0 {
				tm.registers[PC_REG] = a
			}
		case "JGT":
			if tm.registers[r] > 0 {
				tm.registers[PC_REG] = a
			}
		case "JEQ":
			if tm.registers[r] == 0 {
				tm.registers[PC_REG] = a
			}
		case "JNE":
			if tm.registers[r] != 0 {
				tm.registers[PC_REG] = a
			}
		}
	}

	tm.handleCpuState()
}

func (tm *TinyMachine) handleCpuState() {
	switch tm.cpustate {
	case cpuOK:
		break
	case cpuDIV_ZERO:
		fmt.Println("Divide by zero error.")
		fallthrough
	case cpuIMEM_ERR:
		fmt.Println("Instruction memory access violation.")
		fallthrough
	case cpuDMEM_ERR:
		fmt.Println("Data memory access violation.")
		fallthrough
	case cpuHALTED:
		fallthrough
	default:
		fmt.Println("Program halted.")
	}
}

func (tm *TinyMachine) runProgram() {
	for {
		tm.stepProgram()
		if tm.cpustate != cpuOK {
			break
		}
	}
}

func (tm *TinyMachine) loadProgram(progfilename string) bool {
	var (
		i       int
		linenum int = 0
	)

	tm.initializeMachine(true)

	programfile, err := os.Open(progfilename)

	if err != nil {
		fmt.Println("Error:", err)
		return false
	} else {
		defer programfile.Close()
		reader := bufio.NewReader(programfile)
		fmt.Println("Reading program from", progfilename)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println("Error reading program:", err)
					return false
				}
			} else {
				linenum++
				if strings.Index(line, "*") == 0 {
					// Comments are lines starting with an asterisk
					continue
				} else {
					// TODO(bdwalton): Skip over blank lines.
					instruction, err := parseInstruction(line[:len(line)-1])

					if err != nil {
						fmt.Println(err)
						fmt.Printf("Error parsing program at line %d: %s\n", linenum, line)
						return false
					} else {
						tm.instruction_memory[i], i = instruction, i+1
					}
				}
			}
		}
	}

	return true
}

func (tm *TinyMachine) dumpRegisters() {
	fmt.Println("Current Tiny Machine register values:")

	for i := 0; i < NUM_REGS; i++ {
		switch i {
		case PC_REG:
			fmt.Println("PC:", tm.registers[i])
		default:
			fmt.Printf("%2d: %d\n", i, tm.registers[i])
		}
	}
}

func (tm *TinyMachine) dumpMemory(start_addr, end_addr int) {
	fmt.Println("Dumping data memory from address %d to %d", start_addr, end_addr)

	for i := start_addr; i <= end_addr; i++ {
		fmt.Printf("%04d: %d\n", i, tm.data_memory[i])
	}
}

func (tm *TinyMachine) dumpProgram(start_addr, end_addr int) {
	fmt.Printf("Dumping instruction memory from address %d to %d.\n", start_addr, end_addr)

	for i := start_addr; i <= end_addr; i++ {
		fmt.Printf("%04d: %v\n", i, tm.instruction_memory[i])
	}
}

func (tm *TinyMachine) readNumber(prompt string, def int) int {
	for {
		fmt.Printf("%s: ", prompt)
		input, err := tm.stdin.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Returning default", def)
			break
		} else {
			num, err := strconv.Atoi(input[:len(input)-1])
			if err != nil {
				fmt.Println("Error converting input. Returning default", def)
				break
			} else {
				return num
			}
		}
	}

	return def
}

func (tm *TinyMachine) Interact() {
	menu := []struct {
		key      string // The key to press to activate this option.
		helptext string // The help text displayed
	}{
		{"c", "clear machine state"},
		{"g", "go - run program to halt state"},
		{"h", "display this help text"},
		{"d", "display data memory"},
		{"i", "display instruction memory"},
		{"q", "quit the tiny machine simulator"},
		{"r", "dump register contents"},
		{"s", "step program forward by one instruction"},
		{"t", "toggle execution tracing"},
	}

	fmt.Println("Tiny Machine simulation (enter h for help)")
interactive:
	for {
		fmt.Printf("Enter command: ")
		input, err := tm.stdin.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Fake up a real "q" entry so we handle eof the same way as a normal
				// exit.
				fmt.Println()
				input = "q\n"
			} else {
				// This will be handled with the unknown case below.
				input = "ijustmashedthekeyboard"
			}
		}

		command := input[:len(input)-1]

		switch command {
		case "c":
			tm.resetState()
		case "g":
			tm.runProgram()
		case "h":
			for _, menuitem := range menu {
				fmt.Printf("%s: %s\n", menuitem.key, menuitem.helptext)
			}
		case "d":
			start_addr := tm.readNumber("Starting Address", 0)
			end_addr := tm.readNumber("Ending Address", MEM_SIZE-1)
			if start_addr > end_addr || start_addr < 0 {
				fmt.Println("Invalid memory region.")
			}

			if end_addr >= MEM_SIZE {
				fmt.Println("Invalid memory region.")
			} else {
				tm.dumpMemory(start_addr, end_addr)
			}
		case "i":
			start_addr := tm.readNumber("Starting Address", 0)
			end_addr := tm.readNumber("Ending Address", MEM_SIZE-1)
			if start_addr > end_addr || start_addr < 0 {
				fmt.Println("Invalid memory region.")
			}

			if end_addr >= MEM_SIZE {
				fmt.Println("Invalid memory region.")
			} else {
				tm.dumpProgram(start_addr, end_addr)
			}
		case "q":
			fmt.Println("Exiting.")
			break interactive
		case "r":
			tm.dumpRegisters()
		case "s":
			tm.stepProgram()
		case "t":
			tm.trace = !tm.trace
			fmt.Println("Execution tracing is now", tm.trace)
		default:
			fmt.Println("Not implemented yet. Try 'h' for help.")
		}
	}

}

func main() {
	var tm TinyMachine

	if len(os.Args) < 2 {
		fmt.Println("You must supply a program as the first argument.")
	} else {
		if tm.loadProgram(os.Args[1]) {
			tm.Interact()
		} else {
			fmt.Println("Error loading program from:", os.Args[1])
		}
	}
}