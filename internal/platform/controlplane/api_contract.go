package controlplane

import "slices"
type APIContract struct{Version string `json:"version"`;Resources []string `json:"resources"`;Events []string `json:"events"`}
func DefaultAPIContract()APIContract{return APIContract{Version:"v1",Resources:[]string{"devices","servers","networks","dns","mesh","fiber","customers","services","deployments","changes","audit"},Events:[]string{"device.state","mesh.topology","deployment.state","change.state","alert"}}}
func(c APIContract)Clone()APIContract{c.Resources=slices.Clone(c.Resources);c.Events=slices.Clone(c.Events);return c}
