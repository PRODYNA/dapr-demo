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
