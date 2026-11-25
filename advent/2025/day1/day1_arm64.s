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
#define b               R2  // current input byte
#define s1r             R3
#define s2r             R4
#define dial            R5
#define start           R6
#define sign            R7
#define digit           R8  // digit from *input for itoa
#define num             R9  // current number
#define k10             R10 // constant 10
#define tmp             R16 // scratch var


TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr
    MOVD input_len+8(FP), n

    // init counters
    MOVD $50, dial
    MOVD dial, start
    MOVD ZR, s1r
    MOVD ZR, s2r
    MOVD ZR, num
    MOVD $10, k10
    MOVD $1, sign

loop:
    CBZ n, done // compare-branch-zero if n is down to 0

    // Read a byte
    MOVBU (ptr), b // b = *input

    ADD $1, ptr, ptr // input++
    SUB $1, n, n // n--

    // L moves left, R moves right
    CMP $'R', b 
    BEQ right
    CMP $'L', b 
    BEQ left

    // check for 0-9 chars, anything else is a newline
    SUB $'0', b, digit // digit = *input - '0'
    CMP $9, digit
    BHI process // if digit > 9 goto process

    MADD k10, digit, num, num  // num = num*10 + digit
    
    B loop

process:
    MOVD dial, start // start = dial
    MUL sign, num, num // num = num * multipler (i.e. -1/1)
    ADDS num, dial, dial // dial += num

    BLT rollleft
    CMP $100, dial
    BGE rollright
    
    // didn't rollover, count zeros
    B countzeros

rollleft:
    // break out of the loop if dial is > 0
    CMP $0, dial
    BGE countzeros

    ADD $100, dial, dial // dial += 100
    MOVD start, tmp // store start
    MOVD dial, start // start = dial
    
     // if tmp == 0, loop again, otherwise tmp != 0 so increment s2
    CBZ tmp, rollleft
    ADD $1, s2r, s2r // else s2++

    // loop again
    B rollleft

rollright:
    // brek out of the loop when we are < 100
    CMP $100, dial
    BLT countzeros
    SUBS $100, dial, dial // dial = dial -100
    CINC NE, s2r, s2r // s2r++ if dial != 0
    B rollright
left:
    MOVD $-1, sign // sign = -1
    B loop
right:
    MOVD $1, sign // sign = 1
    B loop

countzeros:
    // zero out the number to start the loop again
    MOVD ZR, num

    // if dial != 0 loop again
    CBNZ dial, loop

    // dial is at 0, incrememnt our solutions and loop
    ADD $1, s1r, s1r
    ADD $1, s2r, s2r
    B loop

done:
    // ret
    MOVD s1r, s1+24(FP)
    MOVD s2r, s2+32(FP)
    RET

