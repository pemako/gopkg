#!/bin/bash

set -eou pipefail

# shellcheck disable=SC1091
source ./scripts/libs.sh
# shellcheck disable=SC1091
source ./scripts/mod.sh

# --mode 指定运行模式支持两种 check delete push replace tidy
# --version 为要操作的版本号
# --nextVersion 为要操作的新版本号
# --commitMsg 提交到 git 仓库的 commit 信息
while [ $# != 0 ]; do
  case "$1" in
  --mode | -m)
    mode=$2
    shift
    ;;
  --version | -v)
    version=$2
    shift
    ;;
  --nextVersion | -n)
    nextVersion=$2
    shift
    ;;
  --commitMsg | -c)
    commitMsg=$2
    shift
    ;;
  -*) ;;
  esac
  shift
done

main() {
  if [[ "$mode" != "tidy" && "$mode" != "check" ]]; then
    if [[ "$version" == "" ]]; then
      log_error "--version|-v is required"
      exit 0
    fi
  fi

  case "$mode" in
  check)
    ___check_command
    ;;
  tidy)
    ___go_mod_tidy
    ;;
  delete)
    ___del_spec_tag "$version"
    ;;
  push)
    ___push_spec_tag "$version" "$commitMsg"
    ;;
  replace)
    ___replace_spec_tag "$version" "$nextVersion"
    ;;
  *)
    echo "--mode|-m params is required"
    echo "you can use [check|tidy|delete|push|replace]"
    ;;
  esac
}

main "$@"
