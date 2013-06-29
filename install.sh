#!/bin/bash

dir=$(dirname `readlink -f $0`)
go=$(which go &>/dev/null)

if [ -z "$go"]; then
  if [ -f /usr/local/go/bin/go ]; then
    go=/usr/local/go/bin/go
  else
   echo cant find go
   exit 1
  fi
fi
  
for pkg in $(cat $dir/Packages); do
  $go get $pkg && $go install $pkg
  if [ $? -ne 0 ]; then
    echo error while installing packages
    exit 1
  fi 
done

$go build scridx
