package main

import (
	// "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	// "fmt"

	"github.com/joeldsouza28/terraform-provider-reddit/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	// token, err := provider.GetAccessToken("dLM-rMhSgRqGDeaHB6-GJw", "2UAcsKmg31k7e0wceT6rX7_fD0o2WQ", "Sensitive-Cake-1569", "Frankcastle@9")
	// if err != nil {
	// 	fmt.Errorf("failed to get access token: %s", err)
	// }
	// provider.SubmitPost(token, "kubernetes", "Hello there", "Hi there")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
