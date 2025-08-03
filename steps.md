


1. /system/bin/pm list packages  --user 0   | rg "$userinput"  | sed 's/package://'

2. /system/bin/pm resolve-activity --brief  --user 0 $fstep | sed ' 1d' | sed 's/ //'


3./system/bin/am  start --user 0  -n  $sstep



