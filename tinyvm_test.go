package main

import (
	"math"
	"reflect"
	"testing"
)

func TestParseRMop(t *testing.T) {
	cases := []struct {
		in       string
		want     []int
		want_err string
	}{
		{"0,0(1)", []int{0, 0, 1}, ""},
		{"10,10(10)", []int{10, 10, 10}, ""},
		{"10,a(1)", nil, "Invalid arguments: 10,a(1)"},
		{"a,10(1)", nil, "Invalid arguments: a,10(1)"},
		{",10(1)", nil, "Invalid arguments: ,10(1)"},
		{"1,(1)", nil, "Invalid arguments: 1,(1)"},
		{"1,", nil, "Invalid arguments: 1,"},
		{"1", nil, "Invalid arguments: 1"},
		{"", nil, "Invalid arguments: "},
	}
	for _, c := range cases {
		got, got_err := parseRMop(c.in)
		if got_err != nil {
			if c.want_err == "" {
				t.Errorf("Unexpected error raised for parseRMop(%q): %q", c.in, got_err.Error())
			} else if c.want_err != got_err.Error() {
				t.Errorf("Expected error '%q' but got '%q'.", c.want_err, got_err.Error())
			}
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseRMop(%q) == %q, want %q.", c.in, got, c.want)
		}
	}
}

func TestParseROop(t *testing.T) {
	cases := []struct {
		in       string
		want     []int
		want_err string
	}{
		{"0,0,1", []int{0, 0, 1}, ""},
		{"10,10,10", []int{10, 10, 10}, ""},
		{"10,a,1", nil, "Invalid arguments: 10,a,1"},
		{"a,10,1", nil, "Invalid arguments: a,10,1"},
		{",10,1", nil, "Invalid arguments: ,10,1"},
		{"1,,1", nil, "Invalid arguments: 1,,1"},
		{"1,", nil, "Invalid arguments: 1,"},
		{"1", nil, "Invalid arguments: 1"},
		{"", nil, "Invalid arguments: "},
	}
	for _, c := range cases {
		got, got_err := parseROop(c.in)
		if got_err != nil {
			if c.want_err == "" {
				t.Errorf("Unexpected error raised for parseROop(%q): %q", c.in, got_err.Error())
			} else if c.want_err != got_err.Error() {
				t.Errorf("Expected error '%q' but got '%q'", c.want_err, got_err.Error())
			}
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("parseROop(%q) == %q, want %q", c.in, got, c.want)
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
		{"HALT   0,0,1", TinyInstruction{"HALT", []int{0, 0, 1}, iopRO}, ""},
		{"IN     0,0,1", TinyInstruction{"IN", []int{0, 0, 1}, iopRO}, ""},
		{"OUT    0,0,0", TinyInstruction{"OUT", []int{0, 0, 0}, iopRO}, ""},
		{"ADD    0,0,0", TinyInstruction{"ADD", []int{0, 0, 0}, iopRO}, ""},
		{"SUB    0,0,0", TinyInstruction{"SUB", []int{0, 0, 0}, iopRO}, ""},
		{"MUL    0,0,0", TinyInstruction{"MUL", []int{0, 0, 0}, iopRO}, ""},
		{"DIV    0,0,0", TinyInstruction{"DIV", []int{0, 0, 0}, iopRO}, ""},
		// Valid RM instructions
		{"LD     0,0(0)", TinyInstruction{"LD", []int{0, 0, 0}, iopRM}, ""},
		{"ST     0,0(0)", TinyInstruction{"ST", []int{0, 0, 0}, iopRM}, ""},
		// Valid RA instructions
		{"LDA    0,0(0)", TinyInstruction{"LDA", []int{0, 0, 0}, iopRA}, ""},
		{"LDC    0,0(0)", TinyInstruction{"LDC", []int{0, 0, 0}, iopRA}, ""},
		{"JLT    0,0(0)", TinyInstruction{"JLT", []int{0, 0, 0}, iopRA}, ""},
		{"JLE    0,0(0)", TinyInstruction{"JLE", []int{0, 0, 0}, iopRA}, ""},
		{"JGT    0,0(0)", TinyInstruction{"JGT", []int{0, 0, 0}, iopRA}, ""},
		{"JGE    0,0(0)", TinyInstruction{"JGE", []int{0, 0, 0}, iopRA}, ""},
		{"JEQ    0,0(0)", TinyInstruction{"JEQ", []int{0, 0, 0}, iopRA}, ""},
		{"JNE    0,0(0)", TinyInstruction{"JNE", []int{0, 0, 0}, iopRA}, ""},
		// Garbage spaces are handled properly
		{"   HALT  0,0,1   ", TinyInstruction{"HALT", []int{0, 0, 1}, iopRO}, ""},
		{"   LD  0,0(1)   ", TinyInstruction{"LD", []int{0, 0, 1}, iopRM}, ""},
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

	tm.cpustate = cpuHALTED
	tm.instruction_memory[0] = TinyInstruction{"LDC", []int{1, 1, 1}, iopRA}
	tm.instruction_memory[MEM_SIZE-1] = TinyInstruction{"ADD", []int{1, 1, 1}, iopRO}
	tm.data_memory[0] = 1
	tm.data_memory[MEM_SIZE-1] = 100
	tm.registers[PC_REG] = 1

	tm.resetState()

	if tm.cpustate != cpuOK {
		t.Errorf("Resetting machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"LDC", []int{1, 1, 1}, iopRA},
		tm.instruction_memory[0]) {
		t.Errorf("Resetting machine cleared instructions.")
	} else if !reflect.DeepEqual(TinyInstruction{"ADD", []int{1, 1, 1}, iopRO},
		tm.instruction_memory[MEM_SIZE-1]) {
		t.Errorf("Resetting machine cleared instructions.")
	} else if tm.data_memory[0] != MEM_SIZE-1 {
		t.Errorf("Resetting machine didn't reset memory state.")
	} else if tm.registers[PC_REG] != 0 {
		t.Errorf("Initializing machine didn't reset the program counter.")
	}
}

func TestInitializeMachine(t *testing.T) {
	var tm TinyMachine

	tm.cpustate = cpuDIV_ZERO
	tm.instruction_memory[0] = TinyInstruction{"LDC", []int{1, 1, 1}, iopRA}
	tm.instruction_memory[MEM_SIZE-1] = TinyInstruction{"ADD", []int{1, 1, 1}, iopRO}
	tm.data_memory[0] = 1
	tm.data_memory[MEM_SIZE-1] = 100
	tm.registers[PC_REG] = 1

	tm.initializeMachine(true)

	if tm.cpustate != cpuOK {
		t.Errorf("Initializing machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int{0, 0, 0}, iopRO}, tm.instruction_memory[0]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int{0, 0, 0}, iopRO}, tm.instruction_memory[MEM_SIZE-1]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if tm.data_memory[0] != MEM_SIZE-1 {
		t.Errorf("Initializing machine didn't reset memory state.")
	} else if tm.registers[PC_REG] != 0 {
		t.Errorf("Initializing machine didn't reset the program counter.")
	}
}

func TestHALTInstruction(t *testing.T) {
	var tm TinyMachine

	cases := []struct {
		expected_pc  int
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
	tm.registers = [NUM_REGS]int{0, -1, 10, 2, 2, math.MinInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"SUB", []int{0, 2, 3}, iopRO}
	tm.instruction_memory[1] = TinyInstruction{"SUB", []int{0, 3, 6}, iopRO}
	// Not necessary, but include for completeness. Machine is initialized with
	// HALT instructions.
	tm.instruction_memory[2] = TinyInstruction{"HALT", []int{0, 0, 0}, iopRO}

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
	tm.registers = [NUM_REGS]int{0, 1, 10, 2, 2, 10, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"DIV", []int{2, 2, 3}, iopRO} // 10 / 2 -> reg2
	tm.instruction_memory[1] = TinyInstruction{"DIV", []int{4, 4, 5}, iopRO} // 2 / 10 -> reg4
	tm.instruction_memory[2] = TinyInstruction{"DIV", []int{0, 1, 0}, iopRO} // 1 / 0  -> reg0

	cases := []struct {
		expected_reg int
		expected_val int
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
	tm.registers = [NUM_REGS]int{0, -1, 10, 2, 4, -5, -7, 0}

	tm.instruction_memory[0] = TinyInstruction{"MUL", []int{2, 2, 3}, iopRO} // 10 * 2  -> reg2
	tm.instruction_memory[1] = TinyInstruction{"MUL", []int{4, 4, 5}, iopRO} // 4 * -5  -> reg4
	tm.instruction_memory[2] = TinyInstruction{"MUL", []int{0, 1, 0}, iopRO} // 0 * -1  -> reg0
	tm.instruction_memory[3] = TinyInstruction{"MUL", []int{0, 5, 6}, iopRO} // -5 * -7 -> reg0

	cases := []struct {
		expected_reg int
		expected_val int
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
	tm.registers = [NUM_REGS]int{0, 1, 10, 2, 2, math.MaxInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"ADD", []int{0, 2, 3}, iopRO} // 10 + 2  -> reg0
	tm.instruction_memory[1] = TinyInstruction{"ADD", []int{0, 3, 6}, iopRO} // 2 + 5   -> reg0
	tm.instruction_memory[2] = TinyInstruction{"ADD", []int{0, 1, 0}, iopRO} // 1 + 7   -> reg0
	tm.instruction_memory[3] = TinyInstruction{"ADD", []int{0, 1, 5}, iopRO} // 1 + MAX -> reg0

	cases := []struct {
		expected_reg int
		expected_val int
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
	tm.registers = [NUM_REGS]int{0, -1, 10, 2, 2, math.MinInt32, 5, 0}

	tm.instruction_memory[0] = TinyInstruction{"SUB", []int{0, 2, 3}, iopRO} // 10 - 2  -> reg0
	tm.instruction_memory[1] = TinyInstruction{"SUB", []int{0, 3, 6}, iopRO} // 2 - 5   -> reg0
	tm.instruction_memory[2] = TinyInstruction{"SUB", []int{0, 1, 0}, iopRO} // -1 - -3  -> reg0
	tm.instruction_memory[3] = TinyInstruction{"SUB", []int{0, 1, 5}, iopRO} // -1 - MIN -> reg0

	cases := []struct {
		expected_reg int
		expected_val int
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
