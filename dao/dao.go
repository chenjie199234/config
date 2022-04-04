package dao

import (
	"net"
	"time"

	//"github.com/chenjie199234/config/api"
	//example "github.com/chenjie199234/config/api/deps/example"
	"github.com/chenjie199234/config/config"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/web"
)

//var ExampleCGrpcApi example.ExampleCGrpcClient
//var ExampleCrpcApi example.ExampleCrpcClient
//var ExampleWebApi  example.ExampleWebClient

//NewApi create all dependent service's api we need in this program
func NewApi() error {
	var e error
	_ = e //avoid unuse

	cgrpcc := getCGrpcClientConfig()
	_ = cgrpcc //avoid unuse

	//init cgrpc client below
	//examplecgrpc e = cgrpc.NewCGrpcClient(cgrpcc, api.Group, api.Name, "examplegroup", "examplename")
	//if e != nil {
	//         return e
	//}
	//ExampleCGrpcApi = example.NewExampleCGrpcClient(examplecgrpc)

	crpcc := getCrpcClientConfig()
	_ = crpcc //avoid unuse

	//init crpc client below
	//examplecrpc, e = crpc.NewCrpcClient(crpcc, api.Group, api.Name, "examplegroup", "examplename")
	//if e != nil {
	// 	return e
	//}
	//ExampleCrpcApi = example.NewExampleCrpcClient(examplecrpc)

	webc := getWebClientConfig()
	_ = webc //avoid unuse

	//init web client below
	//exampleweb, e = web.NewWebClient(webc, api.Group, api.Name, "examplegroup", "examplename", "http://examplehost:exampleport")
	//if e != nil {
	// 	return e
	//}
	//ExampleWebApi = example.NewExampleWebClient(exampleweb)

	return nil
}

func getCGrpcClientConfig() *cgrpc.ClientConfig {
	gc := config.GetCGrpcClientConfig()
	return &cgrpc.ClientConfig{
		ConnectTimeout: time.Duration(gc.ConnectTimeout),
		GlobalTimeout:  time.Duration(gc.GlobalTimeout),
		HeartPorbe:     time.Duration(gc.HeartProbe),
		Discover:       cgrpcDNS,
	}
}

func cgrpcDNS(group, name string, manually <-chan *struct{}, client *cgrpc.CGrpcClient) {
	tker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tker.C:
		case <-manually:
			tker.Reset(time.Second * 10)
		}
		result := make(map[string]*cgrpc.RegisterData)
		addrs, e := net.LookupHost(name + "-service-headless" + "." + group)
		if e != nil {
			log.Error(nil, "[cgrpc.dns] get:", name+"-service-headless", "addrs error:", e)
			continue
		}
		for i := range addrs {
			addrs[i] = addrs[i] + ":10000"
		}
		dserver := make(map[string]struct{})
		dserver["dns"] = struct{}{}
		for _, addr := range addrs {
			result[addr] = &cgrpc.RegisterData{DServers: dserver}
		}
		for len(tker.C) > 0 {
			<-tker.C
		}
		client.UpdateDiscovery(result)
	}
}

func getCrpcClientConfig() *crpc.ClientConfig {
	rc := config.GetCrpcClientConfig()
	return &crpc.ClientConfig{
		ConnectTimeout: time.Duration(rc.ConnectTimeout),
		GlobalTimeout:  time.Duration(rc.GlobalTimeout),
		HeartPorbe:     time.Duration(rc.HeartProbe),
		Discover:       crpcDNS,
	}
}

func crpcDNS(group, name string, manually <-chan *struct{}, client *crpc.CrpcClient) {
	tker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tker.C:
		case <-manually:
			tker.Reset(time.Second * 10)
		}
		result := make(map[string]*crpc.RegisterData)
		addrs, e := net.LookupHost(name + "-service-headless" + "." + group)
		if e != nil {
			log.Error(nil, "[crpc.dns] get:", name+"-service-headless", "addrs error:", e)
			continue
		}
		for i := range addrs {
			addrs[i] = addrs[i] + ":9000"
		}
		dserver := make(map[string]struct{})
		dserver["dns"] = struct{}{}
		for _, addr := range addrs {
			result[addr] = &crpc.RegisterData{DServers: dserver}
		}
		for len(tker.C) > 0 {
			<-tker.C
		}
		client.UpdateDiscovery(result)
	}
}

func getWebClientConfig() *web.ClientConfig {
	wc := config.GetWebClientConfig()
	return &web.ClientConfig{
		ConnectTimeout: time.Duration(wc.ConnectTimeout),
		GlobalTimeout:  time.Duration(wc.GlobalTimeout),
		IdleTimeout:    time.Duration(wc.IdleTimeout),
		HeartProbe:     time.Duration(wc.HeartProbe),
		MaxHeader:      1024,
	}
}
