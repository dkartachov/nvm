#!/bin/bash

function nvm() {
  function trim_path() {
    # Delete path by parts so we can never accidentally remove sub paths
    if [ "$PATH" == "$1" ]; then 
      PATH=""
    fi

    PATH=${PATH//":$1:"/":"} # delete any instances in the middle
    PATH=${PATH/#"$1:"/} # delete any instance at the beginning
    PATH=${PATH/%":$1"/} # delete any instance in the at the end
  }

  if [ -z "$NVM_HOME" ]; then
    export NVM_HOME="$HOME/.nvm"
    export NVM_NODE="$NVM_HOME/node_versions"

    CURRENT_NODE="$(cat $NVM_HOME/current.txt)"

    [[ -z "$CURRENT_NODE" ]] && return

    PATH="$HOME/.nvm/node_versions/$CURRENT_NODE/bin:$PATH"

    return
  fi

  go run main.go "$@"

  if [ -n "$CURRENT_NODE" ] && [ $? -eq 0 ]; then
    trim_path "$HOME/.nvm/node_versions/$CURRENT_NODE/bin"

    CURRENT_NODE="$(cat $NVM_HOME/current.txt)"
    PATH="$HOME/.nvm/node_versions/$CURRENT_NODE/bin:$PATH"
  fi
}