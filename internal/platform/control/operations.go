package control

import("errors";"strings")
type Operation string
const(OpReadHealth Operation="read.health";OpReadServices Operation="read.services";OpReadLogs Operation="read.logs";OpDeploy Operation="app.deploy";OpRollback Operation="app.rollback";OpDNSRead Operation="dns.read";OpDNSWrite Operation="dns.write";OpNetworkRead Operation="network.read";OpNetworkChange Operation="network.change";OpServerRestart Operation="server.restart")
var allowed=map[Operation]bool{OpReadHealth:true,OpReadServices:true,OpReadLogs:true,OpDeploy:true,OpRollback:true,OpDNSRead:true,OpDNSWrite:true,OpNetworkRead:true,OpNetworkChange:true,OpServerRestart:true}
func Validate(op Operation)error{op=Operation(strings.TrimSpace(string(op)));if !allowed[op]{return errors.New("operation is not allowlisted")};return nil}
func RequiresApproval(op Operation)bool{switch Operation(strings.TrimSpace(string(op))){case OpDNSWrite,OpNetworkChange,OpServerRestart:return true;default:return false}}
