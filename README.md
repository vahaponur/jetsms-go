# jetsms-go

Basit bir Go client paketi ile JetSMS REST API kullanımı.

Paket ismi:

```go
import jetsms "github.com/vahaponur/jetsms-go"
```

> Not: Kimlik doğrulama, her istekte gövdeye `user` ve `password` alanları konularak yapılıyor (header değil).

## Kurulum

```bash
go get github.com/vahaponur/jetsms-go
```

## Client oluşturma

```go
client, err := jetsms.NewClient(jetsms.Config{
    Username:   "API_USERNAME",
    Password:   "API_PASSWORD",
    Originator: "MYTITLE",                  // varsayılan başlık
    // BaseURL: "https://smsapi.jetsms.com.tr", // boş bırakırsan default bu
})
if err != nil {
    log.Fatal(err)
}
```

## SMS Gönderimi

### 1-1 gönderim (`SendOne`)

Tek numaraya, tek mesaj:

```go
resp, err := client.SendOne(ctx, "9053XXXXXXXX", "Merhaba, bu bir test", nil)
if err != nil {
    // network / HTTP hatası
}
fmt.Println(resp.ResponseCode, resp.ResponseMessage)
```

### 1-n gönderim (`SendMany`)

Aynı mesajı birden fazla numaraya:

```go
recipients := []string{"9053XXXXXXXX", "9053YYYYYYYY"}
resp, err := client.SendMany(ctx, recipients, "Toplu test mesajı", &jetsms.SendOptions{
    Channel: "TRKC", // isteğe bağlı; dokümandaki kanal kodlarından biri
})
```

### Düşük seviye gönderim (`SendSMS`)

Doğrudan `SendSmsRequest` ile çalışmak istersen:

```go
req := jetsms.SendSmsRequest{
    Originator: "MYTITLE",
    SmsMessages: []jetsms.SmsMessage{
        {Recipient: "9053XXXXXXXX", MessageText: "Detaylı istek"},
    },
}
resp, err := client.SendSMS(ctx, req)
```

## SMS İptali

`/api/CancelSms` için:

```go
cancelReq := jetsms.CancelSmsRequest{
    ProcessID: "GROUP_ID",   // SendSms responseMessage ile dönen id
    // MessageID: "OPTIONAL_MESSAGE_ID",
}
cancelResp, err := client.CancelSMS(ctx, cancelReq)
```

## Raporlama

### Detaylı SMS raporu (`ReportSmsDetail`)

Tarih aralığına göre detaylı rapor:

```go
detailReq := jetsms.ReportSmsDetailRequest{
    StartDate: "01012025000000", // ddMMyyyyhhmmss
    EndDate:   "02012025000000",
}
detailResp, err := client.ReportSmsDetail(ctx, detailReq)
```

### Gönderim raporu (`ReportSms`)

Tek bir groupId için mesaj bazlı rapor:

```go
repReq := jetsms.ReportSmsRequest{
    GroupID: "GROUP_ID",
}
repResp, err := client.ReportSms(ctx, repReq)
```

### Gönderim özet raporu (`ReportSmsSummary`)

```go
sumReq := jetsms.ReportSmsSummaryRequest{
    GroupID: "GROUP_ID",
}
sumResp, err := client.ReportSmsSummary(ctx, sumReq)
```

## Başlık Listeleme

Kullanıcının mesaj başlıklarını (originator/title) listeler:

```go
origResp, err := client.ListOriginators(ctx)
for _, o := range origResp.Originators {
    fmt.Println(o.Channel, o.MessageTitle)
}
```

## Testler

- Entegrasyon testleri `tests/` klasöründe ve `.gitignore` içinde olduğu için commit edilmiyor.
- `tests/sms_send_test.go` içinde kendi `testUsername`, `testPassword`, `testOriginator`, `testRecipient` (ve istersen `testRecipient2`) değerlerini doldurup:

```bash
go test ./tests -run TestSend
```

komutuyla gerçek SMS göndererek fonksiyonları test edebilirsin.  
Bu testler, yanlışlıkla CI/CD’de koşup SMS fırlatmamak için özel dizinde tutuluyor.

