// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"terraform-provider-kuma/internal/kuma"
	"terraform-provider-kuma/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name kuma

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "registry.terraform.io/kenli/kuma",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func main1() {
	// host := "http://127.0.0.1:8000"
	// username := "admin"
	// password := "admin"
	// client, err := kuma.NewClient(&host, &username, &password)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	var item kuma.Monitor
	p := `{
	"id": 20,
	"name": "demo_monitor",
	"pathName": "demo_monitor",
	"url": "https://google.com",
	"method": "GET",
	"port": 53,
	"maxretries": 5,
	"weight": 2000,
	"active": true,
	"type": "http",
	"interval": 60,
	"retryInterval": 20,
	"expiryNotification": true,
	"ignoreTls": true,
	"packetSize": 56,
	"maxredirects": 10,
	"accepted_statuscodes": [
		"200-299"
	],
	"dns_resolve_type": "A",
	"dns_resolve_server": "1.1.1.1",
	"gamedigGivenPortOnly": true,
	"httpBodyEncoding": "json",
	"includeSensitiveData": true
}`

	json.Unmarshal([]byte(p), &item)
	var plan provider.MonitorResourceModel

	// resp, err := client.CreateMonitor(kuma.Monitor{
	// 	Name: "demo_monitor",
	// 	Url:  "https://google.com",
	// 	Type: "http",
	// })
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	if err := provider.ConvertStruct(item, &plan, true); err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v", plan)
}
