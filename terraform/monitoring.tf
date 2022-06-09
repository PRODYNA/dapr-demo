resource "kubernetes_namespace" "monitoring" {
  metadata {
    name = "monitoring"
  }
}

resource "helm_release" "kube-prometheus-stack" {
  name       = "monitoring"
  namespace  = kubernetes_namespace.monitoring.id
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"

  values = [
    file("helm/kube-prometheus-stack.yaml")
  ]
}
