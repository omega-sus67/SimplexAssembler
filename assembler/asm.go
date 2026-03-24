package main
//SHRUT GAUTAM
//2401CS62
//declaration of authorship-- This code is completely written by me shrut gautam

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
)

var labelRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)

// Maps mnemonics to a boolean indicating if they require an operand. 
// This helps us quickly validate instruction syntax during the first pass.
var opcodes = map[string]bool{
    "data": true, "ldc": true, "adc": true, "ldl": true, "stl": true,
    "ldnl": true, "stnl": true, "adj": true, "call": true,
    "brz": true, "brlz": true, "br": true, "SET": true,
    "add": false, "sub": false, "shl": false, "shr": false,
    "a2sp": false, "sp2a": false, "return": false, "HALT": false,
}

// Stores the base 8-bit opcode values for each SIMPLEX instruction 
// to be used during machine code generation in the second pass.
var opcodeVals = map[string]uint32{
    "ldc": 0, "adc": 1, "ldl": 2, "stl": 3, "ldnl": 4, "stnl": 5,
    "add": 6, "sub": 7, "shl": 8, "shr": 9, "adj": 10, "a2sp": 11,
    "sp2a": 12, "call": 13, "return": 14, "brz": 15, "brlz": 16,
    "br": 17, "HALT": 18,
}

// purify strips away comments and trims whitespace so the parser 
// only has to deal with the actual instruction.
func purify(s string) string {
    parts := strings.SplitN(s, ";", 2)
    comless := strings.TrimSpace(parts[0])
    return comless
}

// labelCheck determines if the current line defines a new label.
// It separates the label from the instruction and validates the label's syntax.
func labelCheck(l string, lineNum int) (string, string, error) {
    // Quick check: if there's no colon, it definitely isn't a label declaration.
    if !strings.Contains(l, ":") {
        return "", l, nil
    }

    // Isolate the label name from the actual instruction that might follow on the same line.
    spl := strings.SplitN(l, ":", 2)
    if len(spl) == 1 {
        return "", l, nil
    }

    label := strings.TrimSpace(spl[0])
    command := strings.TrimSpace(spl[1])

    // Catch edge cases where the line starts with a colon (e.g., ": ldc 5") 
    // to prevent malformed labels.
    if label == "" {
        return "", "", fmt.Errorf("line %d: missing label name before colon", lineNum)
    }

    // Enforce the SIMPLEX assembly rule that valid labels must be alphanumeric 
    // and start with a letter.
    if !labelRegex.MatchString(label) {
        return "", "", fmt.Errorf("line %d: bogus label name '%s'", lineNum, label)
    }

    return label, command, nil
}

// firstpass scans the entire file to build the symbol table (resolving label addresses) 
// and catch early syntax errors without outputting any machine code.
func firstpass(file *os.File) (map[string]int, []string, []error, []error) {
    // The symbol table holds the name and PC location of every label and SET variable.
    symtable := make(map[string]int) 
    var parsedLines []string
    var errList []error
    var warnList []error 

    pc := 0
    
    // Keep track of the actual line number for accurate error reporting to the user.
    currline := 0 

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        currline += 1
        line := purify(scanner.Text()) 
        
        // Skip blank lines or pure comment lines so they don't consume PC addresses.
        if line == "" {
            continue
        }

        label, command, err := labelCheck(line, currline)

        if err != nil {
            errList = append(errList, err)
        } else if label != "" {
            // Register valid, unique labels into the symbol table. 
            // We flag duplicate declarations immediately to prevent ambiguous branch targets.
            if _, ok := symtable[label]; ok {
                errList = append(errList, fmt.Errorf("line %d: duplicate label definition '%s'", currline, label))
            } else {
                symtable[label] = pc
            }
        }

        if command != "" {
            parts := strings.Fields(command)
            mnemonic := parts[0]

            takesOperand, isValid := opcodes[mnemonic]
            
            // Verify the mnemonic exists in the SIMPLEX architecture. 
            // If not, flag it and skip further parsing for this line.
            if !isValid {
                errList = append(errList, fmt.Errorf("line %d: bogus mnemonic '%s'", currline, mnemonic))
                continue
            }

            if takesOperand {
                // Enforce the strict 1-operand rule for this specific instruction type.
                if len(parts) == 1 {
                    errList = append(errList, fmt.Errorf("line %d: missing operand for '%s'", currline, mnemonic))
                } else if len(parts) > 2 {
                    errList = append(errList, fmt.Errorf("line %d: extra on end of line", currline))
                } else {
                    operand := parts[1]
                    
                    // Determine if the operand is a numeric literal or a label reference.
                    _, errParse := strconv.ParseInt(operand, 0, 32)
                    isLabel := labelRegex.MatchString(operand)

                    if errParse != nil && !isLabel {
                        errList = append(errList, fmt.Errorf("line %d: not a number / invalid operand '%s'", currline, operand))
                    }
                    
                    // Handle the SET pseudo-instruction separately since it assigns 
                    // a custom value to a label instead of the current PC.
                    if mnemonic == "SET" {
                        if label == "" {
                            errList = append(errList, fmt.Errorf("line %d: SET instruction requires a label", currline))
                        } else if errParse != nil {
                            errList = append(errList, fmt.Errorf("line %d: SET operand must be a valid number", currline))
                        } else {
                            val, _ := strconv.ParseInt(operand, 0, 32)
                            symtable[label] = int(val)
                        }
                    }
                }
            } else {
                // Ensure zero-operand instructions (like HALT or a2sp) don't have stray arguments trailing them.
                if len(parts) > 1 {
                    errList = append(errList, fmt.Errorf("line %d: unexpected operand for '%s'", currline, mnemonic))
                }
            }

            // SET is a pseudo-instruction, so it doesn't take up space in the final binary. 
            // We only increment the PC for actual machine instructions.
            if mnemonic != "SET" {
                parsedLines = append(parsedLines, command)
                pc += 1
            }
        }
    }

    // Catch standard file I/O errors that might occur during the scan.
    if err := scanner.Err(); err != nil {
        errList = append(errList, fmt.Errorf("file read error: %v", err))
    }

    return symtable, parsedLines, errList, warnList
}

// isBranch is a helper to identify branch instructions, as their operand values 
// need to be calculated as relative PC offsets rather than absolute addresses.
func isBranch(mnemonic string) bool {
    return mnemonic == "br" || mnemonic == "brz" || mnemonic == "brlz" || mnemonic == "call"
}

// secondpass uses the built symbol table to generate the final machine code, 
// calculate relative branch offsets, and write the output files.
func secondpass(symtable map[string]int, parsedLines []string, filename string) ([]error, []error) {
    var errList []error
    var warnList []error
    usedLabels := make(map[string]bool)

    baseName := strings.TrimSuffix(filename, ".asm")
    objName := baseName + ".o"
    lstName := baseName + ".lst"

    objFile, err := os.Create(objName)
    if err != nil {
        return append(errList, fmt.Errorf("failed to create object file: %v", err)), warnList
    }
    defer objFile.Close()

    lstFile, err := os.Create(lstName)
    if err != nil {
        return append(errList, fmt.Errorf("failed to create listing file: %v", err)), warnList
    }
    defer lstFile.Close()

    pc := 0

    for _, line := range parsedLines {
        parts := strings.Fields(line)
        mnemonic := parts[0]

        var operandVal int32 = 0

        if len(parts) > 1 {
            operandStr := parts[1]
            
            // First, try parsing the operand as a raw number. 
            // If it fails, treat it as a label reference.
            val, errParse := strconv.ParseInt(operandStr, 0, 32)
            if errParse == nil {
                operandVal = int32(val)
            } else {
                addr, exists := symtable[operandStr]
                
                // If the label wasn't found in pass 1, we must fail it here. 
                // We assign a dummy 0 value to maintain file alignment.
                if !exists {
                    errList = append(errList, fmt.Errorf("PC %04d: no such label '%s'", pc, operandStr))
                    operandVal = 0 
                } else {
                    usedLabels[operandStr] = true

                    // For branches, we calculate the relative offset. 
                    // Remember, the PC implicitly increments before execution, so we subtract (pc+1).
                    if isBranch(mnemonic) {
                        operandVal = int32(addr) - int32(pc+1)
                    } else {
                        operandVal = int32(addr)
                    }
                }
            }
        }

        var machineCode uint32
        
        // The 'data' directive writes the raw 32-bit value directly. 
        // For normal instructions, we pack the 8-bit opcode and the 24-bit masked operand into a single 32-bit word.
        if mnemonic == "data" {
            machineCode = uint32(operandVal)
        } else {
            opCode := opcodeVals[mnemonic]
            maskedOperand := uint32(operandVal) & 0xFFFFFF
            machineCode = (maskedOperand << 8) | opCode
        }

        // Write the instruction to the binary file in little-endian format as required.
        err = binary.Write(objFile, binary.LittleEndian, machineCode)
        if err != nil {
            errList = append(errList, fmt.Errorf("failed to write to binary file at PC %d", pc))
        }

        fmt.Fprintf(lstFile, "%08X %08X %s\n", pc, machineCode, line)

        pc++
    }

    // Post-generation check: alert the user about any labels they defined but never jumped to or referenced.
    for label := range symtable {
        if !usedLabels[label] {
            warnList = append(warnList, fmt.Errorf("unused label '%s'", label))
        }
    }

    return errList, warnList
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run asm.go <filename.asm>")
        os.Exit(1)
    }

    filename := os.Args[1]
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Failed to open file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    symTable, parsedLines, errListPass1, _ := firstpass(file)

    errListPass2, warnListPass2 := secondpass(symTable, parsedLines, filename)

    allErrors := append(errListPass1, errListPass2...)
    
    fmt.Printf("--- Assembler Results for %s ---\n", filename)
    
    if len(allErrors) > 0 {
        fmt.Println("\nERRORS:")
        for _, e := range allErrors {
            fmt.Println("-", e)
        }
    } else {
        fmt.Println("\nAssembly completed successfully. No errors.")
    }

    if len(warnListPass2) > 0 {
        fmt.Println("\nWARNINGS:")
        for _, w := range warnListPass2 {
            fmt.Println("-", w)
        }
    }
}