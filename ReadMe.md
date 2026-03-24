# CS2206 MiniProject: SIMPLEX Assembler & Emulator Usage Guide
**Author:** Shrut Gautam (2401CS62)

## Declaration of Authorship
I certify that the code provided in `asm.go` and `emu.go` is entirely my own work, written and submitted for the CS2206 MiniProject.

## Prerequisites
* **Go Environment:** The assembler and emulator are implemented in Go. Ensure you have Go installed on your system to compile and execute the source files.

## Compilation Instructions
While the project specification references `gcc` for C programs, this implementation uses Go. You can compile the source files into standalone executables using the `go build` command.

To build the executables manually:
go build asm.go
go build emu.go

---

## 1. The Assembler (`asm.go`)
The assembler processes a SIMPLEX assembly language text file in two passes to assign label values, instruction opcodes, and generate machine code.

**Execution:**
./asm <filename.asm>

# Alternatively, run directly from source without building:
go run asm.go <filename.asm>

**Outputs Produced:**
* **`<filename>.o`**: A binary object file containing the translated machine code. The execution code starts at address zero.
* **`<filename>.lst`**: A human-readable listing file displaying the memory address, the 32-bit machine code (formatted as 8 hex characters), and the original mnemonic and operand for easy debugging.

---

## 2. The Emulator (`emu.go`)
The emulator loads the binary object file into its memory and executes the SIMPLEX instructions.

**Execution:**
./emu <filename.o>

# Alternatively, run directly from source without building:
go run emu.go <filename.o>

**Outputs Produced:**
* **Execution Trace File (`<filename>_trace.txt`)**: This file serves as the execution output log for your tests. It provides a step-by-step trace of the program, recording the Program Counter (PC), the current instruction, and the states of registers (A, B, and SP) just before each instruction is executed.
* **Console Memory Dump**: Upon successfully reaching the `HALT` instruction, the memory is printed inside the terminal