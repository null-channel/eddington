# api service
docker_build("api-service","./api")
k8s_yaml('./api/deploy/deploy.yaml')
k8s_resource(
  workload='api-service',
  port_forwards=9000
)


docker_build("nullchannel/eddington-container-builder","./container-builder")
k8s_yaml('./deployment/container-builder/config.yaml')
k8s_yaml('./deployment/container-builder/deployment.yaml')
k8s_yaml('./deployment/container-builder/service.yaml')

# container running service
docker_build("controller","./null-operator/")
k8s_yaml(kustomize('./null-operator/config/default'))
