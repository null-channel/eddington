package controllers

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getApplication(name string, namespace string, image string) *unstructured.Unstructured {
	application := &unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
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
				"name":    name,
				"image":   image,
				"cpu":     "100m",
				"memory":  "128Mi",
				"storage": "1Gi",
			},
		}}}
	return application
}

// getPatchData will return difference between original and modified document
func getPatchData(originalObj, modifiedObj interface{}) ([]byte, error) {
	originalData, err := json.Marshal(originalObj)
	if err != nil {
		return nil, errors.Wrapf(err, "failed marshal original data")
	}
	modifiedData, err := json.Marshal(modifiedObj)
	if err != nil {
		return nil, errors.Wrapf(err, "failed marshal original data")
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
