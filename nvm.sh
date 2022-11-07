#!/bin/bash

function nvm() {
  function trim_path() {
    # Delete path by parts so we can never accidentally remove sub paths
    if [ "$PATH" == "$1" ] ; then PATH="" ; fi

    PATH=${PATH//":$1:"/":"} # delete any instances in the middle
    PATH=${PATH/#"$1:"/} # delete any instance at the beginning
    PATH=${PATH/%":$1"/} # delete any instance in the at the end
  }

  if [ $1 = "init" ]
  then
    export NVM_HOME="$HOME/.nvm"
    export NVM_NODE="$NVM_HOME/versions/node"

    CURRENT_NODE="$(cat $NVM_HOME/current.txt)"

    [[ -z "$CURRENT_NODE" ]] && return

    PATH="$HOME/.nvm/versions/node/$CURRENT_NODE/bin:$PATH"

    return
  fi

  go run main.go "$@"

  if [ $? -eq 0 ]
  then
    trim_path $
    NEW_VERSION="$(cat $NVM_HOME/current.txt)"

    echo "current: $CURRENT_NODE"
    echo "new: $NEW_VERSION"
  fi
}