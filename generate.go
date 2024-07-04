package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name zonefile

//go:generate rm -rf ./release/sources
//go:generate go run github.com/google/go-licenses save . --save_path=./release/sources --ignore=github.com/ahamlinman/terraform-provider-zonefile
