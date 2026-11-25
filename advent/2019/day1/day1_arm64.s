//go:build arm64 && !purego

#include "textflag.h"
// func solve(input []byte) (s1 int, s2 int)
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
// R0 input
// R1 i = len(input)
// R2 b byte
// R3 num
// R4 digit
// R5 fuel
// R6 s1
// R7 s2
// R10 10 (digit constant)

TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), R0    // pointer into R0
    MOVD input_len+8(FP), R1    // length into R1

    // init counters
    MOVD ZR, R3 // num = 0
    MOVD ZR, R4 // digit = 0
    MOVD ZR, R6 // s1 = 0
    MOVD ZR, R7 // s2 = 0
    MOVD $10, R10 // atoi 10x multiplier
    MOVD $3, R11 // fuel / 3 const

loop:
    CBZ R1, done // compare-branch-zero if R1 is down to 0

    // Read a byte
    MOVBU (R0), R2 // b = *input

    ADD $1, R0, R0 // input++
    SUB $1, R1, R1 // R1--

    SUB $'0', R2, R4 // digit = b - '0'
    CMP $9, R4
    BHI calcfuel // if digit > 9 goto calcfuel

    MADD R10, R4, R3, R3 // num = num*10 + digit
    B loop

calcfuel:
    CBZ R3, loop // if num == 0 (like a newline), skip it

    UDIV R11, R3, R5 // fuel = num / 3
    SUB $2, R5, R5 // fuel = fuel - 2

    ADD R5, R6, R6 // s1 += fuel
    ADD R5, R7, R7 // s2 += fuel

    CMP $6, R5
    BGT calcfuelfuel

    MOVD ZR, R3 // num = 0
    B loop

calcfuelfuel:
    UDIV R11, R5, R5 // fuel = fuel / 3
    SUB $2, R5, R5 // fuel = fuel - 2
    ADD R5, R7, R7 // s2 += fuel

    // keep looping while fuel > 6
    CMP $6, R5
    BGT calcfuelfuel // if fuel > 6 goto calccalcfuel

    MOVD ZR, R3 // num = 0
    B loop
done:
    CBNZ R3, calcfuel   // flush final number if input lacks trailing \n
    MOVD R6, s1+24(FP)
    MOVD R7, s2+32(FP)
    RET

