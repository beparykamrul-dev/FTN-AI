package main

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"sync"
)

type FTNBFDObservation struct { Peer string; State FTNBFDState }
type FTNBFDSource interface { ListBFD(context.Context) ([]FTNBFDObservation,error) }
type FTNBFDReconciler struct { mu sync.Mutex; source FTNBFDSource; bridge *FTNRoutedEventBridge; previous map[string]FTNBFDState }
func NewFTNBFDReconciler(source FTNBFDSource,bridge *FTNRoutedEventBridge)*FTNBFDReconciler{return &FTNBFDReconciler{source:source,bridge:bridge,previous:make(map[string]FTNBFDState)}}
func(r *FTNBFDReconciler)Reconcile(ctx context.Context)error{if r==nil||r.source==nil||r.bridge==nil{return fmt.Errorf("BFD source and event bridge are required")};if ctx==nil{return fmt.Errorf("context is required")};select{case<-ctx.Done():return ctx.Err();default:};r.mu.Lock();defer r.mu.Unlock();observations,err:=r.source.ListBFD(ctx);if err!=nil{return err};seen:=make(map[string]struct{},len(observations));for _,o:=range observations{peer:=strings.TrimSpace(o.Peer);if _,err:=netip.ParseAddr(peer);err!=nil{continue};switch o.State{case FTNBFDUp,FTNBFDDown,FTNBFDUnknown:default:continue};seen[peer]=struct{}{};previous,ok:=r.previous[peer];if !ok||previous!=o.State{if err:=r.bridge.BFDState(ctx,peer,o.State);err!=nil{return err}};r.previous[peer]=o.State};for peer:=range r.previous{if _,ok:=seen[peer];ok{continue};if err:=r.bridge.BFDState(ctx,peer,FTNBFDUnknown);err!=nil{return err};delete(r.previous,peer)};return nil}
