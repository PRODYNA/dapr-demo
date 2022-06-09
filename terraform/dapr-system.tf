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

resource "kubernetes_ingress_v1" "dapr-dashboard" {
  metadata {
    name      = "dapr"
    namespace = "dapr-system"
    annotations = {
      "kubernetes.io/ingress.class" = "nginx"
    }
  }

  spec {
    rule {
      host = "dapr.minikube"
      http {
        path {
          path = "/"
          backend {
            service {
              name = "dapr-dashboard"
              port {
                number = "8080"
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_manifest" "podmonitor-dapr-operator" {
  depends_on = [
    helm_release.kube-prometheus-stack
  ]

  for_each  = { "dapr-operator" = {}, "dapr-placement-server" = {}, "dapr-sentry" = {}, "dapr-sidecar-injector" = {} }


  manifest = {
    "apiVersion" = "monitoring.coreos.com/v1"
    "kind" : "PodMonitor"
    "metadata" = {
      "name"      = each.key
      "namespace" = "dapr-system"
    }
    "spec" = {
      "namespaceSelector" = {
        "matchNames" = [
          "dapr-system"
        ]
      }
      "podMetricsEndpoints" = [
        {
          "interval" = "15s"
          "path"     = "/metrics"
          "port"     = "metrics"
        }
      ]
      "selector" = {
        "matchLabels" = {
          "app" = each.key
        }
      }
    }
  }
}

/*
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: dapr-operator
  namespace: dapr-system
spec:
  namespaceSelector:
    matchNames:
    - dapr-system
  podMetricsEndpoints:
  - interval: 15s
    path: /metrics
    port: metrics
  selector:
    matchLabels:
      app: dapr-operator
*/