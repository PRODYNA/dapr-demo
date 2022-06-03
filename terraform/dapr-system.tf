resource "kubernetes_namespace" "dapr-system" {
  metadata {
    name = "dapr-system"
  }
}

resource "helm_release" "dapr-system" {
  name       = "dapr-system"
  chart      = "dapr"
  repository = "https://dapr.github.io/helm-charts"
  namespace  = "dapr-system"
  values = [
    file("helm/dapr-system.yaml")
  ]
}