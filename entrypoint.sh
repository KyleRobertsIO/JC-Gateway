#!/bin/bash

eval $(printenv | sed -n "s/^\([^=]\+\)=\(.*\)$/export \1=\2/p" | sed 's/"/\\\"/g' | sed '/=/s//="/' | sed 's/$/"/' >> /etc/profile)
echo 'Starting Job Container Manager Application'
go run ./main.go
echo 'Closing Job Container Manager Application'
