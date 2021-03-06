# Jordan Lewis
# Matrix multiplication assignment

    .data
    .align 2


m1: .word 1,2,3,4,5,6,7,8,9
m2: .word 11,12,13,14,15,16,17,18,19
m3: .word 0,0,0,0,0,0,0,0,0

    .text

main:
    # Constants
    ori $t3,$zero,3 # sizeof col
    ori $t4,$zero,4 # sizeof word

    # Initialize loop variables
    ori  $t0,$zero,0 # i
ILOOP:
    ori  $t1,$zero,0 # j
JLOOP:
    ori  $t2,$zero,0 # k
KLOOP:
    # Want to implement following calculation
    # m3[i][j] += m1[i][k] * m2[k][j]

    # Calculate m1 offset, store m1[i][k] in t7
    la $t5,m1
    mult $t0,$t3    # i * 3
    mflo $t6        # i * 3
    add $t6,$t6,$t2 # (i * 3) + k
    mult $t6,$t4    # (i * 3) + k) * 4
    mflo $t6        # (i * 3) + k) * 4
    add $t6,$t6,$t5 # t6 = m1 + ((i * 3) + k) * 4)
    lw $t7,0($t6)    # t7 = m1[i][k]

    # Calculate m2 offset, store m2[k][j] in t8
    la $t5,m2
    mult $t2,$t3    # k * 3
    mflo $t6        # k * 3
    add $t6,$t6,$t1 # (k * 3) + j
    mult $t6,$t4    # (k * 3) + j) * 4
    mflo $t6        # (k * 3) + j) * 4
    add $t6,$t6,$t5 # t6 = m2 + ((k * 3) + j) * 4)
    lw $t8,0($t6)    # t8 = m2[k][j]

    # Calculate m1[i][k] * m2[k][j], put it in t7
    mult $t7,$t8
    mflo $t7

    # Calculate m3 offset, put it in t6
    la $t5,m3
    mult $t0,$t3    # i * 3
    mflo $t6        # i * 3
    add $t6,$t6,$t1 # (i * 3) + j
    mult $t6,$t4    # (i * 3) + j) * 4
    mflo $t6        # (i * 3) + j) * 4
    add $t6,$t6,$t5 # t6 = m3 + ((i * 3) + j) * 4)  # t6 has the offset



    # Add m3[i][j] to m1[i][k] * m2[k][j], put it in t9
    lw $t8,0($t6)     # t8 = m3[i][j]
    add $t9,$t7,$t8

    # Put the result of the above back into m3[i][j]
    sw $t9,0($t6)


    # Done with calculation, now deal with jumps
    addi $t2,$t2,1 # k++
    blt  $t2,$t3,KLOOP   # go to kloop if k < 3
    addi $t1,$t1,1 # j++
    blt  $t1,$t3,JLOOP   # go to jloop if j < 3
    addi $t0,$t0,1 # i++
    blt  $t0,$t3,ILOOP   # go to iloop if i < 3

