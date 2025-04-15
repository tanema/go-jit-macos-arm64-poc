# JIT in Go on Apple Silicon
This is not straight forward and there are not really any examples out there so
this will hopefully help anyone else trying to do this.

- [reference](https://github.com/zeusdeux/jit-example-macos-arm64)
- [JIT in C](https://medium.com/@gamedev0909/jit-in-c-injecting-machine-code-at-runtime-1463402e6242)
- [porting JIT to apple silicon](https://developer.apple.com/documentation/apple-silicon/porting-just-in-time-compilers-to-apple-silicon?language=objc)
- [extract flat (pure) binary](https://stackoverflow.com/a/13306947)
- [Making system calls from Assembly in Mac OS X](https://filippo.io/making-system-calls-from-assembly-in-mac-os-x/)
- [arm64 syscalls](https://stackoverflow.com/questions/56985859/ios-arm64-syscalls)

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
