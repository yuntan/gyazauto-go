#!/bin/bash

build () {
  os=$1
  arch=$2
  echo "==> for ${os} ${arch}"

  os_name=$os
  if [[ $os == "darwin" ]]; then
    os_name=macos
  fi

  arch_name=64bit
  if [[ $arch == "386" ]]; then
    arch_name=32bit
  fi

  bin=bin/gyazauto-$(git describe)-${os_name}-${arch_name}
  if [[ $os == "windows" ]]; then
    bin=$bin.exe
  fi

  ldflags="-X main.version=$(git describe)"
  if [[ $os == "windows" ]]; then
    ldflags="${ldflags} -H=windowsgui"
  fi

  pkg=github.com/yuntan/gyazauto-go/cmd/gyazauto

  GOOS=$os GOARCH=$arch go build -o $bin -ldflags "$ldflags" $pkg \
    && echo "--> generated $bin"
}

if [ $# == 2 ]; then
  build $1 $2
  exit 0
fi

for os in linux darwin windows; do
  build $os amd64
done

for os in linux windows; do
  build $os 386
done
