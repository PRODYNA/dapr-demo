resource "kubernetes_namespace" "backend" {
  metadata {
    name = "backend"
  }
}

# Deploy the backend services
resource "helm_release" "service" {
  for_each = {"servicea" = {}, "serviceb" = {}, "servicec" = {}}
  name = each.key
  chart = "../charts/service"
  namespace = "backend"
  values = [
    file("helm/${each.key}.yaml")
  ]
  recreate_pods = true
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
      "name" = "pubsub"
      "namespace" = kubernetes_namespace.backend.id
    }
    "spec" = {
      "type" = "pubsub.redis"
      "version" = "v1"
      "metadata" = [
        {
          "name"  = "redisHost"
          "value" = "redis-master.messaging:6379"
        },
        {
          "name" = "consumerID"
          "value" = "eCommerce"
        },
        {
          "name"  = "redisType"
          "value" = "node"
        },
        {
          "name" = "redisUsername",
          "value" = "redis"
        },
        {
          "name" = "redisPassword",
          "value" = "redis"
        },
        {
          "name" = "enableTLS"
          "value" = "false"
        }
      ]
    }
  }
}
