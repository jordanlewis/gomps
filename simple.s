.data
foo: .word 1,2,3,4,5,6,7,8,9
bar: .word 11,12,13,14,15,16,17,18,19

.text

ori $t0,$zero,0
ori $t1,$zero,1
ori $t2,$zero,100
LOOP:
add $t0,$t0,$t1
sw  $t0,0($zero)
blt $t0,$t2,LOOP

lw  $t3,0($zero)
mult $t3,$t3
mflo $t4
addi $t4,$t4,1
la $t5,bar
ori $t1,$zero,3
ori $t4,$zero,4
mult $t1,$t4
mflo $t6
add $t6,$t6,$t5

lw $t7,0($t6)
