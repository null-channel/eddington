package templates

var Deployment = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.NullApplicationName}}-{{.AppName}}
  namespace: {{.CustomerID}}
spec:
  selector:
    matchLabels:
      app: {{.NullApplicationName}}-{{.AppName}}
  template:
    metadata:
      labels:
        app: {{.NullApplicationName}}-{{.AppName}}
    spec:
      containers: 
      - name: {{.NullApplicationName}}-{{.AppName}}
        image: {{.Image}}
        resources:
          limits:
            memory: '128Mi'
            cpu: '500m'
        ports:
        - containerPort: 3000`

type DeploymentTemplate struct {
	NullApplicationName string
	AppName             string
	CustomerID          string
	Image               string
}
