package main

import (
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
		{"HALT   0,0,1", TinyInstruction{"HALT", []int{0, 0, 1}}, ""},
		{"IN     0,0,1", TinyInstruction{"IN", []int{0, 0, 1}}, ""},
		{"OUT    0,0,0", TinyInstruction{"OUT", []int{0, 0, 0}}, ""},
		{"ADD    0,0,0", TinyInstruction{"ADD", []int{0, 0, 0}}, ""},
		{"SUB    0,0,0", TinyInstruction{"SUB", []int{0, 0, 0}}, ""},
		{"MUL    0,0,0", TinyInstruction{"MUL", []int{0, 0, 0}}, ""},
		{"DIV    0,0,0", TinyInstruction{"DIV", []int{0, 0, 0}}, ""},
		// Valid RM instructions
		{"LD     0,0(0)", TinyInstruction{"LD", []int{0, 0, 0}}, ""},
		{"ST     0,0(0)", TinyInstruction{"ST", []int{0, 0, 0}}, ""},
		{"LDA    0,0(0)", TinyInstruction{"LDA", []int{0, 0, 0}}, ""},
		{"LDC    0,0(0)", TinyInstruction{"LDC", []int{0, 0, 0}}, ""},
		{"JLT    0,0(0)", TinyInstruction{"JLT", []int{0, 0, 0}}, ""},
		{"JLE    0,0(0)", TinyInstruction{"JLE", []int{0, 0, 0}}, ""},
		{"JGT    0,0(0)", TinyInstruction{"JGT", []int{0, 0, 0}}, ""},
		{"JGE    0,0(0)", TinyInstruction{"JGE", []int{0, 0, 0}}, ""},
		{"JEQ    0,0(0)", TinyInstruction{"JEQ", []int{0, 0, 0}}, ""},
		{"JNE    0,0(0)", TinyInstruction{"JNE", []int{0, 0, 0}}, ""},
		// Garbage spaces are handled properly
		{"   HALT  0,0,1   ", TinyInstruction{"HALT", []int{0, 0, 1}}, ""},
		{"   LD  0,0(1)   ", TinyInstruction{"LD", []int{0, 0, 1}}, ""},
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
	tm.instruction_memory[0] = TinyInstruction{"LDC", []int{1, 1, 1}}
	tm.instruction_memory[MEM_SIZE-1] = TinyInstruction{"ADD", []int{1, 1, 1}}
	tm.data_memory[0] = 1
	tm.data_memory[MEM_SIZE-1] = 100
	tm.registers[PC_REG] = 1

	tm.resetState()

	if tm.cpustate != cpuOK {
		t.Errorf("Resetting machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"LDC", []int{1, 1, 1}},
		tm.instruction_memory[0]) {
		t.Errorf("Resetting machine cleared instructions.")
	} else if !reflect.DeepEqual(TinyInstruction{"ADD", []int{1, 1, 1}},
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
	tm.instruction_memory[0] = TinyInstruction{"LDC", []int{1, 1, 1}}
	tm.instruction_memory[MEM_SIZE-1] = TinyInstruction{"ADD", []int{1, 1, 1}}
	tm.data_memory[0] = 1
	tm.data_memory[MEM_SIZE-1] = 100
	tm.registers[PC_REG] = 1

	tm.initializeMachine(true)

	if tm.cpustate != cpuOK {
		t.Errorf("Initializing machine didn't clear halt state.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int{0, 0, 0}}, tm.instruction_memory[0]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if !reflect.DeepEqual(TinyInstruction{"HALT", []int{0, 0, 0}}, tm.instruction_memory[MEM_SIZE-1]) {
		t.Errorf("Initializing machine didn't clear instruction memory.")
	} else if tm.data_memory[0] != MEM_SIZE-1 {
		t.Errorf("Initializing machine didn't reset memory state.")
	} else if tm.registers[PC_REG] != 0 {
		t.Errorf("Initializing machine didn't reset the program counter.")
	}
}

func TestDIVInstruction(t *testing.T) {
	var tm TinyMachine

	tm.initializeMachine(true)
	// Stuff some values into the registers
	tm.registers = [NUM_REGS]int{0, 1, 10, 2, 2, 10, 0, 0}

	tm.instruction_memory[0] = TinyInstruction{"DIV", []int{2, 2, 3}} // 10 / 2 -> reg2
	tm.instruction_memory[1] = TinyInstruction{"DIV", []int{4, 4, 5}} // 2 / 10 -> reg4
	tm.instruction_memory[2] = TinyInstruction{"DIV", []int{0, 1, 0}} // 1 / 0  -> reg0

	tm.stepProgram()
	if tm.registers[2] != 5 {
		t.Errorf("DIV 10/2 didn't work.")
		if tm.cpustate != cpuOK {
			t.Errorf("DIV 10/2 worked, but cpustate is invalid.")
		}
	}

	tm.stepProgram()
	if tm.registers[4] != 0 {
		t.Errorf("DIV 2/10 didn't work.")
		if tm.cpustate != cpuOK {
			t.Errorf("DIV 2/10 worked, but cpustate is invalid.")
		}
	}

	tm.stepProgram()
	if tm.cpustate != cpuDIV_ZERO {
		t.Errorf("DIV by 0 didn't set correct cpu state.")
	}
}
