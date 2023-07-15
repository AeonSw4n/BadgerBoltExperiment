# BadgerDb and BoltDb Memory Experiment

This repo runs a simple experiment to test the performance of BadgerDB and BoltDB. In the experiment, I test writing key-value pairs to the DB in batches of various sizes, totalling 5GB of written data. For each write batch, I also added a bunch of read, delete and iterate operations to simulate regular DB usage.

Throughout the experiment, the program measures time and memory allocation. At the end of each test, the code prints the maximum and mean of the benchmarked stats. To reproduce results yourself, just run `go test -v .` Below are my results for each test configuration. 

In the Badger tests, there are various parameters and approaches that bear significant impact on the performance of the DB. For instance, we can write data to the DB by using a standard Update transaction, or by using a WriteBatch. Moreover, we can use the "Default" or "Performance" config for Badger. Many of these settings are benchmarked in the experiment.

### Conclusion
BadgerDB has a memory leak problem, where allocated memory instantaneously spikes to as much as 40-50GB. This can occur regardless of the size of the write batch. The memory leak is more frequent when using "Performance" configuration and when using the WriteBatch. The memory leak happens least frequently when using "Default" configuration and writing in batches of 10MB using an Update transaction. In this case, allocated memory doesn't exceed 8GB. It's worth noting that the memory leak appears randomly, and it's unclear what exactly triggers it. But it has occured for every test configuration at least once, even the "Default" config.

BoltDB is much more stable, and does not have a memory leak problem. The allocated memory never exceeded 1GB, and usually stays in lower 100MBs. However, it is a little slower than BadgerDB, about 1/2 slower. E.g. Bolt's 10MB batch test took 51 seconds, while Badger's 10MB batch test took 28 seconds. But writing in batches of 100MB is much faster for Bolt, taking only 28s, compared to Badger's 19s.

## Benchmark Results
### BoltDB - 5GB Total - 10MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (51.184827684999995)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 50 MiB 	TotalAlloc = 14131 MiB	Sys = 151 MiB	NumGC = 639
        MAX STATS 	|	 Alloc = 138 MiB	TotalAlloc = 28555 MiB	Sys = 153 MiB	NumGC = 1259
```
### BoltDB - 5GB Total - 25MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (36.55428825600001)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 126 MiB	TotalAlloc = 41848 MiB	Sys = 281 MiB	NumGC = 1485
        MAX STATS 	|	 Alloc = 264 MiB	TotalAlloc = 55207 MiB	Sys = 284 MiB	NumGC = 1715
```
### BoltDB - 5GB Total - 100MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (28.230409374000008)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 486 MiB	TotalAlloc = 68111 MiB	Sys = 912 MiB	NumGC = 1777
        MAX STATS 	|	 Alloc = 579 MiB	TotalAlloc = 80440 MiB	Sys = 972 MiB	NumGC = 1830
```
### BadgerDB - Default Config - 5GB Total - 10MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (28.125933275999987)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 2737 MiB	TotalAlloc = 119058 MiB	Sys = 4724 MiB	NumGC = 1862
        MAX STATS 	|	 Alloc = 5958 MiB	TotalAlloc = 165680 MiB	Sys = 6988 MiB	NumGC = 1885
```
### BadgerDB - Performance Config - 5GB Total - 10MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (18.586051209999997)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 11930 MiB	TotalAlloc = 195696 MiB	Sys = 15550 MiB	NumGC = 1893
        MAX STATS 	|	 Alloc = 27450 MiB	TotalAlloc = 217015 MiB	Sys = 27940 MiB	NumGC = 1896
```
### BadgerDB - Performance Config - 5GB Total - 25MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (18.934279655)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 13270 MiB	TotalAlloc = 240102 MiB	Sys = 33152 MiB	NumGC = 1898
        MAX STATS 	|	 Alloc = 24514 MiB	TotalAlloc = 257357 MiB	Sys = 33152 MiB	NumGC = 1901
```
### BadgerDB - Performance Config - 5GB Total - 100MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (18.974326708)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 10448 MiB	TotalAlloc = 276811 MiB	Sys = 38258 MiB	NumGC = 1903
        MAX STATS 	|	 Alloc = 23019 MiB	TotalAlloc = 292394 MiB	Sys = 38258 MiB	NumGC = 1905
```
### BadgerDB - Default Config - WriteBatch - 5GB Total - 10MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (34.115560625000015)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 4891 MiB	TotalAlloc = 337949 MiB	Sys = 38258 MiB	NumGC = 1927
        MAX STATS 	|	 Alloc = 32487 MiB	TotalAlloc = 387991 MiB	Sys = 38259 MiB	NumGC = 1946
```
### BadgerDB - Default Config - WriteBatch - 5GB Total - 25MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (26.763459406)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 3571 MiB	TotalAlloc = 423994 MiB	Sys = 38259 MiB	NumGC = 1963
        MAX STATS 	|	 Alloc = 8320 MiB	TotalAlloc = 463785 MiB	Sys = 38259 MiB	NumGC = 1980
```
### BadgerDB - Default Config - WriteBatch - 5GB Total - 100MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (26.595065420999997)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 3231 MiB	TotalAlloc = 503746 MiB	Sys = 38259 MiB	NumGC = 2002
        MAX STATS 	|	 Alloc = 8405 MiB	TotalAlloc = 545919 MiB	Sys = 38259 MiB	NumGC = 2020
```
### BadgerDB - Performance Config - WriteBatch - 5GB Total - 10MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (18.961160376000002)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 11971 MiB	TotalAlloc = 574700 MiB	Sys = 38259 MiB	NumGC = 2027
        MAX STATS 	|	 Alloc = 26908 MiB	TotalAlloc = 595849 MiB	Sys = 38259 MiB	NumGC = 2030
```
### BadgerDB - Performance Config - WriteBatch - 5GB Total - 25MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (25.305671966000013)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 22099 MiB	TotalAlloc = 618571 MiB	Sys = 43088 MiB	NumGC = 2031
        MAX STATS 	|	 Alloc = 42619 MiB	TotalAlloc = 636188 MiB	Sys = 43381 MiB	NumGC = 2033
```
### BadgerDB - Performance Config - WriteBatch - 5GB Total - 100MB Batches
```
    main_test.go:428: Timer results:
        Timer.End: event (Experiment) total elapsed time (16.707491454)
    main_test.go:429: Profiler results:
        MEAN STATS 	|	 Alloc = 18434 MiB	TotalAlloc = 655966 MiB	Sys = 43381 MiB	NumGC = 2034
        MAX STATS 	|	 Alloc = 23153 MiB	TotalAlloc = 670976 MiB	Sys = 43382 MiB	NumGC = 2036
```
