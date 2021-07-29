geoips prints geographic locations of IP addresses. The output is sorted by
country and then by city.

```
go build
./geoips 1.1.1.1 8.8.8.8
cat ips.txt | ./geoips
```
