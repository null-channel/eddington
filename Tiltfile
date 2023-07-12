# # api service
# docker_build("api-service","./api/")
# k8s_yaml('./api/deploy/deploy.yaml')
# k8s_resource(
#   workload='api-service',
#   port_forwards=9000
# )


# run nats.io container
nats = {'nats': {'image': 'nats', 'ports': ['4222:4222']}}

docker_compose([encode_yaml({'services': nats})])


# # container running service
# docker_build("controller","./application/container-runner/")
# # k8s_yaml(kustomize('./application/container-runner/config/crd'))
# k8s_yaml(kustomize('./application/container-runner/config/default'))