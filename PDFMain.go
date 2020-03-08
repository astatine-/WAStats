package main

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var pdf *gofpdf.Fpdf

//see: https://github.com/jung-kurt/gofpdf
func genPDFStart() {
	pdf = gofpdf.New("P", "mm", "A4", "")
	pdf.SetAuthor("Generated", false)
	pdf.SetCreator("WAStats v"+verString, false)
	runTimestamp = time.Now().Format("Mon 02/01/06, 15:04:05")
}

func genPDFAddHeader() {
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 6, "Statistical analysis of WhatsApp messages by WAStats", "", 0, "C", false, 0, "") //print text
	pdf.Ln(-1)
	pdf.CellFormat(0, 6, docTitle, "B", 0, "C", false, 0, "")

}

func genPDFAddFooter() {
	pdf.SetFont("Helvetica", "", 8)
	pdf.MoveTo(10, 270)
	pdf.CellFormat(0, 5, "Generated at "+runTimestamp+" by WAStats v"+verString, "T", 0, "C", false, 0, "")
}

func genPDFEnd(suffix string) {
	err := pdf.OutputFileAndClose("WAStats " + suffix + ".pdf")
	if err != nil {
		fmt.Print(err)
	}
}
