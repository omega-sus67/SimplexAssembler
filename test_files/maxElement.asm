; ==========================================
; FIND MAXIMUM ELEMENT IN AN ARRAY
; Scans 'array' and stores the largest value in 'max'
; ==========================================

        ; Initialize max = array[0]
        ldc array
        ldnl 0      ; A = array[0]
        ldc max
        stnl 0      ; max = array[0]

        ; Initialize ptr = array + 1 (point to second element)
        ldc array
        adc 1       ; A = address of array[1]
        ldc ptr
        stnl 0      ; ptr = array + 1

loop:
        ; Check if count == 0
        ldc count
        ldnl 0      ; A = count
        brz done    ; If count is 0, we have checked everything

        ; Load current array element (*ptr)
        ldc ptr
        ldnl 0      ; A = address currently in ptr
        ldnl 0      ; A = memory[ptr]
        
        ; Compare *ptr and max
        ldc max
        ldnl 0      ; A = max, B = *ptr
        sub         ; A = B - A (which is *ptr - max)
        
        ; If (*ptr - max) < 0, *ptr is smaller. Skip the update.
        brlz skip
        ; If (*ptr - max) == 0, they are equal. Skip the update.
        brz skip

        ; --- Update max ---
        ; If we get here, *ptr is strictly greater than max
        ldc ptr
        ldnl 0
        ldnl 0      ; A = *ptr
        ldc max
        stnl 0      ; max = *ptr

skip:
        ; Increment pointer (ptr++)
        ldc ptr
        ldnl 0
        adc 1
        ldc ptr
        stnl 0

        ; Decrement count (count--)
        ldc count
        ldnl 0
        adc -1
        ldc count
        stnl 0

        br loop     ; Go to next element

done:
        HALT

; --- VARIABLES ---
count:  data 4      ; 4 remaining elements to check
max:    data 0
ptr:    data 0

; --- ARRAY ---
array:  
        data 12
        data 45
        data 7
        data 89     ; <--- This should be the final max
        data 23