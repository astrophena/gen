#!/usr/bin/env zsh

_gen() {
  local -a opts
  opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")

  _describe 'values' opts

  return
}

compdef _gen gen