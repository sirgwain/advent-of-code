//go:build arm64 && !purego

#include "textflag.h"
// func solve(input []byte) (s1r int, s2r int)
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
#define ptr         R0  // input pointer
#define n           R1  // remaining input bytes to parse
#define x           R2  // current x loc
#define y           R3  // current y loc
#define hx          R4  // x heading
#define hy          R5  // y heading
#define s1r         R6  // s1 accumulator
#define s2r         R7  // s2 accumulator
#define b           R8  // current byte from input
#define digit       R9  // digit from *input for itoa
#define dist        R10 // dist travelled from R*/L* input
#define gridptr     R11 // pointer into the grid for tracking movement
#define found       R12 // 1 if we found shortest overlap
#define stepX       R13 
#define stepY       R14
#define stride      R15
#define gridLoc     R19
#define word        R20
#define mask        R21
#define j           R22
#define k1          R23 // constant 1 for bit shifting
#define k10         R24 // constant 10 value for MADD itoa calls
#define tmp         R16

// helper to ABS(reg) to set the reg to its absolute value
#define ABS(r) CMP $0, r; CNEG LT, r, r

TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), ptr // ptr = input
    MOVD input_len+8(FP), n    // n = len(input)

    // init counters
    MOVD ZR, x // x = 0
    MOVD ZR, y // y = 0
    MOVD ZR, hx // hx = 0
    MOVD $1, hy // hy = 1 - facing north
    MOVD ZR, s1r // s1r = 0
    MOVD ZR, s2r // s2r = 0
    MOVD $10, k10 // atoi 10x multiplier
    MOVD ZR, dist // atoi = 0
    MOVD ZR, found // found = 0

    MOVD $·grid(SB), gridptr
    MOVD gridptr, tmp
    MOVD $512, j
clear:
    // clear(grid[:])
    // store zero register twice, then increment gridptr by 16
    // do this 4 times to clear 64 bytes, then decrement the counter and loop
    STP.P (ZR, ZR), 16(tmp)
    STP.P (ZR, ZR), 16(tmp)
    STP.P (ZR, ZR), 16(tmp)
    STP.P (ZR, ZR), 16(tmp)
    SUBS $1, j, j
    BNE clear

    // setup starting grid
    // we want 
    // gridLoc := 256*512 + 256
    // grid[gridLoc>>6] = 1
    // but since gridLoc starts as a constant, we just want grid[2052] = 1
    // because (256*512 + 256)/64 = 2052
    // we still need to keep gridLoc initialized for stepping later
    MOVD $131328, gridLoc // gridLoc = 256*512 + 256
    LSR $6, gridLoc, word // word = gridLoc>>6
    MOVD $1, k1 // k1 = 1
    MOVD k1, ·grid+16416(SB)    // grid[2052] = 1

loop:
    CBZ n, done // compare-branch-zero if n is down to 0

    // Read a byte
    MOVBU.P 1(ptr), b // b = *input++

    SUB $1, n, n // n--
 
    CMP $'R', b
    BEQ right

    CMP $'L', b
    BEQ left

    CMP $',', b
    BEQ process

    // check for 0-9 chars, skip anything else
    SUB $'0', b, digit // digit = *input - '0'
    CMP $9, digit
    BHI loop // if digit > 9 goto loop

    MADD k10, digit, dist, dist  // dist = dist*10 + digit
    
    B loop

right:
    // rotate right, hx, hy = hy, -hx
    NEG hx, digit // digit = -hx
    MOVD hy, hx // hx = hy
    MOVD digit, hy // hy = digit
    B loop
left:
    // rotate left, hx, hy = -hy, hx
    NEG hy, digit // digit = -hy
    MOVD hx, hy // hy = hx
    MOVD digit, hx // hx = digit
    B loop
process:
    CBZ found, step // if !found goto step

processPostStep:

    MADD hx, x, dist, x   // x = dist*hx + x
    MADD hy, y, dist, y   // y = dist*hy + y
    MOVD ZR, dist // zero out the counter
    B loop

step:
    // step through a movement and record it in the grid
    MOVD x, stepX // stepX = x
    MOVD y, stepY // stepY = y

    ADD hy<<9, hx, stride // stride = hy*512 + hx
    MOVD dist, j // j = dist
    CBZ  j, processPostStep   // nothing to walk

steploop:
    ADD hx, stepX, stepX // stepX = stepX + hx
    ADD hy, stepY, stepY // stepY = stepY + hy
    ADD stride, gridLoc, gridLoc // gridLoc = gridLoc + stride

    LSR $6, gridLoc, word // word = gridLoc>>6
    LSL gridLoc, k1, mask // mask = 1 << gridLoc&63 (registers already drop the bits we AND away)

    MOVD (gridptr)(word<<3), tmp // tmp = grid[word]
    TST mask, tmp // if tmp & mask
    BNE foundOverlap
    ORR  mask, tmp, tmp        // tmp = tmp|mask
    MOVD tmp, (gridptr)(word<<3) // grid[word] = tmp

    SUBS $1, j, j // j = j - 1
    BNE steploop
    B processPostStep

foundOverlap: 

    // s2r = |stepX| + |stepY|
    ABS(stepX)
    ABS(stepY)
    ADD stepX, stepY, s2r // put s2r into s2r

    MOVD k1, found // found = 1
    B processPostStep

done:
    CBNZ dist, process    // last move had no ',' after it, so walk it now
    
    // ret s1r = |x| + |y|
    ABS(x)
    ABS(y)
    ADD x, y, s1r // put s1r into s1r

    MOVD s1r, s1+24(FP)
    MOVD s2r, s2+32(FP)
    RET
    