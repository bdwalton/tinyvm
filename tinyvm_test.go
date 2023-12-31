package main

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestParseRMop(t *testing.T) {
	cases := []struct {
		in       string
		want     []int32
		want_err string
	}{
		{"0,0(1)", []int32{0, 0, 1}, ""},
		{"2,12(2)", []int32{2, 12, 2}, ""},
		{"1,a(1)", nil, "Invalid arguments: 1,a(1)"},
		{"a,10(1)", nil, "Invalid arguments: a,10(1)"},
		{",10(1)", nil, "Invalid arguments: ,10(1)"},
		{"1,(1)", nil, "Invalid arguments: 1,(1)"},
		{"1,", nil, "Invalid arguments: 1,"},
		{"1", nil, "Invalid arguments: 1"},
		{"", nil, "Invalid arguments: "},
		{"10,1(1)", nil, "Invalid arguments. Bad register: 10"},
		{"1,1(12)", nil, "Invalid arguments. Bad register: 12"},
	}
	for i, c := range cases {
		got, got_err := parseRMop(c.in)
		if c.want == nil {
			if got_err == nil {
				t.Errorf("%d: Expected invalid result when calling parseRMop(%q).",
					i, c.in)
			} else if c.want_err != got_err.Error() {
				t.Errorf("%d: Expected error '%q' but got '%q'.",
					i, c.want_err, got_err.Error())
			}
		} else {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%d: parseRMop(%q) == %v, want %v.", i, c.in, got, c.want)
			}
		}
	}
}

func TestParseROop(t *testing.T) {
	cases := []struct {
		in       string
		want     []int32
		want_err string
	}{
		{"0,0,1", []int32{0, 0, 1}, ""},
		{"2,2,2", []int32{2, 2, 2}, ""},
		{"2,a,1", nil, "Invalid arguments: 2,a,1"},
		{"a,10,1", nil, "Invalid arguments: a,10,1"},
		{",10,1", nil, "Invalid arguments: ,10,1"},
		{"1,,1", nil, "Invalid arguments: 1,,1"},
		{"1,", nil, "Invalid arguments: 1,"},
		{"1", nil, "Invalid arguments: 1"},
		{"", nil, "Invalid arguments: "},
		{"12,1,1", nil, "Invalid arguments. Bad register: 12"},
		{"2,13,1", nil, "Invalid arguments. Bad register: 13"},
		{"2,1,14", nil, "Invalid arguments. Bad register: 14"},
	}
	for i, c := range cases {
		got, got_err := parseROop(c.in)

		if c.want == nil {
			if got_err == nil {
				t.Errorf("%d: Expected invalid result when calling parseROop(%q).",
					i, c.in)
			} else if c.want_err != got_err.Error() {
				t.Errorf("%d: Expected error '%q' but got '%q'.",
					i, c.want_err, got_err.Error())
			}
		} else {
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%d: parseROop(%q) == %v, want %v.", i, c.in, got, c.want)
			}
		}
	}
}

func TestParseInstruction(t *testing.T) {
	cases := []struct {
		in       string
		want     TinyInstruction
		want_err string
	}{
		// Valid RO instructions
		{"HALT   0,0,1", TinyInstruction{"HALT", []int32{0, 0, 1}, iopRO}, ""},
		{"IN     0,0,1", TinyInstruction{"IN", []int32{0, 0, 1}, iopRO}, ""},
		{"OUT    0,0,0", TinyInstruction{"OUT", []int32{0, 0, 0}, iopRO}, ""},
		{"ADD    0,0,0", TinyInstruction{"ADD", []int32{0, 0, 0}, iopRO}, ""},
		{"SUB    0,0,0", TinyInstruction{"SUB", []int32{0, 0, 0}, iopRO}, ""},
		{"MUL    0,0,0", TinyInstruction{"MUL", []int32{0, 0, 0}, iopRO}, ""},
		{"DIV    0,0,0", TinyInstruction{"DIV", []int32{0, 0, 0}, iopRO}, ""},
		// Valid RM instructions
		{"LD     0,0(0)", TinyInstruction{"LD", []int32{0, 0, 0}, iopRM}, ""},
		{"ST     0,0(0)", TinyInstruction{"ST", []int32{0, 0, 0}, iopRM}, ""},
		// Valid RA instructions
		{"LDA    0,0(0)", TinyInstruction{"LDA", []int32{0, 0, 0}, iopRA}, ""},
		{"LDC    0,0(0)", TinyInstruction{"LDC", []int32{0, 0, 0}, iopRA}, ""},
		{"JLT    0,0(0)", TinyInstruction{"JLT", []int32{0, 0, 0}, iopRA}, ""},
		{"JLE    0,0(0)", TinyInstruction{"JLE", []int32{0, 0, 0}, iopRA}, ""},
		{"JGT    0,0(0)", TinyInstruction{"JGT", []int32{0, 0, 0}, iopRA}, ""},
		{"JGE    0,0(0)", TinyInstruction{"JGE", []int32{0, 0, 0}, iopRA}, ""},
		{"JEQ    0,0(0)", TinyInstruction{"JEQ", []int32{0, 0, 0}, iopRA}, ""},
		{"JNE    0,0(0)", TinyInstruction{"JNE", []int32{0, 0, 0}, iopRA}, ""},
		// Garbage spaces are handled properly
		{"   HALT  0,0,1   ", TinyInstruction{"HALT", []int32{0, 0, 1}, iopRO}, ""},
		{"   LD  0,0(1)   ", TinyInstruction{"LD", []int32{0, 0, 1}, iopRM}, ""},
		// RM format for RO opcode
		{"IN    0,0(1)", TinyInstruction{}, "Invalid arguments for opcode IN: '0,0(1)'"},
		// RO format for RM opcode
		{"LD    0,0,0", TinyInstruction{}, "Invalid arguments for opcode LD: '0,0,0'"},
		// Missing opcode
		{"   0,0,1   ", TinyInstruction{}, "Invalid instruction: '0,0,1'"},
		{"   0,0(1)   ", TinyInstruction{}, "Invalid instruction: '0,0(1)'"},
		// Missing operands
		{"OPCODE   ", TinyInstruction{}, "Invalid instruction: 'OPCODE'"},
		// Invalid opcode
		{"OPCODE 0,0,1   ", TinyInstruction{}, "Invalid opcode: 'OPCODE'"},
		{"OPCODE 0,0(1)  ", TinyInstruction{}, "Invalid opcode: 'OPCODE'"},
		// Garbage inputs
		{"IN 0,a,1   ", TinyInstruction{}, "Invalid arguments for opcode IN: '0,a,1'"},
		{"ST 0,a(1)   ", TinyInstruction{}, "Invalid arguments for opcode ST: '0,a(1)'"},
	}
	for _, c := range cases {
		got, got_err := parseInstruction(c.in)
		if got_err != nil {
			if c.want_err == "" {
				t.Errorf("Unexpected error raised for parseInstruction(%q): %q.", c.in, got_err.Error())
			} else if c.want_err != got_err.Error() {
				t.Errorf("Expected error '%q' but got '%q'.", c.want_err, got_err.Error())
			}
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseInstruction(%q) == %q, want %q.", c.in, got, c.want)
		}
	}
}

func TestResetState(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)

	tm.cpustate = cpuHALTED
	tm.instruction_memory[0] = TinyInstruction{"LDC", []int32{1, 1, 1}, iopRA}
	tm.instruction_memory[DEF_MEM_SIZE-1] = TinyInstruction{"ADD", []int32{1, 1, 1}, iopRO}
	tm.data_memory[0] = 1
	tm.data_memory[DEF_MEM_SIZE-1] = 100
	tm.registers[PC_REG] = 1

	tm.resetState()

	if tm.cpustate != cpuOK {
		t.Errorf("Resetting machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"LDC", []int32{1, 1, 1}, iopRA},
		tm.instruction_memory[0]) {
		t.Errorf("Resetting machine cleared instructions.")
	} else if !reflect.DeepEqual(TinyInstruction{"ADD", []int32{1, 1, 1}, iopRO},
		tm.instruction_memory[DEF_MEM_SIZE-1]) {
		t.Errorf("Resetting machine cleared instructions.")
	} else if tm.data_memory[0] != DEF_MEM_SIZE-1 {
		t.Errorf("Resetting machine didn't reset memory state.")
	} else if tm.registers[PC_REG] != 0 {
		t.Errorf("Initializing machine didn't reset the program counter.")
	}
}

func TestLoadProgram(t *testing.T) {
	var tm TinyMachine

	cases := []struct {
		prog     string
		valid    bool
		imem_pos []int
		ti       []TinyInstruction
	}{
		// Comment lines ignored
		{"LDC 1,1(0)\n* This is a comment\nADD 1,1,1\n",
			true, []int{0, 1}, []TinyInstruction{{"LDC", []int32{1, 1, 0}, iopRA},
				{"ADD", []int32{1, 1, 1}, iopRO}}},
		{"ST 1,1(0)\nSUB 1,1,1\n",
			true, []int{1}, []TinyInstruction{{"SUB", []int32{1, 1, 1}, iopRO}}},
		// Blank lines ignored.
		{"ST 1,1(0)\n\nSUB 1,1,1\n",
			true, []int{1}, []TinyInstruction{{"SUB", []int32{1, 1, 1}, iopRO}}},
		// Invalid instruction
		{"STORE 1,1(0)\nSUB 1,1,1\n",
			false, []int{}, []TinyInstruction{}},
		// Empty program
		{"",
			true, []int{0}, []TinyInstruction{{"HALT", []int32{0, 0, 0}, iopRO}}},
	}

	for i, c := range cases {
		program := bytes.NewBufferString(c.prog)
		ok := tm.loadProgram(fmt.Sprintf("test-%d", i), program)

		if ok != c.valid {
			t.Errorf("%d: Expected %t load, but didn't get it.", i, c.valid)
		} else if ok {
			// Only test the loaded instructions if the program is valid
			for x, v := range c.imem_pos {
				inst := tm.instruction_memory[v]
				expected := c.ti[x]
				if !reflect.DeepEqual(expected, inst) {
					t.Errorf("Expected instruction '%s' in location %d. Got '%s'.", expected, x, inst)
				}
			}
		}
	}
}

func TestInitializeMachine(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)

	if tm.cpustate != cpuOK {
		t.Errorf("Initializing machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int32{0, 0, 0}, iopRO}, tm.instruction_memory[0]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int32{0, 0, 0}, iopRO}, tm.instruction_memory[DEF_MEM_SIZE-1]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if tm.data_memory[0] != DEF_MEM_SIZE-1 {
		t.Errorf("Initializing machine didn't reset memory state.")
	} else if tm.registers[PC_REG] != 0 {
		t.Errorf("Initializing machine didn't reset the program counter.")
	}
}

func TestHALTInstruction(t *testing.T) {
	var tm TinyMachine

	cases := []struct {
		expected_pc  int32
		expected_cpu TinyCPUState
	}{
		{1, cpuOK},
		{2, cpuOK},
		{3, cpuHALTED},
		// Verify that running the machine when halted doesn't advance PC,
		// change state
		{3, cpuHALTED},
	}

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, -1, 10, 2, 2, math.MinInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"SUB", []int32{0, 2, 3}, iopRO}
	tm.instruction_memory[1] = TinyInstruction{"SUB", []int32{0, 3, 6}, iopRO}
	// Not necessary, but include for completeness. Machine is initialized with
	// HALT instructions.
	tm.instruction_memory[2] = TinyInstruction{"HALT", []int32{0, 0, 0}, iopRO}

	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("PC invalid. Expected %d, got %d",
				c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("PC register moved, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestDIVInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, 1, 10, 2, 2, 10, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"DIV", []int32{2, 2, 3}, iopRO} // 10 / 2 -> reg2
	tm.instruction_memory[1] = TinyInstruction{"DIV", []int32{4, 4, 5}, iopRO} // 2 / 10 -> reg4
	tm.instruction_memory[2] = TinyInstruction{"DIV", []int32{0, 1, 0}, iopRO} // 1 / 0  -> reg0

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{2, 5, cpuOK},
		{4, 0, cpuOK},
		{0, 0, cpuDIV_ZERO},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("DIV instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("DIV instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestMULInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, -1, 10, 2, 4, -5, -7, 0}

	tm.instruction_memory[0] = TinyInstruction{"MUL", []int32{2, 2, 3}, iopRO} // 10 * 2  -> reg2
	tm.instruction_memory[1] = TinyInstruction{"MUL", []int32{4, 4, 5}, iopRO} // 4 * -5  -> reg4
	tm.instruction_memory[2] = TinyInstruction{"MUL", []int32{0, 1, 0}, iopRO} // 0 * -1  -> reg0
	tm.instruction_memory[3] = TinyInstruction{"MUL", []int32{0, 5, 6}, iopRO} // -5 * -7 -> reg0

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{2, 20, cpuOK},
		{4, -20, cpuOK},
		{0, 0, cpuOK},
		{0, 35, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("MUL instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("MUL instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestADDInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, 1, 10, 2, 2, math.MaxInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"ADD", []int32{0, 2, 3}, iopRO} // 10 + 2  -> reg0
	tm.instruction_memory[1] = TinyInstruction{"ADD", []int32{0, 3, 6}, iopRO} // 2 + 5   -> reg0
	tm.instruction_memory[2] = TinyInstruction{"ADD", []int32{0, 1, 0}, iopRO} // 1 + 7   -> reg0
	tm.instruction_memory[3] = TinyInstruction{"ADD", []int32{0, 1, 5}, iopRO} // 1 + MAX -> reg0

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{0, 12, cpuOK},
		{0, 7, cpuOK},
		{0, 8, cpuOK},
		{0, -2147483648, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("ADD instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("ADD instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestSUBInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, -1, 10, 2, 2, math.MinInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"SUB", []int32{0, 2, 3}, iopRO} // 10 - 2  -> reg0
	tm.instruction_memory[1] = TinyInstruction{"SUB", []int32{0, 3, 6}, iopRO} // 2 - 5   -> reg0
	tm.instruction_memory[2] = TinyInstruction{"SUB", []int32{0, 1, 0}, iopRO} // -1 - -3  -> reg0
	tm.instruction_memory[3] = TinyInstruction{"SUB", []int32{0, 1, 5}, iopRO} // -1 - MIN -> reg0

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{0, 8, cpuOK},
		{0, -3, cpuOK},
		{0, 2, cpuOK},
		{0, 2147483647, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("SUB instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("SUB instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestLDInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, DEF_MEM_SIZE - 3, 0, 0, 0, 0, 0, 0}
	tm.data_memory[DEF_MEM_SIZE-4] = 54321
	tm.data_memory[DEF_MEM_SIZE-1] = 12345
	tm.instruction_memory[0] = TinyInstruction{"LD", []int32{0, 0, 0}, iopRM}  // Load DEF_MEM_SIZE
	tm.instruction_memory[1] = TinyInstruction{"LD", []int32{0, 2, 1}, iopRM}  // Load 12345
	tm.instruction_memory[2] = TinyInstruction{"LD", []int32{0, -1, 1}, iopRM} // Load 54321

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{0, 1023, cpuOK},
		{0, 12345, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("LD instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("LD instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestSTInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{DEF_MEM_SIZE + 1, DEF_MEM_SIZE - 3, 0, 0, 0, 0, 0, 0}
	tm.data_memory[DEF_MEM_SIZE-4] = 54321
	tm.data_memory[DEF_MEM_SIZE-1] = 12345
	tm.instruction_memory[0] = TinyInstruction{"ST", []int32{0, 1, 2}, iopRM} // ST DEF_MEM_SIZE+1 -> 1
	tm.instruction_memory[1] = TinyInstruction{"ST", []int32{1, 2, 1}, iopRM} // Load 12345

	cases := []struct {
		expected_addr int32
		expected_aval int32
		expected_cpu  TinyCPUState
	}{
		{1, DEF_MEM_SIZE + 1, cpuOK},
		{DEF_MEM_SIZE - 1, DEF_MEM_SIZE - 3, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.data_memory[c.expected_addr] != c.expected_aval {
			t.Errorf("ST instruction didn't work. Expected %d in addr[%d]. Got %d.",
				c.expected_aval, c.expected_addr, tm.data_memory[c.expected_addr])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("LD instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestLDCInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, 0, 0, 0, 0, 0, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"LDC", []int32{0, 100, 0}, iopRA} // 100 -> reg0
	tm.instruction_memory[1] = TinyInstruction{"LDC", []int32{1, -2, 1}, iopRA}  // -2 -> reg1

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{0, 100, cpuOK},
		{1, -2, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("LDC instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("LDC instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestLDAInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, 0, 0, 0, 0, 0, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"LDA", []int32{0, 100, 0}, iopRA} // 100 -> reg0
	tm.instruction_memory[1] = TinyInstruction{"LDA", []int32{3, -2, 0}, iopRA}  // 98 -> reg3
	tm.instruction_memory[2] = TinyInstruction{"LDA", []int32{4, 5, 3}, iopRA}   // 103 -> reg4

	cases := []struct {
		expected_reg int32
		expected_val int32
		expected_cpu TinyCPUState
	}{
		{0, 100, cpuOK},
		{3, 98, cpuOK},
		{4, 103, cpuOK},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[c.expected_reg] != c.expected_val {
			t.Errorf("LDA instruction didn't work. Expected %d in reg[%d]. Got %d.",
				c.expected_val, c.expected_reg, tm.registers[c.expected_reg])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("LDA instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJLTInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{-1, -2, 0, 0, 0, 0, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JLT", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JLT", []int32{4, 5, 3}, iopRA}   // !(pcreg -> 0)
	tm.instruction_memory[100] = TinyInstruction{"JLT", []int32{1, 3, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for _, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("JLT instruction didn't work. Expected PC to be %d. Got %d.",
				c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("JLT instruction fine, but cpuState invalid. Wanted %d, got %d.",
				c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJLEInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{-1, 0, 0, 1, 0, 1, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JLE", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JLE", []int32{5, 5, 3}, iopRA}   // !(pcreg -> 6)
	tm.instruction_memory[100] = TinyInstruction{"JLE", []int32{1, 3, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for i, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: JLE instruction didn't work. Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: JLE instruction fine, but cpuState invalid. Wanted %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJGEInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{1, 0, 0, 1, 0, -11, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JGE", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JGE", []int32{5, 5, 3}, iopRA}   // !(pcreg -> 6)
	tm.instruction_memory[100] = TinyInstruction{"JGE", []int32{1, 1, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for i, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: JGE instruction didn't work. Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: JGE instruction fine, but cpuState invalid. Wanted %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJGTInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{1, 100, 0, 1, 0, -11, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JGT", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JGT", []int32{5, 5, 3}, iopRA}   // !(pcreg -> 6)
	tm.instruction_memory[100] = TinyInstruction{"JGT", []int32{1, 1, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for i, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: JGT instruction didn't work. Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: JGT instruction fine, but cpuState invalid. Wanted %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJEQInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{0, 0, 0, 1, 0, -11, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JEQ", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JEQ", []int32{5, 5, 3}, iopRA}   // !(pcreg -> 6)
	tm.instruction_memory[100] = TinyInstruction{"JEQ", []int32{1, 2, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for i, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: JEQ instruction didn't work. Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: JEQ instruction fine, but cpuState invalid. Wanted %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
	}
}

func TestJNEInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int32{1, -145, 0, 1, 0, 0, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"JNE", []int32{0, 100, 2}, iopRA} // pcreg -> 100
	tm.instruction_memory[2] = TinyInstruction{"JNE", []int32{5, 5, 3}, iopRA}   // !(pcreg -> 6)
	tm.instruction_memory[100] = TinyInstruction{"JNE", []int32{1, 1, 0}, iopRA} // pcreg -> 2

	cases := []struct {
		expected_pc  int32        // Expected PC value
		expected_cpu TinyCPUState // Expected CPU state
	}{
		{100, cpuOK},
		{2, cpuOK},
		{3, cpuOK},
		{4, cpuHALTED},
	}
	for i, c := range cases {
		tm.stepProgram()
		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: JNE instruction didn't work. Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: JNE instruction fine, but cpuState invalid. Wanted %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
	}
}

func TestDMEM_ERR_State(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)

	cases := []struct {
		given_inst   TinyInstruction // The instruction to execute
		expected_pc  int32           // Expected PC value
		expected_cpu TinyCPUState    // Expected CPU state
	}{
		{TinyInstruction{"LD", []int32{0, DEF_MEM_SIZE, 1}, iopRM}, 1, cpuDMEM_ERR},
		{TinyInstruction{"LD", []int32{0, -1, 1}, iopRM}, 1, cpuDMEM_ERR},
		{TinyInstruction{"ST", []int32{0, 0, 0}, iopRM}, 1, cpuDMEM_ERR},
		{TinyInstruction{"ST", []int32{0, -1, 1}, iopRM}, 1, cpuDMEM_ERR},
	}
	for i, c := range cases {
		// Stuff some values into the registers
		tm.registers = [NUM_REGS]int32{DEF_MEM_SIZE, 0, 0, 0, 0, 0, 0, 0}
		tm.instruction_memory[0] = c.given_inst // Load the instruction that should be a memory violation

		tm.stepProgram()

		if tm.registers[PC_REG] != c.expected_pc {
			t.Errorf("%d: Expected PC to be %d. Got %d.",
				i, c.expected_pc, tm.registers[PC_REG])
		}
		if tm.cpustate != c.expected_cpu {
			t.Errorf("%d: Instruction didn't trigger DMEM_ERR. %d, got %d.",
				i, c.expected_cpu, tm.cpustate)
		}
		tm.resetState() // Reset so the next test instruction has a clean start
	}
}

func TestIMEM_ERR_State(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)

	cases := []int32{
		-1,
		DEF_MEM_SIZE,
	}
	for i, pc := range cases {
		// Stuff some values into the registers
		tm.registers = [NUM_REGS]int32{0, 0, 0, 0, 0, 0, 0, pc}

		tm.stepProgram()

		if tm.cpustate != cpuIMEM_ERR {
			t.Errorf("%d: Expected cpu state to be %d. Got %d.",
				i, cpuIMEM_ERR, tm.cpustate)
		}

		tm.resetState() // Reset so the next test instruction has a clean start
	}
}
