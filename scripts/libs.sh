#!/bin/bash

COLOR_RED='\033[0;31m'
COLOR_ORANGE='\033[0;33m'
COLOR_GREEN='\033[0;32m'
COLOR_LIGHTCYAN='\033[0;36m'
COLOR_BLUE='\033[0;94m'
COLOR_MAGENTA='\033[95m'
COLOR_BOLD='\033[1m'
COLOR_NONE='\033[0m' # No Color

function log_error {
  >&2 echo -n -e "${COLOR_BOLD}${COLOR_RED}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_warning {
  >&2 echo -n -e "${COLOR_ORANGE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_callout {
  >&2 echo -n -e "${COLOR_LIGHTCYAN}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_cmd {
  >&2 echo -n -e "${COLOR_BLUE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_success {
  >&2 echo -n -e "${COLOR_GREEN}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_info {
  >&2 echo -n -e "${COLOR_NONE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

if [[ $(uname) == 'Darwin' ]]; then
  sed() {
    gsed "$@"
  }
fi

# 比较版本号的函数 v1 v2 比较结果如下
# v1 = v2 return 0
# v1 > v2 return 1
# v1 < v2 return 2
function compare_versions() {
  local v1=$1
  local v2=$2

  # 保存原始的 IFS 值
  local old_ifs=$IFS
  # 使用 . 分隔版本号的各个部分，并保存到数组中
  IFS='.' read -r -a v1_arr <<<"$v1"
  IFS='.' read -r -a v2_arr <<<"$v2"

  # 恢复原始的 IFS 值
  IFS=$old_ifs

  local eq=0
  local v1_gt_v2=1
  local v1_lt_v2=2

  # 使用数组比较版本号的各个部分大小
  if ((${v1_arr[0]} != ${v2_arr[0]})); then
    [[ ${v1_arr[0]} -lt ${v2_arr[0]} ]] && echo "${v1_lt_v2}" || echo "${v1_gt_v2}"
  elif ((${v1_arr[1]} != ${v2_arr[1]})); then
    [[ ${v1_arr[1]} -lt ${v2_arr[1]} ]] && echo "${v1_lt_v2}" || echo "${v1_gt_v2}"
  elif ((${v1_arr[2]} != ${v2_arr[2]})); then
    [[ ${v1_arr[2]} -lt ${v2_arr[2]} ]] && echo "${v1_lt_v2}" || echo "${v1_gt_v2}"
  else
    echo "$eq"
  fi
}

# 检查需要用到 命令是否存在
___check_command() {
  local level=0
  if [[ $(uname) == 'Darwin' ]]; then
    if ! command -v gsed &>/dev/null; then
      log_error "gsed command not exist on mac os, you can install use brew install gnu-sed"
      level=1
    fi
  else
    if ! command -v sed &>/dev/null; then
      log_error "sed command not exist."
      level=1
    fi
  fi

  if ! command -v awk &>/dev/null; then
    log_error "awk command not exist."
    level=1
  fi

  # 检查 git 命令
  if ! command -v git &>/dev/null; then
    log_error "git command not exist."
    level=1
  fi

  # 检查 golang 的版本
  if ! command -v go &>/dev/null; then
    log_error "go command not exist."
    level=1
  else # 存在的情况下检查 go 的版本
    modVersion=$(grep '^go [0-9]*' go.mod | awk '{print $NF}')
    goVersion=$(go version | awk '{print $3}' | awk -F'go' '{print $NF}')
    if (($(compare_versions "$modVersion" "$goVersion") == 1)); then
      log_error "go version=${goVersion} lt go.mod use go version ${modVersion}, please update go"
      level=1
    fi
  fi

  if (("$level" == 0)); then
    log_info "Congratulations, your environment ok."
  fi
  echo "$level"
}
