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
    name = "dapr"
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