#!/bin/bash

function nvm() {
  if [ $1 = "init" ]
  then
    CURRENT_NODE="$(cat $HOME/.nvm/current.txt)"

    [[ -z "$CURRENT_NODE" ]] && return

    PATH="$HOME/.nvm/versions/node/$CURRENT_NODE/bin:$PATH"

    return
  fi

  OUTPUT=$(go run main.go "$@")

  if [ $? -eq 0 ]
  then
    NEW_PATH=$(echo $OUTPUT | grep "PATH")

    if [ ! -z "$NEW_PATH" ]
    then
      eval $NEW_PATH
    else
      printf "$OUTPUT"
    fi
  fi
}