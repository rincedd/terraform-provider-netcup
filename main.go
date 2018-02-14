package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rincedd/terraform-provider-netcup/netcup"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return netcup.Provider()
		},
	})
	//client := api.Client{os.Getenv("NETCUP_WS_USER"), os.Getenv("NETCUP_WS_PASSWORD")}
	//ips, err := client.GetVServerIPs("v22018025690561161")
	//if err != nil {
	//	fmt.Printf("Fatal %s\n", err)
	//	return
	//}
	//fmt.Printf("Got IPs %s\n", ips)
}
