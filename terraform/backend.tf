resource "kubernetes_namespace" "backend" {
  metadata {
    name = "backend"
  }
}
