package main

import (
	"math/rand"
	"os"
	"sort"

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
	//opt "github.com/AaronGonsior/polygonioClient"
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

type spread struct{
	num int
	calls []callfunc
	weights []float64
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


//goland:noinspection ALL
func main(){

	riskAndTimePlottesting := false
	optimalTransporttesting := false
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
			cumSpline := pdistSplines[d].IntegrateDUMB()
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

	if apitesting {

		/* User Inputs */
		update := false
		ticker := "TSLA"

		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		var pdistX map[string][]float64 = make(map[string][]float64,0)
		var pdistY map[string][]float64 = make(map[string][]float64,0)
		var pdistDates []string

		pdistDates = append(pdistDates,"2024-06-01")
		pdistX["2024-06-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  , 600,1000}
		pdistY["2024-06-01"] = []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1    , 0.2  ,0}


		pdistDates = append(pdistDates,"2025-01-01")
		pdistX["2025-01-01"] = []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  , 600,1000}
		pdistY["2025-01-01"] = []float64{0, 2	, 5	, 7	  , 15	, 17   , 17   , 15   , 12   , 10   , 7   , 5     , 1  ,0}

		var pdistSplines map[string]my_spline
		pdistSplines = make(map[string]my_spline,0)

		for _,d := range pdistDates {
			//fmt.Println(pdistX[d],pdistY[d])
			s := NewSpline(splinetype,pdistX[d],pdistY[d])
			ns := NewNormedSpline(s)
			pdistSplines[d] = ns
		}

		apiKey := opt.LoadJson("apiKey.json")

		var optreq opt.OptionURLReq
		var options []opt.Option

		optreq = opt.OptionURLReq{
			Ticker:      ticker,
			ApiKey:      apiKey,
			StrikeRange: []int{0,1000},
			DateRange:   /*[]string{"2024-06-01","2024-07-01"}*/[]string{"2023-06-01","2027-01-01"},
			Contract_type: "call",
		}


		/* End User Inputs */

		content := "SetDirectory[NotebookDirectory[]]\n"

		url := "https://api.polygon.io/v2/aggs/ticker/C:USDEUR/prev?adjusted=true&apiKey="+apiKey
		fmt.Println("url: ",url)
		_,body,err := opt.APIRequest(url,1)
		check(err)
		body = strings.Split(body,"\"c\":")[1]
		body = strings.Split(body,",")[0]
		fmt.Println(body)

		usdtoeur,err = strconv.ParseFloat(body,64)
		check(err)
		eurtousd = 1/usdtoeur


		var share_price float64
		url = "https://api.polygon.io/v2/aggs/ticker/"+ticker+"/prev?adjusted=true&apiKey="+apiKey
		fmt.Println("url: ",url)
		_,body,err = opt.APIRequest(url,1)
		check(err)
		body = strings.Split(body,"\"c\":")[1]
		body = strings.Split(body,",")[0]

		share_price,err = strconv.ParseFloat(body,64)
		check(err)
		fmt.Println("share_price(",ticker,"): ",share_price)


		//how many successive requests at most; -1 is Inf
		nMax := -1

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


		debug := true

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


		/*
		//25Q1
		x := []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		y := []float64{0, 2	, 5	, 7	  , 15	, 17   , 17   , 15   , 12   , 10   , 7   , 5     }


		//24Q2
		x := []float64{0, 25, 50, 100 , 150	, 200  , 250  , 300  , 350  , 400  , 450  , 500  }
		y := []float64{0, 2	, 6	, 7	  , 15	, 17   , 15   , 12   , 8    , 6    , 3    , 1     }


		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		s := NewSpline(splinetype,x,y)
		ns := NewNormedSpline(s)
		pdist := ns
		 */




		mathCode := "SetDirectory[NotebookDirectory[]]\n"
		mathCodeSigma := "SetDirectory[NotebookDirectory[]]\n"
		//dx := 0.01

		path, err := os.Getwd()
		check(err)
		fmt.Println(path)
		currentTime := time.Now()
		live := currentTime.Format("2006-01-02")

		err = os.Mkdir(path+"\\tmp\\"+live, 0755)
		path = path+"\\tmp\\"+live+"\\"

		var strikes []float64
		var costs []float64

		mathematicaCompressionLevel := ".75"
		mathematicaImageResolution := "250"

		for _,d := range optionsDates {

			pdist := pdistSplines[pdistDates[0]] //careful: date should eventually be optimal transported.

			folderName := ticker+d+"(live data from "+live+")"
			err = os.Mkdir(path+folderName, 0755)
			check(err)
			optionsList := optionsMap[d]
			callList := callListMap[d]
			callList = callList[len(addToAll):len(callList)]

			fmt.Println(optionsList)

			tmp,idPdist := pdist.PrintMathematicaCode()
			mathCode = tmp
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\pdist.png\", " + fmt.Sprintf("Show[fctplot%v]",idPdist) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			bestcall, bestE := findBestCall(pdist, callList)
			fmt.Println("Best Call:", bestcall, "\nwith expected return:", bestE)
			mathCode = bestcall.PrintMathematicaCode()
			fmt.Println(mathCode)

			content += fmt.Sprintf("msg1 := Text[\"Assuming the probability distribution (left) for the date %v, the call with strike %.1f has the highest expected return out of all calls options available with %.1f %% expected return. Owning the underlying asset (%v) has an expected return of %.1f %%.  \"];\n\n", callList[0].date, bestcall.base, bestE, ticker, long.ExpectedReturn(pdist))
			content += mathCode
			content += "Export[\"" + folderName + "\\-bestCall.png\", {msg1 \n , "+fmt.Sprintf("Show[fctplot%v]",idPdist) +", Show[call,long]}, \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("owning $TSLA has an expected return of: ", long.ExpectedReturn(pdist))

			fmt.Println("\nPrint all calls:\n")
			mathCode = PrintMathematicaCode(callList, share_price)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\allCalls.png\", Show[s], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("\nDistribution Chart for Call-Long intersections:\n")
			mathCode = MathematicaCodeLongIntersection(callList, share_price)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallLongIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("\nDistribution Chart for Call-Zero intersections:\n")
			mathCode = MathematicaCodeZeroIntersection(callList)
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallZeroIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("\nDistribution Chart for Call-Zero-Volumes intersections:\n")
			mathCode = MathematicaCodeZeroIntersectionVolumes(optionsList)
			content += mathCode
			content += "Export[\"" + folderName + "\\CallZeroVolumesIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("\nExpected returns for each strike:")
			mathCode = MathematicaPrintExpectedReturns(pdistSplines[pdistDates[0]], callList) //careful: date should eventually be optimal transported.
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\expected_returns_strike.png\", Show[xy], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			strikes = make([]float64,0)
			costs = make([]float64,0)
			for _, opt := range optionsMap[d] {
				strikes = append(strikes, float64(opt.Strike_price))
				costs = append(costs, (opt.Close))
			}
			mathCode = MathematicaXYPlot(strikes, costs)
			fmt.Println("\nPlot strike vs cost:\n")
			fmt.Println(mathCode)
			content += mathCode
			content += "Export[\"" + folderName + "\\strike_price.png\", Show[xy], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			//still causes some indexing bug


			bestCallSpline := bestcall.ToSpline(min(pdist.x),max(pdist.x))
			tmp,id := bestCallSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_bestCallSpline.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			//only testing
			/*
			call1, call2 := UnionXYCC(callList[1].ToSpline(min(pdist.x),max(pdist.x)),callList[0].ToSpline(min(pdist.x),max(pdist.x)))

			tmp,id = call1.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_call1.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			tmp,id = call2.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_call2.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			callmult := call1.SplineMultiply(call2)
			tmp, id = callmult.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_callMult.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			 */

			//more testing with constants
			const1 := my_spline{
				deg:        1,
				splineType: splinetype,
				x:          []float64{0,100,200},
				y:          []float64{2,2,1},
				coeffs:     []float64{0,2,0.01,1},
				unique:     false,
			}
			const2 := my_spline{
				deg:        0,
				splineType: splinetype,
				x:          []float64{0,100,300},
				y:          []float64{3,3,4},
				coeffs:     []float64{3,4},
				unique:     false,
			}
			fmt.Println("constant splines before UnionXYCC():")
			fmt.Println("len(const1.x)=",len(const1.x)," ; len(const2.x)=",len(const2.x))
			fmt.Println("len(const1.coeffs)=",len(const1.coeffs)," ; len(const2.coeffs)=",len(const2.coeffs))


			/*
			tmp,id = const1.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_const1.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			tmp,id = const2.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_const2.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			const1, const2 = UnionXYCC(const1, const2)

			callmult := const1.SplineMultiply(const2)
			tmp, id = callmult.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_constMult.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			const1Integrate := const1.Integrate()
			tmp, id = const1Integrate.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_const1Integrate.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			fmt.Println("constMult.FullIntegralSpline()=",callmult.FullIntegralSpline())
			 */







			//bestCallSpline, pdist = UnionXYCC(bestCallSpline,pdist)
			/*
			pdist,bestCallSpline = UnionXYCC(pdist,bestCallSpline)

			tmp,id = bestCallSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_bestCallSpline.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			tmp,id = pdist.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\TEST_Unionized_pdist.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			 */


			fmt.Print("probReturn...")
			probReturn := pdist.SplineMultiply(long.ToSpline(0,max(pdist.x)))
			tmp,id = probReturn.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\probReturn.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			fmt.Println(" done.")

			fmt.Print("probReturnIntegral...")
			probReturnIntegral := probReturn.Integrate()
			tmp,id = probReturnIntegral.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\probReturnIntegral.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			fmt.Println(" done.")

			fmt.Print("pdistIntegral...")
			pdistIntegrate := pdist.Integrate()
			tmp,id = pdistIntegrate.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\pdistIntegrate.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			fmt.Println(" done.")

			//risk evaluation
			fmt.Print("riskSpline...")
			eps := 0.1
			yMin := min(probReturnIntegral.y) +eps
			yMax := max(probReturnIntegral.y) -eps
			dy := 10.0
			var ys,probs []float64
			for y := yMin ; y <= yMax ; y += dy {
				fmt.Print(y," ")
				ys = append(ys,y)
				probs = append(probs, pdist.SplineMultiply(probReturnIntegral.OneBelow(y)).FullIntegralSpline() )
			}
			fmt.Print(" ... ")
			riskSpline := NewSpline(splinetype,probs,ys)
			tmp,id = riskSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\riskSpline.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			fmt.Println(" done.")



			//careful: sometimes out of memory
			var spreads []spread
			//only neigboring calls
			
			for i := 0 ; i < len(callList)-1 ; i++ {
				spreads = append(spreads,spread{
					num:     2,
					calls:   []callfunc{callList[i],callList[i+1]},
					weights: []float64{0.5,-0.5},
				})
				spreads = append(spreads,spread{
					num:     2,
					calls:   []callfunc{callList[i],callList[i+1]},
					weights: []float64{-0.5,0.5},
				})
			}


			/*
			// all 2-combinations of calls and 50-50% weighing in both directions (buy&sell)
			ws := []float64{0.5,0.75,0.25}
			for i := 0 ; i < len(callList)-1 ; i++ {
				for j := i ; j < len(callList)-1 ; j++{
					for _,w := range ws {
						spreads = append(spreads,spread{
							num:     2,
							calls:   []callfunc{callList[i],callList[j]},
							weights: []float64{w,-(1.0-w)},
						})
						spreads = append(spreads,spread{
							num:     2,
							calls:   []callfunc{callList[i],callList[j]},
							weights: []float64{w,1.0-w},
						})
						spreads = append(spreads,spread{
							num:     2,
							calls:   []callfunc{callList[i],callList[j]},
							weights: []float64{-w,(1.0-w)},
						})
					}
				}
			}
			 */


			bestSpread,bestSpreadExp := FindBestSpread(pdist,spreads)
			/*
			// Expected returns of all spreads
			dx = 0.1
			var SpreadsExpReturns []float64
			for i := range spreads {
				SpreadsExpReturns = append(SpreadsExpReturns,spreads[i].ExpectedReturn(pdist,dx))
			}

			// Find spread with highest Exp Return* (Later include risk)
			var bestSpread spread = spreads[0]
			var bestSpreadExp float64 = SpreadsExpReturns[0]
			for i,spExp := range SpreadsExpReturns[1:] {
				if spExp > bestSpreadExp {
					bestSpread = spreads[i]
					bestSpreadExp = spExp
				}
			}
			 */

			// make all spreads to splines
			var spreadSplines []my_spline
			for _,s := range spreads {
				spreadSplines = append(spreadSplines,s.ToSpline(0,max(pdist.x)))
			}

			// make best spread to spline
			bestSpreadSpline := bestSpread.ToSpline(0,max(pdist.x))

			// make a .png for every spread
			/*
			err = os.Mkdir(path+folderName+"/spreads", 0755)
			for i,s := range spreadSplines {
				tmp,id = s.PrintMathematicaCode()
				mathCode += tmp+"\n"
				mathCode += "Export[\"" + folderName+"\\spreads" + "\\"+strconv.Itoa(i)+".png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
				content += mathCode
			}
			 */

			//make a .png for bestSpread
			content += fmt.Sprintf("msg1 := Text[\"Assuming the probability distribution (left) for the date %v, the spread %v has the highest expected return out of all spreads (2-50/50) available with %.1f %% expected return. Owning the underlying asset (%v) has an expected return of %.1f %%.  \"];\n\n", callList[0].date, fmt.Sprint(bestSpread), bestSpreadExp, ticker, long.ExpectedReturn(pdist))
			tmp,id = bestSpreadSpline.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName+"\\-bestSpread.png\",{msg1,"+fmt.Sprintf("Show[fctplot%v]",idPdist)+"," + fmt.Sprintf("Show[{s%v,long}]",id) + "}, \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode

			fmt.Sprintf("Show[fctplot%v]",idPdist)


			/*
			integralProbReturn := probReturn.IntegrateDUMB()
			tmp,id = integralProbReturn.PrintMathematicaCode()
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\IntegralProbReturn.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
			 */



			//FindSigmas
			/*
			levels := []float64{0,0.125,0.25,0.5,0.75,0.875,1}

			cumSpline := pdist.IntegrateDUMB()
			tmp, id = cumSpline.PrintMathematicaCode()
			mathCodeSigma += tmp+"\n"
			mathCodeSigma += fmt.Sprintf("s%v\n",id)
			sigmas := pdist.FindSigmas(levels)

			for i,s := range sigmas {
				fmt.Println("expected return at ", levels[i]*100, "% : ", bestcall.At(s))
			}

			 */

			/*
			for risk metric, use %(0-100) where the investment is breakeven. For that, implement my_spline Multiply() to
			multiply pdist and call and use NewtonRoot()
			 */


			a:=0.0;b:=1000.0;
			fmt.Println(fmt.Sprintf("pdist.IntegralSpline(%v,%v)=",a,b),pdist.IntegralSpline(a,b))

		}



		WriteFile("sigmas.nb",mathCodeSigma,"/")




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

		//dx := 0.0001

		//fmt.Println(s.Integral(min(x),max(x),dx))

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
		bestcall, bestE := findBestCall(pdist, callList)
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



		/*
		dx := 0.0001
		area := s.Integral(min(x),max(x),dx)

		fmt.Println("Integral:", area)


		ns := NewNormedSpline(s)

		nsarea := ns.Integral(min(x),max(x),dx)
		fmt.Println("nsIntegral:", nsarea)

		tmp,_ = ns.PrintMathematicaCode()
		mathCode = tmp
		fmt.Println(mathCode)
		 */


	}



}




// ------------------------------- spread specific functions -------------------------------

func (sp spread) ToSpline(a,b float64) my_spline {
	result := sp.calls[0].ToSpline(a,b).Factor(sp.weights[0])
	for s := 1 ; s < len(sp.calls) ; s++ {
		result = result.Add(sp.calls[s].ToSpline(a,b).Factor(sp.weights[s]))
	}
	return result
}

func (sp spread) ExpectedReturn(pdist my_spline) float64 {
	var expReturns float64
	for i,c := range sp.calls {
		expReturns += sp.weights[i]*c.ExpectedReturn(pdist)
	}
	return expReturns
}

func FindBestSpread(pdist my_spline, spreads []spread) (spread,float64) {
	// Expected returns of all spreads
	//dx := 0.1
	var SpreadsExpReturns []float64
	for i := range spreads {
		SpreadsExpReturns = append(SpreadsExpReturns,spreads[i].ExpectedReturn(pdist))
	}
	var bestSpread spread = spreads[0]
	var bestSpreadExp float64 = SpreadsExpReturns[0]
	for i,spExp := range SpreadsExpReturns[1:] {
		if spExp > bestSpreadExp {
			bestSpread = spreads[i]
			bestSpreadExp = spExp
		}
	}
	return bestSpread,bestSpreadExp
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

func (ms my_spline) PrintMathematicaCode() (string,string){


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
	for i := 0 ; i < (ms.deg+1)*(len(ms.x)-1) ; i += ms.deg+1 {
		result += fmt.Sprint("{")
		for d := ms.deg ; d >= 0 ; d-- {
			if ms.coeffs[i+(ms.deg-d)] >= 0 {
				result += fmt.Sprint("+")
			}
			result += fmt.Sprintf("%.20fx^%v",ms.coeffs[i+(ms.deg-d)],d)
		}
		result += fmt.Sprint(",")
		result += fmt.Sprintf("%.3f",ms.x[i/(ms.deg+1)])
		result += fmt.Sprint("<=x<=")
		result += fmt.Sprintf("%.3f",ms.x[i/(ms.deg+1)+1])
		result += fmt.Sprint("}")
		if i<(ms.deg+1)*(len(ms.x)-1)-(ms.deg+1) {
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

/*
func (ms my_spline) PrintMathematicaCode() (string,string) {

	if ms.deg != 3 {
		fmt.Println("my_spline.PrintMathematicaCode: for this print ms.deg has to be 3 atm. Degree here is ",ms.deg," still going to try.")
	}

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
 */

func (ms my_spline) At (x float64) float64{
	//debug := true

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
	if (ms.deg+1)*(splineNr+1)+1 < len(coeffs) {
		coeffs = coeffs[(ms.deg+1)*(splineNr):(ms.deg+1)*(splineNr+1)+1]
	} else {
		coeffs = coeffs[(ms.deg+1)*(splineNr):]
	}

	if len(coeffs) == 0 {
		return 0
	}

	result := 0.0
	for deg := 0 ; deg <= ms.deg ; deg++ {
		result += coeffs[deg]*math.Pow(x,float64(ms.deg-deg))
	}
	return result
}

/*
func (ms my_spline) Integral(a float64, b float64, dx float64) float64{
	var err error
	f := make([]float64,int((b-a)/dx))
	for i := 0 ; i < len(f) ; i++ {
		f[i] = ms.At(a+float64(i)*dx)
		check(err)
	}
	return Integral(f, dx)
}
 */

func (ms my_spline) IntegralSplineOld(a,b float64) float64 {
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
		newCoeffs = append(newCoeffs,ms.coeffs[(ms.deg+1)*j+d])
	}
	for i := j ; ms.x[i] < b && i < len(ms.x)-1 ; i++ {
		newX = append(newX, ms.x[i])
		newY = append(newY, ms.At(ms.x[i]))
		for d := 0 ; d < ms.deg+1 ; d++ {
			newCoeffs = append(newCoeffs,ms.coeffs[(ms.deg+1)*i+d])
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

func (ms my_spline) IntegralSpline(a,b float64) float64 {
	integral := 0.0
	i1 := 0
	i2 := len(ms.x)-1
	for ms.x[i1] < a {i1++};
	for ms.x[i2] > b {i2--};
	if i1 < 1 {i1=1}
	if i2 > len(ms.x)-2 {i2=len(ms.x)-2}


	//a-ms.x[i1]
	for d := 0 ; d <= ms.deg ; d++ {
		integral += (ms.coeffs[(ms.deg+1)*(i1-1)+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i1],float64(ms.deg-d)+1)
		integral -= (ms.coeffs[(ms.deg+1)*(i1-1)+d]/(float64(ms.deg-d)+1)) * math.Pow(a,float64(ms.deg-d)+1)
	}

	//ms.x[i2]-b
	for d := 0 ; d <= ms.deg ; d++ {
		integral += (ms.coeffs[(ms.deg+1)*i2+d]/(float64(ms.deg-d)+1)) * math.Pow(b,float64(ms.deg-d)+1)
		integral -= (ms.coeffs[(ms.deg+1)*i2+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i2],float64(ms.deg-d)+1)
	}


	for i := i1 ; i <= i2 ; i++ {
		for d := 0 ; d <= ms.deg ; d++ {
			integral += (ms.coeffs[(ms.deg+1)*i+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i+1],float64(ms.deg-d)+1)
			integral -= (ms.coeffs[(ms.deg+1)*i+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i],float64(ms.deg-d)+1)
		}
	}
	return integral
}

//has some bug especially regarding coeffs
func UnionXYCC (ms1, ms2 my_spline) (my_spline , my_spline) {

	debug := false

	if isUnionized(ms1,ms2){
		return ms1,ms2
	}


	newX := []float64{}

	i1 := 0
	i2 := 0
	//min := math.Max(min(ms1.x),min(ms2.x))
	//max := math.Min(max(ms1.x),max(ms2.x))
	for i := 0 ; i < len(ms1.x)+len(ms2.x) ; i++ {

		//if outside of shared definition space, stop
		/*
		if i1 < len(ms1.x) && (ms1.x[i1] < min || ms1.x[i1] > max) {i++;continue}
		if i2 < len(ms2.x) && (ms2.x[i2] < min || ms2.x[i2] > max) {i++;continue}
		 */

		//if x1's are at the end, add all x2's that are not already in there
		eps := 0.01
		if i1 >= len(ms1.x) {
			for i2 < len(ms2.x) {
				if !containsFloat(newX,ms2.x[i2],eps){
					newX = append(newX,ms2.x[i2])
				}
				i2++
				//i++
			}
			break
		}

		//if x2's are at the end, add all x1's that are not already in there
		if i2 >= len(ms2.x) {
			for i1 < len(ms1.x) {
				if !containsFloat(newX,ms1.x[i1],eps){
					newX = append(newX,ms1.x[i1])
				}
				i1++
				//i++
			}
			break
		}

		if ms1.x[i1] < ms2.x[i2] {
			if !containsFloat(newX,ms1.x[i1],eps){
				newX = append(newX,ms1.x[i1])
			}
			i1++
		} else {
			if !containsFloat(newX,ms2.x[i2],eps){
				newX = append(newX,ms2.x[i2])
			}
			i2++
		}


		/*
		if !containsFloat(newX,x) {
			newX = append(newX,x)
		}

		 */
	}

	if debug {
		fmt.Println("len(newX): ",len(newX))
	}

	degMax := int(math.Max( float64(ms1.deg) , float64(ms2.deg) ))

	if debug{
		fmt.Println("degMax: ", degMax)
	}

	//add new Y's
	var newY1, newY2 []float64
	for _,nx := range newX {
		newY1 = append(newY1,ms1.At(nx))
		newY2 = append(newY2,ms2.At(nx))
	}

	var newC1, newC2 []float64
	for i,nx := range newX {
		i1 := 0
		i2 := 0

		//increase i1,i2 s.t. ms1.x and ms2.x are just under nx
		for ms1.x[i1] <= nx && i1<len(ms1.x)-1 {i1++};i1--
		for ms2.x[i2] <= nx && i2<len(ms2.x)-1 {i2++};i2--
		if i1<0{i1=0}
		if i2<0{i2=0}

		if debug {
			fmt.Println("For the first spline, index ", i1, "is just small enough s.t. ms1.x[i1]<nx: ( ",ms1.x[i1],"<",nx," ; ms1.x[i1+1]=",ms1.x[i1+1],")")
			fmt.Println("For the second spline, index ", i2, "is just small enough s.t. ms2.x[i2]<nx: ( ",ms2.x[i2],"<",nx," ; ms2.x[i2+1]=",ms2.x[i2+1],")")
		}

		//ms1 update
		//add new C
		if i != len(newX)-1 {
			for j := 0 ; j < degMax-ms1.deg ;  j++ {
				newC1 = append(newC1,0.0)
			}
			for j := 0 ; j <= ms1.deg ; j++ {
				newC1 = append(newC1,ms1.coeffs[i1*(ms1.deg+1)+j])
			}
		}


		/*
		for j := i1 * (degMax+1) ; j < (i1+1) * (degMax+1) ; j++ {
			if degMax > ms1.deg + j {
				newC1 = append(newC1,0.0)
			} else {
				newC1 = append(newC1,ms1.coeffs[i1-degMax+ms1.deg])
			}
		}
		 */


		/*
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
			for j := i * ms1.deg ; j < i * (ms1.deg+1) ; j++ {
				newC1 = append(newC1,ms1.coeffs[j])
			}
		}
		 */

		//ms2 update
		//add new C
		if i != len(newX)-1 {
			for j := 0 ; j < degMax-ms2.deg ;  j++ {
				newC2 = append(newC2,0.0)
			}
			for j := 0 ; j <= ms2.deg ; j++ {
				newC2 = append(newC2,ms2.coeffs[i2*(ms2.deg+1)+j])
			}
		}

		/*
		for j := i2 * (degMax+1) ; j < (i2+1) * (degMax+1) ; j++ {
			if degMax > ms2.deg + j {
				newC2 = append(newC2,0.0)
			} else {
				newC2 = append(newC2,ms2.coeffs[j-degMax+ms2.deg])
			}
		}
		 */


		/*
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
		 */
	}

	newms1 := my_spline{
		deg:        degMax,
		splineType: ms1.splineType,
		x:          newX,
		y:          newY1,
		coeffs:     newC1,
		unique:     false,
	}
	newms2 := my_spline{
		deg:        degMax,
		splineType: ms2.splineType,
		x:          newX,
		y:          newY2,
		coeffs:     newC2,
		unique:     false,
	}

	//check
	if !isUnionized(newms1,newms2){
		fmt.Errorf("bug in UnionXYCC! Not unionized at the end.")
		fmt.Println("bug in UnionXYCC! Not unionized at the end.")
		os.Exit(1)
	}

	if debug {
		fmt.Println("len(newX)=",len(newX)," , len(newC1)=",len(newC1)," , len(newC2)=",len(newC2))
	}


	return newms1,newms2

}

func isUnionized (ms1, ms2 my_spline) bool {
	if len(ms1.x) != len(ms2.x){return false}
	for i := range ms1.x {
		if ms1.x[i] != ms2.x[i]{
			//fmt.Println("isUnionized check: len(ms1.x), len(ms2.x): ",len(ms1.x),len(ms2.x))
			return false
		}
	}
	if len(ms1.coeffs) != len(ms2.coeffs){
		//fmt.Println("isUnionized check: len(ms1.coeffs), len(ms2.coeffs): ",len(ms1.coeffs),len(ms2.coeffs))
		return false
	}
	return true
}

// multiplies two splines into the product spline
func (ms1 my_spline) SplineMultiply(ms2 my_spline) my_spline {

	debug := false

	ms1, ms2 = UnionXYCC(ms1, ms2)
	degSum := ms1.deg + ms2.deg
	newC := []float64{}

	for iSp := 0 ; iSp < len(ms1.x)-1 ; iSp++ {
		//if iSp == len(ms1.x)-1{break}

		//for each segment (between x's)
		for d := degSum ; d >= 0 ; d-- {
			tmp := 0.0
			for d1 := 0 ; d1 <= ms1.deg && d1 <= d ; d1++ {
				d2 := d - d1
				if d2 > ms2.deg {
					continue
				}
				//if d2 < 0 {break}
				tmp += ms1.coeffs[iSp*(ms1.deg+1)+ms1.deg-d1] * ms2.coeffs[iSp*(ms2.deg+1)+ms2.deg-d2]
				if debug{
					fmt.Println("SplineMultiply(): For degree ",d," with degree ",d1," from ms1 and degree ",d2," from ms2, use coeffs ",ms1.coeffs[iSp*(ms1.deg+1)+d1], " and ", ms2.coeffs[iSp*(ms2.deg+1)+d2], "to add ",ms1.coeffs[iSp*(ms1.deg+1)+d1] * ms2.coeffs[iSp*(ms2.deg+1)+d2]," to make the coeff for degree ",d," ",tmp)
				}
			}
			newC = append(newC, tmp)
			if debug {
				fmt.Println("SplineMultiply(): addded in newC: ",tmp)
			}
		}
	}

	msMult := my_spline{
		deg:        degSum,
		splineType: ms1.splineType,
		x:          ms1.x,
		y:          []float64{},
		coeffs:     newC,
		unique:     false,
	}

	newY := []float64{}
	for _,x := range msMult.x {
		newY = append(newY,msMult.At(x))
	}

	msMult = my_spline{
		deg:        degSum,
		splineType: ms1.splineType,
		x:          ms1.x,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}

	if debug {
		fmt.Println("SplineMultiply(): len(msMult.x)=",len(msMult.x)," ; len(msMult.y)=",len(msMult.y)," ; len(msMult.coeffs)=",len(msMult.coeffs), " ; msMult.deg=",msMult.deg)
	}

	return msMult
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

	/*
	//careful at different degrees!
	if ms1.deg != ms2.deg {
		fmt.Println("unequal degrees in myspline.Add() not supported yet")
		return my_spline{}
	}
	 */

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
		deg:        int(math.Max(float64(ms1.deg),float64(ms2.deg))),
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
		for d := 0 ; d <= ms.deg ; d++ {
			integral += (ms.coeffs[(ms.deg+1)*i+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i+1],float64(ms.deg-d)+1)
			integral -= (ms.coeffs[(ms.deg+1)*i+d]/(float64(ms.deg-d)+1)) * math.Pow(ms.x[i],float64(ms.deg-d)+1)
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

//bugged
//return Integrated spline, one dim higher
func (ms my_spline) Integrate() my_spline {
	var newC []float64
	for i,_ := range ms.x[:len(ms.x)-1] {
		for d := ms.deg ; d >= 0 ; d-- {
			newC = append(newC,1.0/(float64(d+1))*ms.coeffs[(ms.deg+1)*i+ms.deg-d])
		}
		newC = append(newC,0.0)
	}
	integral := my_spline{
		deg:        ms.deg+1,
		splineType: ms.splineType,
		x:          ms.x,
		y:          []float64{},
		coeffs:     newC,
		unique:     false,
	}
	//find suitable C/constant s.t. integral is continuous
	newC = integral.coeffs
	for i,_ := range ms.x[:len(ms.x)-1] {
		newC[(i+1)*(integral.deg+1)-1] = ms.IntegralSpline(0,ms.x[i]) - integral.At(ms.x[i+1])
	}
	integral = my_spline{
		deg:        integral.deg,
		splineType: integral.splineType,
		x:          integral.x,
		y:          integral.y,
		coeffs:     newC,
		unique:     false,
	}
	var newY []float64
	for _,x := range integral.x {
		newY = append(newY,integral.At(x))
	}
	integral = my_spline{
		deg:        ms.deg+1,
		splineType: ms.splineType,
		x:          ms.x,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}
	return integral
}

//still to be tested
func (ms my_spline) IntegrateDUMB() my_spline {
	cumX := []float64{}
	fmt.Println("Test: FullIntegral: ",ms.FullIntegralSpline())
	for _,x := range ms.x {
		cumX = append(cumX,ms.IntegralSpline(min(ms.x),float64(x)))
		fmt.Printf("Test: Integral from 0 to %v : %v \n",x,ms.IntegralSpline(min(ms.x),float64(x)))
	}
	return NewSpline(ms.splineType,ms.x,cumX)
}

//return derivated spline, one dim lower
//not done.
func (ms my_spline) Derive() my_spline {
	var newC []float64
	for i,_ := range ms.x[:len(ms.x)-1] {
		for d := ms.deg ; d > 0 ; d-- {
			newC = append(newC,float64(d)*ms.coeffs[(ms.deg+1)*i+ms.deg-d])
		}
	}
	derivative := my_spline{
		deg:        ms.deg-1,
		splineType: ms.splineType,
		x:          ms.x,
		y:          []float64{},
		coeffs:     newC,
		unique:     false,
	}
	var newY []float64
	for _,x := range derivative.x {
		newY = append(newY,derivative.At(x))
	}
	derivative = my_spline{
		deg:        ms.deg-1,
		splineType: ms.splineType,
		x:          ms.x,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}
	return derivative
}

//returns the value of the derivative of ms at x
func (ms my_spline) D (x float64) float64 {
	//find ix s.t. ms.x[ix] just under x
	ix:=0;for ms.x[ix]<x {ix++};ix--;
	if ix<0{ix=0}
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

//returns 1_{ms<=y} as my_spline
func (ms my_spline) OneBelow (y float64) my_spline {
	debug := true
	intersections := ms.Intersections(y)
	if debug {
		fmt.Println(intersections)
	}
	newX := ms.x
	for _,ix := range intersections {
		if !containsFloat(newX,ix,1){
			newX = append(newX, ix)
		}
	}
	sort.Float64s(newX)
	newY := []float64{}
	newC := []float64{}
	for _,x := range newX {
		if ms.At(x) <= y {
			newY = append(newY,1.0)
			newC = append(newC, 1.0)
		} else {
			newY = append(newY,0.0)
			newC = append(newC, 0.0)
		}
	}
	result := my_spline{
		deg:        0,
		splineType: []string{""},
		x:          newX,
		y:          newY,
		coeffs:     newC,
		unique:     false,
	}
	return result
}


//not completed!
//still to be tested
func (ms my_spline) Intersections(y float64) []float64 {
	if ms.deg > 2 { //when cubic implemented >3
		//use NewtonRoots here
		return ms.NewtonRoots(y,0.01,10)
	}

	var intersections []float64

	if ms.deg == 1 {
		for i,_ := range ms.x[:len(ms.x)-1] {
			if (ms.x[i] > y && ms.x[i+1] < y) || ((ms.x[i] < y && ms.x[i+1] > y)) {
				//linear
				intersections = append(intersections,(y-ms.y[i])*(ms.x[i+1]-ms.x[i])/(ms.y[i+1]-ms.y[i]))
			}
		}
	}

	if ms.deg == 2 {
		for i,_ := range ms.x[:len(ms.x)-1] {
			if (ms.x[i] > y && ms.x[i+1] < y) || ((ms.x[i] < y && ms.x[i+1] > y)) {
				//squadratic
				coeffs := ms.coeffs[i*(ms.deg+1):i*(ms.deg+2)-1]
				if len(coeffs) != 3 {fmt.Println("Error: coeffs in my_spline.Intersections(float64) should be length 3 for squadratic!");os.Exit(2)}
				p := coeffs[1]/coeffs[0]
				q := coeffs[2]/coeffs[0]
				intersectionCandidates := []float64{p/2.0-math.Sqrt(math.Pow(p/2.0,2))-q ,p/2.0+math.Sqrt(math.Pow(p/2.0,2))-q }
				if intersectionCandidates[0] < ms.x[i+1] && intersectionCandidates[0] > ms.x[i] {
					intersections = append(intersections,  intersectionCandidates[0] )
				} else if intersectionCandidates[1] < ms.x[i+1] && intersectionCandidates[1] > ms.x[i] {
					intersections = append(intersections,  intersectionCandidates[1] )
				}


				//check which one (of the two) is in range [ms.x[i],ms.x[i+1]]
			}
		}
	}

	if ms.deg == 3 {
		for i,_ := range ms.x[:len(ms.x)-1] {
			if (ms.x[i] > y && ms.x[i+1] < y) || ((ms.x[i] < y && ms.x[i+1] > y)) {
				//cubic

				//check which one (of the three) is in range [ms.x[i],ms.x[i+1]]

			}
		}
	}

	return intersections

}


func (ms my_spline) NewtonRoots (y float64, tol float64, n int) []float64 {
	debug := true
	span := ms.x[0]-ms.x[len(ms.x)-1]
	dx := span / float64(n)
	var intersections []float64
	for x := ms.x[0] ; x < ms.x[len(ms.x)-1] ; x += dx {
		root,err := ms.NewtonRoot(x,y,tol)
		if err != nil {continue}
		if debug {
			fmt.Println("NewtonRoots: root:",root)
		}
		if !containsFloat(intersections,root,tol){
			intersections = append(intersections,root)
		}
	}
	return intersections
}

//For degrees <=3, should be replaced by pq or cubic root formula
//also do for result []float64
//finds roots (y=0) of ms, starting at xo with a tolerance of 0<tol. For other y's it doesn't find roots but where ms is y.
//implement Derive() and calculate it once instead of using D() multiple times
func (ms my_spline) NewtonRoot(x0 float64, y float64, tol float64) (float64,error) {
	debug := true
	derivative := ms.Derive()
	if debug{
		fmt.Println("calculated derivative.")
	}
	xn := x0
	skip := 100
	for math.Abs(y-ms.At(xn)) > tol{
		skip--
		if skip < 0 {return 0,fmt.Errorf("newton couldn't find root")}
		//fmt.Println("old xn: ",xn," , old yn: ", ms.At(xn), " , D(xn)=",ms.D(xn))
		xn = math.Min(max(ms.x),math.Max(min(ms.x),		xn+(y-ms.At(xn))/derivative.At(xn)		))
		if debug{
			fmt.Println(xn,":",ms.At(xn) , " , difference to ",y," is ",ms.At(xn)-y)
		}
		time.Sleep(1)
	}
	return xn,nil
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
	cumSpline := ms.IntegrateDUMB()
	for _,l := range levels {
		root,err := cumSpline.NewtonRoot(m,l,tol)
		if err != nil {continue}
		intersections = append(intersections,root)
	}
	return intersections
}





// ------------------------------- call specific functions -------------------------------

func findBestCall(pdist my_spline, calllist []callfunc) (callfunc, float64){
	best := calllist[0]
	bestE := best.ExpectedReturn(pdist)
	for _,c := range calllist {
		cE := c.ExpectedReturn(pdist)
		if cE > bestE{
			best = c
			bestE = cE
		}
	}
	return best, bestE
}

func ExpectedReturns(pdist my_spline, calllist []callfunc) ([]callfunc,[]float64) {
	var expReturns []float64
	var calls []callfunc
	for _,c := range calllist {
		expReturns = append(expReturns , c.ExpectedReturn(pdist) )
		calls = append(calls,c)
	}
	return calls,expReturns
}

func MathematicaPrintExpectedReturns(pdist my_spline, calllist []callfunc) string {
	calls, expReturns := ExpectedReturns(pdist,calllist)
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

//still not the same - possibly error in FullIntegralSpline()!!
func (call callfunc) ExpectedReturn (pdist my_spline) float64{
	debug := false
	var ref float64
	var dx float64
	if debug{
		dx = 0.1
		ref = call.ExpectedReturnDX(pdist,dx)
	}
	//result := pdist.SplineMultiply(call.ToSpline(min(pdist.x),max(pdist.x))).Integral(min(pdist.x),max(pdist.x),dx)
	result := pdist.SplineMultiply(call.ToSpline(min(pdist.x),max(pdist.x))).FullIntegralSpline()
	if debug{
		if math.Abs(ref - result) > 0.5*math.Abs(ref) {
			fmt.Println("ExpectedReturn(): error: ref not close enough. Ref is ",ref," and the new result is ",result)
			os.Exit(1)
		}
	}
	return result
}

func (call callfunc) ExpectedReturnDX (pdist my_spline,dx float64) float64{
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
		splineType: []string{"3","2","=Sl","=Cv","EQSl"},
		x:          []float64{a,call.base,b},
		y:          []float64{-100,-100,call.At(b)},
		coeffs:     []float64{0,-100,call.factor/call.cost*100,-100-100*call.base*call.factor/call.cost},
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
			cost:   optt.Close,
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

func containsFloat(list []float64, item float64, eps float64) bool {
	//eps := 0.01
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
				cost:   optt.Close,
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
				cost:   optt.Close,
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
