#!/bin/sh

set -e

ctlptl apply -f ./hack/ctlptl-kind-config.yaml

kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.1/manifests/tigera-operator.yaml
kubectl apply -f ./hack/calico.yaml

kubectl create namespace eddington-admin
kubectl create namespace eddington
kubectl create serviceaccount eddington-admin
kubectl create clusterrolebinding eddington-admin-crb --clusterrole=cluster-admin --serviceaccount=eddington-admin:eddington-admin

cat << EOF | kubectl create -f -
apiVersion: v1
kind: Secret
metadata:
  name: eddington-admin-sa
  namespace: eddington-admin
  annotations:
    kubernetes.io/service-account.name: eddington-admin
type: kubernetes.io/service-account-token
EOF

sleep 20

istioctl install -f ./hack/kind-istio-install.yaml --skip-confirmation

kubectl apply -f ./hack/kpack.yaml
