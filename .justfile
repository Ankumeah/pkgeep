default: run

run *args:
  go run . {{ args }}

ci ci = "":
  #! /bin/sh

  if test -n '{{ ci }}'; then
    echo 'CI = true'
    export CI=true
  fi


  ./.ci.sh
