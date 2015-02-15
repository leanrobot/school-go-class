#!/bin/bash

for i in {1..20}; do
	#curl localhost:8080 &
	curl localhost:8080/time --write-out "%{http_code}\n" --silent --output /dev/null &
done

for job in `jobs -p`
do
    wait $job
done