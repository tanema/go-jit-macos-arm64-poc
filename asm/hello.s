// small program to write a constant string to stdout
.global _start // Provide program starting address to linker
.align  2      // Make sure everything is aligned properly

_start:
	mov X0, #1      // 1 = StdOut
	adr X1, msg     // string to print
	mov X2, msg_len // length of our string
	mov X16, #4     // Unix write system call
	svc 0           // Call kernel to output the string
	ret

msg:
	.ascii "Hello World From Go!\n"
	.equ   msg_len, . - msg
