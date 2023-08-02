package controllers

import (
	"encoding/json"

	"errors"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/null-channel/eddington/api/app/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getApplication(app models.NullApplication) *unstructured.Unstructured {
	application := &unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      app.Name,
			"namespace": app.Namespace,
			"valueFrom": map[string]interface{}{
				"fieldRef": map[string]interface{}{
					"fieldPath": "metadata.namespace",
				},
			},
		},
		"spec": map[string]interface{}{
			"name":       "name",
			"appVersion": "v1",
			"apps": map[string]interface{}{
				"name":    app.NullApplicationService[0].Name,
				"image":   app.NullApplicationService[0].Image,
				"cpu":     "100m",
				"memory":  "128Mi",
				"storage": "1Gi",
			},
		}}}
	return application
}

// getPatchData will return difference between original and modified document
//
//nolint:golint,unused
func getPatchData(originalObj, modifiedObj interface{}) ([]byte, error) {
	originalData, err := json.Marshal(originalObj)
	if err != nil {
		return nil, errors.Wrapf(err, "failed marshal original data")
	} *pb.ContainerServiceClient
}
	// Using strategicpatch package can cause below error
	// Error: CreateTwoWayMergePatch failed: unable to find api field in struct Unstructured for the json field "spec"
	//patchBytes, err := strategicpatch.CreateTwoWayMergePatch(originalData, modifiedData, originalObj)
	// if err != nil {
	// 	return nil, errors.Errorf("CreateTwoWayMergePatch failed: %v", err)
	// }

	patchBytes, err := jsonpatch.CreateMergePatch(originalData, modifiedData)
	if err != nil {
		return nil, errors.Errorf("CreateTwoWayMergePatch failed: %v", err)
	}
	return patchBytes, nil
}
