; --- Initialization ---
ldc 0x1000      ; Load initial stack pointer address into Accumulator (A) [cite: 115]
a2sp            ; Transfer A to Stack Pointer (SP) [cite: 116]
adj -1          ; Allocate 1 word on the stack for the result pointer [cite: 117]
ldc result      ; Load the memory address of the 'result' array into A [cite: 118]
stl 0           ; Store the result pointer at SP[0] [cite: 119]
ldc count       ; Load the memory address of the 'count' variable into A [cite: 120]
ldnl 0          ; Dereference A to load the actual value of 'count' into A [cite: 121]
call main       ; Call the main subroutine (B = count, A = return address) [cite: 122]
adj 1           ; Clean up the stack allocation after execution finishes [cite: 123]
HALT            ; Terminate emulator execution [cite: 124]

; --- Main Subroutine ---
main:           ; Main entry point for the loop [cite: 125]
adj -3          ; Allocate 3 words on the stack for local variables [cite: 126]
stl 1           ; Save return address to SP[1] [cite: 127]
stl 2           ; Save 'count' argument to SP[2] [cite: 128]
ldc 0           ; Load 0 into A [cite: 129]
; zero accumulator ; Zero out the accumulator (original comment) [cite: 130]
stl 0           ; Initialize loop counter 'i' to 0 at SP[0] [cite: 131]

loop:           ; Loop entry point [cite: 132]
adj -1          ; Allocate 1 word on stack to pass argument to subroutine [cite: 133]
ldl 3           ; Load 'count' into A [cite: 134]
stl 0           ; Save it to SP[0] [cite: 135]
ldl 1           ; Load loop counter 'i' into A [cite: 136]
call triangle   ; Call the triangle subroutine [cite: 137]
adj 1           ; Clean up argument from stack after return [cite: 138]

ldl 3           ; Load 'result' array pointer into A [cite: 139]
stnl 0          ; Store the subroutine's returned value into the array [cite: 140]
ldl 3           ; Reload the array pointer [cite: 141]
adc 1           ; Increment the pointer to point to the next array element [cite: 142]
stl 3           ; Save the updated array pointer back to the stack [cite: 143]

ldl 0           ; Load loop counter 'i' [cite: 144]
adc 1           ; Increment 'i' by 1 [cite: 145]
stl 0           ; Save the updated 'i' to the stack [cite: 146]

ldl 0           ; Load the updated 'i' into A [cite: 147]
ldl 2           ; Load 'count' into A, shifting 'i' into B [cite: 148]
sub             ; Subtract to check loop bounds (i - count) [cite: 149]
brlz loop       ; Branch back to 'loop' if i < count (result is negative) [cite: 150]

ldl 1           ; Load return address [cite: 151]
adj 3           ; Deallocate local variables from the stack [cite: 152]
; reload it     ; Reload return address (original comment) [cite: 153]
; get return address ; Get return address (original comment) [cite: 154]
return          ; Return to the caller [cite: 155]

; --- Triangle Subroutine ---
triangle:       ; Recursive subroutine entry point [cite: 156]
adj -3          ; Allocate 3 words for local variables [cite: 156]
stl 1           ; Save return address [cite: 157]
stl 2           ; Save argument 'n' [cite: 158]
ldc 1           ; Load constant 1 [cite: 159]
shl             ; Bitwise shift left [cite: 160]
ldl 3           ; Load local variable from stack [cite: 161]
sub             ; Subtract A from B [cite: 162]
brlz skip       ; Branch if less than zero to 'skip' [cite: 163]

ldl 3           ; Load local variable [cite: 164]
ldl 2           ; Load argument 'n' [cite: 165]
sub             ; Subtract [cite: 166]
stl 2           ; Store updated value to stack [cite: 167]

skip:           ; Branch target [cite: 168]
ldl 2           ; Load value from stack [cite: 168]
brz one         ; Branch if zero to 'one' [cite: 169]
ldl 3           ; Load local variable [cite: 170]
adc -1          ; Add -1 [cite: 171]
stl 0           ; Store to SP[0] [cite: 172]

; --- THE FIX: The orphaned base case instructions were moved here ---
one:            ; Base case branch target [cite: 173]
ldc 1           ; Base case return value (1) [cite: 193]
ldl 1           ; Load return address [cite: 194]
adj 3           ; Clean up stack frame [cite: 195]
return          ; Return from base case [cite: 196]
; --------------------------------------------------------------------

adj -1          ; Allocate 1 word on stack (note: count: label removed from here) [cite: 176]
ldl 1           ; Load value from stack [cite: 177]
stl 0           ; Store to SP[0] [cite: 178]
ldl 3           ; Load local variable [cite: 179]
adc -1          ; Add -1 (decrement for recursion) [cite: 180]
call triangle   ; First recursive call [cite: 181]

ldl 1           ; Load value from stack [cite: 182]
stl 0           ; Store to SP[0] [cite: 183]
stl 1           ; Store to SP[1] [cite: 184]
ldl 3           ; Load local variable [cite: 185]
call triangle   ; Second recursive call [cite: 186]

adj 1           ; Deallocate stack space [cite: 187]
ldl 0           ; Load result [cite: 188]
add             ; Add results of recursive calls together [cite: 189]
ldl 1           ; Load return address [cite: 190]
adj 3           ; Clean up stack frame [cite: 191]
return          ; Return combined result [cite: 192]

; --- Data Section ---
count: data 10  ; The variable 'count', initialized to 10 (label moved here) [cite: 175, 197]
result: data 0  ; Base address to start storing the result array [cite: 198]