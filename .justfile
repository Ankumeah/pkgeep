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

  just format_spec

format_spec:
  #! /bin/sh

  set -eu

  cd specs/

  for file in *.json; do
    echo "${file}"

    jq . "${file}" > "${file}.tmp"
    mv "${file}.tmp" "${file}"
  done
