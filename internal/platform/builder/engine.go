package builder

import (
	"errors"
	"sync"
	"time"
)

type Engine struct {
	mu sync.RWMutex
	projects map[string]Project
	jobs map[string]BuildJob
}

func NewEngine() *Engine { return &Engine{projects: map[string]Project{}, jobs: map[string]BuildJob{}} }

func (e *Engine) CreateProject(p Project) error {
	if p.ID == "" || p.Name == "" { return errors.New("project id and name are required") }
	if p.Target != TargetWeb && p.Target != TargetAndroid && p.Target != TargetReact && p.Target != TargetFlutter { return errors.New("unsupported target") }
	if p.Backend.Language == "" || p.Backend.Framework == "" || p.Frontend.Framework == "" { return errors.New("backend and frontend specifications are required") }
	p.CreatedAt = time.Now().UTC()
	e.mu.Lock(); defer e.mu.Unlock(); e.projects[p.ID] = p
	return nil
}

func (e *Engine) GetProject(id string) (Project, bool) { e.mu.RLock(); defer e.mu.RUnlock(); p, ok := e.projects[id]; return p, ok }

func (e *Engine) CreateBuild(projectID, jobID string) (BuildJob, error) {
	e.mu.RLock(); _, ok := e.projects[projectID]; e.mu.RUnlock()
	if !ok { return BuildJob{}, errors.New("project not found") }
	job := BuildJob{ID: jobID, ProjectID: projectID, Status: "queued"}
	e.mu.Lock(); e.jobs[jobID] = job; e.mu.Unlock()
	return job, nil
}

func (e *Engine) UpdateBuild(job BuildJob) error {
	e.mu.Lock(); defer e.mu.Unlock()
	if _, ok := e.jobs[job.ID]; !ok { return errors.New("build job not found") }
	e.jobs[job.ID] = job
	return nil
}
