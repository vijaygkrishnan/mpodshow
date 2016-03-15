#mpodshow 
tool to fetch and decode route-maps from ACI switches 

#example
mpod.go:

```
package main
import "github.com/vijaygkrishnan/mpodshow"
mpodshow.RunCli()
```

go build mpod.go

#usage
mpod swmp5-spine1 showrun rpm

output:
```
[isisInterLeakP: interleak-ospf-default-1] InterLeak from ospf to isis
route-map interleak_rtmap_prefix_remote_pod_teps permit 2
   match ip prefix-list prefix_all_ifcs_tep_range
      ip prefix-list prefix_all_ifcs_tep_range seq 1 permit 10.0.0.0/27 from 32 to 32
route-map interleak_rtmap_prefix_remote_pod_teps permit 1
   match ip prefix-list prefix_remote_pod_teps
      ip prefix-list prefix_remote_pod_teps seq 2 permit 10.20.58.2/32
      ip prefix-list prefix_remote_pod_teps seq 1 permit 10.10.58.5
      ip prefix-list prefix_remote_pod_teps seq 3 permit 10.1.0.0/16
   match ip prefix-list prefix_ipn_remote_subnets
      ip prefix-list prefix_ipn_remote_subnets seq 2 permit 20.0.0.0/8 from 0 to 32
      ip prefix-list prefix_ipn_remote_subnets seq 1 permit 50.3.50.1/32
   set metric 100
```

