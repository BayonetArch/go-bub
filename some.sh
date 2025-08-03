#!/usr/bin/env sh
sm=$1
am start --user 0 -n "$(
  /system/bin/pm resolve-activity --brief --user 0 "$(
    /system/bin/pm list packages --user 0 | rg "$sm" | sed 's/package://'
  )" | sed '1d; s/ //'
)"
