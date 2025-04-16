// small program to write a user inputted string to stdout
.global _start // Provide program starting address to linker
.align  2      // Make sure everything is aligned properly

_start:
	mov x2, x1  // save msg length received from caller in x4
	mov x1, x0  // save msg received from caller in x3
	mov x0, #1  // 1 -> stdout
	mov x16, #4 // 4 -> write syscall
	svc 0       // Call kernel to output the string
	ret
