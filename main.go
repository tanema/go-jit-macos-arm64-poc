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
	debug = os.Getenv("DEBUG") != ""
)

func printDbg(msg string, args ...any) {
	if debug {
		fmt.Printf(msg+"\n", args...)
	}
}

func check(err error, label string) {
	if err != nil {
		log.Fatal("ERR "+label+" : ", err)
	}
}

func main() {
	// If a message was supplied use the write program instead
	xargs := os.Args[1:]
	var code []byte
	if len(xargs) > 0 {
		code = write[:24] // 6 instructions at 32 bits each this is how small it can be
	} else {
		code = hello[:45] // 6 instructions + 21 characters in our string
	}

	// map page of memory for us to load our binary code into
	// Protections: protections can only be RW on MacOS because the hardened runtime prevents
	// mapping memory that is both writable and executable
	// Flags: We need to use private and anon because we are not mapping a file.
	//		Private would mean that it would not alter the underlying file but since
	//			we are not using a file it should definitely not be shared
	//		Anon: indicates we are operating on a virtual file and not a real file.
	//		JIT: This is hinting to the OS that we are using this for JIT purposes.
	printDbg("mapping %v bytes", len(code))
	fn, err := syscall.Mmap(-1, 0, len(code), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANON|syscall.MAP_JIT)
	check(err, "mmap")

	// put our code into the mmaped memory
	n := copy(fn, code)
	printDbg("copied %v bytes into memory", n)

	// Change the protections of our memory from RW to RX now that we no longer need
	// to modify it. This allows us to execute the code in this memory.
	printDbg("making memory executable")
	err = syscall.Mprotect(fn, syscall.PROT_READ|syscall.PROT_EXEC)
	check(err, "mprotect")

	// Convert our memory from a byte array to an executable function with some
	// very unsafe memory handling.
	printDbg("calling jit function")
	printDbg("====================")
	unsafeFunc := (uintptr)(unsafe.Pointer(&fn))
	if len(xargs) == 0 {
		f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
		f()
	} else {
		msg := strings.Join(xargs, " ") + "\n"
		f := *(*writeFunc)(unsafe.Pointer(&unsafeFunc))
		f(msg, uint(len(msg)))
	}
	printDbg("====================")

	// Unmap the memory that we mapped to release it from the application.
	printDbg("unmapping memory")
	err = syscall.Munmap(fn)
	check(err, "munmap")
	printDbg("done.")
}
