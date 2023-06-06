resource "google_cloudbuild_trigger" "cloudbuild_trigger" {
  filename = "cloudbuild.yaml"
  github {
    owner = "hsmtkk"
    name  = var.project_name
    push {
      branch = "main"
    }
  }
}