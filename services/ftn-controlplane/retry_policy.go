package controlplane

type RetryPolicy struct{MaxAttempts int;BackoffSeconds int}
func(p RetryPolicy)Valid()bool{return p.MaxAttempts>=0&&p.MaxAttempts<=20&&p.BackoffSeconds>=0&&p.BackoffSeconds<=3600}
func(p RetryPolicy)CanRetry(attempt int)bool{return p.Valid()&&attempt>=0&&attempt<p.MaxAttempts}
