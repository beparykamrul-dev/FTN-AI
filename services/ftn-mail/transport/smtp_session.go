package transport

import (
 "bufio"
 "context"
 "crypto/tls"
 "errors"
 "fmt"
 "net"
 "strings"
)

type SessionConfig struct { Hostname string; TLSConfig *tls.Config; RequireTLS bool; MaxMessageSize int64 }
type Authenticator interface { Authenticate(context.Context,string,string)(string,error) }
type Delivery interface { Deliver(context.Context,string,string,[]string,[]byte) error }

func ServeSession(ctx context.Context, conn net.Conn, cfg SessionConfig, auth Authenticator, delivery Delivery) error {
 if conn==nil || cfg.MaxMessageSize<=0 || delivery==nil { return errors.New("invalid SMTP session configuration") }
 defer conn.Close(); br:=bufio.NewReader(conn); bw:=bufio.NewWriter(conn)
 write:=func(code int,text string) error { if _,e:=fmt.Fprintf(bw,"%d %s\r\n",code,text); e!=nil{return e}; return bw.Flush() }
 if e:=write(220,cfg.Hostname+" FTN Mail Service Ready"); e!=nil{return e}
 tlsActive,authenticated:=false,false; identityID,mailFrom:="",""; var recipients []string
 for { select{case <-ctx.Done():return ctx.Err();default:}
  line,e:=br.ReadString('\n'); if e!=nil{return e}; cmdLine:=strings.TrimSpace(line); parts:=strings.Fields(cmdLine)
  if len(parts)==0 { if e:=write(500,"Invalid command");e!=nil{return e};continue }; cmd:=strings.ToUpper(parts[0])
  switch cmd {
  case "EHLO","HELO": if e:=write(250,cfg.Hostname+"\r\n250-AUTH PLAIN\r\n250 STARTTLS");e!=nil{return e}
  case "STARTTLS":
   if tlsActive||cfg.TLSConfig==nil {if e:=write(454,"TLS unavailable");e!=nil{return e};continue}; if e:=write(220,"Ready to start TLS");e!=nil{return e}
   tc:=tls.Server(conn,cfg.TLSConfig);if e:=tc.HandshakeContext(ctx);e!=nil{return e};conn=tc;br=bufio.NewReader(conn);bw=bufio.NewWriter(conn);tlsActive=true;authenticated=false;identityID="";mailFrom="";recipients=nil
  case "AUTH":
   if auth==nil {if e:=write(454,"Authentication unavailable");e!=nil{return e};continue}; if cfg.RequireTLS&&!tlsActive {if e:=write(538,"Encryption required");e!=nil{return e};continue}; if authenticated {if e:=write(503,"Already authenticated");e!=nil{return e};continue}
   args:=strings.Fields(cmdLine);if len(args)<2||!strings.EqualFold(args[1],"PLAIN"){if e:=write(504,"Only AUTH PLAIN is supported");e!=nil{return e};continue};encoded:=""
   if len(args)>=3 {encoded=args[2]} else {if e:=write(334,"");e!=nil{return e};r,e:=br.ReadString('\n');if e!=nil{return e};encoded=strings.TrimSpace(r)}
   id,e:=(AuthExchange{Auth:auth}).AuthenticatePLAIN(ctx,encoded);if e!=nil {if e2:=write(535,"Authentication failed");e2!=nil{return e2};continue};authenticated=true;identityID=id;if e:=write(235,"Authentication successful");e!=nil{return e}
  case "MAIL":
   if cfg.RequireTLS&&!tlsActive {if e:=write(530,"Must issue STARTTLS first");e!=nil{return e};continue};if !authenticated {if e:=write(530,"Authentication required");e!=nil{return e};continue};arg:=strings.TrimSpace(cmdLine[len(parts[0]):]);if len(arg)<5||!strings.EqualFold(arg[:5],"FROM:"){if e:=write(501,"Invalid sender");e!=nil{return e};continue};mailFrom=strings.TrimSpace(arg[5:]);if mailFrom==""{if e:=write(501,"Invalid sender");e!=nil{return e};continue};recipients=nil;if e:=write(250,"Sender OK");e!=nil{return e}
  case "RCPT":
   if !authenticated {if e:=write(530,"Authentication required");e!=nil{return e};continue};if mailFrom==""{if e:=write(503,"Need MAIL FROM first");e!=nil{return e};continue};arg:=strings.TrimSpace(cmdLine[len(parts[0]):]);if len(arg)<3||!strings.EqualFold(arg[:3],"TO:"){if e:=write(501,"Invalid recipient");e!=nil{return e};continue};rcpt:=strings.TrimSpace(arg[3:]);if rcpt==""{if e:=write(501,"Invalid recipient");e!=nil{return e};continue};recipients=append(recipients,rcpt);if e:=write(250,"Recipient OK");e!=nil{return e}
  case "DATA":
   if !authenticated {if e:=write(530,"Authentication required");e!=nil{return e};continue};if mailFrom==""||len(recipients)==0{if e:=write(503,"Need sender and recipient");e!=nil{return e};continue};if e:=write(354,"End data with <CRLF>.<CRLF>");e!=nil{return e};raw,e:=ReadMessage(br,cfg.MaxMessageSize);if e!=nil{if e2:=write(552,"Message too large or invalid");e2!=nil{return e2};continue};if e=delivery.Deliver(ctx,identityID,mailFrom,recipients,raw);e!=nil{if e2:=write(451,"Temporary delivery failure");e2!=nil{return e2};continue};if e:=write(250,"Message accepted");e!=nil{return e};mailFrom="";recipients=nil
  case "RSET": mailFrom="";recipients=nil;if e:=write(250,"Reset OK");e!=nil{return e}
  case "NOOP": if e:=write(250,"OK");e!=nil{return e}
  case "QUIT": _=write(221,"Bye");return nil
  default: if e:=write(502,"Command not implemented");e!=nil{return e}
  }
 }
}
