# api service
docker_build("api-service","./api/")
k8s_yaml('./api/deploy/deploy.yaml')
k8s_resource(
  workload='api-service',
  port_forwards=9000
)


# container running service
docker_build("controller","./application/container-runner/")
# k8s_yaml(kustomize('./application/container-runner/config/crd'))
k8s_yaml(kustomize('./application/container-runner/config/default'))