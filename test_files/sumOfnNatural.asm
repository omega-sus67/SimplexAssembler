; ==========================================
; SUM OF FIRST N NATURAL NUMBERS
; Calculates N + (N-1) + ... + 1 and stores in 'sum'
; ==========================================

loop:
        ; Check if N == 0
        ldc n
        ldnl 0      ; A = n
        brz done    ; If N reached 0, the sum is complete

        ; sum = sum + n
        ldc sum
        ldnl 0      ; A = sum
        ldc n
        ldnl 0      ; A = n, B = sum
        add         ; A = B + A (sum + n)
        ldc sum     ; A = address of sum, B = new sum
        stnl 0      ; memory[sum] = B

        ; Decrement n (n--)
        ldc n
        ldnl 0      ; A = n
        adc -1      ; A = n - 1
        ldc n
        stnl 0      ; memory[n] = n - 1

        br loop     ; Repeat the loop

done:
        HALT

; --- VARIABLES ---
n:      data 10     ; The N value. Sum of 1 to 10 is 55.
sum:    data 0      ; This will hold 55 (0x37) when the program halts