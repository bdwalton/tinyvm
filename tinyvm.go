package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// If --mem_size isn't passed, the default size of data and instruction memory.
const (
	DEF_MEM_SIZE = 1024
	NUM_REGS     = 8 // The total number of registers available.
	PC_REG       = 7 // The registered used as the program counter.
)

var (
	mem_size = flag.Uint64("mem_size", DEF_MEM_SIZE, "This size of program and data memory.")
)

type menuAction struct {
	desc   string
	action func(tm *TinyMachine)
}

type TinyInstructionType int

const (
	iopRO TinyInstructionType = iota // Register-only
	iopRM TinyInstructionType = iota // Register-memory
	iopRA TinyInstructionType = iota // Register-address
)

// Instructions are composed of one operation and up to three
// arguments.
type TinyInstruction struct {
	iop     string
	iargs   []int32
	ioptype TinyInstructionType
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
	stdin              *bufio.Reader     // To handle data input
	registers          [NUM_REGS]int32   // 8 registers
	mem_size           int32             // How many memory slots
	data_memory        []int32           // Data memory
	instruction_memory []TinyInstruction // Instruction memory
	trace              bool              // Output instructions as they're executed
	cpustate           TinyCPUState      // See cpu* constants above
}

func (ti TinyInstruction) String() string {
	var s string

	switch ti.ioptype {
	case iopRO:
		s = fmt.Sprintf("%-4s %d,%d,%d", ti.iop, ti.iargs[0], ti.iargs[1], ti.iargs[2])
	default:
		s = fmt.Sprintf("%-4s %d,%d(%d)", ti.iop, ti.iargs[0], ti.iargs[1], ti.iargs[2])
	}

	return s
}

// Operands are of the form r,s,t where r, s and t are all integers
func parseROop(args string) ([]int32, error) {
	string_args := strings.Split(args, ",")
	converted_args := make([]int32, 3)

	if len(string_args) != 3 {
		return nil, errors.New("Invalid arguments: " + args)
	} else {
		for i := 0; i < 3; i++ {
			num, err := strconv.ParseInt(string_args[i], 10, 32)
			if err != nil {
				return nil, errors.New("Invalid arguments: " + args)
			} else {
				// Ensure that all operands are valid registers
				if num < 0 || num >= NUM_REGS {
					return nil, errors.New("Invalid arguments. Bad register: " + string_args[i])
				} else {
					converted_args[i] = int32(num)
				}
			}
		}
	}

	return converted_args, nil
}

// Operands are of the form r,s(t) where r, s and t are all integers
func parseRMop(args string) ([]int32, error) {
	converted_args := make([]int32, 3)

	x := strings.Index(args, ",")
	y := strings.Index(args, "(")
	z := strings.Index(args, ")")

	if x < 1 || y < x || z < y {
		return nil, errors.New("Invalid arguments: " + args)
	} else {
		indexes := [][]int{[]int{0, x}, []int{x + 1, y}, []int{y + 1, z}}

		for i, bounds := range indexes {
			str_num := args[bounds[0]:bounds[1]]
			num, err := strconv.ParseInt(str_num, 10, 32)

			if err != nil {
				return nil, errors.New("Invalid arguments: " + args)
			} else {
				// Ensure that the 1st and 3rd operands are valid registers
				if (i == 0 || i == 2) && (num < 0 || num >= NUM_REGS) {
					return nil, errors.New("Invalid arguments. Bad register: " + str_num)
				} else {
					converted_args[i] = int32(num)
				}
			}
		}
	}

	return converted_args, nil
}

func parseInstruction(line string) (TinyInstruction, error) {
	var args []int32
	var err error
	var ti TinyInstruction
	var ioptype TinyInstructionType

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
			ioptype = iopRO
		case "LD", "ST":
			ioptype = iopRM
			args, err = parseRMop(line_parts[1])
		case "LDA", "LDC", "JLT", "JLE", "JGT", "JGE", "JEQ", "JNE":
			args, err = parseRMop(line_parts[1])
			ioptype = iopRA
		default:
			return ti, errors.New("Invalid opcode: '" + line_parts[0] + "'")
		}

		if err != nil {
			m := "Invalid arguments for opcode " + line_parts[0] + ": '" + line_parts[1] + "'"
			return ti, errors.New(m)
		} else {
			ti.iop = line_parts[0]
			ti.iargs = args
			ti.ioptype = ioptype
		}
	}

	return ti, nil
}

func (tm *TinyMachine) speak(saywhat ...interface{}) {
	fmt.Println(saywhat...)
}

func (tm *TinyMachine) initializeMachine(clearprogram bool) {
	tm.mem_size = int32(*mem_size)
	tm.data_memory = make([]int32, tm.mem_size)

	for i := 0; i < NUM_REGS; i++ {
		tm.registers[i] = 0
	}

	for i := 0; i < int(tm.mem_size); i++ {
		tm.data_memory[i] = 0
	}

	if clearprogram {
		tm.instruction_memory = make([]TinyInstruction, tm.mem_size)
		for i := 0; i < int(tm.mem_size); i++ {
			tm.instruction_memory[i] = TinyInstruction{"HALT", []int32{0, 0, 0}, iopRO}
		}
	}

	// Store the size of the memory in the first memory element.
	tm.data_memory[0] = tm.mem_size - 1
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
	if pc < 0 || pc > tm.mem_size-1 {
		tm.cpustate = cpuIMEM_ERR
	} else {
		// Step the program counter
		tm.registers[PC_REG] = pc + 1

		instruction := tm.instruction_memory[pc]
		if tm.trace {
			tm.speak("Executing:", instruction)
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
			tm.speak(tm.registers[r])
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
			tm.registers[r] = a
		case "LDC":
			tm.registers[r] = s
		case "LD":
			if a < 0 || a >= tm.mem_size {
				tm.cpustate = cpuDMEM_ERR
			} else {
				tm.registers[r] = tm.data_memory[a]
			}
		case "ST":
			if a < 0 || a >= tm.mem_size {
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
		tm.speak("Divide by zero error. Program halted.")
	case cpuIMEM_ERR:
		tm.speak("Instruction memory access violation. Program halted.")
	case cpuDMEM_ERR:
		tm.speak("Data memory access violation. Program halted.")
	case cpuHALTED:
		tm.speak("Program halted.")
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

func (tm *TinyMachine) loadProgram(progname string, fh io.Reader) bool {
	var (
		i       int
		linenum int = 0
	)

	tm.initializeMachine(true)

	reader := bufio.NewReader(fh)
	tm.speak("Reading program from", progname)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				tm.speak("Error reading program:", err)
				return false
			}
		} else {
			linenum++
			chomped_line := line[:len(line)-1] // Strip the
			r := regexp.MustCompile("[[:alnum:]]")
			has_content := r.MatchString(chomped_line)

			if !has_content || strings.Index(chomped_line, "*") == 0 {
				continue // Comments blank and comment lines
			} else {
				instruction, err := parseInstruction(chomped_line)

				if err != nil {
					tm.speak(err)
					tm.speak(fmt.Sprintf("Error parsing program at line %d: %s", linenum, chomped_line))
					return false
				} else {
					tm.instruction_memory[i], i = instruction, i+1
				}
			}
		}
	}

	return true
}

func (tm *TinyMachine) dumpRegisters() {
	tm.speak("Current Tiny Machine register values:")

	regs_even := ""
	regs_odd := ""

	for i := 0; i < NUM_REGS; i += 2 {
		regs_even += fmt.Sprintf("%2d: %011d  ", i, tm.registers[i])
		regs_odd += fmt.Sprintf("%2d: %011d  ", i+1, tm.registers[i+1])
	}

	tm.speak(regs_even + "\n" + regs_odd)
}

func (tm *TinyMachine) dumpMemory(start_addr, end_addr int32) {
	tm.speak("Dumping data memory from address %d to %d", start_addr, end_addr)

	for i := start_addr; i <= end_addr; i++ {
		tm.speak(fmt.Sprintf("%04d: %d", i, tm.data_memory[i]))
	}
}

func (tm *TinyMachine) dumpProgram(start_addr, end_addr int32) {
	fmt.Printf("Dumping instruction memory from address %d to %d.\n", start_addr, end_addr)

	for i := start_addr; i <= end_addr; i++ {
		fmt.Printf("%04d: %v\n", i, tm.instruction_memory[i])
	}
}

func (tm *TinyMachine) readNumber(prompt string, def int32) int32 {
	for {
		fmt.Printf("%s: ", prompt)
		input, err := tm.stdin.ReadString('\n')
		if err != nil {
			tm.speak("Error reading input. Returning default", def)
			break
		} else {
			num, err := strconv.ParseInt(input[:len(input)-1], 10, 32)
			if err != nil {
				tm.speak("Error converting input. Returning default", def)
				break
			} else {
				return int32(num)
			}
		}
	}

	return def
}

func handleClear(tm *TinyMachine) {
	tm.resetState()
}

func handleDataMemoryDump(tm *TinyMachine) {
	start_addr := tm.readNumber("Starting Address", 0)
	end_addr := tm.readNumber("Ending Address", tm.mem_size-1)
	if start_addr > end_addr || start_addr < 0 {
		tm.speak("Invalid memory region")
	}

	if end_addr >= tm.mem_size {
		tm.speak("Invalid memory region.")
	} else {
		tm.dumpMemory(start_addr, end_addr)
	}
}

func handleInstructionMemoryDump(tm *TinyMachine) {
	start_addr := tm.readNumber("Starting Address", 0)
	end_addr := tm.readNumber("Ending Address", tm.mem_size-1)
	if start_addr > end_addr || start_addr < 0 {
		tm.speak("Invalid memory region.")
	}

	if end_addr >= tm.mem_size {
		tm.speak("Invalid memory region.")
	} else {
		tm.dumpProgram(start_addr, end_addr)
	}
}

func handleGo(tm *TinyMachine) {
	tm.runProgram()
}

func handleQuit(tm *TinyMachine) {
	tm.speak("Exiting.")
	os.Exit(0)
}

func handleRegDump(tm *TinyMachine) {
	tm.dumpRegisters()
}

func handleStep(tm *TinyMachine) {
	tm.stepProgram()
}

func handleTrace(tm *TinyMachine) {
	tm.trace = !tm.trace
	tm.speak("Execution tracing is now", tm.trace)
}

func (tm *TinyMachine) Interact() {
	menu := map[string]menuAction{
		"?": menuAction{"display this help text", nil},
		"c": menuAction{"clear machine state", handleClear},
		"d": menuAction{"display data memory", handleDataMemoryDump},
		"g": menuAction{"run program to halt state", handleGo},
		"h": menuAction{"display this help text", nil},
		"i": menuAction{"display instruction memory", handleInstructionMemoryDump},
		"q": menuAction{"quit the tiny machine simulator", handleQuit},
		"r": menuAction{"dump register contents", handleRegDump},
		"s": menuAction{"step program forward by one instruction", handleStep},
		"t": menuAction{"toggle execution tracing", handleTrace},
	}

	tm.speak("Tiny Machine simulation (enter h for help)")

	for {
		fmt.Printf("Enter command: ")
		input, err := tm.stdin.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Fake up a real "q" entry so we handle eof the same way as a normal
				// exit.
				tm.speak()
				input = "q\n"
			} else {
				// This will be handled with the unknown case below.
				input = "ijustmashedthekeyboard\n"
			}
		}

		command := input[:len(input)-1]
		if menuitem, ok := menu[command]; ok {
			switch menuitem.action {
			case nil:
				// Show the help text if the menu key has no action
				for k, m := range menu {
					fmt.Printf("%s: %s\n", k, m.desc)
				}
			default:
				menuitem.action(tm)
			}
		} else {
			tm.speak("Not implemented yet. Try 'h' for help.")
		}
	}
}

func main() {
	var tm TinyMachine

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("You must supply a program as the first argument.")
	}

	programfile, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatalf("Error reading from %s: %s\n", flag.Args()[0], err)
	}
	defer programfile.Close()

	if tm.loadProgram(flag.Args()[0], programfile) {
		tm.Interact()
	} else {
		log.Fatalf("Error loading program from:", flag.Args()[0])
	}
}
