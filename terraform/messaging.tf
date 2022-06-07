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
  chart = "redis"
  namespace = "messaging"
  name = "redis"

  values = [
    file("helm/redis.yaml")
  ]
}

# Component for redis
resource "kubernetes_manifest" "redis-pubsub" {

  depends_on = [
    helm_release.dapr-system
  ]

  manifest = {
    "apiVersion" = "dapr.io/v1alpha1"
    "kind" = "Component"
    "metadata" = {
      "name" = "redis-pubsub"
      "namespace" = "messaging"
    }
    "spec" = {
      "type" = "pubsub.redis"
      "version" = "v1"
      "metadata" = [
        {
          "name"  = "redisHost"
          "value" = "redis-master:6379"
        },
        {
          "name" = "consumerID"
          "value" = "eCommerce"
        },
        {
          "name" = "redisType"
          "value" = "node"
        },
        {
          "name" = "redisPassword",
          "value" = "redis"
        }
      ]
    }
  }
}
