#!/bin/bash

___modules() {
  modules=(
    "guid"
    "ctxlog"
    "config"
    "envload"
    "logger"
    "lumberjack"
    "rotatelogs")
  echo "${modules[@]}"
}

# 删除本地及仓库上指定的 tag 版本，主要用于用户打的版本号出现了错误进行清理
___del_spec_tag() {
  local version="$1"
  if ! git tag -d "$version"; then
    log_error "git tag delete $version error"
  fi

  if ! git push origin :refs/tags/"$version"; then
    log_error "delete remote tag :refs/tags/$version error"
  fi

  for tag in $(___modules); do
    log_info "should delete tag ${tag}/${version}"
    if ! git tag -d "${tag}/${version}"; then
      log_error "delete tag ${tag}/${version} error"
    fi
    log_info "deleted local tag ${tag}/${version}"
    if ! git push origin :refs/tags/"${tag}/${version}"; then
      log_error "delete remote tag :refs/tags/${tag}/${version} error"
    fi
    log_info "delete romote tag ${tag}/${version}"
  done
}

# 推送新的 tag 版本到 git 仓库
___push_spec_tag() {
  local version="$1"
  local msg="$2"
  for tag in $(___modules); do
    git tag -a "${tag}/${version}" -m "${tag}/${version} $msg"
  done

  git tag -a "$version" -m "$version $msg"

  git push --tags
}

# 设置 sdk go.mod 中依赖的版本为新版本
___replace_spec_tag() {
  local version newVersion
  version="$1"
  newVersion="$2"
  repo="github.com/pemako/gopkg"
  for tag in $(___modules); do
    _sp="$repo/$tag $version"
    _np="$repo/$tag $newVersion"
    while read -r line; do
      log_info "$line" "$version" "$newVersion"
      sed -i "s|${_sp}|${_np}|g" "$line"
    done < <(grep "${_sp}" -rl * | grep -v 'go.sum')
  done

  # 替换项目根目录下的 go.mod
  sed -i "s|$repo $version|$repo $newVersion|g" go.mod
  # 替换 verion.go 中的上报上报版本
  # sed -i "s|$version|$newVersion|g" version.go
}

___go_mod_tidy() {
  for tag in $(___modules); do
    pushd "$tag" >/dev/null || return 1
    rm -rf go.sum
    go mod tidy
    popd >/dev/null || return 1
  done
}
