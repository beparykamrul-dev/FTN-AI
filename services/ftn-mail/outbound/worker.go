package outbound

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"
)

type DeliveryAttempt struct { ID string; Sender string; Recipients []string; Raw []byte; Attempts int }
type RemoteSender interface { Send(context.Context, string, []string, []byte) error }
type Worker struct { DB *sql.DB; Sender RemoteSender; PollInterval time.Duration; MaxAttempts int; BaseBackoff time.Duration; ProcessingTTL time.Duration }

func (w *Worker) Run(ctx context.Context) error {
	if w == nil || w.DB == nil || w.Sender == nil { return errors.New("invalid outbound worker") }
	if w.PollInterval <= 0 { w.PollInterval = 5*time.Second }; if w.MaxAttempts <= 0 { w.MaxAttempts=8 }; if w.BaseBackoff<=0 { w.BaseBackoff=30*time.Second }; if w.ProcessingTTL<=0 { w.ProcessingTTL=15*time.Minute }
	t:=time.NewTicker(w.PollInterval); defer t.Stop()
	for { if err:=w.recoverStuck(ctx); err!=nil{return err}; if err:=w.processOne(ctx); err!=nil{return err}; select{case <-ctx.Done():return nil;case <-t.C:} }
}
func (w *Worker) processOne(ctx context.Context) error {
	row:=w.DB.QueryRowContext(ctx, `UPDATE ftn_mail_outbound_queue SET status='processing',updated_at=now() WHERE id=(SELECT id FROM ftn_mail_outbound_queue WHERE status IN ('queued','retry') AND next_attempt_at<=now() ORDER BY next_attempt_at,created_at FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING id,sender,recipients,raw_message,attempts`)
	var m DeliveryAttempt; if err:=row.Scan(&m.ID,&m.Sender,&m.Recipients,&m.Raw,&m.Attempts);err!=nil{if errors.Is(err,sql.ErrNoRows){return nil};return err}
	if err:=w.Sender.Send(ctx,m.Sender,m.Recipients,m.Raw);err==nil{_,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='sent',updated_at=now() WHERE id=$1`,m.ID);return e}
	m.Attempts++; if m.Attempts>=w.MaxAttempts{_,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='failed',attempts=$2,last_error=$3,updated_at=now() WHERE id=$1`,m.ID,m.Attempts,"remote delivery failed");return e}
	backoff:=time.Duration(float64(w.BaseBackoff)*math.Pow(2,float64(m.Attempts-1)));if backoff>24*time.Hour{backoff=24*time.Hour}
	_,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='retry',attempts=$2,last_error=$3,next_attempt_at=now()+$4::interval,updated_at=now() WHERE id=$1`,m.ID,m.Attempts,"remote delivery failed",backoff.String());return e
}
func (w *Worker) recoverStuck(ctx context.Context) error { _,err:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='retry',next_attempt_at=now(),updated_at=now(),last_error='recovered stale processing job' WHERE status='processing' AND updated_at < now()-$1::interval`,w.ProcessingTTL.String());return err }
