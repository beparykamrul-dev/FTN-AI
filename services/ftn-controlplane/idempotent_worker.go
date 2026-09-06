package controlplane

import("errors";"fmt";"strings";"sync";"time")

var ErrLeaseRequired=errors.New("valid lease required")
var ErrExecutorRequired=errors.New("job executor required")
type JobExecutor func(DurableJob)error
type WorkerResult struct{Job DurableJob;Event JournalEvent;Executed bool}
type IdempotentWorker struct{Jobs JobStore;Events EventJournal;Leases LeaseStore;ExecuteJob JobExecutor;mu sync.Mutex}

func(w *IdempotentWorker)Execute(jobID,workerID,leaseKey string,token uint64,eventType string,now time.Time)(WorkerResult,error){if w==nil{return WorkerResult{},errors.New("worker dependencies are required")};w.mu.Lock();defer w.mu.Unlock();if w.Jobs==nil||w.Events==nil||w.Leases==nil{return WorkerResult{},errors.New("worker dependencies are required")};job,err:=w.Jobs.Get(strings.TrimSpace(jobID));if err!=nil{return WorkerResult{},err};if err=w.Leases.Validate(leaseKey,workerID,token,now);err!=nil{return WorkerResult{},ErrLeaseRequired};if job.State==JobSucceeded{return WorkerResult{Job:job,Executed:false},nil};if w.ExecuteJob==nil{return WorkerResult{},ErrExecutorRequired};if job.State==JobPending{job,err=TransitionJob(job,JobRunning,now);if err!=nil{return WorkerResult{},err};if err=w.Jobs.Update(job);err!=nil{return WorkerResult{},err}};if err=w.ExecuteJob(job);err!=nil{failed,_:=TransitionJob(job,JobFailed,now);failed.LastError=strings.TrimSpace(err.Error());_ = w.Jobs.Update(failed);return WorkerResult{Job:failed},err};eventType=strings.TrimSpace(eventType);if eventType==""{eventType="job.succeeded"};event,err:=w.Events.Append(JournalEvent{ID:fmt.Sprintf("job:%s:attempt:%d",job.ID,job.Attempt),TenantID:job.TenantID,Type:eventType,CorrelationID:job.ID,AggregateID:job.ID});if err!=nil{return WorkerResult{},err};job,err=TransitionJob(job,JobSucceeded,now);if err!=nil{return WorkerResult{},err};if err=w.Jobs.Update(job);err!=nil{return WorkerResult{},err};return WorkerResult{Job:job,Event:event,Executed:true},nil}

func(w *IdempotentWorker)ReconcilePending(tenantID,workerID,leaseKey string,token uint64,now time.Time,ids []string)error{if w==nil{return errors.New("worker is required")};tenantID=strings.TrimSpace(tenantID);for _,id:=range ids{job,err:=w.Jobs.Get(id);if err!=nil||job.TenantID!=tenantID||job.State!=JobPending{continue};if _,err=w.Execute(id,workerID,leaseKey,token,"job.reconciled",now);err!=nil{return err}};return nil}
