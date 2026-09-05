package builder

import("strings";"time")

type Target string
const(TargetWeb Target="web";TargetAndroid Target="android";TargetReact Target="react";TargetFlutter Target="flutter")
type Project struct{ID string `json:"id"`;Name string `json:"name"`;Target Target `json:"target"`;Description string `json:"description,omitempty"`;Backend BackendSpec `json:"backend"`;Frontend FrontendSpec `json:"frontend"`;CreatedAt time.Time `json:"created_at"`}
type BackendSpec struct{Language string `json:"language"`;Framework string `json:"framework"`;Database string `json:"database,omitempty"`;API string `json:"api"`}
type FrontendSpec struct{Framework string `json:"framework"`;Platform string `json:"platform"`}
type BuildJob struct{ID string `json:"id"`;ProjectID string `json:"project_id"`;Status string `json:"status"`;Logs []string `json:"logs,omitempty"`;Artifact string `json:"artifact,omitempty"`}

func(p Project)Valid()bool{return strings.TrimSpace(p.ID)!=""&&strings.TrimSpace(p.Name)!=""&&(p.Target==TargetWeb||p.Target==TargetAndroid||p.Target==TargetReact||p.Target==TargetFlutter)&&strings.TrimSpace(p.Backend.Language)!=""&&strings.TrimSpace(p.Backend.Framework)!=""&&strings.TrimSpace(p.Frontend.Framework)!=""}
func(j BuildJob)Valid()bool{switch strings.ToLower(strings.TrimSpace(j.Status)){case "queued","running","succeeded","failed","cancelled":return strings.TrimSpace(j.ID)!=""&&strings.TrimSpace(j.ProjectID)!="";default:return false}}
