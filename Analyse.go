package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/wcharczuk/go-chart"
)

type nvList struct {
	name  string
	value int
}

func countWAM() {
	var wamcount string
	err := wamdb.QueryRow("select count(*) from wam").Scan(&wamcount)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Records count is ", wamcount)

	var sendercount string
	err = wamdb.QueryRow("select  count (distinct Msender) from wam ").Scan(&sendercount)
	if err != nil {
		log.Println(err)
		return
	}
	//fmt.Println("Sender count is ", sendercount)

	_ = fmt.Sprint(",Message count,", wamcount, ",Sender count,", sendercount, "\n")

	genPDFAddGeneralStats(wamcount, sendercount)
}

func hourWiseReport() map[int]int {
	hourmap := make(map[int]int)
	rows, err := wamdb.Query("select substr(MTime,1,2) as hour, count (*) as hourcnt from wam group by substr(MTime,1,2)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			hour    string
			hourcnt int
		)
		err := rows.Scan(&hour, &hourcnt)
		if err != nil {
			log.Fatal(err)
		}
		hournumber, err := strconv.Atoi(hour)
		if err == nil {
			hourmap[hournumber] = hourcnt
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(hourmap)
	hourChart(hourmap)
	return hourmap
}

func monthWiseReport() []nvList {
	monthmap := make(map[string]int)
	rows, err := wamdb.Query("select substr(MDate,6,2) as month, count (*) as monthcnt from wam group by substr(MDate,6,2);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			month    string
			monthcnt int
		)
		err := rows.Scan(&month, &monthcnt)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(month, monthcnt)
		//monthnumber, err := strconv.Atoi(month)
		if err == nil {
			monthmap[month] = monthcnt
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(monthmap)
	monthlist := organizeBy("Month", monthmap)
	monthChart(monthlist)
	return monthlist
}

func dayWiseReport() []nvList {
	daymap := make(map[string]int)
	rows, err := wamdb.Query("select substr(MDate,12,3) as day, count (*) as daycnt from wam group by substr(MDate,12,3);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			day    string
			daycnt int
		)
		err := rows.Scan(&day, &daycnt)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(month, monthcnt)
		//monthnumber, err := strconv.Atoi(day)
		if err == nil {
			daymap[day] = daycnt
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(monthmap)
	daylist := organizeBy("Day", daymap)
	dayChart(daylist)
	return daylist
}

type senderCount struct {
	sender     string
	count      int
	mediacount int
	linkcount  int
	textcount  int
}

type byPosts []senderCount

func (a byPosts) Len() int      { return len(a) }
func (a byPosts) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPosts) Less(i, j int) bool {
	return (a[i].mediacount + a[i].linkcount) < (a[j].mediacount + a[j].linkcount)
}

func senderReportSlice() {
	senderlist := make([]senderCount, 0)
	rows, err := wamdb.Query("select MSender as sender, count (*) as sendercnt from wam group by MSender order by count(MSender) desc;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			sender    string
			sendercnt int
		)
		err := rows.Scan(&sender, &sendercnt)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(sender, sendercnt)
		var asender senderCount
		asender.sender = sender
		asender.count = sendercnt
		senderlist = append(senderlist, asender)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for i := range senderlist {
		senderlist[i].mediacount, senderlist[i].linkcount, senderlist[i].textcount =
			enrichSenderWithMedia(senderlist[i].sender)
	}
	//fmt.Println(senderlist)

	genPDFAddSenderReport(senderlist)

	sort.Sort(sort.Reverse(byPosts(senderlist)))
	genPDFAddTopPostsReport(senderlist)

}

func enrichSenderWithMedia(sender string) (int, int, int) {
	//select MMedia, count (*) as mcnt from wam where MSender="Arun Tanksali" group by MMedia order by MMedia desc;
	rows, err := wamdb.Query("select MMedia, count (*) as mcnt from wam where MSender='" + sender + "' group by MMedia order by MMedia desc;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		MediaCnt int
		LinkCnt  int
		TextCnt  int
	)
	for rows.Next() {
		var (
			media string
			mcnt  int
		)
		err := rows.Scan(&media, &mcnt)
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case media == "Media":
			{
				MediaCnt = mcnt
			}
		case media == "Link":
			{
				LinkCnt = mcnt
			}
		case media == "Text":
			{
				TextCnt = mcnt
			}
		default:
			{
				TextCnt += mcnt
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return MediaCnt, LinkCnt, TextCnt
}

func mediaWiseReport() map[string]int {
	mediamap := make(map[string]int)
	rows, err := wamdb.Query("select MMedia as media, count (*) as mediacnt from wam group by MMedia")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			media    string
			mediacnt int
		)
		err := rows.Scan(&media, &mediacnt)
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			mediamap[media] = mediacnt
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(mediamap)
	mediaChart(mediamap)
	return mediamap
}

// see : https://github.com/wcharczuk/go-chart
func hourChart(hourdata map[int]int) {

	var hourChartData [24]chart.Value

	for i := 0; i < 24; i++ {
		hourChartData[i].Value = 0.0
		hourChartData[i].Label = fmt.Sprint(i)
	}

	for l, v := range hourdata {
		hourChartData[l].Value = float64(v)
	}

	graph := chart.BarChart{
		Title:      "Hourly Message Statistics",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:     512,
		BarWidth:   50,
		BarSpacing: 10,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: hourChartData[0:24],
	}

	//buffer := bytes.NewBuffer([]byte{})
	file, _ := os.Create("wh1.png")
	defer file.Close()
	w := bufio.NewWriter(file)
	_ = graph.Render(chart.PNG, w)
	//err = file.Write(buffer.Read)
	w.Flush()
	//fmt.Println(err)

}

func monthChart(monthdata []nvList) {
	var monthChartData [12]chart.Value

	for l, v := range monthdata {
		monthChartData[l].Value = float64(v.value)
		monthChartData[l].Label = v.name
	}

	graph := chart.BarChart{
		Title:      "Monthly Message Statistics",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: monthChartData[0:12],
	}

	//buffer := bytes.NewBuffer([]byte{})
	file, _ := os.Create("wh2.png")
	defer file.Close()
	w := bufio.NewWriter(file)
	_ = graph.Render(chart.PNG, w)
	//err = file.Write(buffer.Read)
	w.Flush()
	//fmt.Println(err)

}

func mediaChart(mediadata map[string]int) {

	var mediaChartData []chart.Value

	for l, v := range mediadata {
		var aMedia chart.Value
		aMedia.Value = float64(v)
		aMedia.Label = l
		mediaChartData = append(mediaChartData, aMedia)
	}

	graph := chart.PieChart{
		Title:      "Media Type Statistics",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 5,
			},
		},
		Height: 512,
		Values: mediaChartData[0:len(mediaChartData)],
	}

	//buffer := bytes.NewBuffer([]byte{})
	file, _ := os.Create("wh3.png")
	defer file.Close()
	w := bufio.NewWriter(file)
	_ = graph.Render(chart.PNG, w)
	//err = file.Write(buffer.Read)
	w.Flush()
	//fmt.Println(err)

}

func dayChart(daydata []nvList) {

	var dayChartData []chart.Value

	for _, v := range daydata {
		var aMedia chart.Value
		aMedia.Value = float64(v.value)
		aMedia.Label = v.name
		dayChartData = append(dayChartData, aMedia)
	}

	graph := chart.BarChart{
		Title:      "Day-wise Message Statistics",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: dayChartData[0:len(dayChartData)],
	}

	//buffer := bytes.NewBuffer([]byte{})
	file, _ := os.Create("wh4.png")
	defer file.Close()
	w := bufio.NewWriter(file)
	_ = graph.Render(chart.PNG, w)
	w.Flush()
}

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var dayNames = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

func organizeBy(flag string, datamap map[string]int) []nvList {
	orderedlist := make([]nvList, 0)
	var anv nvList
	switch flag {
	case "Day":
		{
			for i := 0; i < len(dayNames); i++ {

				if val, ok := datamap[dayNames[i]]; ok {
					anv = nvList{dayNames[i], val}
				} else {
					anv = nvList{dayNames[i], 0}
				}
				orderedlist = append(orderedlist, anv)
			}
		}
	case "Month":
		{
			for i := 0; i < len(monthNames); i++ {
				monthnums := fmt.Sprintf("%02d", i+1)
				if val, ok := datamap[monthnums]; ok {
					anv = nvList{monthNames[i], val}
				} else {
					anv = nvList{monthNames[i], 0}
				}
				orderedlist = append(orderedlist, anv)
			}

		}
	}
	return orderedlist
}

func allStats() {
	countWAM()

	hourData := hourWiseReport()

	dayData := dayWiseReport()

	monthData := monthWiseReport()

	mediaData := mediaWiseReport()

	genPDFAddHourMonth(hourData, dayData, monthData, mediaData)

	//fmt.Println("-----------sender wise report--------------")
	senderReportSlice()

}
