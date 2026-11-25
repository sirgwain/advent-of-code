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
//  s1+24(FP)
//  s2+32(FP)
// 
// registers:
#define ptr             R0  // input pointer
#define n               R1  // remaining input bytes to parse
#define b               R2   // current input byte
#define floor           R3  // current floor
#define pos             R4  // current position
#define basement        R5  // basement


TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr    // pointer into ptr
    MOVD input_len+8(FP), n    // length into n

    // init counters
    MOVD ZR, floor // currentFloor = 0
    MOVD ZR, pos // position = 0
    MOVD ZR, basement // firstBasementPosition = 0

loop:
    CBZ n, done // compare-branch-zero if n is down to 0

    // Read a byte
    MOVBU (ptr), b // b = *input

    ADD $1, ptr, ptr // input++
    SUB $1, n, n // n--
    ADD $1, pos, pos // pos++

    // check characters
    CMP $'(', b 
    BEQ up // if b == '(' goto up

    // check characters
    CMP $')', b 
    BNE loop // if b != ')' goto loop
    // fallthrough to down
down:
    SUB $1, floor, floor // floor-- (currentFloor--)

    // if we already found the basement, skip this
    CBNZ basement, loop // if basement != 0 (firstBasementPosition != 0)

    // is the floor -1?
    CMP $-1, floor 
    BNE loop // if floor != -1 goto loop (currentFloor != -1)

    // save the first basement
    MOVD pos, basement // basement = pos
    B loop
up:
    ADD $1, floor, floor // floor = floor + 1
    B loop

done:
    // ret
    MOVD floor, s1+24(FP)
    MOVD basement, s2+32(FP)
    RET

