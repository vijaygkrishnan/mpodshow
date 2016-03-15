package mpodshow

import (
          "encoding/json"
          "fmt"
       )

type BgpPeerT struct {
    Addr      string   `json:"addr"`
    AdminSt   string   `json:"adminSt"`
    Asn       string   `json:"asn"`
    SrcIf     string   `json:"srcIf"`
    Type      string   `json:"type"`
}

func (p BgpPeerT) Print() {
    fmt.Printf(" addr:%s\t asn:%s\t type:%s\n",p.Addr,p.Asn,p.Type) 
}

func BgpShowOutput (str string) error {
    var data []string
    var err error
    data,err = GetClassAttributes("error", str)
    if err != nil {
        if len(data) >= 1 {
            println(data[0])
        }
        return err
    }    
    var peer BgpPeerT
    data,err = GetClassAttributes("bgpPeer", str)
    if err != nil {
        return err
    }    
    for _,peerData := range(data) {
        err = json.Unmarshal([]byte(peerData),&peer)
        if err != nil {
            return err
        }    
        peer.Print()
    }
    return nil
}

