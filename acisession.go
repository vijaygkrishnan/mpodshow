package mpodshow

import (
          "net/http"
          "net/http/cookiejar"
          "bytes"
          "io/ioutil"
          "crypto/tls"
          "encoding/json"
       )

type Session struct {
    url_prefix  string
    jar         http.CookieJar
}

func NewSession(switch_ip string) (s *Session) {
    return &Session{ url_prefix : "https://"+switch_ip+"/api/" }
}

// Login:
//  POST aaaLogin MO and save the cookie in CookieJar
func (s *Session) Login() {
    url := s.url_prefix + "aaaLogin.json"
    Debug.Println("URL:>", url)

    var jsonStr = []byte(`{"aaaUser": {"attributes": {"pwd": "ins3965!", "name": "admin"}}}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    s.jar, _ = cookiejar.New(nil)

    tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client := &http.Client{Transport: tr, Jar: s.jar}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    Debug.Println("response Status:", resp.Status)
    Debug.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    Debug.Println("response Body:", string(body))
}

// GetMo:
//   Given the URL, GET the MO from the switch
//   Use the cookie jar saved during Login and cookie is 
//   validated for every request
func (s Session) GetMo (url string) (string, error) {
    //url := s.url_prefix + "class/" + name + ".json?query-target=self"
    Debug.Println("URL:>", url)

    req, err := http.NewRequest("GET", url, nil)
    req.SetBasicAuth("admin","ins3965!")

    tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
    client := &http.Client{Transport: tr, Jar: s.jar}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
        return "",err
    }
    defer resp.Body.Close()

    Debug.Println("response Status:", resp.Status)
    Debug.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    Debug.Println("response Body:", string(body))
    return string(body),nil
}

type Resp struct {
    TotalCount   string   `json:"totalCount"`
    Imdata     []dataT    `json:"imdata"`
}
type dataT  map[string] json.RawMessage

// GetClassAttributes:
//   Utility function to split json string into multiple 
//   strings matching the given MO class name
func GetClassAttributes (className,jsonStr string) ([]string, error) {
    var resp Resp
    var dat  dataT
    err := json.Unmarshal([]byte(jsonStr), &resp)
    if err != nil {
        return nil,err
    }
    var classList []string
    for _,imdata := range(resp.Imdata) {
        if data,ok := imdata[className]; ok {
            err = json.Unmarshal([]byte(string(data)), &dat)
            if err == nil {
                if data,ok = dat["attributes"]; ok {
                    classList = append(classList, string(data))
                }
            } else {
                println("Error parsing JSON data");
                return nil,err
            }
        }
    }
    return classList,nil
}

// Test 
func main() {
    switch_ip := "swmp5-spine1.insieme.local"
    s := NewSession(switch_ip)
    s.Login()
    url := s.url_prefix + "class/bgpPeer.json?query-target=self"
    str,_ := s.GetMo(url)
    println(str)
}
