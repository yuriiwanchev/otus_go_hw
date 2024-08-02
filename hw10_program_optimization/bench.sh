go test -benchmem -run=^$ -count=10 -bench ^BenchmarkGetDomainStat$ > old.txt
go test -benchmem -run=^$ -count=10 -bench ^BenchmarkGetDomainStat$ > new.txt
benchstat old.txt new.txt > benchstat.txt