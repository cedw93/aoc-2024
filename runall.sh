#!/bin/bash

for filename in d[0-9]*; do
    cd "$filename"
    go run main.go < input.txt
    cd ..
    exit
done