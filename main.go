package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type (
	execFunc  func()
	writeFunc func(msg string, size uint)
)

var (
	//go:embed build/hello_flat.bin
	hello []byte
	//go:embed build/write_flat.bin
	write []byte
)

func main() {
	xargs := os.Args[1:]
	msg := strings.Join(xargs, " ") + "\n"

	// If a message was supplied use the write program instead
	code := hello
	if len(xargs) > 0 {
		code = write
	}

	// map page of memory for us to load our binary code into
	// Protections: protections can only be RW on MacOS because the hardened runtime prevents
	// mapping memory that is both writable and executable
	// Flags: We need to use private and anon because we are not mapping a file.
	//		Private would mean that it would not alter the underlying file but since
	//			we are not using a file it should definitely not be shared
	//		Anon: indicates we are operating on a virtual file and not a real file.
	//		JIT: This is hinting to the OS that we are using this for JIT purposes.
	fmt.Printf("mapping %v bytes\n", len(code))
	fn, err := syscall.Mmap(-1, 0, len(code), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANON|syscall.MAP_JIT)
	if err != nil {
		log.Fatalf("mmap err: %v", err)
	}

	// put our code into the mmaped memory
	n := copy(fn, code)
	fmt.Printf("copied %v bytes into memory\n", n)

	// Change the protections of our memory from RW to RX now that we no longer need
	// to modify it. This allows us to execute the code in this memory.
	fmt.Println("making memory executable")
	if err = syscall.Mprotect(fn, syscall.PROT_READ|syscall.PROT_EXEC); err != nil {
		log.Fatalf("mprotect err: %v", err)
	}

	// Convert our memory from a byte array to an executable function with some
	// very unsafe memory handling.
	fmt.Println("calling jit function")
	fmt.Println("====================")
	unsafeFunc := (uintptr)(unsafe.Pointer(&fn))
	if len(xargs) == 0 {
		f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
		f()
	} else {
		f := *(*writeFunc)(unsafe.Pointer(&unsafeFunc))
		f(msg, uint(len(msg)))
	}
	fmt.Println("====================")

	// Unmap the memory that we mapped to release it from the application.
	fmt.Println("unmapping memory")
	if err := syscall.Munmap(fn); err != nil {
		log.Fatalf("munmap err: %v", err)
	}
	fmt.Println("done.")
}
