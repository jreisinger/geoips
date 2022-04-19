geoips prints geographic locations of IP addresses. The output is sorted by
country and then by city.

```
$ ./geoips 1.1.1.1 4.4.8.8 8.8.8.8
Australia      Sydney         1.1.1.1
United States  Mountain View  8.8.8.8
United States  Nashville      4.4.8.8
```

```
go build
cat ips.txt | ./geoips
```
