package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var runTimestamp string
var docTitle string

const verString = "1.09"

func main() {
	var allContent string

	fmt.Print("WAStats v" + verString + " - WhatsApp message statistics\n(c) 2018-20, Arun Tanksali. All rights reserved\n\n")
	osptr := flag.String("os", "Android", "-os Android or -os iOS")
	fptr := flag.String("db", "", "WhatsApp data export file, with full path")
	titleptr := flag.String("title", "none", "-title document-title")
	verptr := flag.Int("ver", 0, "-ver dataver either 18 or 20")
	fromptr := flag.String("from", "01/01/2000", "-from dd/mm/yyyy")
	toptr := flag.String("to", "01/01/2100", "-to dd/mm/yyyy")

	flag.Parse()

	if *fptr == "" {
		println("Usage:\nWAStats -db <filename> -os <Phone OS> -title <title> -ver <Android ver> -from <fromdate> -to <todate>")
		println("-db <filename> : Mandatory: the filename is that of the WhasApp exported file, without media")
		println("-os <Phone OS> : Optional: this is the OS of the phone on which the file above is exported. Values Android or iOS")
		println("-title <title> : Optional : this is the title that will be used in the report and in the output filename")
		println("-ver <Android ver> : Optional : Specify 18 as the value if and Android file is not processed")
		println("-from and -to : Optional : Select messages from a period to process. Format DD/MM/YYYY")
		return
	}

	if *titleptr == "none" {
		docTitle = *fptr
	} else {
		docTitle = *titleptr
	}

	version := 0
	if *verptr == 0 {
		version = 20
	} else {
		version = *verptr
	}

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Usage: WAStats -db WhatsApp exported file -os iOS or Android -title document-title")
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	rand.Seed(time.Now().UTC().UnixNano())
	openWAM()

	fromdate, _ := time.Parse("02/01/2006", *fromptr)
	todate, _ := time.Parse("02/01/2006", *toptr)

	s := bufio.NewScanner(f)
	allok := false
	for s.Scan() {
		var msgDate, msgTime, msgSender, msgType, msgContent string
		if *osptr == "Android" {
			allok, msgDate, msgTime, msgSender, msgType, msgContent = parseAndroid(s.Text(), version, fromdate, todate)
		} else {
			allok, msgDate, msgTime, msgSender, msgType, msgContent = parseiOS(s.Text(), fromdate, todate)
		}
		if allok && (msgType == "Text" || msgType == "Media" || msgType == "Link") {
			insertWAM("1", msgDate, msgTime, msgType, msgSender)
			if msgType == "Text" {
				allContent = allContent + " " + msgContent
			}
		}
		//fmt.Println(msgContent)
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	genPDFStart()
	allStats()
	generateWordcloud(allContent)
	genPDFEnd(docTitle)

	fmt.Print("Done. Results available in file: WAStats " + docTitle + ".pdf\n")
}
