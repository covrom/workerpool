#!/bin/bash

go build .
./workerpool -pool -c 300 -head > res.csv
./workerpool -pool -c 3000 >> res.csv
./workerpool -pool -c 30000 >> res.csv
./workerpool -pool -c 300000 >> res.csv
./workerpool -pool -c 3000000 >> res.csv
./workerpool -pool -c 30000000 >> res.csv
./workerpool -go -c 300 >> res.csv
./workerpool -go -c 3000 >> res.csv
./workerpool -go -c 30000 >> res.csv
./workerpool -go -c 300000 >> res.csv
./workerpool -go -c 3000000 >> res.csv
./workerpool -go -c 30000000 >> res.csv
