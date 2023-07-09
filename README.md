# BadgerBoltExperiment

Simple experiment to test the performance of BadgerDB and BoltDB. In the experiment, we insert key-value pairs in each DB in batches of 25MB, totalling 1GB of random data.

Throughout, we measure the time taken and memory allocation. To reproduce results yourself, just run `go test -v .`

## Results
### BoltDB
```
    main_test.go:66: Timer results:
        Timer.End: event (Experiment) total elapsed time (7.118351249000001)
    main_test.go:67: Profiler results:
        MEAN STATS 	|	 Alloc = 123 MiB	TotalAlloc = 2663 MiB	Sys = 255 MiB	NumGC = 48
        MAX STATS 	|	 Alloc = 159 MiB	TotalAlloc = 5159 MiB	Sys = 287 MiB	NumGC = 88
```
### BadgerDB
```
    main_test.go:66: Timer results:
        Timer.End: event (Experiment) total elapsed time (5.817741123000003)
    main_test.go:67: Profiler results:
        MEAN STATS 	|	 Alloc = 5754 MiB	TotalAlloc = 16935 MiB	Sys = 7211 MiB	NumGC = 90
        MAX STATS 	|	 Alloc = 7689 MiB	TotalAlloc = 30419 MiB	Sys = 8196 MiB	NumGC = 94
```