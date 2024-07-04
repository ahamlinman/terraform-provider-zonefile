data "zonefile_records" "example" {
  origin  = "terraform-provider-zonefile.example."
  content = file("terraform-provider-zonefile.example.zone")
}
