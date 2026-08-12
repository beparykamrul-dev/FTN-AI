package outbound

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type SMTPSender struct { Resolver MXResolver; TLSConfig *tls.Config; ConnectTimeout time.Duration; Hostname string }

func (s *SMTPSender) Send(ctx context.Context, sender string, recipients []string, raw []byte) error {
	if s == nil || s.Resolver == nil || sender == "" || len(recipients) == 0 || len(raw) == 0 { return errors.New("invalid SMTP sender") }
	_, domain, ok := strings.Cut(sender, "@"); if !ok || domain == "" { return errors.New("invalid sender address") }
	mxs, err := s.Resolver.LookupMX(ctx, domain); if err != nil || len(mxs) == 0 { return errors.New("MX lookup failed") }
	for _, mx := range mxs { if err := s.sendMX(ctx, strings.TrimSuffix(mx.Host, "."), sender, recipients, raw); err == nil { return nil } }
	return errors.New("all MX delivery attempts failed")
}

func (s *SMTPSender) sendMX(ctx context.Context, host, sender string, recipients []string, raw []byte) error {
	conn, err := (&net.Dialer{Timeout:s.ConnectTimeout}).DialContext(ctx, "tcp", net.JoinHostPort(host,"25")); if err != nil { return err }; defer conn.Close()
	if deadline,ok:=ctx.Deadline();ok{_ = conn.SetDeadline(deadline)}
	br:=bufio.NewReader(conn); bw:=bufio.NewWriter(conn)
	read:=func()(string,error){line,e:=br.ReadString('\n');if e!=nil{return "",e};return strings.TrimSpace(line),nil}
	write:=func(v string)error{_,e:=fmt.Fprintf(bw,"%s\r\n",v);if e==nil{e=bw.Flush()};return e}
	g,e:=read();if e!=nil||len(g)<3||g[:3]!="220"{return errors.New("invalid remote SMTP greeting")}
	if e=write("EHLO "+s.Hostname);e!=nil{return e};if _,e=read();e!=nil{return e}
	if s.TLSConfig!=nil {if e=write("STARTTLS");e!=nil{return e};r,e:=read();if e!=nil||!strings.HasPrefix(r,"220"){return errors.New("remote STARTTLS rejected")};tc:=tls.Client(conn,&tls.Config{ServerName:host,MinVersion:s.TLSConfig.MinVersion});if e=tc.HandshakeContext(ctx);e!=nil{return e};conn=tc;br=bufio.NewReader(conn);bw=bufio.NewWriter(conn);if e=write("EHLO "+s.Hostname);e!=nil{return e};if _,e=read();e!=nil{return e}}
	if e=write("MAIL FROM:<"+sender+">");e!=nil{return e};r,e:=read();if e!=nil||len(r)<3||r[:3]!="250"{return errors.New("remote sender rejected")}
	for _,rcpt:=range recipients{if e=write("RCPT TO:<"+rcpt+">");e!=nil{return e};r,e=read();if e!=nil||len(r)<3||r[:3]!="250"{return errors.New("remote recipient rejected")}}
	if e=write("DATA");e!=nil{return e};r,e=read();if e!=nil||!strings.HasPrefix(r,"354"){return errors.New("remote DATA rejected")};if _,e=bw.Write(raw);e!=nil{return e};if !strings.HasSuffix(string(raw),"\r\n"){if _,e=bw.WriteString("\r\n");e!=nil{return e}};if _,e=bw.WriteString(".\r\n");e!=nil{return e};if e=bw.Flush();e!=nil{return e};r,e=read();if e!=nil||len(r)<3||r[:3]!="250"{return errors.New("remote message rejected")};_ = write("QUIT");return nil
}
