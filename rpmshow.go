package mpodshow

import (
          "encoding/json"
          "fmt"
       )

/*
 * Used http://mholt.github.io/json-to-go/ to generate the go structure
 */
type RibLeakResponse struct {
    TotalCount string `json:"totalCount"`
    Imdata []struct {
        OspfRibLeakP struct {
            Attributes struct {
                Always string `json:"always"`
                RtMap string  `json:"rtMap"`
            } `json:"attributes"`
        } `json:"ospfRibLeakP"`
    } `json:"imdata"`
}
// RpmShowInterLeakP:
//   Decode and print inter-leak info
func RpmShowRibLeakP (s *Session, to_proto string, recursive bool) error {
    var url string
    if to_proto == "ospf" {
        url = s.url_prefix + "class/ospfRibLeakP.json?rsp-subtree=full"
    } else {
        fmt.Println("Only OSPF ribLeak is supported")
        return nil
    }
    json_str,err := s.GetMo(url)
    if err != nil {
        return err
    }
    var data []string
    data,err = GetClassAttributes("error", json_str)
    if err != nil {
        if len(data) >= 1 {
            println(data[0])
        }
        return err
    }    
    var resp RibLeakResponse
    err = json.Unmarshal([]byte(json_str), &resp)
    if err != nil { return err }    

    for _,d := range(resp.Imdata) {
        fmt.Printf("ospfRibLeakP\n")
        fmt.Printf("   rtMap : %s\n", d.OspfRibLeakP.Attributes.RtMap)
        fmt.Printf("   always: %s\n", d.OspfRibLeakP.Attributes.Always)
        if recursive {
            err = RpmShowRouteMap(s, d.OspfRibLeakP.Attributes.RtMap, recursive)
            if err != nil { return nil }
            fmt.Println()
        }
    }
    return nil
}

type InterLeakPT struct {
    Attributes struct {
        Proto string  `json:"proto"`
        Inst string   `json:"inst"`
        Asn string    `json:"asn"`
        RtMap string  `json:"rtMap"`
        Scope string  `json:"scope"`
        Name string   `json:"name"`
    } `json:"attributes"`
} 
type InterLeakResponse struct {
    TotalCount string `json:"totalCount"`
    Imdata []struct {
        OspfInterLeakP  *InterLeakPT `json:"ospfInterLeakP"`
        IsisInterLeakP  *InterLeakPT `json:"isisInterLeakP"`
    } `json:"imdata"`
}
// RpmShowInterLeakP:
//   Decode and print inter-leak info
func RpmShowInterLeakP (s *Session, to_proto string, recursive bool) error {
    var url string
    if to_proto == "ospf" {
        url = s.url_prefix + "class/ospfInterLeakP.json?rsp-subtree=full"
    } else {
        url = s.url_prefix + "class/isisInterLeakP.json?rsp-subtree=full"
    }
    json_str,err := s.GetMo(url)
    if err != nil {
        return err
    }
    var data []string
    data,err = GetClassAttributes("error", json_str)
    if err != nil {
        if len(data) >= 1 {
            println(data[0])
        }
        return err
    }    
    var resp InterLeakResponse
    err = json.Unmarshal([]byte(json_str), &resp)
    if err != nil {
        return err
    }    
    for _,d := range(resp.Imdata) {
        var InterLeakP *InterLeakPT
        if to_proto == "ospf" {
            InterLeakP = d.OspfInterLeakP
        } else {
            InterLeakP = d.IsisInterLeakP
        }
        if InterLeakP == nil {
            return err
        }
        interLeakName := fmt.Sprintf("interleak-%s-%s-%s", 
                    InterLeakP.Attributes.Proto,
                    InterLeakP.Attributes.Inst,
                    InterLeakP.Attributes.Asn)
        fmt.Printf("[%sInterLeakP: %s] InterLeak from %s to %s\n",
                   to_proto, interLeakName,
                   InterLeakP.Attributes.Proto, to_proto)
        if recursive {
            err = RpmShowRouteMap(s, InterLeakP.Attributes.RtMap, recursive)
            if err != nil { return nil }
            fmt.Println()
        } else {
            fmt.Printf("   proto: %s\n", InterLeakP.Attributes.Proto)
            fmt.Printf("   inst : %s\n", InterLeakP.Attributes.Inst)
            fmt.Printf("   asn  : %s\n", InterLeakP.Attributes.Asn)
            fmt.Printf("   rtmap: %s\n", InterLeakP.Attributes.RtMap)
        }
    }
    return nil
}

type RtmapMatchRtTypeT struct {
    Attributes struct {
        RouteT string `json:"routeT"`
        Type   string `json:"type"`
    } `json:"attributes"`
}
type RtmapMatchRtDstT struct {
    Attributes struct {
        Type string `json:"type"`
    } `json:"attributes"`
    Children []struct {
        RtmapRsRtDstAtt struct {
            Attributes struct {
                ForceResolve string `json:"forceResolve"`
                State        string `json:"state"`
                TCl          string `json:"tCl"`
                TSKey        string `json:"tSKey"`
            } `json:"attributes"`
        } `json:"rtmapRsRtDstAtt"`
    } `json:"children"`
}
type RtmapSetNhT struct {
    Attributes struct {
        Addr string `json:"addr"`
        Ctrl string `json:"ctrl"`
        Type string `json:"type"`
    } `json:"attributes"`
}
type RtmapSetMetricT struct {
    Attributes struct {
        Metric string `json:"metric"`
        Type   string `json:"type"`
    } `json:"attributes"`
}
type RtmapSetPathSelectT struct {
    Attributes struct {
        Criteria string  `json:"criteria"`
        Status   string  `json:"status"`
        Type     string  `json:"type"`
    } `json:"attributes"`
}
type RtmapResponse struct {
    TotalCount string  `json:"totalCount"`
    Imdata []struct {
        RtmapRule struct {
            Attributes struct {
                Name string `json:"name"`
            } `json:"attributes"`
            Children []struct {
                RtmapEntry struct {
                    Attributes struct {
                        Action string `json:"action"`
                        Name   string `json:"name"`
                        Order  string `json:"order"`
                    } `json:"attributes"`
                    Children []struct {
                        RtmapMatchRtDst    *RtmapMatchRtDstT    `json:"rtmapMatchRtDst,omitempty"`
                        RtmapMatchRtType   *RtmapMatchRtTypeT   `json:"rtmapMatchRtType,omitempty"`
                        RtmapSetMetric     *RtmapSetMetricT     `json:"rtmapSetMetric,omitempty"`
                        RtmapSetNh         *RtmapSetNhT         `json:"rtmapSetNh,omitempty"`
                        RtmapSetPathSelect *RtmapSetPathSelectT `json:"rtmapSetPathSelect,omitempty"`
                    } `json:"children"`
                } `json:"rtmapEntry"`
            } `json:"children"`
        } `json:"rtmapRule"`
    } `json:"imdata"`
}

// RpmShowRouteMap:
//   Decode and print route-map info
func RpmShowRouteMap (s *Session, name string, recursive bool) error {
    var url string
    if name == "all" {
        url = s.url_prefix + "class/rtmapRule.json?rsp-subtree=full"
    } else {
        url = s.url_prefix + "mo/sys/rpm/rtmap-"+name+".json?rsp-subtree=full"
    }
    json_str,err := s.GetMo(url)
    if err != nil {
        return err
    }
    var data []string
    data,err = GetClassAttributes("error", json_str)
    if err != nil {
        if len(data) >= 1 {
            println(data[0])
        }
        return err
    }    
    var resp RtmapResponse
    err = json.Unmarshal([]byte(json_str), &resp)
    if err != nil {
        fmt.Println("ERROR:", err)
        return err
    }    
    for _,d := range(resp.Imdata) {
        for _,c := range(d.RtmapRule.Children) {
            //rtmapEntry
            fmt.Printf("route-map %s %s %s\n", d.RtmapRule.Attributes.Name, 
                       c.RtmapEntry.Attributes.Action, c.RtmapEntry.Attributes.Order)
            for _,ent := range(c.RtmapEntry.Children) {
                if ent.RtmapMatchRtDst != nil {
                    for _,dst := range(ent.RtmapMatchRtDst.Children) {
                        var pfx string
                        if dst.RtmapRsRtDstAtt.Attributes.State == "unformed" {
                            pfx = "unformed"
                        } else {
                            pfx = dst.RtmapRsRtDstAtt.Attributes.TSKey
                        }
                        fmt.Printf("   match ip prefix-list %s\n", pfx)
                        if recursive && pfx != "unformed" {
                            err = RpmShowPrefixList(s, pfx, recursive)
                            if err != nil { return nil }
                        }

                    }
                }
                if ent.RtmapMatchRtType != nil {
                    fmt.Printf("   match %s %s\n", 
                               ent.RtmapMatchRtType.Attributes.Type,
                               ent.RtmapMatchRtType.Attributes.RouteT)
                }
                if ent.RtmapSetMetric != nil {
                    fmt.Printf("   set %s %s\n", 
                               ent.RtmapSetMetric.Attributes.Type, 
                               ent.RtmapSetMetric.Attributes.Metric)
                }
                if ent.RtmapSetNh != nil {
                    fmt.Printf("   set ip next-hop %s %s\n", 
                               ent.RtmapSetNh.Attributes.Addr, 
                               ent.RtmapSetNh.Attributes.Ctrl)
                }
                if ent.RtmapSetPathSelect != nil {
                    fmt.Printf("   set path-selection %s\n",
                               ent.RtmapSetPathSelect.Attributes.Criteria)
                }
            }
        }
    }
    return nil
}

type RtpfxRtRtDstAttT  struct {
    Attributes struct {
        ParentSKey string `json:"parentSKey"`
        TCl string        `json:"tCl"`
        TDn string        `json:"tDn"`
    } `json:"attributes"`
}
type RtpfxEntryT struct {
    Attributes struct {
        Name string       `json:"name"`
        Action string     `json:"action"`
        Criteria string   `json:"criteria"`
        Order string      `json:"order"`
        Pfx string        `json:"pfx"`
        FromPfxLen string `json:"fromPfxLen"`
        ToPfxLen string   `json:"toPfxLen"`
    } `json:"attributes"`
}
type PfxListResponse struct {
    TotalCount string `json:"totalCount"`
    Imdata []struct {
        RtpfxRule struct {
            Attributes struct {
                Name string `json:"name"`
            } `json:"attributes"`
            Children []struct {
                RtpfxRtRtDstAtt *RtpfxRtRtDstAttT `json:"rtpfxRtRtDstAtt,omitempty"`
                RtpfxEntry      *RtpfxEntryT      `json:"rtpfxEntry,omitempty"`
            } `json:"children"`
        } `json:"rtpfxRule"`
    } `json:"imdata"`
}
// RpmShowPrefixList:
//   Decode and print prefix-list info
func RpmShowPrefixList (s *Session, name string, recursive bool) error {
    var url string
    if name == "all" {
        url = s.url_prefix + "class/rtpfxRule.json?rsp-subtree=full"
    } else {
        url = s.url_prefix + "mo/sys/rpm/pfxlist-"+name+".json?rsp-subtree=full"
    }
    json_str,err := s.GetMo(url)
    if err != nil {
        return err
    }
    var data []string
    data,err = GetClassAttributes("error", json_str)
    if err != nil {
        if len(data) >= 1 {
            println(data[0])
        }
        return err
    }    
    var resp PfxListResponse
    err = json.Unmarshal([]byte(json_str), &resp)
    if err != nil {
        return err
    }    
    for _,d := range(resp.Imdata) {
        for _,ent := range(d.RtpfxRule.Children) {
            if ent.RtpfxEntry != nil {
                if recursive { fmt.Printf("      ") }
                fmt.Printf("ip prefix-list %s seq %s %s %s ",
                           d.RtpfxRule.Attributes.Name,
                           ent.RtpfxEntry.Attributes.Order,
                           ent.RtpfxEntry.Attributes.Action,
                           ent.RtpfxEntry.Attributes.Pfx)
                if ent.RtpfxEntry.Attributes.Criteria == "inexact" {
                    fmt.Printf("from %s to %s\n",
                           ent.RtpfxEntry.Attributes.FromPfxLen,
                           ent.RtpfxEntry.Attributes.ToPfxLen)
                } else {
                    fmt.Printf("\n")
                }
            }
        }
    }
    return nil
}

