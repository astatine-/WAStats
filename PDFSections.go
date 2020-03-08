package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func genPDFAddGeneralStats(wamcount, sendercount string) {

	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 30)

	pdf.MoveTo(0, 50)
	pdf.CellFormat(0, 20, docTitle, "BT", 0, "C", false, 0, "")

	pdf.SetFont("Helvetica", "B", 10)

	pdf.MoveTo(0, 80)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(0, 20, "Total participants in group: "+sendercount, "B", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(0, 20, "Total messages in group: "+wamcount, "B", 0, "C", false, 0, "")
	pdf.MoveTo(0, 150)
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 20, "WAStats Analysis generated at: "+runTimestamp, "T", 0, "C", false, 0, "")

	genPDFAddFooter()
}

func genPDFAddHourMonth(hourdata map[int]int, daydata []nvList, monthdata []nvList, mediadata map[string]int) {

	genPDFAddHeader()

	pdf.MoveTo(10, 30)
	genPDFMonthDataTable(monthdata)
	pdf.MoveTo(60, 30)
	pdf.Image("wh2.png", 60, 30, 130, 0, true, "", 0, "")
	os.Remove("wh2.png")

	pdf.MoveTo(10, 160)
	genPDFDayDataTable(daydata)
	pdf.MoveTo(60, 160)
	pdf.Image("wh4.png", 60, 160, 130, 0, true, "", 0, "")
	os.Remove("wh4.png")

	genPDFAddFooter()
	genPDFAddHeader()

	pdf.MoveTo(10, 30)
	genPDFHourDataTable(hourdata)
	pdf.MoveTo(60, 80)
	pdf.Image("wh1.png", 60, 80, 130, 0, true, "", 0, "")
	os.Remove("wh1.png")

	pdf.MoveTo(10, 200)
	genPDFMediaDataTable(mediadata)
	pdf.MoveTo(60, 180)
	pdf.Image("wh3.png", 60, 180, 100, 0, true, "", 0, "")
	os.Remove("wh3.png")

	genPDFAddFooter()

}

func genPDFHourDataTable(hourdata map[int]int) {
	w := []float64{16.0, 16.0}
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(w[0], 6, "Hour", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Count", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	for hour := 0; hour < 24; hour++ {
		pdf.CellFormat(w[0], 6, fmt.Sprint(hour), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, fmt.Sprint(hourdata[hour]), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

}

func genPDFDayDataTable(dayData []nvList) {
	w := []float64{16.0, 16.0}

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(w[0], 6, "Day", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Count", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	for _, s := range dayData {
		pdf.CellFormat(w[0], 6, fmt.Sprint(s.name), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, fmt.Sprint(s.value), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

}

func genPDFMonthDataTable(monthdata []nvList) {
	w := []float64{16.0, 16.0}

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(w[0], 6, "Month", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Count", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	for _, s := range monthdata {
		pdf.CellFormat(w[0], 6, fmt.Sprint(s.name), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, fmt.Sprint(s.value), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

}

func genPDFMediaDataTable(mediadata map[string]int) {
	w := []float64{16.0, 16.0}
	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(w[0], 6, "Type", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Count", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	for m, s := range mediadata {
		pdf.CellFormat(w[0], 6, fmt.Sprint(m), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, fmt.Sprint(s), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

}

func genPDFAddSenderReport(scl []senderCount) {
	w := []float64{16.0, 106.0, 16.0, 16.0, 16.0, 16.0}
	genPDFAddHeader()
	pdf.SetFont("Helvetica", "B", 14)
	//pdf.Text(10, 25, "Sender Name/Number - Count - Media - Link - Text\n")
	pdf.MoveTo(10, 30)
	pdf.CellFormat(w[0], 6, "Rank", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Sender Name/Number", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[2], 6, "Text", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[3], 6, "Media", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[4], 6, "Link", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[5], 6, "Total", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	rows := 0
	rank := 1
	for _, s := range scl {
		sendername := strings.Trim(s.sender, " \n\t\xe2\x80\xaa\u00ac")
		senderlenmax := math.Min(float64(len(sendername)), 45)
		pdf.SetFont("Helvetica", "", 12)
		pdf.CellFormat(w[0], 6, fmt.Sprint(rank), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, fmt.Sprint(sendername[0:int(senderlenmax)]), "LRBT", 0, "", false, 0, "")
		pdf.SetFont("Helvetica", "", 14)
		pdf.CellFormat(w[2], 6, fmt.Sprint(s.textcount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[3], 6, fmt.Sprint(s.mediacount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[4], 6, fmt.Sprint(s.linkcount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[5], 6, fmt.Sprint(s.count), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
		rows++
		rank++

		if ((rows + 1) % 36) == 0 {
			genPDFAddFooter()
			genPDFAddHeader()
			pdf.SetFont("Helvetica", "B", 14)
			pdf.MoveTo(10, 30)
			pdf.CellFormat(w[0], 6, "Rank", "LRBT", 0, "", false, 0, "")
			pdf.CellFormat(w[1], 6, "Sender Name/Number", "LRBT", 0, "", false, 0, "")
			pdf.CellFormat(w[2], 6, "Text", "LRBT", 0, "", false, 0, "")
			pdf.CellFormat(w[3], 6, "Media", "LRBT", 0, "", false, 0, "")
			pdf.CellFormat(w[4], 6, "Link", "LRBT", 0, "", false, 0, "")
			pdf.CellFormat(w[5], 6, "Total", "LRBT", 0, "", false, 0, "")
			pdf.SetFont("Helvetica", "", 14)
			pdf.Ln(-1)
			rows = 0
		}

	}
	genPDFAddFooter()
}

func genPDFAddTopPostsReport(scl []senderCount) {
	w := []float64{16.0, 106.0, 16.0, 16.0, 16.0, 16.0}
	genPDFAddHeader()
	pdf.SetFont("Helvetica", "B", 14)
	pdf.MoveTo(0, 30)
	pdf.CellFormat(0, 6, "Top 25 Posters", "", 0, "C", false, 0, "")
	//pdf.Text(10, 25, "Sender Name/Number - Count - Media - Link - Text\n")
	pdf.MoveTo(10, 40)
	pdf.CellFormat(w[0], 6, "Rank", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[1], 6, "Sender Name/Number", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[2], 6, "Text", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[3], 6, "Media", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[4], 6, "Link", "LRBT", 0, "", false, 0, "")
	pdf.CellFormat(w[5], 6, "Posts", "LRBT", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 14)
	pdf.Ln(-1)
	rows := 0
	rank := 1
	for _, s := range scl {
		pdf.CellFormat(w[0], 6, fmt.Sprint(rank), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, strings.Trim(s.sender, " \n\t\xe2\x80\xaa\u00ac"), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 6, fmt.Sprint(s.textcount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[3], 6, fmt.Sprint(s.mediacount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[4], 6, fmt.Sprint(s.linkcount), "LRBT", 0, "", false, 0, "")
		pdf.CellFormat(w[5], 6, fmt.Sprint(s.mediacount+s.linkcount), "LRBT", 0, "", false, 0, "")
		pdf.Ln(-1)
		rows++
		rank++

		if ((rows + 1) % 26) == 0 {
			break
		}

	}
	genPDFAddFooter()
}

func genPDFAddWordCloud() {

	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 30)

	pdf.Image("wc.png", 10, 10, 210, 0, true, "", 0, "")
	os.Remove("wc.png")

	genPDFAddFooter()
}
