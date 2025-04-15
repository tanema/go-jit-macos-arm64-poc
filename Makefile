.PHONY: asm

run: asm
	@go run main.go

asm:
	@mkdir -p ./build
	@as -o ./build/hello.o ./asm/hello.s
	@as -o ./build/write.o ./asm/write.s
	@ld -o ./build/hello.bin ./build/hello.o -e _start -arch arm64
	@ld -o ./build/write.bin ./build/write.o -e _start -arch arm64
	@otool -l ./build/hello.bin \
		| grep -A4 "sectname __text" \
		| tail -1 \
		| grep -o "\d*" \
		| xargs -n1 -I% dd if=./build/hello.bin of=./build/hello_flat.bin ibs=% skip=1 2>/dev/null
	@otool -l ./build/write.bin \
		| grep -A4 "sectname __text" \
		| tail -1 \
		| grep -o "\d*" \
		| xargs -n1 -I% dd if=./build/write.bin of=./build/write_flat.bin ibs=% skip=1 2>/dev/null

# not needed right now
jit: entitlements
	@go build -o jit .
	@codesign -s - -f --entitlements jit.entitlements jit

entitlements:
	@rm -f ./jit.entitlements && /usr/libexec/PlistBuddy -c "Add :com.apple.security.cs.allow-jit bool true" jit.entitlements
