// Package jetsms provides a small client for the
// JetSMS REST API.
//
// It currently implements:
//   - SMS sending (1-1 via Client.SendOne, 1-n via Client.SendMany)
//   - SMS cancellation (Client.CancelSMS)
//   - Reporting (Client.ReportSmsDetail, Client.ReportSms,
//     Client.ReportSmsSummary, Client.ListOriginators)
//
// Example:
//
//	client, err := jetsms.NewClient(jetsms.Config{
//		Username:   "api-user",
//		Password:   "api-password",
//		Originator: "MYTITLE",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	resp, err := client.SendOne(context.Background(), "9053XXXXXXXX", "Merhaba!", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("response code=%s message=%s", resp.ResponseCode, resp.ResponseMessage)
package jetsms

