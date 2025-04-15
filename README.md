# JIT in Go on Apple Silicon
Writing a JIT compiler on apple silicon can be extra hard because of the hardened
runtime, and there are no specific exampe that exist for go so I wanted to publish
one to make it easier for others to figure out. This will not get you to the point
of code generation though as this example uses the assembler from the macos toolchain
to generate and link the binary code used in this JIT.

## Instructions
- `make help` show all make targets
- `make asm` to build assembly and generate flat binaries
- `make run` to build assembly, generate flat binaries, and run the go code
- `go run main.go "Hello from Tim"` (only after building asm) to call the write.s code
- `make jit` run the app as a signed binary (see macos jit notes)

## Reference

- [C repo used for some insight on how to do this in Go](https://github.com/zeusdeux/jit-example-macos-arm64)
- [JIT in C article](https://medium.com/@gamedev0909/jit-in-c-injecting-machine-code-at-runtime-1463402e6242)
- [porting JIT to apple silicon](https://developer.apple.com/documentation/apple-silicon/porting-just-in-time-compilers-to-apple-silicon?language=objc)
- [extract flat (pure) binary](https://stackoverflow.com/a/13306947)
- [Making system calls from Assembly in Mac OS X](https://filippo.io/making-system-calls-from-assembly-in-mac-os-x/)
- [arm64 syscalls](https://stackoverflow.com/questions/56985859/ios-arm64-syscalls)

## MacOS JIT notes
- `mmap` cannot be mapped with RWX permissions [because of the hardened runtime](https://stackoverflow.com/questions/74124485/mmap-rwx-page-on-macos-arm64-architecture)
  so it needs to be mapped with RW, the data copied in, and then the memory protections
  should be toggled to RX to allow for execution.
- Regarding entitlements, a lot of people say you cannot execute JIT code without
  the `com.apple.security.cs.allow-jit` entitlement signed into the binary. I have
  found that this repo runs without it but please let me know if you have problems.
  There is a make target to build this repo with those entitlements just in case
  and you can try it with `make jit`

## extract flat (pure) binary:
This is just a shortcut so that you can use an assembler to get the raw binary but
use free standing environment code in your jit without the OS header

  1. `otool -l hello.bin | grep -A4 "sectname __text" | tail -1`  (offset field is in decimal not hex btw)
    - Take the offset, convert to hex and verify code starts there in the hexdump view of the compiled binary
  2. `dd if=hello.bin of=hello_flat.bin ibs=<offset> skip=1`

```bash
  otool -l hello.bin |\
    grep -A4 "sectname __text" |\
    tail -1 |\
    grep -o "\d+" |\
    xargs -n1 -I% dd if=hello.bin of=hello_flat.bin ibs=% skip=1
```
