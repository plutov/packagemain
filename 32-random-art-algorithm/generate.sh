#!/bin/bash

for i in {1..100}; do
	random_string=$(LC_ALL=C.UTF-8 tr -dc '[:alnum:]' </dev/urandom | head -c 16)
	go run *.go --depth 5 --phrase "$random_string" --out "out-$i.png"
done

echo "Finished generating 100 images."
