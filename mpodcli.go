package mpodshow

import (
        "os"
        "fmt"
        "log"
        "io/ioutil"
        "strings"
        "github.com/codegangsta/cli"
       )

// global variable to store the switch ip
// Didnt find a way to store user info in 
// cli.App or cli.Context
var switch_ip string

var Debug *log.Logger

// RunCli:
//   Uses codegangsta/cli module to implement the mpod cli
func RunCli() {
    Debug = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	app := cli.NewApp()
	app.Usage = "Show commands for multipod"
	app.UsageText = app.Name+" <switch-name/ip> [global options] commands"
    app.EnableBashCompletion = true
	app.Version = "1.0"
	app.Commands = []cli.Command{
		{
			Name:        "showrun",
			Usage:       "running-config",
			Description: "running-config command",
            Subcommands: []cli.Command{
              {
                Name:  "bgp",
                Usage: "bgp-peers",
			    Description: "show configured BGP peers",
                Action: bgpShow,
              },
              {
                Name:  "rpm",
                Usage: "rpm [recursive]",
                Action: rpmShow,
              },
            },
		}, 
        {
			Name:        "route-map",
			Usage:       "route-map [rtmap-name] [recursive]",
			Description: "show route-map command",
			Action:       rtmapShow,
		},
        {
			Name:        "prefix-list",
			Usage:       "prefix-list [pfxList-name] [recursive]",
			Description: "show prefix-list command",
			Action:       pfxlistShow,
		},
        {
			Name:        "interleak",
			Usage:       "interleak [to_proto] [recursive]",
			Description: "show interleak command",
			Action:       interleakShow,
		},
	}

    args := os.Args
    // If the last parameter is -d, run in debug mode
    if len(args) >= 1 && args[len(args)-1] == "-d" {
        // Also, remove the last parameter from arg list
        os.Args = args[:len(args)-1]
        args = os.Args
        // Redirect debug logs to Stdout
        Debug.SetOutput(os.Stdout)
    }
    // The first parameter is the switch_ip
    if len(os.Args) >= 2 {
        switch_ip = os.Args[1]
        if strings.Count(switch_ip, ".") < 2 {
            switch_ip = switch_ip+".insieme.local"
        }
        // remove the parameter from arg list
        args = os.Args[1:]
    } 
    // Since most use of this tool is to run 'showrun rpm'
    // make that as the default command if none specified
    if len(os.Args) == 2 {
        fmt.Println("Using default command 'showrun rpm'\n")
        args = []string{"cmd", "showrun", "rpm"}
    }
	app.Run(args)
}

// bgpShow:
//   Connect to switch and get bgpPeer MO
func bgpShow (c *cli.Context) {
    s := NewSession(switch_ip)
    s.Login()
    url := s.url_prefix + "class/bgpPeer.json?query-target=self"
    str,err := s.GetMo(url)
    if err == nil {
        BgpShowOutput(str)
    } else {
        println(err)
    }
}

// rpmShow:
//   Connect to switch and get rpm running config including
//   interleak, route-map and prefix-list. In recursive mode
//   print each interleak followed by its route-map and prefix-list
func rpmShow (c *cli.Context) {
    var recursive bool

    s := NewSession(switch_ip)
    s.Login()

    // default is recursive. Print each interleak/ribleak 
    // followed by route-maps and prefix-lists it refers to
    recursive = true
    if len(c.Args()) >= 1 {
        if c.Args()[0] == "all" {
            recursive = false
        }
    }

    var err error
    err = RpmShowRibLeakP(s, "ospf", recursive)
    if err != nil {
        fmt.Println(err)
    }
    err = RpmShowInterLeakP(s, "ospf", recursive)
    if err != nil {
        fmt.Println(err)
    }
    err = RpmShowInterLeakP(s, "isis", recursive)
    if err != nil {
        fmt.Println(err)
    }
    if !recursive {
        err = RpmShowRouteMap(s, "all", false)
        if err != nil {
            fmt.Println(err)
        }
        err = RpmShowPrefixList(s, "all", false)
        if err != nil {
            fmt.Println(err)
        }
    }
}

// ribleakshow:
//   connect to switch and get ribleakp mo
func ribleakShow (c *cli.Context) {
    recursive := false
    name := "all"

    s := NewSession(switch_ip)
    s.Login()

    if len(c.Args()) >= 1 {
        argList := strings.Join(c.Args(), " ")
        if strings.Contains(argList, "recursive") {
            recursive = true
        }
        if c.Args()[0] != "recursive" {
            name = c.Args()[0]
        }
    }

    err := RpmShowRibLeakP(s, name, recursive)
    if err != nil {
        fmt.Println(err)
    }
}

// interleakShow:
//   Connect to switch and get InterLeakP mo
//   Get the specific InterLeakP if specified
func interleakShow (c *cli.Context) {
    recursive := false
    name := "all"

    s := NewSession(switch_ip)
    s.Login()

    if len(c.Args()) >= 1 {
        argList := strings.Join(c.Args(), " ")
        if strings.Contains(argList, "recursive") {
            recursive = true
        }
        if c.Args()[0] != "recursive" {
            name = c.Args()[0]
        }
    }

    var err error
    if name != "all" {
        err = RpmShowInterLeakP(s, name, recursive)
        if err != nil {
            fmt.Println(err)
        }
    } else {
        err = RpmShowInterLeakP(s, "ospf", recursive)
        if err != nil {
            fmt.Println(err)
        }
        err = RpmShowInterLeakP(s, "isis", recursive)
        if err != nil {
            fmt.Println(err)
        }
    }
}

// rtmapShow:
//   Connect to switch and get rtmapRule mo
//   Get the specific route-map if specified
func rtmapShow (c *cli.Context) {
    recursive := false
    name := "all"

    s := NewSession(switch_ip)
    s.Login()

    if len(c.Args()) >= 1 {
        argList := strings.Join(c.Args(), " ")
        if strings.Contains(argList, "recursive") {
            recursive = true
        }
        if c.Args()[0] != "recursive" {
            name = c.Args()[0]
        }
    }

    err := RpmShowRouteMap(s, name, recursive)
    if err != nil {
        fmt.Println(err)
    }
}

// pfxlistShow:
//   Connect to switch and get rtpfxRule mo
//   Get the specific prefix-list if specified
func pfxlistShow (c *cli.Context) {
    recursive := false
    name := "all"

    s := NewSession(switch_ip)
    s.Login()

    if len(c.Args()) >= 1 {
        argList := strings.Join(c.Args(), " ")
        if strings.Contains(argList, "recursive") {
            recursive = true
        }
        if c.Args()[0] != "recursive" {
            name = c.Args()[0]
        }
    }

    err := RpmShowPrefixList(s, name, recursive)
    if err != nil {
        fmt.Println(err)
    }
}
