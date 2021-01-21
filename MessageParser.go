package main

import (
	"strings"
	"time"
)

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
func parseAndroid(msgline string, version int, fromdate, todate time.Time, tslayout string) (bool, string, string, string, string, string) {

	var msgDate, msgTime, msgSender, msgType, msgContent string

	msgDate = ""
	msgTime = ""
	msgSender = ""
	msgType = ""
	msgSender = ""
	msgContent = ""

	//7/18/18, 10:04 PM
	//	Mon Jan 2 15:04:05 -0700 MST 2006
	//	0   1   2  3  4  5
	layout := "02/01/2006, 3:04 pm" // changed from "1/2/06, 3:04 PM" on 12 Jan 2020
	if tslayout == "default" {
		if version == 18 {
			layout = "1/2/06, 3:04 PM" // changed from "1/2/06, 3:04 PM" on 12 Jan 2020
		}
		if version == 01 {
			layout = "2/1/06, 3:04 pm" // changed from "1/2/06, 3:04 PM" on 12 Jan 2020
		}
		if version == 24 {
			layout = "2/1/06, 15:04" // changed from "1/2/06, 3:04 PM" on 12 Jan 2020
		}
	} else {
		layout = tslayout // user has given the layout string in proper format
	}

	dash := strings.Index(msgline, " - ")
	if dash != -1 {
		tst := msgline[:dash]
		t, e := time.Parse(layout, tst)
		if e == nil {
			if !inTimeSpan(fromdate, todate, t) {
				return false, msgDate, msgTime, msgSender, msgType, msgContent
			}
			switch {
			case strings.Contains(msgline, "Messages to this"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "CreateGroup"
				}
			case strings.Contains(msgline, "created group"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "CreateGroup"
				}

			case strings.Contains(msgline, "added"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "MemberAdded"
				}

			case strings.Contains(msgline, "now an admin"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "MemberAdmin"
				}

			default:
				{
					colon := strings.Index(msgline, ": ")
					if colon != -1 {
						sender := msgline[dash+3 : colon]
						msg := msgline[colon+1:]
						msgDate = t.Format("2006-01-02 Mon")
						msgTime = t.Format("15:04:05")
						msgSender = sender
						msgContent = msg

						if strings.Contains(msg, "omitted") {
							msgType = "Media"
						} else if strings.Contains(msg, "http://") || strings.Contains(msg, "https://") {
							msgType = "Link"
						} else {
							msgType = "Text"
						}
						//insertWAM("1", t.Format("2006-01-02"), t.Format("15:04:05"), "text", sender)
					} else {
						//fmt.Println(">>>", msgline)
						//replace with an accumulating kind of operation
					}

				}
			}
		} else {
			//fmt.Println(">>>", msgline)
			//fmt.Println("continuation text")
		}
	} else {
		//fmt.Println(">>>", msgline)
		//fmt.Println("continuation text")
	}

	return true, msgDate, msgTime, msgSender, msgType, msgContent
}

func parseiOS(msgline string, fromdate, todate time.Time, tslayout string) (bool, string, string, string, string, string) {

	var msgDate, msgTime, msgSender, msgType, msgContent string

	msgDate = ""
	msgTime = ""
	msgSender = ""
	msgType = ""
	msgSender = ""
	msgContent = ""

	//[01/08/13, 9:49:19 AM]
	layout := "2/1/06, 3:04:05 PM"
	if tslayout != "default" {
		layout = tslayout // user has given the layout string in proper format
	}
	//fmt.Println(msgline)
	dash := strings.Index(msgline, "] ")
	if dash != -1 {
		tst := msgline[1:dash]
		t, e := time.Parse(layout, tst)
		if e == nil {
			if !inTimeSpan(fromdate, todate, t) {
				return false, msgDate, msgTime, msgSender, msgType, msgContent
			}

			switch {
			case strings.Contains(msgline, "Messages to this"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "CreateGroup"
				}
			case strings.Contains(msgline, "created group"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "CreateGroup"
				}

			case strings.Contains(msgline, "added"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "MemberAdded"
				}

			case strings.Contains(msgline, "now an admin"):
				{
					//fmt.Println("--------Admin Message--------")
					msgType = "MemberAdmin"
				}

			default:
				{
					colon := strings.Index(msgline, ": ")
					if colon != -1 {
						sender := msgline[dash+2 : colon]
						msg := msgline[colon+1:]
						//fmt.Println("[", t, "]", "[", sender, "]", "[", msg, "]")
						msgDate = t.Format("2006-01-02 Mon")
						msgTime = t.Format("15:04:05")
						msgSender = sender
						msgContent = msg

						if strings.Contains(msg, "omitted") {
							msgType = "Media"
						} else if strings.Contains(msg, "http://") || strings.Contains(msg, "https://") {
							msgType = "Link"
						} else {
							msgType = "Text"
						}

					} else {
						//fmt.Println(">>>", msgline)
						//replace with an accumulating kind of operation
					}

				}
			}
		} else {
			//fmt.Println(">>>", msgline)
			//fmt.Println("continuation text")
		}
	} else {
		//fmt.Println(">>>", msgline)
		//fmt.Println("continuation text")
	}

	return true, msgDate, msgTime, msgSender, msgType, msgContent
}
