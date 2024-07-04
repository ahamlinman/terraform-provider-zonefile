data "zonefile_record_sets" "example" {
  origin  = "terraform-provider-zonefile.example."
  content = file("terraform-provider-zonefile.example.zone")
}
