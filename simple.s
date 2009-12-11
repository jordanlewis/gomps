ori $t0,$zero,0
ori $t1,$zero,1
ori $t2,$zero,100
LOOP:
add $t0,$t0,$t1
sw  $t0,0($zero)
blt $t0,$t2,LOOP

lw  $t3,0($zero)
