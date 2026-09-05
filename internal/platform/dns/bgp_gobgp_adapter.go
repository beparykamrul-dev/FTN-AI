package dns
import("context";"fmt";"strings")
type GoBGPAdapter struct{Address string;Enabled bool}
func NewGoBGPAdapter(address string,enabled bool)*GoBGPAdapter{return &GoBGPAdapter{Address:strings.TrimSpace(address),Enabled:enabled}}
func(a *GoBGPAdapter)Validate()error{if a==nil{return fmt.Errorf("GoBGP adapter is required")};if !a.Enabled{return nil};address:=strings.TrimSpace(a.Address);if address==""{return fmt.Errorf("GoBGP address is required")};a.Address=address;return nil}
func(a *GoBGPAdapter)Publish(ctx context.Context,advertisements []BGPAdvertisement)error{if err:=a.Validate();err!=nil{return err};if ctx==nil{return context.Canceled};if err:=ctx.Err();err!=nil{return err};for _,adv:=range advertisements{if err:=ValidateBGPAdvertisement(adv);err!=nil{return err}};return nil}
func(a *GoBGPAdapter)Withdraw(ctx context.Context,advertisements []BGPAdvertisement)error{if err:=a.Validate();err!=nil{return err};if ctx==nil{return context.Canceled};if err:=ctx.Err();err!=nil{return err};for _,adv:=range advertisements{if err:=ValidateBGPAdvertisement(adv);err!=nil{return err}};return nil}
