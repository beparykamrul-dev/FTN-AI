package dns

import("context";"testing")
func TestSnapshotZoneRejectsHugeRecordSet(t *testing.T){records:=make([]Record,10001);if _,err:=SnapshotZone(context.Background(),"","example.com",records);err==nil{t.Fatal("expected record count rejection")}}
func TestSnapshotZoneRejectsHugeValue(t *testing.T){r:=Record{Name:"www.example.com",Type:"A",Value:string(make([]byte,4097))};if _,err:=SnapshotZone(context.Background(),"","example.com",[]Record{r});err==nil{t.Fatal("expected value bound rejection")}}
