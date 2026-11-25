//go:build arm64 && !purego

#include "textflag.h"
// func solve(input []byte) (floor int, basement int)
//
// args:
//  input+0(FP) = data pointer
//  input+8(FP) = length
//  input+16(FP) = capacity
//
// return:
//  floor+24(FP)
//  basement+32(FP)

TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), R0    // pointer into R0
    MOVD input_len+8(FP), R1    // length into R1

    // init counters
    MOVD ZR, R2 // currentFloor = 0
    MOVD ZR, R3 // position = 0
    MOVD ZR, R4 // firstBasementPosition = 0

loop:
    CBZ R1, done // compare-branch-zero if R1 is down to 0

    // Read a byte
    MOVBU (R0), R5 // R5 = *input

    ADD $1, R0, R0 // input++
    SUB $1, R1, R1 // R1--
    ADD $1, R3, R3 // pos++

    // check characters
    CMP $'(', R5 
    BEQ up // if R5 == '(' goto up

    // check characters
    CMP $')', R5 
    BNE loop // if R5 != ')' goto loop
    // fallthrough to down
down:
    SUB $1, R2, R2 // R2-- (currentFloor--)

    // if we already found the basement, skip this
    CBNZ R4, loop // if R4 != 0 (firstBasementPosition != 0)

    // is the floor -1?
    CMP $-1, R2 
    BNE loop // if R2 != -1 goto loop (currentFloor != -1)

    // save the first basement
    MOVD R3, R4 // R4 = R3
    B loop
up:
    ADD $1, R2, R2 // R2 = R2 + 1
    B loop

done:
    // ret
    MOVD R2, floor+24(FP)
    MOVD R4, basement+32(FP)
    RET

