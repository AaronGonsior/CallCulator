package main

import (
	"math/rand"
	"os"

	"bufio"
	//"encoding/json"
	"fmt"
	"github.com/cnkei/gospline"
	"math"
	//"os"
	"strconv"
	"strings"
	//"github.com/cnkei/gospline"

	//"github.com/Arafatk/glot"
	//_ "github.com/gnuplot/gnuplot-old"
	"time"

	opt "github.com/AaronGonsior/optionsscheine2"
)

func check(err error){
	if err != nil{
		fmt.Println(err)
	}
}

/*
func loadJson(path string) string{

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	// Decode the JSON data from the file
	var readStr string
	if err := json.NewDecoder(file).Decode(&readStr); err != nil {
		fmt.Println(err)
		return ""
	}

	// Return the decoded string
	return fmt.Sprint(readStr)

}
 */


type callfunc struct{
	base float64
	cost float64
	factor float64
	date []int
}


type spline_old struct {
	spline gospline.Spline
}
type normalized_spline_old struct{
	spline spline_old
	xrange []float64
	integral float64
}


//Lineares Gleichungs-System (linear equation-system)
type LGS struct{
	n int
	A [][]float64
	y []float64
}

type my_spline struct {
	deg int
	splineType []string
	x []float64
	y []float64
	coeffs []float64
	unique bool
}

var usdtoeur float64
var eurtousd float64



func WriteFile(filename string, content string, pathExt string) {

	/*
	d1 := []byte("hello\ngo\n")
	err := os.WriteFile("/tmp/dat1", d1, 0644)
	check(err)
	 */

	path, err := os.Getwd()

	fmt.Println(path)

	f, err := os.Create(path+pathExt+filename)
	check(err)

	defer f.Close()
	/*
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(content)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)
	 */

	/*
	n3, err := f.WriteString("writes\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n3)

	 */

	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString(content)
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)

	w.Flush()

}

//goland:noinspection ALL
func main(){

	riskAndTimePlottesting := true
	optimalTransporttesting := false
	newIntegraltesting := false
	apitesting := true
	calltesting := false
	splinetesting := false

	if riskAndTimePlottesting {

		var pdistX map[string][]float64 = make(map[string][]float64,0)
		var pdistY map[string][]float64 = make(map[string][]float64,0)
		var pdistDates []string

		pdistDates = append(pdistDates,"2024-06-01")
		pdistX["2024-06-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		pdistY["2024-06-01"] = []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1     }


		pdistDates = append(pdistDates,"2025-01-01")
		pdistX["2025-01-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		pdistY["2025-01-01"] = []float64{0, 2	, 5	, 7	  , 15	, 17   , 17   , 15   , 12   , 10   , 7   , 5     }

		var pdistSplines map[string]my_spline
		pdistSplines = make(map[string]my_spline,0)
		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}

		for _,d := range pdistDates {
			//fmt.Println(pdistX[d],pdistY[d])
			s := NewSpline(splinetype,pdistX[d],pdistY[d])
			ns := NewNormedSpline(s)
			pdistSplines[d] = ns
		}

		mathCode := "SetDirectory[NotebookDirectory[]]\n"

		//FindSigmas
		var sigmasMap map[string][]float64
		sigmasMap = make(map[string][]float64,0)
		levels := []float64{0,0.125,0.25,0.5,0.75,0.875,1}
		for _,d := range pdistDates {
			cumSpline := pdistSplines[d].Integrate()
			tmp,id := cumSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)
			sigmasMap[d] = pdistSplines[d].FindSigmas(levels)
		}



		//Test Print
		for _,d := range pdistDates {
			fmt.Println(d+" : ")
			for _,s := range sigmasMap[d] {
				fmt.Printf("%.6f , ",s)
			}
			fmt.Println("")
		}


		WriteFile("sigmas.nb",mathCode,"/")



	}

	if optimalTransporttesting{


		var pdistX map[string][]float64 = make(map[string][]float64,0)
		var pdistY map[string][]float64 = make(map[string][]float64,0)
		var pdistDates []string

		pdistDates = append(pdistDates,"2024-06-01")
		pdistX["2024-06-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		pdistY["2024-06-01"] = []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1     }


		pdistDates = append(pdistDates,"2025-01-01")
		pdistX["2025-01-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		pdistY["2025-01-01"] = []float64{0, 2	, 5	, 7	  , 15	, 17   , 17   , 15   , 12   , 10   , 7   , 5     }

		var pdistSplines map[string]my_spline
		pdistSplines = make(map[string]my_spline,0)
		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}

		for _,d := range pdistDates {
			//fmt.Println(pdistX[d],pdistY[d])
			s := NewSpline(splinetype,pdistX[d],pdistY[d])
			ns := NewNormedSpline(s)
			pdistSplines[d] = ns
		}

		mathCode := "SetDirectory[NotebookDirectory[]]\n"


		for _,d := range pdistDates {
			tmp,id := pdistSplines[d].PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)
		}

		var cumSplines []my_spline
		var invCumSplines []my_spline

		for _,d := range pdistDates {
			cumX := []float64{}
			fmt.Println("Test: FullIntegral: ",pdistSplines[d].FullIntegralSpline())
			for _,x := range pdistSplines[d].x {
				cumX = append(cumX,pdistSplines[d].IntegralSpline(0,float64(x)))
				fmt.Printf("Test: Integral from 0 to %v : %v \n",x,pdistSplines[d].IntegralSpline(0,float64(x)))
			}
			cumSpline := NewSpline(splinetype,pdistSplines[d].x,cumX)
			cumSplines = append(cumSplines,cumSpline)
			invIntSpline := NewSpline(splinetype,cumX,pdistSplines[d].x)
			invCumSplines = append(invCumSplines,invIntSpline)

			tmp,id := cumSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)

			tmp,id = invIntSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)
		}

		//optimal transport between dates
		//this is not enough
		//for ref, see: https://math.nyu.edu/~tabak/publications/Kuang_Tabak.pdf
		//cumSplines[0].At(invCumSplines[1].At(0.5))

		test := cumSplines[0].Subtract(cumSplines[1])
		tmp,id := test.PrintMathematicaCode()
		mathCode += tmp+"\n"
		mathCode += fmt.Sprintf("s%v\n",id)


		var transportMapFloat []float64
		//cumSplines[0], cumSplines[1] = UnionXYCC(cumSplines[0],cumSplines[1])
		//dx := 0.1
		for _,x := range test.x {
			transportMapFloat = append(transportMapFloat,cumSplines[0].Subtract(cumSplines[1]).IntegralSpline(0,x))
		}
		transportMapSpline := NewSpline(splinetype,cumSplines[0].x,transportMapFloat)
		tmp,id = transportMapSpline.PrintMathematicaCode()
		mathCode += tmp+"\n"
		mathCode += fmt.Sprintf("s%v\n",id)






		WriteFile("optTransport.nb",mathCode,"/")


	}

	if newIntegraltesting {

		x := []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		y := []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1     }

		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		s := NewSpline(splinetype,x,y)

		dx := 0.001

		fmt.Println(s.Integral(min(x),max(x),dx))

		tmp,_ := s.PrintMathematicaCode()
		mathCode := tmp
		fmt.Println(mathCode)


		ns := NewNormedSpline(s)

		fmt.Println("New Spline Integral Test:")
		fmt.Println("Old normed spline Integral:")
		fmt.Println(ns.Integral(ns.x[0],ns.x[len(ns.x)-1],dx))
		fmt.Println("New Spline Full Integral:")
		fmt.Println(ns.FullIntegralSpline())
		//bugged
		fmt.Println("New Spline Integral in bound but max bounds:")
		fmt.Println(ns.IntegralSpline(ns.x[0],ns.x[len(ns.x)-1]))
		a:=120.0
		b:=325.0
		fmt.Printf("New Spline Integral in bound with bounds %v and %v:\n",a,b)
		fmt.Println(ns.IntegralSpline(a,b))
		fmt.Println("Old Integral for same range:")
		fmt.Println(ns.Integral(a,b,dx))

		tmp,_ = ns.PrintMathematicaCode()
		mathCode = tmp
		fmt.Println(mathCode)

	}

	if apitesting {


		content := "SetDirectory[NotebookDirectory[]]\n"


		update := true

		apiKey := opt.LoadJson("apiKey.json")

		url := "https://api.polygon.io/v2/aggs/ticker/C:USDEUR/prev?adjusted=true&apiKey="+apiKey
		fmt.Println("url: ",url)
		_,body,err := opt.APIRequest(url,1)
		check(err)
		body = strings.Split(body,"\"vw\":")[1]
		body = strings.Split(body,",")[0]
		fmt.Println(body)

		usdtoeur,err = strconv.ParseFloat(body,64)
		check(err)
		eurtousd = 1/usdtoeur

		ticker := "TSLA"

		var share_price float64
		url = "https://api.polygon.io/v2/aggs/ticker/"+ticker+"/prev?adjusted=true&apiKey="+apiKey
		fmt.Println("url: ",url)
		_,body,err = opt.APIRequest(url,1)
		check(err)
		body = strings.Split(body,"\"vw\":")[1]
		body = strings.Split(body,",")[0]

		share_price,err = strconv.ParseFloat(body,64)
		check(err)
		fmt.Println("share_price(",ticker,"): ",share_price)


		nMax := -1

		var optreq opt.OptionURLReq
		var options []opt.Option

		optreq = opt.OptionURLReq{
			Ticker:      ticker,
			ApiKey:      apiKey,
			StrikeRange: []int{0,1000},
			DateRange:   /*[]string{"2024-06-01","2024-07-01"}*/[]string{"2023-06-01","2027-01-01"},
			Contract_type: "call",
		}

		if update {

			log := ""
			var msg string

			options, msg = opt.GetOptions(optreq,nMax)
			log += msg

			for _,opt := range options {
				fmt.Println(opt.Print())
			}

			opt.WriteJson("log.json",log)
			opt.WriteJson("options.json",fmt.Sprint(options))

		}


		if !update{

			readStr := opt.LoadJson("options.json")
			readStr = strings.Replace(readStr,"} {","\n",-1)
			readStr = strings.Replace(readStr,"}]","",-1)
			readStr = strings.Replace(readStr,"[{","",-1)
			fmt.Println(readStr)

			options = opt.JsonToOptions("options.json")
			fmt.Println("loaded options: \n",options)
		}

		debug := true


		long := callfunc{
			base:   0,
			cost:   share_price,
			factor: 1,
			date:   nil,
		}
		var addToAll []callfunc
		addToAll = append(addToAll,long)



		optionsDates, optionsMap, callListMap := OptionsToOptionsDates(options, addToAll)

		/*
		var optionsMap map[string][]opt.Option
		optionsMap = make(map[string][]opt.Option)
		var optionsDates []string
		var callListMap map[string][]callfunc
		callListMap = make(map[string][]callfunc)

		callList := []callfunc{}
		for _,optt := range options {

			dateStr := strings.Split(optt.Expiration_date,"-")
			dateInt := []int{}
			for i:=0;i<3;i++ {
				tmp,_ := strconv.Atoi(dateStr[i])
				dateInt = append(dateInt,tmp)
			}

			if len(optionsMap[optt.Expiration_date])>0 {
				optionsMap[optt.Expiration_date] = append(optionsMap[optt.Expiration_date],optt)
				callListMap[optt.Expiration_date] = append(callListMap[optt.Expiration_date],callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Vw,
					factor: 1,
					date:   dateInt,
				})
			} else {
				optionsDates = append(optionsDates,optt.Expiration_date)
				tmp := make([]opt.Option,1)
				tmp[0] = optt
				optionsMap[optt.Expiration_date] = tmp
				tmpp := make([]callfunc,2)
				tmpp[0] = callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Vw,
					factor: 1,
					date:   dateInt,
				}
				tmpp[1] = long
				callListMap[optt.Expiration_date] = tmpp
			}


			callList = append(callList,callfunc{
				base:   float64(optt.Strike_price),
				cost:   optt.Vw,
				factor: 1,
				date:  dateInt ,
			})

		}
		 */



		if debug {
			fmt.Println("len of optionsDates2: ", len(optionsMap), " aka. for how many different dates call options got loaded.\n These are all the dates:")
			for _,d := range optionsDates {
				fmt.Println(d,":")
				fmt.Println("callListMap:")
				for _,c := range callListMap[d] {
					fmt.Println("   ",c)
				}
				fmt.Println("OptionsMap:")
				for _,o := range optionsMap[d] {
					fmt.Println("   ",o)
				}
			}
		}

		//add long

		/*
		callList = append(callList,long)
		fmt.Println(callList)
		 */

		// var shift float64 = +0

		/*
		//25Q1
		x := []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		y := []float64{0, 2	, 5	, 7	  , 15	, 17   , 17   , 15   , 12   , 10   , 7   , 5     }
		 */

		//24Q2
		x := []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		y := []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1     }

		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		s := NewSpline(splinetype,x,y)
		ns := NewNormedSpline(s)
		pdist := ns

		//fmt.Println(s.Integral(min(x),max(x),dx))

		/*
		fmt.Println("New Spline Integral Test:")
		fmt.Println("Old normed spline Integral:")
		fmt.Println(ns.Integral(ns.x[0],ns.x[len(ns.x)-1],dx))
		fmt.Println("New Spline Full Integral:")
		fmt.Println(ns.FullIntegralSpline())
		//bugged
		fmt.Println("New Spline Integral in bound but max bounds:")
		fmt.Println(ns.IntegralSpline(ns.x[0],ns.x[len(ns.x)-1]))
		fmt.Println("New Spline Integral in bound with bounds 0 and 300:")
		fmt.Println(ns.IntegralSpline(0,300))

		 */

		var mathCode string
		dx := 0.01

		path, err := os.Getwd()
		check(err)
		fmt.Println(path)
		currentTime := time.Now()
		live := currentTime.Format("2006-01-02")

		err = os.Mkdir(path+"\\tmp\\"+live, 0755)
		path = path+"\\tmp\\"+live+"\\"

		var strikes []float64
		var costs []float64

		for _,d := range optionsDates {
			folderName := ticker+d+"(live data from "+live+")"
			err = os.Mkdir(path+folderName, 0755)
			check(err)
			optionsList := optionsMap[d]
			callList := callListMap[d]
			callList = callList[len(addToAll):len(callList)]

			fmt.Println(optionsList)

			tmp,id := pdist.PrintMathematicaCode()
			mathCode = tmp
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\pdist.png\", " + fmt.Sprintf("Show[fctplot%v]",id) + ", \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			bestcall, bestE := findBestCall(pdist, callList, dx)
			fmt.Println("Best Call:", bestcall, "\nwith expected return:", bestE)
			mathCode = bestcall.PrintMathematicaCode()
			fmt.Println(mathCode)

			content += fmt.Sprintf("msg1 := Text[\"Assuming the probability distribution (left) for the date %v, the call with strike %.1f has the highest expected return out of all calls options available with %.1f %% expected return. Owning the underlying asset (%v) has an expected return of %.1f %%.  \"];\n\n", callList[0].date, bestcall.base, bestE, ticker, long.ExpectedReturn(ns, dx))
			content += mathCode
			content += "Export[\"" + folderName + "\\-bestCall.png\", {msg1 \n , "+fmt.Sprintf("Show[fctplot%v]",id) +", Show[call,long]}, \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			fmt.Println("owning $TSLA has an expected return of: ", long.ExpectedReturn(ns, dx))

			fmt.Println("\nPrint all calls:\n")
			mathCode = PrintMathematicaCode(callList, share_price)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\allCalls.png\", Show[s], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			fmt.Println("\nDistribution Chart for Call-Long intersections:\n")
			mathCode = MathematicaCodeLongIntersection(callList, share_price)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallLongIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			fmt.Println("\nDistribution Chart for Call-Zero intersections:\n")
			mathCode = MathematicaCodeZeroIntersection(callList)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallZeroIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			fmt.Println("\nDistribution Chart for Call-Zero-Volumes intersections:\n")
			mathCode = MathematicaCodeZeroIntersectionVolumes(optionsList)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallZeroVolumesIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			fmt.Println("\nExpected returns for each strike:")
			mathCode = MathematicaPrintExpectedReturns(pdist, callList, dx)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\expected_returns_strike.png\", Show[xy], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

			strikes = make([]float64,0)
			costs = make([]float64,0)
			for _, opt := range optionsMap[d] {
				strikes = append(strikes, float64(opt.Strike_price))
				costs = append(costs, (opt.Vw))
			}
			mathCode = MathematicaXYPlot(strikes, costs)
			fmt.Println("\nPlot strike vs cost:\n")
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\strike_price.png\", Show[xy], \"CompressionLevel\" -> .25, \n ImageResolution -> 300];\n"

		}


		WriteFile("output.nb",content,"/tmp/"+live+"/")

	}


	if calltesting{


		//cur := 125.0

		//call := callfunc{850.0,5.95*eurtousd,0.01,[]int{1,1,2022}}

		date1 := []int{17,6,2025}
		factor := 0.1
		callList := []callfunc{
			{
				base:   1000,
				cost:   0.31*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   450,
				cost:   1.04*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   300,
				cost:   1.83*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   150,
				cost:   4.3*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   260,
				cost:   2.22*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   200,
				cost:   3.13*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   320,
				cost:   1.67*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   500,
				cost:   0.9*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   600,
				cost:   0.67*eurtousd,
				factor: factor,
				date: date1,
			},
			{
				base:   170,
				cost:   3.8*eurtousd,
				factor: factor,
				date: date1,
			},
		}

		fmt.Println(callList)


		//x := 111.111111111111111111111111111111111111111111
		//fmt.Println(call_v(x, call))
		//fmt.Println("gain:",call_gain_perc(x, call) , "%" )

		/*
		fmt.Println(call)
		fmt.Println("breakeven_ground:", call_breakeven_ground(call))
		fmt.Println("breakeven_base:" , call_breakeven_base(call,cur))
		 */



		//splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		var shift float64 = +0
		x := []float64{0	, 50+shift	, 100+shift	, 150+shift	, 200+shift	, 250+shift	, 300+shift	, 350+shift	, 400+shift	, 450+shift	, 500+shift	}
		y := []float64{0	, 5			, 7			, 15		, 10		, 10		, 15		, 15		, 10		, 7			, 5			}


		//var xrange = []float64{min(x),max(x)}

		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		s := NewSpline(splinetype,x,y)

		dx := 0.0001

		fmt.Println(s.Integral(min(x),max(x),dx))

		tmp, _ := s.PrintMathematicaCode()
		mathCode := tmp
		fmt.Println(mathCode)


		ns := NewNormedSpline(s)

		tmp, _ = ns.PrintMathematicaCode()
		mathCode = tmp
		fmt.Println(mathCode)

		/*
		E := call.ExpectedReturn(ns, dx)
		fmt.Println("ExpectedReturn:", E)

		res := []int{40,100}
		call.ASCIIPlot(xrange,res)

		mathCode = call.PrintMathematicaCode()
		fmt.Println(mathCode)

		 */


		pdist := ns
		bestcall, bestE := findBestCall(pdist, callList, dx)
		fmt.Println("Best Call:",bestcall,"\nwith expected return:", bestE)
		mathCode = bestcall.PrintMathematicaCode()
		fmt.Println(mathCode)


	}

	if splinetesting{

		//splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		splinetype := []string{"2","2","=Sl","=Cv","EQSl"}


		//x := []float64{1,2,3,4,5,6}
		//y := []float64{1,2,4,2,3,7}

		//x := []float64{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
		//y := []float64{1,2,3,5,5,3,0,1,1,5,7,9,9,8,5}

		x := []float64{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
		y := []float64{1,2,3,5,5,2,0,.5,1,5,7,9,9,8,5,3,2,1,1,0}

		/*
		fmt.Println("splinetype:",splinetype)
		fmt.Println("x =",x)
		fmt.Println("y =",y)
		 */



		s := NewSpline(splinetype,x,y)
		s = s.Init(splinetype,x,y)

		tmp,_ := s.PrintMathematicaCode()
		mathCode := tmp
		fmt.Println(mathCode)



		dx := 0.0001
		area := s.Integral(min(x),max(x),dx)

		fmt.Println("Integral:", area)


		ns := NewNormedSpline(s)

		nsarea := ns.Integral(min(x),max(x),dx)
		fmt.Println("nsIntegral:", nsarea)

		tmp,_ = ns.PrintMathematicaCode()
		mathCode = tmp
		fmt.Println(mathCode)


	}



}









// ------------------------------- my spline specific functions -------------------------------

func NewSpline(splineType []string, x []float64, y []float64) my_spline{
	ms := my_spline{
		splineType: nil,
		x:          nil,
		y:          nil,
		coeffs:     nil,
		unique:     false,
	}
	ms = ms.Init(splineType, x, y)
	return ms
}

func (ms my_spline) Init(splineType []string, x []float64, y []float64) my_spline{
	ms.splineType = splineType
	ms.x = x
	ms.y = y
	lgs, err := SplineLGSInit(splineType, x, y)
	check(err)
	ms.unique = lgs.GaussElimination()
	if !ms.unique{
		fmt.Println("Caution: Solution not unique!")
	}
	ms.coeffs = lgs.y
	tmp, err := strconv.ParseFloat(splineType[0],64)
	check(err)
	ms.deg = int(tmp)
	return ms
}

func SplineLGSInit(splineType []string, x []float64, y []float64) (LGS,error){

	inmethodprint := false

	//spline_func_deg := 3
	tmp ,err := strconv.ParseFloat(splineType[0],64)
	check(err)
	spline_func_deg := int(tmp)
	tmp ,err = strconv.ParseFloat(splineType[1],64)
	check(err)
	lamda := int(tmp)
	if lamda != 2{
		return LGS{},fmt.Errorf("spline type only supported with lamda=2")
	}
	deg := spline_func_deg+1 //=4

	constraints := splineType[2:]


	//set value vectors
	x_var := make([][]float64,len(x))
	for i:=0 ; i < len(x_var) ; i++ {
		x_var[i] = make([]float64,deg)
		for j:=0 ; j < len(x_var[0]) ; j++ {
			x_var[i][j] = math.Pow(x[i],float64(deg-j-1))
		}
	}
	if inmethodprint{
		fmt.Println("x_var =",x_var)
	}

	//set slope vectors
	x_slope := make([][]float64,len(x))
	for i:=0 ; i < len(x_slope) ; i++ {
		x_slope[i] = make([]float64,deg)
		for j:=0 ; j < len(x_slope[0])-1 ; j++ {
			//fmt.Println(deg-1-j,"xi^",deg-2-j)
			x_slope[i][j] = float64((deg-1-j))*math.Pow(x[i],float64(deg-2-j))
		}
		x_slope[i][len(x_slope[0])-1] = 0
	}

	if inmethodprint {
		fmt.Println("x_slope =", x_slope)
	}

	//set curvature vectors
	x_curv := make([][]float64,len(x))
	for i:=0 ; i < len(x_curv) ; i++ {
		x_curv[i] = make([]float64,deg)
		for j:=0 ; j < len(x_curv[0])-2 ; j++ {
			//fmt.Println(factorial(3-j)/factorial(deg-3-j),"xi^",deg-3-j)
			x_curv[i][j] = float64(factorial(3-j)/factorial(deg-3-j))*math.Pow(x[i],float64(deg-3-j))
		}
		x_curv[i][len(x_curv[0])-1] = 0
	}
	if inmethodprint {
		fmt.Println("x_curv =", x_curv)
	}



	m := deg*(len(x)-1)
	zeromatrix := make([][]float64,m)
	for i := 0 ; i < m ; i++{
		zeromatrix[i] = make([]float64,m)
	}
	zerovector := make([]float64,m)
	M := LGS{m, zeromatrix, zerovector}


	cur := 0

	//x_var left
	for i := 0 ; i < len(x_var)-1 ; i++ {
		//fmt.Println("fct val cond. left")
		M.SetRow(cur,x_var[i],4*i,y[i])
		cur++
	}

	//x_var right
	for i := 0 ; i < len(x_var)-1 ; i++ {
		//fmt.Println("fct val cond. right")
		M.SetRow(cur,x_var[i+1],4*i,y[i+1])
		cur++
	}

	//=Sl
	if contains(constraints,"=Sl"){
		for i := 0 ; i < len(x_slope)-2 ; i++ {
			//fmt.Println("=Sl")
			row := floatlist_cat(x_slope[i+1],floatlist_negation_compwise(x_slope[i+1]))
			//fmt.Println(row)
			M.SetRow(cur,row,4*i,0)
			cur++
		}
	}

	//0Sl
	if contains(constraints,"0Sl"){
		for i := 0 ; i < len(x_slope)-2 ; i++ {
			//fmt.Println("0Sl")
			M.SetRow(cur,x_slope[i],4*i,0)
			cur++
		}
	}

	//=Cv
	if contains(constraints,"=Cv"){
		for i := 1 ; i < len(x_curv)-1 ; i++ {
			//fmt.Println("=Cv")
			row := floatlist_cat(x_curv[i],floatlist_negation_compwise(x_curv[i]))
			//fmt.Println(row)
			M.SetRow(cur,row,4*(i-1),0)
			cur++
		}
	}

	//0Cv
	if contains(constraints,"0Cv"){
		for i := 1 ; i < len(x_curv)-1 ; i++ {
			//fmt.Println("0Cv")
			//row := floatlist_cat(x_curv[i],floatlist_negation_compwise(x_curv[i]))
			//M.AddRow(cur,row,4*(i-1),0)
			M.SetRow(cur,x_curv[i],4*(i-1),0)
			cur++
		}
	}

	//E0Sl
	if contains(constraints,"E0Sl") {
		//fmt.Println("E0Sl")
		//first
		M.SetRow(cur,x_slope[0],4*0,0)
		cur++
		//last
		M.SetRow(cur,x_slope[len(x_slope)-1],4*(len(x_slope)-2),0)
		cur++
	}

	//E0Cv
	if contains(constraints,"E0Cv"){
		//fmt.Println("E0Cv")
		//first
		M.SetRow(cur,x_curv[0],4*(len(x_curv)-2),0)
		cur++
		//last
		M.SetRow(cur,x_curv[len(x_curv)-1],4*(len(x_curv)-2),0)
		cur++
	}

	//EQSl
	if contains(constraints,"EQSl") {
		//fmt.Println("EQSl")
		//first
		M.SetRow(cur,x_slope[0],4*0,(y[1]-y[0])/(x[1]-x[0]))
		cur++
		//last
		M.SetRow(cur,x_slope[len(x_slope)-1],4*(len(x_slope)-2),(y[len(y)-1]-y[len(y)-2])/(x[len(y)-1]-x[len(y)-2]))
		cur++
	}


	return M,nil
}

func (ms my_spline) PrintMathematicaCode() (string,string) {

	id := fmt.Sprint(rand.Intn(999))

	result := ""
	//result += fmt.Sprintln("Mathematica Code to visualize:\n\n")

	//x={x[0],...,x[n]};
	result += fmt.Sprintf("x%v={",id)
	result += fmt.Sprint(ms.x[0])
	for i := 1 ; i < len(ms.x) ; i++ {
		result += fmt.Sprint(",",ms.x[i])
	}
	result += fmt.Sprintln("};")

	//y={y[0],...,y[n]};
	result += fmt.Sprintf("y%v={",id)
	result += fmt.Sprint(ms.y[0])
	for i := 1 ; i < len(ms.y) ; i++ {
		result += fmt.Sprint(",",ms.y[i])
	}
	result += fmt.Sprintln("};")

	//xyPlot
	result += fmt.Sprintf("xy%v:=ListPlot[Transpose[{x%v, y%v}], PlotStyle -> {AbsolutePointSize[8]},ImageSize -> Large, PlotRange -> Automatic];\n",id,id,id)

	//piecewisePlot
	result += fmt.Sprintf("fct%v[x_]:=Piecewise[{",id)
	for i := 0 ; i < 4*(len(ms.x)-1) ; i+=4 {
		result += fmt.Sprint("{")
		if ms.coeffs[i]>=0{
			result += fmt.Sprint("+")
		}
		result += fmt.Sprintf("%.20fx^3",ms.coeffs[i])
		if ms.coeffs[i+1]>=0{
			result += fmt.Sprint("+")
		}
		result += fmt.Sprintf("%.20fx^2",ms.coeffs[i+1])
		if ms.coeffs[i+2]>=0{
			result += fmt.Sprint("+")
		}
		result += fmt.Sprintf("%.20fx^1",ms.coeffs[i+2])
		if ms.coeffs[i+3]>=0{
			result += fmt.Sprint("+")
		}
		result += fmt.Sprintf("%.20f",ms.coeffs[i+3])
		result += fmt.Sprint(",")
		result += fmt.Sprintf("%.3f",ms.x[i/4])
		result += fmt.Sprint("<=x<=")
		result += fmt.Sprintf("%.3f",ms.x[i/4+1])
		result += fmt.Sprint("}")
		if i<4*(len(ms.x)-1)-4-1 {
			result += fmt.Sprint(",")
		}
	}
	result += fmt.Sprint("}];\n")



	result += fmt.Sprintf("fctplot%v := Plot[fct%v[x]",id,id)
	result += ",{x,"
	result += fmt.Sprintf("%.3f",ms.x[0])
	result += fmt.Sprint(",")
	result += fmt.Sprintf("%.3f",ms.x[len(ms.x)-1])
	result += fmt.Sprint("},ImageSize->Large, PlotRange -> Automatic];\n")

	//Show
	result += fmt.Sprintf("s%v:=Show[fctplot%v, xy%v];\n\n",id,id,id)

	return result,id
}

func (ms my_spline) At (x float64) float64{
	splineNr := 0
	if x > max(ms.x) || x < min(ms.x) {
		fmt.Errorf("x not in range")
		return 0
	}
	for i := 0 ; i < len(ms.x) ; i++ {
		if i+1<len(ms.x){
			if x >= ms.x[i] && x <= ms.x[i+1]{
				splineNr = i
				break
			}
		} else {
			splineNr = i-1
		}
	}
	coeffs := ms.coeffs
	if (ms.deg+1)*(splineNr+1)+1 < len(coeffs){
		coeffs = coeffs[(ms.deg+1)*(splineNr):(ms.deg+1)*(splineNr+1)+1]
	} else {
		coeffs = coeffs[(ms.deg+1)*(splineNr):]
	}

	result := 0.0
	for deg := 0 ; deg <= ms.deg ; deg++ {
		result += coeffs[deg]*math.Pow(x,float64(ms.deg-deg))
	}
	return result
}

func (ms my_spline) Integral(a float64, b float64, dx float64) float64{
	var err error
	f := make([]float64,int((b-a)/dx))
	for i := 0 ; i < len(f) ; i++ {
		f[i] = ms.At(a+float64(i)*dx)
		check(err)
	}
	return Integral(f, dx)
}

func (ms my_spline) IntegralSpline(a,b float64) float64 {
	debug := false

	if debug {
		fmt.Printf("original spline: deg: %v , len(x)=%v , len(y)=%v, len(coeffs)=%v \n",ms.deg, len(ms.x), len(ms.y), len(ms.coeffs))
	}

	if b <= a {
		return 0
	}
	var newX []float64
	var newY []float64
	var newCoeffs []float64
	if a >= ms.x[0] {
		newX = append(newX,a)
		newY = append(newY,ms.At(a))
	} else {
		newX = append(newX,ms.x[0])
		newY = append(newY,ms.y[0])
	}

	j:=0
	for ms.x[j] <= a {
		j++
	}
	for d := 0 ; d < ms.deg+1 ; d++ {
		newCoeffs = append(newCoeffs,ms.coeffs[4*j+d])
	}
	for i := j ; ms.x[i] < b && i < len(ms.x)-1 ; i++ {
		newX = append(newX, ms.x[i])
		newY = append(newY, ms.At(ms.x[i]))
		for d := 0 ; d < ms.deg+1 ; d++ {
			newCoeffs = append(newCoeffs,ms.coeffs[4*i+d])
		}
	}
	if b <= ms.x[len(ms.x)-1] {
		newX = append(newX, b)
		newY = append(newY, ms.At(b))
	}


	var newSpline my_spline = my_spline{
		deg:        ms.deg,
		splineType: ms.splineType,
		x:          newX,
		y:          newY,
		coeffs:     newCoeffs,
		unique:     false,
	}

	if debug {
		fmt.Printf("newSpline: deg: %v , len(x)=%v , len(y)=%v, len(coeffs)=%v \n",newSpline.deg, len(newSpline.x), len(newSpline.y), len(newSpline.coeffs))
	}

	return newSpline.FullIntegralSpline()

}


func UnionXYCC (ms1, ms2 my_spline) (my_spline , my_spline) {
	if isUnionized(ms1,ms2){
		return ms1,ms2
	}
	newX := ms1.x
	for _, x := range ms2.x {
		if !containsFloat(newX,x) {
			newX = append(newX,x)
		}
	}


	var newY1, newY2, newC1, newC2 []float64
	for i,nx := range newX {
		i1 := 0
		i2 := 0

		//increase i1,i2 s.t. ms1.x and ms2.x are just under nx
		for ms1.x[i1] < nx {i1++};i1--
		for ms2.x[i2] < nx {i2++};i2--

		//ms1 update
		if !containsFloat(ms1.x,nx){
			//add new Y
			newY1 = append(newY1,ms1.At(nx))

			//add new C
			for j := ms1.deg*i1 ; j < i1 + ms1.deg+1 ; j++ {
				newC1 = append(newC1,ms1.coeffs[j])
			}

		} else {
			//add old Y
			newY1 = append(newY1,ms1.y[i])
			//add old C
			for j := ms1.deg*i ; j < i + ms1.deg+1 ; j++ {
				newC1 = append(newC1,ms1.coeffs[j])
			}
		}

		//ms2 update
		if !containsFloat(ms2.x,nx){
			//add new Y
			newY2 = append(newY2,ms2.At(nx))

			//add new C
			for j := ms2.deg*i2 ; j < i2 + ms2.deg+1 ; j++ {
				newC2 = append(newC2,ms2.coeffs[j])
			}

		} else {
			//add old Y
			newY2 = append(newY2,ms2.y[i])
			//add old C
			for j := ms2.deg*i ; j < i + ms2.deg+1 ; j++ {
				newC2 = append(newC2,ms2.coeffs[j])
			}
		}
	}

	newms1 := my_spline{
		deg:        ms1.deg,
		splineType: ms1.splineType,
		x:          newX,
		y:          newY1,
		coeffs:     newC1,
		unique:     false,
	}
	newms2 := my_spline{
		deg:        ms2.deg,
		splineType: ms2.splineType,
		x:          newX,
		y:          newY2,
		coeffs:     newC2,
		unique:     false,
	}
	return newms1,newms2


}

func isUnionized (ms1, ms2 my_spline) bool {
	if len(ms1.x) != len(ms2.x){return false}
	for i := range ms1.x {
		if ms1.x[i] != ms2.x[i]{return false}
	}
	if len(ms1.coeffs) != len(ms2.coeffs){return false}
	return true
}

//need UnionXYC first
func (ms1 my_spline) SplineMultiply(ms2 my_spline) my_spline {
	return my_spline{}
}

func (ms my_spline) Factor (factor float64) my_spline {
	var newY []float64
	var newC []float64
	for _,y := range ms.y {
		newY = append(newY,y*factor)
	}
	for _,c := range ms.coeffs {
		newC = append(newC,c*factor)
	}
	return my_spline{
		deg:        ms.deg,
		splineType: ms.splineType,
		x:          ms.x,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}
}

func (ms1 my_spline) Add (ms2 my_spline) my_spline {
	ms1,ms2 = UnionXYCC(ms1,ms2)

	//careful at different degrees!
	if ms1.deg != ms2.deg {
		fmt.Println("unequal degrees in myspline.Add() not supported yet")
		return my_spline{}
	}

	/*
	degSmall,degBig := 0,0
	if ms1.deg < ms2.deg{
		msSmall := ms1
		msBig := ms2
		degSmall = ms1.deg
		degBig = ms2.deg
	} else if ms1.deg > ms2.deg {
		msSmall := ms2
		msBig := ms1
		degSmall = ms2.deg
		degBig = ms1.deg
	}


	if degSmall < degBig{
		for
	}
	 */

	newY,err := addFloat(ms1.y,ms2.y)
	check(err)
	newC,err := addFloat(ms1.coeffs,ms2.coeffs)
	check(err)
	return my_spline{
		deg:        0,
		splineType: ms1.splineType,
		x:          ms1.x,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}
}

func (ms1 my_spline) Subtract (ms2 my_spline) my_spline {
	ms2 = ms2.Factor(-1.0)
	//careful at different degrees!
	return ms1.Add(ms2)
}

func (ms my_spline) FullIntegralSpline() float64 {
	integral := 0.0
	for i := 0 ; i < len(ms.x)-1 ; i++ {
		for d := 0 ; d < ms.deg+1 ; d++ {
			integral += (ms.coeffs[4*i+d]/(float64(ms.deg-d)+1))*math.Pow(ms.x[i+1],float64(ms.deg-d)+1) - (ms.coeffs[4*i+d]/(float64(ms.deg-d)+1))*math.Pow(ms.x[i],float64(ms.deg-d)+1)
		}
	}
	return integral
}

func NewNormedSpline(ms my_spline) my_spline{
	//area := ms.Integral(min(ms.x),max(ms.x),precision)
	area := ms.FullIntegralSpline()
	ns_y := make([]float64, len(ms.y))
	for i,y := range ms.y {
		ns_y[i] = y/area
	}
	ns_coeffs := make([]float64,len(ms.coeffs))
	for i,c := range ms.coeffs {
		ns_coeffs[i] = c/area
	}
	return my_spline{
		deg:        ms.deg,
		splineType: ms.splineType,
		x:          ms.x,
		y:          ns_y,
		coeffs:     ns_coeffs,
		unique:     ms.unique,
	}
}

//return Integrated spline, one dim higher
//still to be tested
func (ms my_spline) Integrate() my_spline {
	cumX := []float64{}
	fmt.Println("Test: FullIntegral: ",ms.FullIntegralSpline())
	for _,x := range ms.x {
		cumX = append(cumX,ms.IntegralSpline(min(ms.x),float64(x)))
		fmt.Printf("Test: Integral from 0 to %v : %v \n",x,ms.IntegralSpline(min(ms.x),float64(x)))
	}
	return NewSpline(ms.splineType,ms.x,cumX)
}

//return derivated spline, one dim lower
func (ms my_spline) Derive() my_spline {
	return my_spline{}
}

//returns the value of the derivative of ms at x
func (ms my_spline) D (x float64) float64 {
	//find ix s.t. ms.x[ix] just under x
	ix:=0;for ms.x[ix]<x {ix++};ix--;
	//Test Print
	//fmt.Println("my_spline D() test: out of ",ms.x, "the one that is just under ",x, " is ",ix)

	//extract relevant coeffs
	var coeffs []float64
	for i := 0 ; i < ms.deg+1 ; i++ {
		coeffs = append(coeffs, ms.coeffs[ix*(ms.deg+1)+i])
	}

	result := 0.0
	for i,c := range coeffs[0:len(coeffs)-1] {
		result += c*(float64(ms.deg-i))*math.Pow(x,float64(ms.deg-i)-1)
	}
	return result
}

func (ms my_spline) NewtonRoot(x0 float64, y float64, tol float64) float64 {
	xn := x0
	for math.Abs(y-ms.At(xn))>tol{
		//fmt.Println("old xn: ",xn," , old yn: ", ms.At(xn), " , D(xn)=",ms.D(xn))
		xn = math.Min(max(ms.x),math.Max(min(ms.x),		xn+(y-ms.At(xn))/ms.D(xn)		))
		//fmt.Println(xn,":",ms.At(xn) , " , difference to ",y," is ",ms.At(xn)-y)
		time.Sleep(1)
	}
	return xn
}

//implement Newton-method first
func (ms my_spline) FindSigmas(levels []float64) []float64 {


	l := min(ms.x)
	r := max(ms.x)
	m := (l+r)/2

	/*
	ml := (m+l)/2
	mr := (m+r)/2
	x0s := []float64{l,r,m,ml,mr}
	 */
	tol := 0.0001
	var intersections []float64
	cumSpline := ms.Integrate()
	for _,l := range levels {
		intersections = append(intersections, cumSpline.NewtonRoot(m,l,tol))
	}
	return intersections
}




// ------------------------------- call specific functions -------------------------------

func findBestCall(pdist my_spline, calllist []callfunc, dx float64) (callfunc, float64){
	best := calllist[0]
	bestE := best.ExpectedReturn(pdist, dx)
	for _,c := range calllist {
		cE := c.ExpectedReturn(pdist, dx)
		if cE > bestE{
			best = c
			bestE = cE
		}
	}
	return best, bestE
}

func ExpectedReturns(pdist my_spline, calllist []callfunc, dx float64) ([]callfunc,[]float64) {
	var expReturns []float64
	var calls []callfunc
	for _,c := range calllist {
		expReturns = append(expReturns , c.ExpectedReturn(pdist,dx) )
		calls = append(calls,c)
	}
	return calls,expReturns
}

func MathematicaPrintExpectedReturns(pdist my_spline, calllist []callfunc, dx float64) string {
	calls, expReturns := ExpectedReturns(pdist,calllist,dx)
	code := "x={"
	for i,c := range calls {
		if i==0{
			code += fmt.Sprintf("%0.f",c.base)
			continue
		}
		code += fmt.Sprintf(",%0.f",c.base)
	}
	code += "};\n"

	code += "y={"
	for i,e := range expReturns {
		if i==0{
			code += fmt.Sprintf("%0.f",e)
			continue
		}
		code += fmt.Sprintf(",%0.f",e)
	}
	code += "};\n"
	code += "xy:=ListPlot[Transpose[{x, y}], PlotStyle -> {AbsolutePointSize[8]},ImageSize -> Large, PlotRange -> Automatic,Joined -> True];\n"
	code += "Show[xy];"
	return code
}

func (call callfunc) At (x float64) float64{
	//return call_v(x, call)
	return call_gain_perc(x, call)
}

func call_breakeven_ground(call callfunc) float64{
	return call.base+call.cost*1/call.factor
}

func call_breakeven_base(call callfunc, curbase float64) float64{
	return call.base*call.factor*curbase/(call.factor*curbase-call.cost)
}

func call_gain_perc(x float64, call callfunc) float64{
	return math.Max(-1,x/(call.cost/call.factor)-call.base/(call.cost/call.factor)-1)*100
}

/*
func call_v(x float64, call callfunc) float64{
	return math.Max(0, (x-call.base)*call.factor)
}
 */

func (call callfunc) PrintMathematicaCode() string{
	//fmt.Println("Mathematica Code to visualize call option value\n\n")
	code := ""
	code += fmt.Sprintln("call:=Plot[100*Max[-1,(x/(",call.cost/call.factor,")-",call.base/(call.cost/call.factor),"-1)],{x,0,500},ImageSize->Large, PlotRange->Automatic];")
	code += fmt.Sprintln("Show[call]")
	return code
}

func PrintMathematicaCode(calls []callfunc, share_price float64) string {
	//fmt.Println("Mathematica Code to visualize call option value\n\n")
	xmax := calls[0].base
	for _,call := range calls {
		if xmax < call.base{
			xmax = call.base
		}
	}
	code := "xmax:=1.5*"+fmt.Sprintf("%.0f",xmax)+";\n"
	for i,call := range calls {
		code += fmt.Sprintf("(* strike: %v *)\n",call.base)
		code += fmt.Sprint("call"+strconv.Itoa(i)+":=Plot[100*Max[-1,(x/(",call.cost/call.factor,")-",call.base/(call.cost/call.factor),"-1)],{x,0,xmax},ImageSize->Large, PlotRange->Automatic];\n")
	}

	code += "long := Plot[100*(x - "+fmt.Sprintf("%.0f",share_price)+")/"+fmt.Sprintf("%.0f",share_price)+", {x, 0, xmax}, PlotStyle -> Red];"

	for i := range calls {
		if i==0{
			code += "s:=Show[{call"+strconv.Itoa(i)
			continue
		}
		code += fmt.Sprintln(",call"+strconv.Itoa(i))
	}
	code += ",long}]\n"
	return code
}

func (call callfunc) ASCIIPlot(xrange []float64, res []int){
	max := call.At(xrange[1])
	min := call.At(xrange[0])

	fmt.Println("min:",min," max:",max)

	fmt.Print(repeatstr(" ",4))
	for i:=0 ; i<res[1]/4 ; i++{
		fmt.Printf("%.1f ",min+float64(4*i)*(max-min)/float64(res[1]))
	}
	fmt.Println("")

	var xi float64
	var yi float64
	for i := 0 ; i <= res[0] ; i++{
		xi =float64(i)*(xrange[1]-xrange[0])/float64(res[0])
		yi = (call.At(xi)-min)/(max-min)*float64(res[1])
		fmt.Printf("%.2f:",xi)
		fmt.Println(repeatstr("#", int(yi)+1 ))
		//fmt.Println(xi,":",yi)
	}
}

/*
func ExpectedReturn (call callfunc, pdist normalized_spline_old) float64 {
	precision := 10000.0
	var E float64
	for i:=pdist.xrange[0] ; i<pdist.xrange[1] ; i+=1/precision{
		E += call.At(i)*pdist.At(i)
	}
	E/=precision
	return E
}
 */

func (call callfunc) ExpectedReturn(pdist my_spline,dx float64) float64{
	var E float64
	for x := min(pdist.x) ; x < max(pdist.x) ; x+=dx {
		E += call.At(x)*pdist.At(x)
	}
	E*=dx
	return E
}

func (call callfunc) ToSpline(a,b float64) my_spline {
	return my_spline{
		deg:        1,
		splineType: nil,
		x:          []float64{a,call.base,b},
		y:          []float64{0,0,call.At(b)},
		coeffs:     []float64{0,0,call.factor/call.cost},
		unique:     false,
	}
}

func (call callfunc) LongIntersection(share_price float64) float64 {
	return share_price*(call.factor*call.base+call.cost)/(call.cost+call.factor*share_price)
}

func (call callfunc) ZeroIntersection() float64 {
	return call.base+call.cost/call.factor
}

func LongIntersection(callList []callfunc, share_price float64) []float64 {
	var interList []float64
	for _,call := range callList {
		interList = append(interList,call.LongIntersection(share_price))
	}
	return interList
}

func MathematicaCodeLongIntersection(callList []callfunc, sharePrice float64) string {
	interList := LongIntersection(callList, sharePrice)
	code := "dist:=DistributionChart[{"
	for i,inter := range interList {
		if i==0 {
			code += ""+fmt.Sprintf("%.0f",inter)
			continue
		}
		code += ","+fmt.Sprintf("%.0f",inter)
	}
	code += "}];\n"
	return code
}

func ZeroIntersection(callList []callfunc) []float64 {
	var interList []float64
	for _,call := range callList {
		interList = append(interList,call.ZeroIntersection())
	}
	return interList
}

func ZeroIntersectionVolume (options []opt.Option) []float64 {
	dateStr := strings.Split(options[0].Expiration_date,"-")
	dateInt := []int{}
	for i:=0;i<3;i++ {
		tmp,_ := strconv.Atoi(dateStr[i])
		dateInt = append(dateInt,tmp)
	}
	var callList []callfunc
	var volumes []int
	for _,optt := range options {
		callList = append(callList, callfunc{
			base:   float64(optt.Strike_price),
			cost:   optt.Vw,
			factor: 1,
			date:   dateInt,
		})
		volumes = append(volumes, optt.Volume)
	}
	var interListVol []float64
	for i,call := range callList {
		for v := 0;v < volumes[i] ; v++ {
			interListVol = append(interListVol,call.ZeroIntersection())
		}
	}
	return interListVol
}

func MathematicaCodeZeroIntersectionVolumes(options []opt.Option) string {
	interList := ZeroIntersectionVolume(options)
	code := "dist:=DistributionChart[{"
	for i,inter := range interList {
		if i==0 {
			code += ""+fmt.Sprintf("%.0f",inter)
			continue
		}
		code += ","+fmt.Sprintf("%.0f",inter)
	}
	code += "}];\n"
	return code
}

func MathematicaCodeZeroIntersection(callList []callfunc) string {
	interList := ZeroIntersection(callList)
	code := "dist:=DistributionChart[{"
	for i,inter := range interList {
		if i==0 {
			code += ""+fmt.Sprintf("%.0f",inter)
			continue
		}
		code += ","+fmt.Sprintf("%.0f",inter)
	}
	code += "}];\n"
	return code
}




// ------------------------------- general functions -------------------------------



func containsFloat(list []float64, item float64) bool {
	eps := 0.01
	for _,l := range list {
		if math.Abs(l-item)<eps {
			return true
		}
	}
	return false
}

func addFloat(l1, l2 []float64) ([]float64,error) {
	if len(l1)!=len(l2){return []float64{},fmt.Errorf("list not same length")}
	var newL []float64
	for i := range l1 {
		newL = append(newL,l1[i]+l2[i])
	}
	return newL,nil
}

func OptionsToOptionsDates (options []opt.Option, addToAll []callfunc) ([]string , map[string][]opt.Option , map[string][]callfunc) {


	var optionsMap map[string][]opt.Option
	optionsMap = make(map[string][]opt.Option)
	var optionsDates []string
	var callListMap map[string][]callfunc
	callListMap = make(map[string][]callfunc)

	var dateStr []string
	var dateInt []int

	for _,optt := range options {

		dateStr = strings.Split(optt.Expiration_date,"-")
		dateInt = []int{}
		for i:=0;i<3;i++ {
			tmp,_ := strconv.Atoi(dateStr[i])
			dateInt = append(dateInt,tmp)
		}

		if len(optionsMap[optt.Expiration_date])>0 {
			optionsMap[optt.Expiration_date] = append(optionsMap[optt.Expiration_date],optt)
			callListMap[optt.Expiration_date] = append(callListMap[optt.Expiration_date],callfunc{
				base:   float64(optt.Strike_price),
				cost:   optt.Vw,
				factor: 1,
				date:   dateInt,
			})
		} else {
			optionsDates = append(optionsDates,optt.Expiration_date)
			tmp := make([]opt.Option,1)
			tmp[0] = optt
			optionsMap[optt.Expiration_date] = tmp
			tmpp := make([]callfunc,len(addToAll))
			for i := range addToAll {
				tmpp[i] = addToAll[i]
			}
			tmpp = append(tmpp,callfunc{
				base:   float64(optt.Strike_price),
				cost:   optt.Vw,
				factor: 1,
				date:   dateInt,
			})

			callListMap[optt.Expiration_date] = tmpp
		}

	}

	return optionsDates,optionsMap,callListMap

}

func MathematicaXYPlot(x,y []float64) string {
	code := "x={"
	for i,xx := range x {
		if i!=0 {
			code += ","
		}
		code += fmt.Sprintf("%.0f",xx)
	}
	code += "};\n"

	code += "y={"
	for i,yy := range y {
		if i!=0 {
			code += ","
		}
		code += fmt.Sprintf("%.0f",yy)
	}
	code += "};\n"
	code += "xy=ListPlot[Transpose[{x, y}], PlotStyle -> {AbsolutePointSize[8]},ImageSize -> Large, PlotRange -> Automatic,Joined -> True];\n"
	code += "Show[xy]\n"
	return code
}

//repeats the string s n times, returns as combined string
func repeatstr(s string, n int) string{
	result := ""
	for i := 0 ; i < n ; i++{
		result = result + s
	}
	return result
}

//Integral of f - []f equidistant (dx) values of f between a and b
func Integral(f []float64, dx float64) float64{
	area := 0.0
	for i:=0 ; i < len(f); i++ {
		area += f[i] * dx
	}
	return area
}

func MVProduct(M [][]float64, V []float64) []float64{
	if len(M)<1{
		return nil
	}
	if len(M[0])!=len(V){
		fmt.Errorf("Incompatible dimensions")
		return nil
	}
	result := make([]float64,len(V))

	for i := 0 ; i < len(result) ; i++ {
		temp := 0.0
		for d := 0 ; d < len(M[0]) ; d++ {
			temp += M[i][d]*V[d]
		}
		result[i] = temp
	}

	return result
}

func MMProduct(M [][]float64, N [][]float64) [][]float64{
	if len(M[0])!=len(N){
		fmt.Errorf("Incompatible dimensions")
		return nil
	}
	result := make([][]float64,len(M))
	for i := 0 ; i < len(result) ; i++{
		result[i] = make([]float64,len(N[0]))
	}
	temp := 0.0
	for i := 0 ; i < len(result) ; i++ {
		for j := 0 ; j < len(result[0]) ; j++{
			temp = 0
			for d := 0 ; d < len(N) ; d++ {
				temp += M[i][d]*N[d][j]
			}
			result[i][j] = temp
		}
	}
	return result
}

//negates a []float64 list (*(-1)) component wise
func floatlist_negation_compwise(f []float64) []float64{
	ff := make([]float64,len(f))
	for i := 0 ; i < len(ff) ; i++ {
		ff[i] = - f[i]
	}
	return ff
}

//concatenates two []float64 lists
func floatlist_cat(f1 []float64, f2 []float64) []float64{
	result := make([]float64,len(f1)+len(f2))
	for i := 0 ; i < len(f1) ; i++ {
		result[i] = f1[i]
	}
	for i := len(f1) ; i < len(result) ; i++ {
		result[i] = f2[i-len(f1)]
	}
	return result
}

//returns true if the list contains a string that contains the search string
func containscontains (list []string, search string) bool{
	for _,i := range list{
		if strings.Contains(i,search){
			return true
		}
	}
	return false
}

//returns true if the list contains a string that is equal to the search string
func contains(list []string, search string) bool{
	for _,i := range list{
		if i==search{
			return true
		}
	}
	return false
}

//1*2*3*...*n
func factorial(n int) int{
	if n == 1 || n == 0{
		return 1
	}
	return n*factorial(n-1)
}

//calculates root of ax^2+bx+c
func pq_formula_plus(a float64, b float64, c float64) float64{
	p := b/a
	pp := p/2
	q := c/a
	return -pp+math.Sqrt(pp*pp-q)
}

func f_x2 (x float64) float64{return x*x}

func max(x []float64) float64{
	if len(x) == 0{
		return 0
	}
	max := x[0]
	for _, a := range x{
		if a > max{
			max = a
		}
	}
	return max
}

func min(x []float64) float64{
	if len(x) == 0{
		return 0
	}
	min := x[0]
	for _, a := range x{
		if a < min{
			min = a
		}
	}
	return min
}


// ------------------------------- LGS specific functions -------------------------------


//Solves LGS
func (M LGS) GaussElimination() bool {

	//M := A
	inmethodprint := false

	//identity matrix
	id := make([][]float64,M.n)
	for i := 0 ; i < M.n ; i++{
		id[i] = make([]float64,M.n)
		id[i][i] = 1
	}
	//zero vector
	zerovector := make([]float64,M.n)

	exchanges := make([][]float64,M.n)
	copy(exchanges, id)
	exchangesM := LGS{M.n,exchanges,zerovector}
	/*//not supported yet
	inverse := make([][]float64,M.n)
	copy(inverse,id)
	*/

	unique := true

	//clear lower left triangle
	for i := 0 ; i < M.n ; i ++ {
		if inmethodprint{
			fmt.Println("In method print:")
			M.Print()
		}


		if M.A[i][i] == 0 {
			//check for row exchanges
			for j := i+1 ; j < M.n ; j++ {
				if M.A[j][i] != 0{
					//exchange rows (I) and (J)
					if inmethodprint{
						fmt.Println("In method print: (",i,")<->(",j,")")
					}
					M.RowExchange(i,j)
					exchangesM.RowExchange(i,j)
					if inmethodprint{
						M.Print()
					}
					break
				}
			}
		}

		if M.A[i][i] != 0 {
			if inmethodprint{
				fmt.Println("In method print: (",i,")*=",1/M.A[i][i])
			}
			M.MultiplyRow(i , 1/M.A[i][i])
			if inmethodprint{
				M.Print()
			}
		} else{
			continue
		}

		for j := i+1 ; j < M.n ; j++ {
			if M.A[j][i] != 0 {
				if inmethodprint{
					fmt.Println("In method print: (",j,")+=",-M.A[j][i],"*(",i,")")
				}
				M.OnePlusEqualsXxTwo(j,-M.A[j][i],i)
				if inmethodprint{
					M.Print()
				}
			}
		}
	}

	//clear upper right triangle
	for i := M.n-1 ; i >= 0 ; i-- {
		if M.A[i][i] == 0 {
			unique = false
			M.A[i][i]=-1.0
			continue
		} else {
			for j := i-1 ; j >= 0 ; j-- {
				if M.A[j][i] != 0 {

					if inmethodprint {
						fmt.Println("In method print: (",j,")+=",-M.A[j][i],"*(",i,")")
					}
					M.OnePlusEqualsXxTwo(j,-M.A[j][i],i)
					if inmethodprint{
						M.Print()
					}
				}
			}
		}
	}

	M.y = MVProduct(exchangesM.A,M.y)


	return unique

}

func (M LGS) AddRow(row int, addA []float64, addFrom int, addY float64){
	for col := addFrom; col < addFrom+len(addA) ; col++ {
		M.A[row][col] = M.A[row][col-addFrom] + addA[col-addFrom]
	}
	M.y[row] = M.y[row] + addY
}

func (M LGS) SetRow(row int, setA []float64, setFrom int, setY float64){
	for col := setFrom; col < setFrom + len(setA) ; col++ {
		M.A[row][col] = setA[col - setFrom]
	}
	M.y[row] = M.y[row] + setY
}

//(1)<->(2)
func (A LGS) RowExchange(one int, two int){
	var temp float64
	for col := 0 ; col < A.n ; col++ {
		temp = A.A[one][col]
		A.A[one][col] = A.A[two][col]
		A.A[two][col] = temp
	}
	temp = A.y[one]
	A.y[one] = A.y[two]
	A.y[two] = temp
}

//(row)*=factor
func (A LGS) MultiplyRow(row int, factor float64){
	for col := 0 ; col < A.n ; col++{
		A.A[row][col] = factor * A.A[row][col]
	}
	A.y[row] = factor * A.y[row]
}

// (1)+=x*(2)
func (A LGS) OnePlusEqualsXxTwo (one int, X float64, two int){
	for col:=0 ; col < A.n ; col++{
		A.A[one][col]=A.A[one][col]+X*A.A[two][col]
	}
	A.y[one] = A.y[one]+X*A.y[two]
}

func (A LGS) Print(){
	fmt.Println("LGS Print")
	for i:=0 ; i < A.n ; i++ {
		for j:=0 ; j < A.n ; j++{
			fmt.Printf("%.1f"+"  ",A.A[i][j])
		}
		fmt.Printf("|  "+"%.1f",A.y[i])
		fmt.Println("")
	}
}





// ------------------------------- old spline specific functions -------------------------------
/*

func (ns normalized_spline_old) Integral(a float64, b float64, precision float64) float64{
	dx := precision*(b-a)
	var f = make([]float64,int(1.0/dx))
	for i := 0 ; i < len(f) ; i++{
		f[i] = ns.At(a+float64(i)*dx)
	}
	return Integral(f, dx)

}


func (s spline_old) Integral(a float64, b float64, precision float64) float64{
	dx := precision*(b-a)
	var f = make([]float64,int(1.0/dx))
	for i := 0 ; i < len(f) ; i++{
		f[i] = s.At(a+float64(i)*dx)
	}
	return Integral(f, dx)
}


func (s spline_old) At (x float64) float64{
	return s.spline.At(x)
}


func (ns normalized_spline_old) At (x float64) float64 {
	if x<ns.xrange[1]&&x>ns.xrange[0]{
		return math.Max(0, ns.spline.At(x)/ns.integral)
	} else {
		return 0
	}
}

func normalize_spline (spline spline_old, xrange []float64) normalized_spline_old{
	return normalized_spline_old{
		spline: spline,
		xrange: xrange,
		integral: spline.Integral(xrange[0],xrange[1],0.001),
	}
}

func (ns normalized_spline_old) ASCIIPlot(res []int){
	max := ns.At(0)
	min := ns.At(0)
	precision := 100000.0
	for i := ns.xrange[0] ; i <= ns.xrange[1] ; i+=1/precision{
		if ns.At(i)>max{
			max = ns.At(i)
		}
		if ns.At(i)<min{
			min = ns.At(i)
		}
	}

	fmt.Println("min:",min," max:",max)

	fmt.Print(repeatstr(" ",4))
	for i:=0 ; i<res[1]/4 ; i++{
		fmt.Printf("%.1f ",min+float64(4*i)*(max-min)/float64(res[1]))
	}
	fmt.Println("")

	var xi float64
	var yi float64
	for i := 0 ; i <= res[0] ; i++{
		xi =float64(i)*(ns.xrange[1]-ns.xrange[0])/float64(res[0])
		yi = (ns.At(xi)-min)/(max-min)*float64(res[1])
		fmt.Printf("%.2f:",xi)
		fmt.Println(repeatstr("#", int(yi)+1 ))
		//fmt.Println(xi,":",yi)
	}
}

*/
