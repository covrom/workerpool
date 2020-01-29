#!/bin/bash

go build .

./workerpool -pool -c 300 -head > res.csv
./workerpool -pool -c 3000 >> res.csv
./workerpool -pool -c 30000 >> res.csv
./workerpool -pool -c 300000 >> res.csv
./workerpool -pool -c 1000000 >> res.csv
./workerpool -pool -c 3000000 >> res.csv
./workerpool -pool -c 6000000 >> res.csv
./workerpool -pool -c 8000000 >> res.csv

./workerpool -go -c 300 >> res.csv
./workerpool -go -c 3000 >> res.csv
./workerpool -go -c 30000 >> res.csv
./workerpool -go -c 300000 >> res.csv
./workerpool -go -c 1000000 >> res.csv
./workerpool -go -c 3000000 >> res.csv
./workerpool -go -c 6000000 >> res.csv
./workerpool -go -c 8000000 >> res.csv

./workerpool -fullpool -c 300 >> res.csv
./workerpool -fullpool -c 3000 >> res.csv
./workerpool -fullpool -c 30000 >> res.csv
./workerpool -fullpool -c 300000 >> res.csv
./workerpool -fullpool -c 1000000 >> res.csv
./workerpool -fullpool -c 3000000 >> res.csv
./workerpool -fullpool -c 6000000 >> res.csv
./workerpool -fullpool -c 8000000 >> res.csv

./workerpool -dgo -c 300 >> res.csv
./workerpool -dgo -c 3000 >> res.csv
./workerpool -dgo -c 30000 >> res.csv
./workerpool -dgo -c 300000 >> res.csv
./workerpool -dgo -c 1000000 >> res.csv
./workerpool -dgo -c 3000000 >> res.csv
./workerpool -dgo -c 6000000 >> res.csv
./workerpool -dgo -c 8000000 >> res.csv

./workerpool -fastpool -c 300 >> res.csv
./workerpool -fastpool -c 3000 >> res.csv
./workerpool -fastpool -c 30000 >> res.csv
./workerpool -fastpool -c 300000 >> res.csv
./workerpool -fastpool -c 1000000 >> res.csv
./workerpool -fastpool -c 3000000 >> res.csv
./workerpool -fastpool -c 6000000 >> res.csv
./workerpool -fastpool -c 8000000 >> res.csv

./workerpool -fgo -c 300 >> res.csv
./workerpool -fgo -c 3000 >> res.csv
./workerpool -fgo -c 30000 >> res.csv
./workerpool -fgo -c 300000 >> res.csv
./workerpool -fgo -c 1000000 >> res.csv
./workerpool -fgo -c 3000000 >> res.csv
./workerpool -fgo -c 6000000 >> res.csv
./workerpool -fgo -c 8000000 >> res.csv
