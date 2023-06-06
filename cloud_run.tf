resource "google_service_account" "cloud_run" {
  account_id = "cloud-run-user"
}

resource "google_project_iam_member" "cloud_run_firestore" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.cloud_run.email}"
}

resource "google_cloud_run_v2_service" "cloud_run" {
  location = var.region
  name     = "web-app"
  template {
    containers {
      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }
    service_account = google_service_account.cloud_run.email
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = var.region
  service     = google_cloud_run_v2_service.cloud_run.name
  policy_data = data.google_iam_policy.noauth.policy_data
}

output "cloud_run_uri" {
  value = google_cloud_run_v2_service.cloud_run.uri
}