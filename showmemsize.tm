* The first (0th) memory address is initialized with the size of the
* data memory. Read that value into register 0.
LD 0,0(0)
* Because memory is zero indexed and we want human friendly output,
* store 1 in register 1.
LDC 1,1(0)
* Add the two stored numbers to get a human friendly value.
ADD 0,0,1
* Spit out the answer.
OUT 0,0,0
