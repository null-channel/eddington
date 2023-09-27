#!/bin/sh
#

set -e

kind create cluster --config=./hack/kind-config.yaml

kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.1/manifests/tigera-operator.yaml
kubectl apply -f ./hack/calico.yaml

sleep 20

istioctl install -f ./hack/kind-istio-install.yaml --skip-confirmation
