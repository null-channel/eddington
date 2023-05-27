package templates

var Service = `apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: {{.NullApplicationName}}-{{.AppName}}
  namespace: {{.CustomerID}}
spec:
  hosts:
  - "*"
  gateways:
  - nullcloud-gateway
  http:
  - name: {{.CustomerID}}/{{.AppName}}
    match:
      - name: "{{.CustomerID}}-{{.NullApplicationName}}-{{.AppName}}"
        uri:
          prefix: "/{{.CustomerID}}/{{.NullApplicationName}}/{{.AppName}}/"
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
