package builder

import (
	"errors"
	"sync"
	"time"
)

type Engine struct {
	mu       sync.RWMutex
	projects map[string]Project
	jobs     map[string]BuildJob
}

func NewEngine() *Engine {
	return &Engine{projects: map[string]Project{}, jobs: map[string]BuildJob{}}
}

func (e *Engine) CreateProject(p Project) error {
	if p.ID == "" || p.Name == "" {
		return errors.New("project id and name are required")
	}
	if p.Target != TargetWeb && p.Target != TargetAndroid && p.Target != TargetReact && p.Target != TargetFlutter {
		return errors.New("unsupported target")
	}
	if p.Backend.Language == "" || p.Backend.Framework == "" || p.Frontend.Framework == "" {
		return errors.New("backend and frontend specifications are required")
	}
	p.CreatedAt = time.Now().UTC()
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.projects[p.ID]; exists {
		return errors.New("project already exists")
	}
	e.projects[p.ID] = p
	return nil
}

func (e *Engine) GetProject(id string) (Project, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	p, ok := e.projects[id]
	return p, ok
}

func (e *Engine) CreateBuild(projectID, jobID string) (BuildJob, error) {
	if projectID == "" || jobID == "" {
		return BuildJob{}, errors.New("project id and job id are required")
	}
	e.mu.RLock()
	_, ok := e.projects[projectID]
	e.mu.RUnlock()
	if !ok {
		return BuildJob{}, errors.New("project not found")
	}
	job := BuildJob{ID: jobID, ProjectID: projectID, Status: "queued"}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.jobs[jobID]; exists {
		return BuildJob{}, errors.New("build job already exists")
	}
	e.jobs[jobID] = job
	return job, nil
}

func (e *Engine) UpdateBuild(job BuildJob) error {
	if job.ID == "" {
		return errors.New("build job id is required")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.jobs[job.ID]; !ok {
		return errors.New("build job not found")
	}
	e.jobs[job.ID] = job
	return nil
}
