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
// note R16,R17 are scratch that could be blown away by CALLs
// R18, R27+ are reserved
// ptr input
// n i = len(input)
// b b byte
// num num
// digit digit
// fuel fuel
// s1r s1
// s2r s2
// k10 10 (digit constant)
#define ptr         R0  // input pointer
#define n           R1  // remaining input bytes to parse
#define b           R2  // current byte from input
#define s1r         R3  // solution1 register
#define s2r         R4  // solution2 register
#define num         R5
#define digit       R6
#define fuel        R7
#define k10         R10 // remaining input bytes to parse

TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr    // pointer into ptr
    MOVD input_len+8(FP), n    // length into n

    // init counters
    MOVD ZR, num // num = 0
    MOVD ZR, digit // digit = 0
    MOVD ZR, s1r // s1 = 0
    MOVD ZR, s2r // s2 = 0
    MOVD $10, k10 // atoi 10x multiplier
    MOVD $3, R11 // fuel / 3 const

loop:
    CBZ n, done // compare-branch-zero if n is down to 0

    // Read a byte
    MOVBU (ptr), b // b = *input

    ADD $1, ptr, ptr // input++
    SUB $1, n, n // n--

    SUB $'0', b, digit // digit = b - '0'
    CMP $9, digit
    BHI calcfuel // if digit > 9 goto calcfuel

    MADD k10, digit, num, num // num = num*10 + digit
    B loop

calcfuel:
    CBZ num, loop // if num == 0 (like a newline), skip it

    UDIV R11, num, fuel // fuel = num / 3
    SUB $2, fuel, fuel // fuel = fuel - 2

    ADD fuel, s1r, s1r // s1 += fuel
    ADD fuel, s2r, s2r // s2 += fuel

    CMP $6, fuel
    BGT calcfuelfuel

    MOVD ZR, num // num = 0
    B loop

calcfuelfuel:
    UDIV R11, fuel, fuel // fuel = fuel / 3
    SUB $2, fuel, fuel // fuel = fuel - 2
    ADD fuel, s2r, s2r // s2 += fuel

    // keep looping while fuel > 6
    CMP $6, fuel
    BGT calcfuelfuel // if fuel > 6 goto calccalcfuel

    MOVD ZR, num // num = 0
    B loop
done:
    CBNZ num, calcfuel   // flush final number if input lacks trailing \n
    MOVD s1r, s1+24(FP)
    MOVD s2r, s2+32(FP)
    RET

