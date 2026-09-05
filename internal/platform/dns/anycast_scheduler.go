package dns
import("context";"fmt";"sync";"time")
type AnycastSchedule struct{Interval time.Duration `json:"interval"`;Probe DNSQueryProbe `json:"probe"`;Policy ProbePolicy `json:"policy"`}
func(s AnycastSchedule)Validate()error{if s.Interval<=0{return fmt.Errorf("probe interval must be positive")};if err:=s.Probe.Validate();err!=nil{return err};return s.Policy.Validate()}
type AnycastScheduler struct{mu sync.Mutex;running bool}
func NewAnycastScheduler()*AnycastScheduler{return &AnycastScheduler{}}
func(s *AnycastScheduler)Run(ctx context.Context,schedule AnycastSchedule,reconcile func(context.Context,DNSProbeResult)error)error{if s==nil{return fmt.Errorf("anycast scheduler is required")};if ctx==nil{return context.Canceled};if err:=schedule.Validate();err!=nil{return err};if reconcile==nil{return fmt.Errorf("reconcile callback is required")};s.mu.Lock();if s.running{s.mu.Unlock();return fmt.Errorf("anycast scheduler is already running")};s.running=true;s.mu.Unlock();defer func(){s.mu.Lock();s.running=false;s.mu.Unlock()}();ticker:=time.NewTicker(schedule.Interval);defer ticker.Stop();run:=func()error{if err:=ctx.Err();err!=nil{return err};observed:=schedule.Probe.Probe(ctx);return reconcile(ctx,observed)};if err:=run();err!=nil{return err};for{select{case<-ctx.Done():return ctx.Err();case<-ticker.C:if err:=run();err!=nil{return err}}}}
