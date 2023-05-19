package templates

var Service = `apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.NullApplicationName}}-{{.AppName}}
  namespace: {{.CustomerID}}-{{.NullApplicationName}}
spec:
  hosts:
  - "*"
  gateways:
  - nullcloud-gateway
  http:
  - name: {{.CustomerID}}/{{.AppName}}
    match:
      - uri:
        prefix: "/{{.CustomerID}}/{{.NullApplicationName}}/{{.AppName}}"
    route:
      - destination:
          host: {{.AppName}}
          port:
            number: 80

`

type ServiceTemplate struct {
	NullApplicationName string
	AppName             string
	CustomerID          string
}
