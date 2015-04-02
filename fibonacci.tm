* Store 1 and 1, the first two numbers of the sequence
LDC 0,1(0)
LDC 1,1(0)
* This will be used to decrement the count
LDC 5,1(0)
* Prompt for a number - output this many elements of the fibonacci sequence
* Store value in register 6
IN 6,0,0
* Test first, in case the user entered a value <= 0
JLE 6,7(7)
* Output a number in the sequence
OUT 0,0,0
* Decrement the counter, since we've printed a number
SUB 6,6,5
* Determine the second next number in the sequence
ADD 2,0,1
* Abuse LDA to shift reg1 -> reg0 and reg2 -> reg1
LDA 0,0(1)
LDA 1,0(2)
* Step back 6 instructions to iterate, if we're not done yet.
JNE 6,-6(7)
* Not necessary, but provides a visible target in the program for the
* first JLE.
HALT 0,0,0
