resource "kubernetes_namespace" "backend" {
  metadata {
    name = "backend"
  }
}

# Deploy the backend services
resource "helm_release" "service" {
  for_each = {"checkout" = {}, "number" = {}, "order" = {}}
  name = each.key
  chart = "../charts/service"
  namespace = "backend"
  values = [
    file("helm/${each.key}.yaml")
  ]
  recreate_pods = true
}

# Component pubsub.redis
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

# Component redis.
resource "kubernetes_manifest" "redis-state" {

  depends_on = [
    helm_release.dapr-system
  ]

  manifest = {
    "apiVersion" = "dapr.io/v1alpha1"
    "kind"       = "Component"
    "metadata"   = {
      "name"      = "state"
      "namespace" = kubernetes_namespace.backend.id
    }
    "spec" = {
      "type"     = "state.redis"
      "version"  = "v1"
      "metadata" = [
        {
          "name"  = "redisHost"
          "value" = "redis-master.messaging:6379"
        },
        {
          "name"  = "redisPassword"
          "value" = "redis"
        },
        {
          "name"  = "redisType"
          "value" = "node"
        },
        {
          "name" = "enableTLS"
          "value" = "false"
        }
      ]
    }
  }
}

/*
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: <NAME>
  namespace: <NAMESPACE>
spec:
  type: state.redis
  version: v1
  metadata:
  - name: redisHost
    value: <HOST>
  - name: redisPassword
    value: <PASSWORD>
  - name: enableTLS
    value: <bool> # Optional. Allowed: true, false.
  - name: failover
    value: <bool> # Optional. Allowed: true, false.
  - name: sentinelMasterName
    value: <string> # Optional
  - name: maxRetries
    value: # Optional
  - name: maxRetryBackoff
    value: # Optional
  - name: ttlInSeconds
    value: <int> # Optional
  - name: queryIndexes
    value: <string> # Optional
*/

# Component redis.
resource "kubernetes_manifest" "redis-pubsub-order" {

  depends_on = [
    helm_release.dapr-system
  ]

  manifest = {
    "apiVersion" = "dapr.io/v1alpha1"
    "kind"       = "Subscription"
    "metadata"   = {
      "name"      = "order-subscription"
      "namespace" = kubernetes_namespace.backend.id
    }
    "spec" = {
      "topic" = "checkout"
      "route" = "/order"
      "pubsubname" = "pubsub"
    }
    "scopes" = [
      "order"
    ]
  }
}

