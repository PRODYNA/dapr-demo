resource "kubernetes_namespace" "dapr-system" {
  metadata {
    name = "dapr-system"
  }
}