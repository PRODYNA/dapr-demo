resource "kubernetes_namespace" "messaging" {
  metadata {
    name = "messaging"
  }
}

/*
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
*/

# Install redis
resource "helm_release" "redis" {
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "redis"
  namespace  = "messaging"
  name       = "redis"

  values = [
    file("helm/redis.yaml")
  ]
}

# Subscription to checkout topic
/*
resource "kubernetes_manifest" "checkout-topic" {

  depends_on = [
    helm_release.dapr-system
  ]

  manifest = {
    "apiVersion" = "dapr.io/v1alpha1"
    "kind" = "Subscription"
    "metadata" = {
      "name" = "checkout-topic"
      "namespace" = kubernetes_namespace.messaging.id
    }
    "spec" = {
      "topic" = "checkout"
      "route" = "/order"
      "pubsubname" =  "checkout-topic"
    }
  }
}
*/