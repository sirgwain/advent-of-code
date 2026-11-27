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
#define ptr             R0  // input pointer
#define n               R1  // remaining input bytes to parse
#define b               R2   // current input byte
#define s1r             R3  // solution1 answer
#define s2r             R4  // solution2 answer
#define num1            R5
#define num2            R6
#define first           R7
#define tmp           R16


TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr
    MOVD input_len+8(FP), n
    CBZ n, done

    // init counters
    MOVD ZR, s1r
    MOVD ZR, s2r

    // remove the newline at the end if there
    SUB $1, n, tmp
    MOVBU (ptr)(tmp), b
    CMP $10, b
    BNE nonewline
    SUB $1, n, n // n-- if there is a \n

nonewline:    

    // record the first
    MOVBU (ptr), b // b = *input
    SUB $'0', b, first // first = b - '0'

    ADD $1, ptr, ptr // input++
    SUB $1, n, n // n--
loop:
    CBZ n, checkcircle // compare-branch-zero if n is down to 0

    MOVBU -1(ptr), b // b = *input++
    SUB $'0', b, num1 // num1 = b - '0'

    // Read a byte, increment input pointer
    MOVBU.P 1(ptr), b // b = *input++
    SUB $1, n, n // n--

    SUB $'0', b, num2 // num2 = b - '0'

    // check characters
    CMP num1, num2 
    BNE loop // if num1 != num2 loop

    ADD num1, s1r, s1r
    B loop

checkcircle:
    // check characters
    CMP first, num2 
    BNE done // if first != num2 goto done
    ADD num2, s1r, s1r

done:
    // ret
    MOVD s1r, s1+24(FP)
    MOVD s2r, s2+32(FP)
    RET

