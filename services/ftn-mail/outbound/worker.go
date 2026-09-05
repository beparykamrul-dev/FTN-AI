package outbound

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type DeliveryAttempt struct { ID string; Sender string; Recipients []string; Raw []byte; Attempts int; LeaseToken string }
type RemoteSender interface { Send(context.Context, string, []string, []byte) error }
type Worker struct { DB *sql.DB; Sender RemoteSender; PollInterval time.Duration; MaxAttempts int; BaseBackoff time.Duration; ProcessingTTL time.Duration }

func (w *Worker) Run(ctx context.Context) error {
	if w == nil || w.DB == nil || w.Sender == nil { return errors.New("invalid outbound worker") }
	if ctx == nil { return errors.New("context is required") }
	if w.PollInterval <= 0 { w.PollInterval = 5*time.Second }; if w.MaxAttempts <= 0 { w.MaxAttempts=8 }; if w.BaseBackoff<=0 { w.BaseBackoff=30*time.Second }; if w.ProcessingTTL<=0 { w.ProcessingTTL=15*time.Minute }
	t:=time.NewTicker(w.PollInterval); defer t.Stop()
	for { if err:=w.recoverStuck(ctx); err!=nil{return err}; if err:=w.processOne(ctx); err!=nil{return err}; select{case <-ctx.Done():return nil;case <-t.C:} }
}
func (w *Worker) processOne(ctx context.Context) error {
	if ctx == nil { return errors.New("context is required") }
	leaseToken:=uuid.NewString()
	row:=w.DB.QueryRowContext(ctx, `UPDATE ftn_mail_outbound_queue SET status='processing',lease_token=$1,lease_expires_at=now()+$2::interval,updated_at=now() WHERE id=(SELECT id FROM ftn_mail_outbound_queue WHERE status IN ('queued','retry') AND next_attempt_at<=now() ORDER BY next_attempt_at,created_at FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING id,sender,recipients,raw_message,attempts,lease_token::text`, leaseToken, w.ProcessingTTL.String())
	var m DeliveryAttempt; if err:=row.Scan(&m.ID,&m.Sender,&m.Recipients,&m.Raw,&m.Attempts,&m.LeaseToken);err!=nil{if errors.Is(err,sql.ErrNoRows){return nil};return err}
	if m.LeaseToken != leaseToken { return errors.New("outbound lease token mismatch") }
	if err:=w.Sender.Send(ctx,m.Sender,m.Recipients,m.Raw);err==nil{res,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='sent',lease_token=NULL,lease_expires_at=NULL,updated_at=now() WHERE id=$1 AND status='processing' AND lease_token=$2 AND lease_expires_at>now()`,m.ID,leaseToken);if e!=nil{return e};if n,_:=res.RowsAffected();n!=1{return errors.New("outbound delivery lease lost")};return nil}
	m.Attempts++; if m.Attempts>=w.MaxAttempts{res,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='failed',attempts=$2,last_error=$3,lease_token=NULL,lease_expires_at=NULL,updated_at=now() WHERE id=$1 AND status='processing' AND lease_token=$4 AND lease_expires_at>now()`,m.ID,m.Attempts,"remote delivery failed",leaseToken);if e!=nil{return e};if n,_:=res.RowsAffected();n!=1{return errors.New("outbound delivery lease lost")};return nil}
	backoff:=time.Duration(float64(w.BaseBackoff)*math.Pow(2,float64(m.Attempts-1)));if backoff>24*time.Hour{backoff=24*time.Hour}
	res,e:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='retry',attempts=$2,last_error=$3,next_attempt_at=now()+$4::interval,lease_token=NULL,lease_expires_at=NULL,updated_at=now() WHERE id=$1 AND status='processing' AND lease_token=$5 AND lease_expires_at>now()`,m.ID,m.Attempts,"remote delivery failed",backoff.String(),leaseToken);if e!=nil{return e};if n,_:=res.RowsAffected();n!=1{return errors.New("outbound delivery lease lost")};return nil
}
func (w *Worker) recoverStuck(ctx context.Context) error { if ctx==nil{return errors.New("context is required")};_,err:=w.DB.ExecContext(ctx,`UPDATE ftn_mail_outbound_queue SET status='retry',next_attempt_at=now(),lease_token=NULL,lease_expires_at=NULL,updated_at=now(),last_error='recovered stale processing job' WHERE status='processing' AND lease_expires_at IS NOT NULL AND lease_expires_at < now()`);return err }
