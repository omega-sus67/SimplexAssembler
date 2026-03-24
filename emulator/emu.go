package main

//SHRUT GAUTAM
//2401CS62
//declaration of authorship-- This code is completely written by me shrut gautam
import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

type CPU struct {
	A  int32
	B  int32
	PC int32
	SP int32
}

var opNames = map[int32]string{
	0: "ldc", 1: "adc", 2: "ldl", 3: "stl", 4: "ldnl", 5: "stnl",
	6: "add", 7: "sub", 8: "shl", 9: "shr", 10: "adj", 11: "a2sp",
	12: "sp2a", 13: "call", 14: "return", 15: "brz", 16: "brlz",
	17: "br", 18: "HALT",
}

func dumpMemory(memory []int32, limit int) {
	fmt.Println("\n--- FINAL MEMORY DUMP ---")
	for i := 0; i < limit; i++ {
		// Dump all non-zero memory, or memory within the initial loaded block
		if memory[i] != 0 || i < limit {
			fmt.Printf("0x%04X : %10d (0x%08X)\n", i, memory[i], uint32(memory[i]))
		}
	}
}

func executeInstruction(cpu *CPU, memory []int32, opcode int32, operand int32) error {
	switch opcode {
	case 0:
		cpu.B = cpu.A
		cpu.A = operand
	case 1:
		cpu.A += operand
	case 2:
		cpu.B = cpu.A
		addr := cpu.SP + operand
		if addr < 0 || addr >= int32(len(memory)) {
			return fmt.Errorf("ldl out of bounds at %d", addr)
		}
		cpu.A = memory[addr]
	case 3:
		addr := cpu.SP + operand
		if addr < 0 || addr >= int32(len(memory)) {
			return fmt.Errorf("stl out of bounds at %d", addr)
		}
		memory[addr] = cpu.A
		cpu.A = cpu.B
	case 4:
		addr := cpu.A + operand
		if addr < 0 || addr >= int32(len(memory)) {
			return fmt.Errorf("ldnl out of bounds at %d", addr)
		}
		cpu.A = memory[addr]
	case 5:
		addr := cpu.A + operand
		if addr < 0 || addr >= int32(len(memory)) {
			return fmt.Errorf("stnl out of bounds at %d", addr)
		}
		memory[addr] = cpu.B
	case 6:
		cpu.A = cpu.B + cpu.A
	case 7:
		cpu.A = cpu.B - cpu.A
	case 8:
		cpu.A = cpu.B << cpu.A
	case 9:
		cpu.A = cpu.B >> cpu.A
	case 10:
		cpu.SP += operand
	case 11:
		cpu.SP = cpu.A
		cpu.A = cpu.B
	case 12:
		cpu.B = cpu.A
		cpu.A = cpu.SP
	case 13:
		cpu.B = cpu.A
		cpu.A = cpu.PC
		cpu.PC += operand
	case 14:
		cpu.PC = cpu.A
		cpu.A = cpu.B
	case 15:
		if cpu.A == 0 {
			cpu.PC += operand
		}
	case 16:
		if cpu.A < 0 {
			cpu.PC += operand
		}
	case 17:
		cpu.PC += operand
	default:
		return fmt.Errorf("illegal instruction opcode %d at PC %d", opcode, cpu.PC-1)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run emu.go <filename.o>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open object file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Prepare the trace file for writing
	baseName := strings.TrimSuffix(filename, ".o")
	traceFilename := baseName + "_trace.log"
	traceFile, err := os.Create(traceFilename)
	if err != nil {
		fmt.Printf("Failed to create trace file: %v\n", err)
		os.Exit(1)
	}
	defer traceFile.Close()

	memory := make([]int32, 65536)

	var instruction uint32
	memAddr := 0
	for {
		err := binary.Read(file, binary.LittleEndian, &instruction)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading binary: %v\n", err)
			os.Exit(1)
		}
		memory[memAddr] = int32(instruction)
		memAddr++
	}

	cpu := CPU{A: 0, B: 0, PC: 0, SP: 0}

	// Write header to the trace file
	fmt.Fprintln(traceFile, "===============================================================================")
	fmt.Fprintln(traceFile, "                              SIMPLEX EXECUTION TRACE")
	fmt.Fprintln(traceFile, "===============================================================================")
	fmt.Fprintf(traceFile, "%-10s | %-18s | %-12s | %-12s | %-12s\n", "PC", "INSTRUCTION", "REG A", "REG B", "REG SP")
	fmt.Fprintln(traceFile, "-------------------------------------------------------------------------------")

	for {
		if cpu.PC < 0 || cpu.PC >= int32(len(memory)) {
			errMsg := fmt.Sprintf("FATAL ERROR: PC out of bounds (%d)", cpu.PC)
			fmt.Println(errMsg)
			fmt.Fprintln(traceFile, errMsg)
			dumpMemory(memory, memAddr)
			break
		}

		machineCode := memory[cpu.PC]
		opcode := machineCode & 0xFF
		operand := machineCode >> 8

		mnemonic, ok := opNames[opcode]
		if !ok {
			mnemonic = "DATA"
		}
		instructionStr := fmt.Sprintf("%s %d", mnemonic, operand)

		// Log state BEFORE execution
		fmt.Fprintf(traceFile, "0x%08X | %-18s | %-12d | %-12d | %-12d\n",
			cpu.PC, instructionStr, cpu.A, cpu.B, cpu.SP)

		cpu.PC++

		if opcode == 18 { // HALT
			fmt.Println("PROGRAM HALTED SUCCESSFULLY.")
			fmt.Fprintln(traceFile, "-------------------------------------------------------------------------------")
			fmt.Fprintln(traceFile, "PROGRAM HALTED SUCCESSFULLY.")
			dumpMemory(memory, memAddr)
			break
		}

		err := executeInstruction(&cpu, memory, opcode, operand)
		if err != nil {
			fmt.Printf("FATAL ERROR: %v\n", err)
			fmt.Fprintf(traceFile, "FATAL ERROR: %v\n", err)
			dumpMemory(memory, memAddr)
			break
		}
	}

	fmt.Printf("Execution trace saved to %s\n", traceFilename)
}
