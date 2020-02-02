# workerpool
benchmark: workerpool vs goroutines with params

```
go run ./main.go run --workers=20 --chan=200 -a 300 -a 3000 -a 30000 -a 300000 -a 1000000 -a 3000000 -a 6000000 -a 1000000 -a 20000000 -a 30000000
go run ./main.go chart
```
