default: run

run *args:
  go build -o build/out .
  ./build/out {{ args }}

ci ci = "":
  #! /bin/sh

  if test -n '{{ ci }}'; then
    echo 'CI = true'
    export CI=true
  fi

  ./.ci.sh
