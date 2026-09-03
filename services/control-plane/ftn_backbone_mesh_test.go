package main

import "testing"

func TestSelectFTNLocalServicePrefersLatencyThenCapacity(t *testing.T) {
	got, err := SelectFTNLocalService([]FTNLocalServiceEndpoint{
		{ServiceID:"cache",POPID:"pop-b",Address:"198.51.100.20",Healthy:true,LatencyMS:8,CapacityPercent:70},
		{ServiceID:"cache",POPID:"pop-a",Address:"198.51.100.10",Healthy:true,LatencyMS:5,CapacityPercent:90},
		{ServiceID:"cache",POPID:"pop-c",Address:"198.51.100.30",Healthy:true,LatencyMS:5,CapacityPercent:40},
	})
	if err != nil { t.Fatal(err) }
	if got.POPID!="pop-c" { t.Fatalf("chosen=%+v",got) }
}

func TestSelectFTNLocalServiceRejectsInvalidEndpoints(t *testing.T) {
	if _,err:=SelectFTNLocalService([]FTNLocalServiceEndpoint{{ServiceID:"cache",POPID:"pop-a",Address:"bad",Healthy:true}});err==nil { t.Fatal("expected no healthy service") }
}
