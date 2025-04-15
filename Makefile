.PHONY: asm

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+%?:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: asm ## Run the jit with just the ./asm/hello.s executed.
	@go run main.go

debug: asm ## Run the jit with debug logging enabled
	@DEBUG=1 go run main.go

asm: clean ## build the asm files that are embedded into the JIT
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

clean:
	@rm -rf ./build

# not needed right now as it seems like we are able to execute JIT code without
# needing this entitlement right now.
jit: asm ## build and codesign the application with the com.apple.security.cs.allow-jit entitlement
	@/usr/libexec/PlistBuddy -c "Add :com.apple.security.cs.allow-jit bool true" ./build/jit.entitlements 2>/dev/null
	@go build -o ./build/jit .
	@codesign -s - -f --entitlements ./build/jit.entitlements ./build/jit 2>/dev/null
	@./build/jit
