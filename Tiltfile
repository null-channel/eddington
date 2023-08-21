# api service
docker_build("api-service","./",dockerfile="./api/Dockerfile",only=["./api","./proto"])
k8s_yaml('./api/deploy/deploy.yaml')
k8s_resource(
  workload='api-service',
  port_forwards=9000
)


docker_build("nullchannel/eddington-container-builder","./",dockerfile="./application/container-builder/Dockerfile", only=["./application/container-builder/","./proto"])
k8s_yaml('./deployment/container-builder/deployment.yaml')
k8s_yaml('./deployment/container-builder/service.yaml')


# container running service
docker_build("controller","./application/container-runner/")
# k8s_yaml(kustomize('./application/container-runner/config/crd'))
k8s_yaml(kustomize('./application/container-runner/config/default'))
