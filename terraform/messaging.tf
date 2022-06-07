resource "kubernetes_namespace" "messaging" {
  metadata {
    name = "messaging"
  }
}

resource "helm_release" "nats" {
  chart = "nats"
  name = "nats"
  repository = "https://nats-io.github.io/k8s/helm/charts/"
  create_namespace = false
  namespace = kubernetes_namespace.messaging.id

  values = [
    file("helm/nats.yaml")
  ]
}
