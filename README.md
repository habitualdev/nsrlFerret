# nsrlFerret

![image](ferret.png)

Easiest way to manage way too many NSRL entries
Self updates to most recent Modern NSRL database when run. Rebuilds bloom filters each time, currently takes
approximately 5 minutes to rebuild on a modern system.


Current API endpoints:


```
/stats/buckets - Number of "bucket" files
/stats/uptime - its in the name
/query?hash=<hash/filename> - Query by hash or filename
```