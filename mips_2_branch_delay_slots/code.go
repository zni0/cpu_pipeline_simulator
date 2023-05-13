package main

func loadCode() {
	Memory = []string{
		"ADDI 10 10 1",
		"ADDI 11 11 1",

		"ADDI 1 0 1", // Value of n
		"ADDI 2 2 1",
		"ADD 12 10 11",
		"ADD 11 0 10",
		"ADD 10 0 12",
		"BNE 1 2 -5",
		"NOOP",
		"NOOP",
		"HALT",
		"",
	}
}

/*
Code for sum of first n natural numbers:
		"ADDI 1 0 10", // Value of n
		"ADDI 2 2 1",
		"ADD 3 3 2",
		"BNE 1 2 -3",
		"NOOP",
		"NOOP",
		"HALT",
		"",
*/

/*
Code for n+1 th Fibonacci number:

		t1=0, t2=1
		for i=0;i<10;i++
		nx = t1 + t2
		t1 = t2
		t2 = nx
		------
		"ADDI 10 10 1",
		"ADDI 11 11 1",
		"ADDI 1 0 10", // Value of n
		"ADDI 2 2 1",
		"ADD 12 10 11",
		"ADD 11 0 10",
		"ADD 10 0 12",
		"BNE 1 2 -5",
		"NOOP",
		"NOOP",
		"HALT",
		"",
*/
