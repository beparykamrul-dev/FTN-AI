package controlplane

import "strings"

type TenantScope struct{TenantID string;Allowed bool}
func InTenant(scope TenantScope,tenantID string)bool{return scope.Allowed&&strings.TrimSpace(scope.TenantID)!=""&&strings.TrimSpace(scope.TenantID)==strings.TrimSpace(tenantID)}
