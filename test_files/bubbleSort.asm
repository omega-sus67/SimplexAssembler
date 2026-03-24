; ==========================================
; SIMPLEX BUBBLE SORT 
; Sorts an array in ascending order
; ==========================================

outerloop:
        ; Check if swapped == 0. If yes, the array is sorted.
        ldc swapped
        ldnl 0
        brz endsort

        ; Reset swapped to 0 for this pass
        ldc 0
        ldc swapped
        stnl 0

        ; Reset j to 0
        ldc 0
        ldc j
        stnl 0

innerloop:
        ; Check if j < nlimit (j - nlimit < 0)
        ldc j
        ldnl 0
        ldc nlimit
        ldnl 0
        sub         ; A = B - A (evaluates to j - nlimit)
        brlz docompare
        br endinner

docompare:
        ; --- THE FIX: Safe memory loading ---
        ; 1. Load arr[j+1] and safely stash it in tmp1
        ldc j
        ldnl 0
        ldc array
        add         
        adc 1       
        ldnl 0      ; A = arr[j+1]
        ldc tmp1    ; B = arr[j+1], A = addr(tmp1)
        stnl 0      ; memory[tmp1] = B

        ; 2. Load arr[j] into A
        ldc j
        ldnl 0
        ldc array
        add         
        ldnl 0      ; A = arr[j]

        ; 3. Load arr[j+1] back into A, safely pushing arr[j] into B
        ldc tmp1
        ldnl 0      ; B = arr[j], A = arr[j+1]

        ; 4. Compare arr[j] and arr[j+1]
        sub         ; A = arr[j] - arr[j+1]
        
        ; If arr[j] - arr[j+1] <= 0, they are in order. Skip the swap.
        brlz noswap
        brz noswap

        ; --- SWAP ROUTINE ---
        ; Calculate and store ptr1 (address of arr[j])
        ldc j
        ldnl 0
        ldc array
        add         
        ldc ptr1
        stnl 0

        ; Calculate and store ptr2 (address of arr[j+1])
        ldc j
        ldnl 0
        ldc array
        add
        adc 1
        ldc ptr2
        stnl 0

        ; tmp1 = *ptr1
        ldc ptr1
        ldnl 0      ; Load address stored in ptr1
        ldnl 0      ; Load actual array value
        ldc tmp1
        stnl 0

        ; tmp2 = *ptr2
        ldc ptr2
        ldnl 0      
        ldnl 0      
        ldc tmp2
        stnl 0

        ; *ptr1 = tmp2
        ldc tmp2
        ldnl 0
        ldc ptr1
        ldnl 0
        stnl 0

        ; *ptr2 = tmp1
        ldc tmp1
        ldnl 0
        ldc ptr2
        ldnl 0
        stnl 0

        ; swapped = 1 (flag that a swap occurred)
        ldc 1
        ldc swapped
        stnl 0

noswap:
        ; j++
        ldc j
        ldnl 0
        adc 1
        ldc j
        stnl 0

        br innerloop

endinner:
        ; nlimit-- (The last element is already in place)
        ldc nlimit
        ldnl 0
        adc -1
        ldc nlimit
        stnl 0

        br outerloop

endsort:
        HALT        ; Stop the emulator

; --- VARIABLES & DATA SEGMENT ---
nlimit:  data 4     ; n - 1 (since we have 5 elements)
swapped: data 1     ; Initialize to 1 to start the outer loop
j:       data 0
ptr1:    data 0
ptr2:    data 0
tmp1:    data 0
tmp2:    data 0

; --- THE ARRAY TO SORT ---
array:   
         data 42
         data 7
         data 99
         data 3
         data 18