package builder

import("context";"fmt";"os/exec";"strings";"time")
type BuildRunner struct{Timeout time.Duration}
func NewBuildRunner(timeout time.Duration)*BuildRunner{if timeout<=0{timeout=15*time.Minute};return &BuildRunner{Timeout:timeout}}
// Run executes only an explicitly supplied executable and arguments; generated source is never passed to a shell.
func(r *BuildRunner)Run(ctx context.Context,dir,executable string,args ...string)([]byte,error){if r==nil||r.Timeout<=0{return nil,fmt.Errorf("build runner is unavailable")};if ctx==nil{return nil,fmt.Errorf("context is required")};if strings.TrimSpace(executable)==""{return nil,fmt.Errorf("build executable is required")};if strings.TrimSpace(dir)==""{return nil,fmt.Errorf("build directory is required")};ctx,cancel:=context.WithTimeout(ctx,r.Timeout);defer cancel();cmd:=exec.CommandContext(ctx,executable,args...);cmd.Dir=dir;out,err:=cmd.CombinedOutput();if err!=nil{if ctx.Err()!=nil{return out,fmt.Errorf("build timed out: %w",ctx.Err())};return out,fmt.Errorf("build failed: %w",err)};return out,nil}
