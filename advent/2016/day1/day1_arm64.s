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
// R2 x
// R3 y
// R4 hx
// R5 hy
// R6 s1
// R7 s2
// R8 b byte (input value)
// R9 digit (for itoa)
// R10 10 constant for itoa
// R11 dist (for itoa)
// R12 grid
// R13 512
// R14 found
// R15 stepX
// R16 stepY
// R17 stride
// R19 gridLoc
// R20 word i.e. gridLoc>>6
// R21 1 constant for grid setting
// R22 mask
// R23 w for testing
// R24 j
//
// Note: R18 is reserved, and R16,R17 are scratch registers, so CALL could change them

TEXT ·solve(SB), NOSPLIT, $0-40
    // load the input slice
    MOVD input_base+0(FP), R0    // pointer into R0
    MOVD input_len+8(FP), R1    // length into R1

    // init counters
    MOVD ZR, R2 // x = 0
    MOVD ZR, R3 // y = 0
    MOVD ZR, R4 // hx = 0
    MOVD $1, R5 // hy = 1 - facing north
    MOVD ZR, R6 // s1 = 0
    MOVD ZR, R7 // s2 = 0
    MOVD $10, R10 // atoi 10x multiplier
    MOVD ZR, R11 // atoi = 0
    MOVD ZR, R14 // found = 0

    MOVD $·grid(SB), R12
    MOVD $512, R13
clear:
    // clear(grid[:])
    // store zero register twice, then increment R12 by 16
    // do this 4 times to clear 64 bytes, then decrement the counter and loop
    STP.P (ZR, ZR), 16(R12)
    STP.P (ZR, ZR), 16(R12)
    STP.P (ZR, ZR), 16(R12)
    STP.P (ZR, ZR), 16(R12)
    SUB $1, R13
    CBNZ R13, clear // keep going until we have looped 512 times to clear the whole grid

    // reset the 512 constant, we reuse it while calculating stride
    MOVD $512, R13
    MOVD $·grid(SB), R12  // clear loop advanced R12, reset it

    // setup starting grid
    // we want 
    // gridLoc := 256*512 + 256
    // grid[gridLoc>>6] = 1
    // but since gridLoc starts as a constant, we just want grid[2052] = 1
    // because (256*512 + 256)/64 = 2052
    // we still need to keep gridLoc initialized for stepping later
    MOVD $131328, R19 // gridLoc = 256*512 + 256
    LSR $6, R19, R20 // R20 = gridLoc>>6
    MOVD $1, R21 // R21 = 1
    MOVD R21, ·grid+16416(SB)    // grid[2052] = 1

loop:
    CBZ R1, done // compare-branch-zero if R1 is down to 0

    // Read a byte
    MOVBU.P 1(R0), R8 // b = *input++

    SUB $1, R1, R1 // R1--
 
    CMP $'R', R8
    BEQ right

    CMP $'L', R8
    BEQ left

    CMP $',', R8
    BEQ process

    // check for 0-9 chars, skip anything else
    SUB $'0', R8, R9 // digit = *input - '0'
    CMP $9, R9
    BHI loop // if digit > 9 goto loop

    MADD R10, R9, R11, R11  // dist = dist*10 + digit
    
    B loop

right:
    // rotate right, hx, hy = hy, -hx
    NEG R4, R9 // R9 = -hx
    MOVD R5, R4 // hx = hy
    MOVD R9, R5 // hy = R9
    B loop
left:
    // rotate left, hx, hy = -hy, hx
    NEG R5, R9 // R9 = -hy
    MOVD R4, R5 // hy = hx
    MOVD R9, R4 // hx = R9
    B loop
process:
    CMP $0, R14
    BEQ step // if !found goto step

processPostStep:

    MADD R4, R2, R11, R2   // x = dist*hx + x
    MADD R5, R3, R11, R3   // y = dist*hy + y
    MOVD ZR, R11 // zero out the counter
    B loop

step:
    // step through a movement and record it in the grid
    MOVD R2, R15 // stepX = x
    MOVD R3, R16 // stepY = y

    MADD R13, R4, R5, R17 // stride = hy*512 + hx

    MOVD $0, R24 // j = 0

steploop:
    // steploop loops through each value in dist
    CMP R11, R24
    BGE processPostStep // if j >= dist goto processPostStep

    ADD R4, R15, R15 // stepX = stepX + hx
    ADD R5, R16, R16 // stepY = stepY + hy
    ADD R17, R19, R19 // gridLoc = gridLoc + stride

    LSR $6, R19, R20 // word = gridLoc>>6
    LSL R19, R21, R22 // mask = 1 << gridLoc&63 (registers already drop the bits we AND away)

    MOVD (R12)(R20<<3), R23 // w = grid[word]
    TST R22, R23 // if w & mask
    BNE found
    ORR  R22, R23, R23        // w = w|mask
    MOVD R23, (R12)(R20<<3) // grid[word] = w

    ADD $1, R24, R24 // j = j + 1
    B steploop

found: 

    // s2 = |stepX| + |stepY|
    CMP $0, R15
    CNEG LT, R15, R15 // x = |x|
    CMP $0, R16
    CNEG LT, R16, R16 // y = |y|
    ADD R15, R16, R7 // put s2 into R7

    MOVD R21, R14 // found = 1
    B processPostStep

done:
    CBNZ R11, process    // last move had no ',' after it, so walk it now
    
    // ret s1 = |x| + |y|
    CMP $0, R2
    CNEG LT, R2, R2 // x = |x|
    CMP $0, R3
    CNEG LT, R3, R3 // y = |y|
    ADD R2, R3, R6 // put s1 into R6

    MOVD R6, s1+24(FP)
    MOVD R7, s2+32(FP)
    RET

