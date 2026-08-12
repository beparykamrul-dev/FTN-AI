package outbound

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type DKIMSigner struct { Domain string; Selector string; PrivateKey *rsa.PrivateKey; Headers []string }

func (s *DKIMSigner) Sign(raw []byte) ([]byte,error) {
	if s==nil||s.Domain==""||s.Selector==""||s.PrivateKey==nil||len(raw)==0{return nil,errors.New("invalid DKIM signer configuration")}
	headers,body,err:=splitMessage(raw);if err!=nil{return nil,err};bh:=sha256.Sum256([]byte(relaxedBody(body)));bodyHash:=base64.StdEncoding.EncodeToString(bh[:])
	selected:=s.Headers;if len(selected)==0{selected=[]string{"From","To","Subject","Date","Message-ID"}};h:=parseHeaders(headers);var present []string
	for _,name:=range selected{if _,ok:=h[strings.ToLower(name)];ok{present=append(present,name)}};if len(present)==0{return nil,errors.New("no DKIM headers available")}
	var ch []string;for _,name:=range present{ch=append(ch,relaxedHeader(name,h[strings.ToLower(name)]))}
	dkim:=fmt.Sprintf("v=1; a=rsa-sha256; c=relaxed/relaxed; d=%s; s=%s; h=%s; bh=%s; b=",s.Domain,s.Selector,strings.Join(present,":"),bodyHash)
	sigInput:=strings.Join(ch,"\r\n")+"\r\n"+relaxedHeader("DKIM-Signature",dkim);digest:=sha256.Sum256([]byte(sigInput));sig,err:=rsa.SignPKCS1v15(rand.Reader,s.PrivateKey,crypto.SHA256,digest[:]);if err!=nil{return nil,err}
	return append([]byte("DKIM-Signature: "+dkim+base64.StdEncoding.EncodeToString(sig)+"\r\n"),raw...),nil
}
func splitMessage(raw []byte)(string,string,error){s:=string(raw);p:=strings.Index(s,"\r\n\r\n");sep:="\r\n\r\n";if p<0{p=strings.Index(s,"\n\n");sep="\n\n"};if p<0{return "","",errors.New("message headers not found")};return s[:p],s[p+len(sep):],nil}
func parseHeaders(v string)map[string]string{m:=map[string]string{};sc:=bufio.NewScanner(strings.NewReader(v));var cur string;for sc.Scan(){line:=sc.Text();if(strings.HasPrefix(line," ")||strings.HasPrefix(line,"\t"))&&cur!=""{m[cur]+=" "+strings.TrimSpace(line);continue};i:=strings.IndexByte(line,':');if i<=0{continue};cur=strings.ToLower(line[:i]);m[cur]=strings.TrimSpace(line[i+1:])};return m}
func relaxedHeader(name,value string)string{v:=regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(value)," ");return strings.ToLower(name)+":"+v}
func relaxedBody(v string)string{v=strings.ReplaceAll(v,"\r\n","\n");lines:=strings.Split(v,"\n");for i:=range lines{lines[i]=strings.TrimRight(regexp.MustCompile(`[ \t]+$`).ReplaceAllString(lines[i],"")," ")};for len(lines)>0&&lines[len(lines)-1]==""{lines=lines[:len(lines)-1]};return strings.Join(lines,"\r\n")+"\r\n"}
