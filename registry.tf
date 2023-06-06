resource "google_artifact_registry_repository" "registry" {
  repository_id = "registry"
  format        = "DOCKER"
}