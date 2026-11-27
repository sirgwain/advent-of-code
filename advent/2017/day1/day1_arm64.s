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
#define bPrev           R5
#define bCur            R6
#define pFirst          R7
#define bFirst          R8
#define half            R9
#define pMid            R10
#define pEnd            R11
#define tmp             R16


TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr
    MOVD input_len+8(FP), n

    // init counters
    MOVD ZR, s1r
    MOVD ZR, s2r

    // remove the newline at the pEnd if there
    ADD n, ptr, pEnd    
    MOVBU -1(pEnd), b
    CMP $10, b
    BNE trimmed
    SUB $1, n, n // n-- if there is a \n
    SUB $1, pEnd, pEnd

trimmed:
    MOVD ptr, pFirst // pFirst = &input[0]

    LSR $1, n, half // half = n/2
    ADD half, ptr, pMid // pMid = &input[half]

loop2:
    MOVBU (ptr)(half), bCur // bCur = input[half]
    MOVBU.P 1(ptr), bPrev // bPrev = *input++
    SUB $'0', bCur, tmp // tmp = bCur - '0'
    CMP bPrev, bCur
    CSEL EQ, tmp, ZR, tmp // tmp = (bPrev == bCur) ? tmp : 0
    ADD tmp, s2r, s2r // s2 += tmp
    CMP pMid, ptr
    BLO loop2 // loop until pMid == ptr

    // double the result when done
    LSL $1, s2r, s2r

part1:
    // loop through each bPrev and next entry
    MOVD pFirst, ptr
    MOVBU (ptr), bFirst
    MOVBU.P 1(ptr), bPrev

loop1:
    // load current byte, compare with previous, if equal, add value to s1
    MOVBU.P 1(ptr), bCur
    SUB $'0', bCur, tmp // tmp = bCur - '0'
    CMP bPrev, bCur
    CSEL EQ, tmp, ZR, tmp // tmp = (bPrev == bCur) ? tmp : 0
    ADD tmp, s1r, s1r // s1 += tmp
    MOVD bCur, bPrev // bPrev = bCur

    // loop until we hit the end
    CMP pEnd, ptr
    BLO loop1

checkwrap:
    SUB $'0', bPrev, tmp // tmp = bPrev - '0'
    CMP bPrev, bFirst
    CSEL EQ, tmp, ZR, tmp
    ADD tmp, s1r, s1r


done:
    // ret
    MOVD s1r, s1+24(FP)
    MOVD s2r, s2+32(FP)
    RET

