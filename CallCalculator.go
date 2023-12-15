package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"

	"bufio"
	//"encoding/json"
	"fmt"
	//"github.com/cnkei/gospline"
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


type Investment interface {
	ExpectedReturn(spline my_spline) float64
	At(x float64) float64
}

type nonInvested struct {}

func (nonInvested) At(x float64) float64 {
	return 0
}

func (nonInvested) ExpectedReturn(pdist my_spline) float64 {
	return 0
}

/*
func (*Investment) At(x float64) float64{
	return -1
}

func (Investment) ExpectedReturn(pdist my_spline) float64 {
	return 0
}

func (*Investment) ToSpline(a float64, b float64) my_spline {
	return my_spline{
		deg:        0,
		splineType: nil,
		x:          []float64{a,b},
		y:          []float64{0,0},
		coeffs:     []float64{0},
		unique:     true,
	}
}



func (*nonInvested) ToSpline(a float64, b float64) my_spline {
	return my_spline{
		deg:        0,
		splineType: nil,
		x:          []float64{a,b},
		y:          []float64{0,0},
		coeffs:     []float64{0},
		unique:     true,
	}
}
 */


type callfunc struct{
	base float64
	cost float64
	factor float64
	date []int
	optionType string
}

/*
type spread struct{
	num int
	calls []Investment
	weights []float64
}
 */

type spread struct{
	num int
	calls []callfunc
	weights []float64
}

/*
type spline_old struct {
	spline gospline.Spline
}
type normalized_spline_old struct{
	spline spline_old
	xrange []float64
	integral float64
}
 */


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

	/*
	var Investments []*Investment
	var Cash = &nonInvested{}
	Investments = append(Investments,Cash.Investment)
	 */

	riskAndTimePlottesting := false
	optimalTransporttesting := false
	apitesting := false
	calltesting := false
	splinetesting := false
	newtontesting := false

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

		mathCode := "SetDirectory[NotebookDirectory[]];\n"

		//FindSigmas
		var sigmasMap map[string][]float64
		sigmasMap = make(map[string][]float64,0)
		levels := []float64{0,0.125,0.25,0.5,0.75,0.875,1}
		for _,d := range pdistDates {
			cumSpline := pdistSplines[d].IntegrateDUMB()
			tmp,id := cumSpline.PrintMathematicaCode(true,"Blue","Automatic")
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
			tmp,id := pdistSplines[d].PrintMathematicaCode(false,"Blue","Automatic")
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

			tmp,id := cumSpline.PrintMathematicaCode(true,"Blue","Automatic")
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)

			tmp,id = invIntSpline.PrintMathematicaCode(true,"Blue","Automatic")
			mathCode += tmp+"\n"
			mathCode += fmt.Sprintf("s%v\n",id)
		}

		//optimal transport between dates
		//this is not enough
		//for ref, see: https://math.nyu.edu/~tabak/publications/Kuang_Tabak.pdf
		//cumSplines[0].At(invCumSplines[1].At(0.5))

		test := cumSplines[0].Subtract(cumSplines[1])
		tmp,id := test.PrintMathematicaCode(true,"Blue","Automatic")
		mathCode += tmp+"\n"
		mathCode += fmt.Sprintf("s%v\n",id)


		var transportMapFloat []float64
		//cumSplines[0], cumSplines[1] = UnionXYCC(cumSplines[0],cumSplines[1])
		//dx := 0.1
		for _,x := range test.x {
			transportMapFloat = append(transportMapFloat,cumSplines[0].Subtract(cumSplines[1]).IntegralSpline(0,x))
		}
		transportMapSpline := NewSpline(splinetype,cumSplines[0].x,transportMapFloat)
		tmp,id = transportMapSpline.PrintMathematicaCode(true,"Blue","Automatic")
		mathCode += tmp+"\n"
		mathCode += fmt.Sprintf("s%v\n",id)






		WriteFile("optTransport.nb",mathCode,"/")


	}

	if apitesting {


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

		tmp, _ := s.PrintMathematicaCode(true,"Blue","Automatic")
		mathCode := tmp
		fmt.Println(mathCode)


		ns := NewNormedSpline(s)

		tmp, _ = ns.PrintMathematicaCode(true,"Blue","Automatic")
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
		mathCode = bestcall.PrintMathematicaCode(max(pdist.x))
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

		tmp,_ := s.PrintMathematicaCode(true,"Blue","Automatic")
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

	if newtontesting{
		path, err := os.Getwd()
		check(err)
		promptName := "newtontesting"
		_,pdistDates,pdistX,pdistY,_,_,_,_,_ := LoadPromptEasy(path+"\\prompts\\","prompt_"+promptName+".json")

		splinetype := []string{"3","2","=Sl","=Cv","EQSl"}
		var pdistSplines map[string]my_spline
		pdistSplines = make(map[string]my_spline,0)
		d:=pdistDates[0]
		s := NewSpline(splinetype,pdistX[d],pdistY[d])
		ns := NewNormedSpline(s)
		pdistSplines[d] = ns

		constMinus := constSpline((max(ns.y)+min(ns.y))/2,[]float64{min(ns.x),max(ns.x)})
		testFunc := pdistSplines[d].Subtract(constMinus)
		fmt.Println(testFunc.PrintMathematicaCode(true,"Blue","Automatic"))
		roots := testFunc.NewtonRoots(0,0.01*((max(ns.y)+min(ns.y))/2),10)
		fmt.Println("roots: ",roots)

		negRange,posRange := testFunc.PosNegRange(0,0.001,100)
		fmt.Println("probReturn negRange: ",negRange)
		fmt.Println("probReturn posRange: ",posRange)
		os.Exit(1)
	}


	path, err := os.Getwd()
	check(err)

	//allowSubDirs := false
	files, err := ioutil.ReadDir(path+"/prompts/")
	check(err)
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	for _, file := range files {
		//if file.IsDir() && allowSubDirs && file.Name != "inactive" {}
		if !strings.Contains(file.Name(),"prompt_"){continue}
		promptName := strings.Split(strings.Split(file.Name(),"prompt_")[1],".json")[0]
		fmt.Println(promptName)
		promptSubPath := "\\prompts\\prompt_"+promptName+".json"

		prompt(promptSubPath)
	}

	//portfolios
	/*
	files, err = ioutil.ReadDir(path+"/portfolios/")
	check(err)
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
	for _, file := range files {
		//if file.IsDir() && allowSubDirs && file.Name != "inactive" {}
		if !strings.Contains(file.Name(),"prompt_"){continue}
		promptName := strings.Split(strings.Split(file.Name(),"prompt_")[1],".json")[0]
		fmt.Println(promptName)

		prompt(promptName)
	}
	 */

	os.Exit(1)

}


func prompt(promptSubPath string){
	forceUpdate := false
	debug := false
	info := true
	ASCIIPlots := true
	mathematicaExport := true
	brute := true
	iterativeFindN := true
	riskCompare := true
	selling := true
	optionsOutdatedHoursLimit := 24

	promptName := strings.Split(strings.Split(strings.Split(promptSubPath,"\\")[2],".")[0],"_")[1]

	if info{
		fmt.Println("promptName=",promptName)
		fmt.Println("Starting to execute prompt with name", promptName)
		fmt.Println("forceUpdate=",forceUpdate)
		fmt.Println("debug=",debug)
		fmt.Println("Starting data loading and preparation.")
	}

	startTime := time.Now()

	if info{
		fmt.Print("Loading prompt data and API key...")
	}
	path,err := os.Getwd();check(err)
	currentTime := time.Now()
	live := currentTime.Format("2006-01-02")

	// Loading prompt Json
	ticker,pdistDates,pdistX,pdistY,StrikeRange,DateRange,Contract_type,riskTolX,riskTolY := LoadPromptEasy(path,promptSubPath)

	// Loading API Key Json
	apiKey := opt.LoadJson("apiKey.json")
	if apiKey == "" {
		if info {fmt.Println("No API Key provided - trying to load options data from json file")}
		forceUpdate = false
	}
	if info {
		fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
	}

	if info {
		fmt.Print("Creating pdist and riskTol splines...")
		startTime = time.Now()
	}

	// Pdist splines
	var splinetype []string
	var pdistSplines map[string]my_spline
	pdistSplines = make(map[string]my_spline,0)
	pdistLinearTransitionThreshold := 40
	for _,d := range pdistDates {
		if debug {fmt.Println(pdistX[d],pdistY[d])}
		if len(pdistX[d]) < pdistLinearTransitionThreshold {
			splinetype = []string{"3","2","=Sl","=Cv","EQSl"}
		} else {
			splinetype = []string{"1","2"}
		}
		s := NewSpline(splinetype,pdistX[d],pdistY[d])
		ns := NewNormedSpline(s)
		pdistSplines[d] = ns
		if !s.unique {fmt.Println("Caution: pdist spline not unique!")}
	}

	// RiskTolSpline
	riskTolSpline := NewSpline([]string{"1","2"},riskTolX,riskTolY)

	if info{
		fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
	}

	if info{
		fmt.Print("Creating OptionsRequest for API Calls...")
		startTime = time.Now()
	}

	// Forming API Requests for each contract_type
	var optreq opt.OptionURLReq
	var optreqs []opt.OptionURLReq
	for _,ct := range Contract_type{
		fmt.Println("Contract_type:",ct)
		optreq = opt.OptionURLReq{
			Ticker:      ticker,
			ApiKey:      apiKey,
			StrikeRange: StrikeRange,
			DateRange:   DateRange,
			Contract_type: ct,
		}
		optreqs = append(optreqs,optreq)
	}

	if info{
		fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
	}

	/* End User Inputs */



	if info {
		fmt.Print("Making API Calls to pull necessary data...")
		startTime = time.Now()
	}

	//Pull USDEUR
	usdtoeur = PolygonAPISharePrice("C:USDEUR",apiKey)
	eurtousd = 1.0/usdtoeur
	if info{fmt.Printf("Data pulled: EURUSD=%.3f\n",eurtousd)}


	//Pull share price
	var share_price float64
	share_price = PolygonAPISharePrice(ticker,apiKey)
	if info{fmt.Printf("Data pulled: share_price(%s)=%.3f\n",ticker,share_price)}



	var options []opt.Option

	if !forceUpdate{

		// Try to load existing options json
		pulledDate, Contract_typeLoaded, options_tmp := LoadOptionsJson(path+"\\options\\"+ticker+"\\","options_latest.json")


		if info {
			fmt.Println("Loaded",len(options_tmp),"options from json file.")
		}

		for _,ct := range Contract_type {
			if !contains(Contract_typeLoaded,ct){
				if apiKey == "" {
					fmt.Println("No API Key nor suitable options json provided. Not enough data for a meaningful answer.")
					os.Exit(0)
				} else {
					forceUpdate = true
				}
			}
		}

		// Checking how outdated loaded options are.
		outOfDate := time.Now().Sub(pulledDate).Minutes()
		if outOfDate > float64(optionsOutdatedHoursLimit*60) {
			fmt.Printf("optionsdata %v h %v m old which is too old (%vh limit) - pulling data again through the API\n",(int)(outOfDate/60),int(math.Mod(outOfDate,60)),optionsOutdatedHoursLimit)
			options = options_tmp
			forceUpdate = true
		} else {

			if info {
				fmt.Printf("Loaded old options data from json. (%vh%vm old)\n",(int)(outOfDate/60),int(math.Mod(outOfDate,60)))
			}

			//options = options_tmp
			//filter dates
			if debug{
				fmt.Println("DateRange=",DateRange)
				fmt.Println("options[0].Expiration_date=",options_tmp[0].Expiration_date)
				fmt.Println("DateInRange(DateRange,options_tmp[0].Expiration_date)=",DateInRange(DateRange,options_tmp[0].Expiration_date))
				//os.Exit(610)
			}
			for _,optt := range options_tmp{
				if DateInRange(DateRange,optt.Expiration_date) {
					options = append(options,optt)
				}
			}

		}
	}

	if forceUpdate {

		// Execute all options requests and append the results in options []opt.Option
		nMax := -1 //how many successive requests at most; -1 is Inf
		options = MergeRequests(optreqs,nMax)

		os.Mkdir("options",0755)
		os.Mkdir("options"+"\\"+ticker,0755)
		os.Mkdir("options"+"\\"+ticker+"\\"+live,0755)

		now := time.Now()
		SaveOptionsJson("\\options\\"+ticker+"\\"+live+"\\","options",now,options,Contract_type)
		SaveOptionsJson("\\options\\"+ticker+"\\","options_latest",now,options,Contract_type)

	}

	// Add other investments
	long := callfunc{
		base:   0,
		cost:   share_price,
		factor: 1.0,
		date:   nil,
		optionType: "call",
	}
	/*
	short := callfunc{
		base:   0,
		cost:   share_price,
		factor: -1.0,
		date:   nil,
		optionType: "call",
	}
	 */

	//nonInvested cannot be represented by a call - we want a constant=0 function
	//nonInvested := constSpline(0.0,[]float64{0,1000000})


	//change to []Investment here
	var addToAll []callfunc
	addToAll = append(addToAll,long)
	//addToAll = append(addToAll,short)
	//addToAll = append(addToAll,nonInvested)
	// together with put, also implement short including checking all functions if they can handle short-interested functions and data


	// Group options by dates
	optionsDates, optionsMap, callListMap := OptionsToOptionsDates(options, addToAll)


	// weights for brute force
	weights := []float64{0.0,0.1,0.25,0.5,0.75,0.9,1.0}

	// Estimate calculation time for brute force
	var comparisonApprox int
	if info {
		fmt.Println("Found ",len(optionsDates), "(",optionsDates,") dates.")
		fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		if brute{

			for _,d := range optionsDates {
				comparisonApprox += int(math.Pow(float64(len(optionsMap[d])),2)*float64(len(weights)))
			}
			fmt.Printf("Brute force 2-spread comparison will compare %d spreads for all %v dates.\n Assuming a calculation time of 0.08 ms per spread, this will take approx. %.1f hours.\n",comparisonApprox,len(optionsDates),float64(comparisonApprox)*0.08/1000.0/60.0/24.0)
		}
	}

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


	//debug = true
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



	if info{
		fmt.Println("Iterating through all dates:")
	}

	// Setup for Mathematica Export
	var mathCode string
	content := "SetDirectory[NotebookDirectory[]];\n"
	mathematicaCompressionLevel := ".75"
	mathematicaImageResolution := "250"
	//mathCode := "SetDirectory[NotebookDirectory[]]\n"
	//mathCodeSigma := "SetDirectory[NotebookDirectory[]]\n"

	// Making files for outputs
	err = os.Mkdir(path+"\\tmp\\"+live,0755)
	err = os.Mkdir(path+"\\tmp\\"+live+"\\"+promptName, 0755)
	path = path+"\\tmp\\"+live+"\\"+promptName+"\\"


	var strikes []float64
	var costs []float64

	var overallBestSpread spread
	var overallBestSpreadExpCAGR float64
	var overallBestSpreadTotalCount int
	var overallBestSpreadRiskMatchCount int
	var overallBestSpreadPDist my_spline
	var overallMsg string
	//var overallBestSpreadExp float64 = -1000
	//var overallDate string
	//var overallBestSpreadRiskTolExclusion string


	// Go through all expiry dates d for options
	for j,d := range optionsDates {

		if info {
			fmt.Println("Starting with date ",d,"...")
			//startTime = time.Now()
		}

		//filter dates
		if debug{
			fmt.Println("DateRange=",DateRange)
		}
		if !DateInRange(DateRange,d) {
			if info {
				fmt.Println("Date",d,"not in date range",DateRange,". Continue.")
			}
			continue
		}

		// Making brute force calculation time approximation
		comparisonApprox = 0
		for i := j ; i < len(optionsDates) ; i++ {
			comparisonApprox += int(math.Pow(float64(len(optionsMap[optionsDates[i]])),2)*float64(len(weights)))
		}
		comparisonApprox = int(math.Pow(float64(len(optionsMap[d])),2)*float64(len(weights)))
		if info && brute {
			fmt.Printf("Brute force 2-spread comparison will compare %d spreads for the remaining %v dates.\n Assuming a calculation time of 0.08 ms per spread, this will take approx. %.3f hours.\n",comparisonApprox,len(optionsDates)-j-1,float64(comparisonApprox)*0.08/1000.0/60.0/24.0)
			fmt.Printf("Brute force 2-spread comparison will compare %d spreads for this date (%s).\n Assuming a calculation time of 0.08 ms per spread, this will take approx. %.3f hours.\n",comparisonApprox,d,float64(comparisonApprox)*0.08/1000.0/60.0/24.0)
		}


		dDate, err := time.Parse("2006-01-02", d);check(err)
		hoursToExpiry := float64(dDate.Sub(time.Now()).Hours())
		daysToExpiry := hoursToExpiry/24.0
		yearsToExpiry := float64(daysToExpiry)/365.0

		/*
		fmt.Println("dDate=",dDate)
		fmt.Println("time.Now()=",time.Now())
		fmt.Println("daysToExpiry=",daysToExpiry)
		fmt.Println("yearsToExpiry=",yearsToExpiry)
		*/


		// ---- Find closest pdist ----------
		// !should eventually be optimal transported!
		pdist := FindCloestPDist(d,pdistDates,pdistSplines,debug)


		// ------- Make folder for date --------
		folderName := ticker+d+"(live data from "+live+")"
		err = os.Mkdir(path+folderName, 0755);check(err)


		optionsList := optionsMap[d]
		callList := callListMap[d]
		callList = callList[len(addToAll):len(callList)]
		if info{
			fmt.Println("Found ",len(optionsList)," options for ",ticker," on the date ",d)
		}


		//debugging - print all(?, just one call,put each?) calls and puts through PrintMathematicaCode()
		/*
		debugBasic := false
		if debugBasic{
			testCode := mathCode
			var testCallPut []callfunc
			var callIn, putIn bool = false,false
			for _,c := range callList {
				if len(testCallPut)== 2{
					testCode += testCallPut[0].ToSpline(0.0,1000.0).MathematicaExport2("Red",testCallPut[1].ToSpline(0.0,1000.0),"Red","",false,folderName,"-testCallPut",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
					WriteFile("testCallPut.nb",testCode,"/tmp/"+live+"/"+promptName+"/")
					//"output.nb",content,"/tmp/"+live+"/"+promptName+"/"
					//WriteFile("testCallPut.nb",testCode,"")
					os.Exit(124)
				}
				if c.optionType=="call" && !callIn && c.base>200.0 {
					testCallPut = append(testCallPut,c)
					callIn=true
				}
				if c.optionType=="put" && !putIn && c.base>200.0 {
					testCallPut = append(testCallPut,c)
					putIn=true
				}
			}
		}
		 */



		// Pdist
		var idPdist string
		if info && mathematicaExport{
			fmt.Print("Creating pdist wolfram mathematica export code...")
			startTime = time.Now()

			mathCode,idPdist = pdist.PrintMathematicaCode(false,"Blue","Automatic")
			content += mathCode
			content += "Export[\"" + folderName + "\\pdist.png\", " + fmt.Sprintf("Show[fctplot%v]",idPdist) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		}
		if ASCIIPlots {
			fmt.Println("Pdist Plot")
			fmt.Println(pdist.PrintASCII(min(pdist.x),max(pdist.x),130,35,false))
		}

		// only pdist testing
		//continue



		// All Options
		if info && mathematicaExport {
			fmt.Print("Creating wolfram mathematica export code for all options(",Contract_type,")...")
			startTime = time.Now()

			mathCode = PrintMathematicaCode(callList, share_price,"Blue","Red",true)
			content += mathCode
			content += "Export[\"" + folderName + "\\allCalls.png\", Show[s], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		}



		// Best Call
		if info{
			fmt.Print("Finding best option(",Contract_type,") ")
			if mathematicaExport {
				fmt.Print("and creating wolfram mathematica export code...")
			} else {
				fmt.Print(" ...")
			}
			startTime = time.Now()
		}
		bestcall, bestE := findBestCall(pdist, callList)
		//bestCallCAGR := (math.Pow(bestE/100.0+1.0,1.0/yearsToExpiry)-1.0)*100
		bestCallCAGR := CAGR(bestE,yearsToExpiry)
		if mathematicaExport {
			mathCode = bestcall.PrintMathematicaCode(max(pdist.x))
			content += fmt.Sprintf("msg1 := Text[\"Assuming the probability distribution (left) for the date %v, the %s with strike %.1f has the highest expected return out of all call options available with %.1f %% expected return (%.2f Percent CAGR). Owning the underlying asset (%v) has an expected return of %.1f %% (%.1f Percent CAGR).  \"];\n\n", callList[0].date,bestcall.optionType, bestcall.base, bestE,bestCallCAGR, ticker, long.ExpectedReturn(pdist),(math.Pow(long.ExpectedReturn(pdist)/100.0+1.0,1.0/yearsToExpiry)-1.0)*100.0)
			content += mathCode
			content += "Export[\"" + folderName + "\\-bestCall.png\", {msg1 \n , "+fmt.Sprintf("Show[fctplot%v]",idPdist) +", Show[call,long]}, \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		}
		if info {
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		}


		//Distribution Chart for Call-Long intersections
		/*
		//fmt.Println("\nDistribution Chart for Call-Long intersections:\n")
		mathCode = MathematicaCodeLongIntersection(callList, share_price)
		//fmt.Println(mathCode)
		content += mathCode
		content += "Export[\"" + folderName + "\\CallLongIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		 */

		//Distribution Chart for Call-Zero intersections
		/*
		//fmt.Println("\nDistribution Chart for Call-Zero intersections:\n")
		mathCode = MathematicaCodeZeroIntersection(callList)
		//fmt.Println(mathCode)
		content += mathCode
		content += "Export[\"" + folderName + "\\CallZeroIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		 */

		//Distribution Chart for Call-Zero-Volumes intersections
		/*
		//fmt.Println("\nDistribution Chart for Call-Zero-Volumes intersections:\n")
		mathCode = MathematicaCodeZeroIntersectionVolumes(optionsList)
		content += mathCode
		content += "Export[\"" + folderName + "\\CallZeroVolumesIntersectionDistribution.png\", Show[dist], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		 */


		//Expected returns for each option
		if info {
			fmt.Print("Computing expected returns (given pdist) for each option (",Contract_type,")")
			if mathematicaExport {
				fmt.Print(" and creating wolfram mathematica export code...")
			} else {
				fmt.Print(" ...")
			}
			startTime = time.Now()
		}
		if mathematicaExport {
			mathCode = MathematicaPrintExpectedReturns(pdistSplines[d], callList)
			content += mathCode
			content += "Export[\"" + folderName + "\\expected_returns_strike.png\", Show[xy], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		}
		if info {
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		}


		// Strike prices
		if mathematicaExport {
			if info {
				fmt.Print("Creating wolfram mathematica export code for a strike price plot...")
				startTime = time.Now()
			}

			strikes = make([]float64,0)
			costs = make([]float64,0)
			for _, opt := range optionsMap[d] {
				strikes = append(strikes, float64(opt.Strike_price))
				costs = append(costs, (opt.Close))
			}

			mathCode = MathematicaXYPlot(strikes, costs)
			content += mathCode
			content += "Export[\"" + folderName + "\\strike_price.png\", Show[xy], \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"

			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
			}
		}


		// probReturn
		var probReturn my_spline
		if mathematicaExport{
			if info {
				fmt.Print("Creating wolfram mathematica export code for probReturn...")
				startTime = time.Now()
			}

			probReturn = pdist.SplineMultiply(long.ToSpline(0,max(pdist.x)))
			content += probReturn.MathematicaExport("Blue","",false,folderName,"probReturn",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")

			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
			}
		}


		// probReturnIntegral
		//var probReturnIntegral my_spline
		if mathematicaExport {
			if info {
				fmt.Print("Creating wolfram mathematica export code for probReturnIntegral...")
				startTime = time.Now()
			}

			probReturnIntegral := probReturn.Integrate()
			content += probReturnIntegral.MathematicaExport("Blue","",false,folderName,"probReturnIntegral",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")

			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
			}
			if debug {
				root := probReturnIntegral.NewtonRoots(0,0.01,100)
				fmt.Println("probReturnIntrgral root: ",root)
			}
		}



		// pdistIntegral
		if info {
			fmt.Print("Creating wolfram mathematica export code for pdistIntegral...")
			startTime = time.Now()
		}
		pdistIntegrate := pdist.Integrate()
		tmp,id := pdistIntegrate.PrintMathematicaCode(true,"Blue","Automatic")
		mathCode += tmp+"\n"
		mathCode += "Export[\"" + folderName + "\\pdistIntegrate.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		content += mathCode
		if info {
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
		}

		/*
		fmt.Print("pdistIntegralInverse...")
		pdistIntegrateInverse := pdist.Integrate().Inversion()
		tmp,id = pdistIntegrateInverse.PrintMathematicaCode()
		mathCode += tmp+"\n"
		mathCode += "Export[\"" + folderName + "\\pdistIntegrateInverse.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		content += mathCode
		fmt.Println(" done.")
		 */

		/*
		fmt.Print("Risk eval")
		var riskY []float64
		var ls []float64
		dl := 0.1
		for l:=0.0;l<=1.0;l+=dl {
			ls = append(ls, l)
			riskY = append(riskY, bestcall.At(pdistIntegrateInverse.At(l)) )
		}
		riskProfile := NewSpline(splinetype,ls,riskY)
		tmp,id = riskProfile.PrintMathematicaCode()
		mathCode += tmp+"\n"
		mathCode += "Export[\"" + folderName + "\\riskProfileBestCall.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		content += mathCode
		fmt.Println(" done.")
		 */

		//risk evaluation
		/*
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
		*/

		// riskSpline
		var riskSpline my_spline
		if mathematicaExport {
			if info{
				fmt.Print("Creating wolfram mathematica export code for riskSpline...")
				startTime = time.Now()
			}
			riskSpline, err = bestcall.ToSpread().riskProfile(pdist);check(err)
			if len(riskSpline.x) == 0 || len(riskSpline.y) == 0 {continue}
			content += riskSpline.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-riskSplineBestCall",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
			}
		}




		/*
			riskSpline := bestcallSpread.riskProfile(pdist)
			tmp,id = riskSpline.PrintMathematicaCode(false)
			mathCode += tmp+"\n"
			mathCode += "Export[\"" + folderName + "\\-riskSplineBestCall.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
			content += mathCode
		*/


		// riskTolSpline
		if mathematicaExport {
			if info{
				fmt.Print("Creating wolfram mathematica export code for riskTolSpline...")
				startTime = time.Now()
			}
			content += riskTolSpline.MathematicaExport("Darker[Red]","",false,folderName,"-UserRiskTolSpline",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(riskSpline,riskTolSpline))
			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
			}
		}


		//test nonInvested
		/*
		nonInvestedSpread := nonInvested.ToSpread()
		content += nonInvestedSpread.MathematicaExport(pdist,false,folderName,"nonInvested",mathematicaCompressionLevel,mathematicaImageResolution)

		bestcallNonInvestedSpread := bestcallSpread.Add(nonInvested,0.5)
		content += bestcallNonInvestedSpread.MathematicaExport(pdist,false,folderName,"bestCallNonInvested",mathematicaCompressionLevel,mathematicaImageResolution)

		*/



		var riskTolExclusion string


		// Iterative Finding

		var msg string
		var iterations int
		var starts int
		var weightDelta float64
		iterations = 100000
		starts = 10
		weightDelta = 0.1
		maxSearchLength := 20


		//Best 2-Spread Iterative Find
		/*
		if info{
			fmt.Println("Starting gradient decent search in all 2-combination spreads and creating wolfram mathematica export codes...")
			startTime = time.Now()
		}
		bestSpreadFind2, bestSpreadFind2Exp := BestSpreadIterativeFind2(pdist,callList,weightDelta,iterations,riskTolSpline,starts)
		fmt.Println("BestSpreadIterativeFind2() compared",2*4*iterations*starts," spreads and found",bestSpreadFind2,"with an exp. return of",bestSpreadFind2Exp)

		msg = fmt.Sprintf("Assuming the probability distribution for the date %v, the 2-spread with strikes and weights {(%s,%.1f, %.2f),(%s,%.1f, %.2f)} has the highest expected return out of all call options available with %.1f Percent expected return. Owning the underlying asset (%v) has an expected return of %.1f Percent. %s", bestSpreadFind2.calls[0].date,bestSpreadFind2.calls[0].optionType, bestSpreadFind2.calls[0].base,bestSpreadFind2.weights[0],bestSpreadFind2.calls[1].optionType,bestSpreadFind2.calls[1].base,bestSpreadFind2.weights[1], bestSpreadFind2Exp, ticker, long.ExpectedReturn(pdist),riskTolExclusion)
		longSpline := long.ToSpline(min(pdist.x),max(pdist.x))
		bestSpreadFindSpline := bestSpreadFind2.ToSpline(min(pdist.x),max(pdist.x))
		content += bestSpreadFindSpline.MathematicaExport2("Blue",longSpline,"Red","",false,folderName,"-bestSpreadFind2",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(pdist,bestSpreadFindSpline,longSpline))
		rpFind,err1 := bestSpreadFind2.riskProfile(pdist)
		if len(rpFind.x) == 0 || len(rpFind.y) == 0 {continue}
		if err1 == nil {
			content += rpFind.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-bestSpreadFind2RiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rpFind,riskTolSpline))
		} else {
			fmt.Println("ERROR in riskProfile:",err1)
		}
		if info{
			elapsed := time.Now().Sub(startTime).Milliseconds()
			elapsedPerSpread := float64(elapsed)/float64(2*4*iterations*starts)
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms - ",elapsedPerSpread,"ms per spread)")
		}
		 */



		//Best 4-Spread Iterative Find
		/*
		if info{
			fmt.Println("Starting gradient decent search in all 4-combination spreads and creating wolfram mathematica export codes...")
			startTime = time.Now()
		}
		bestSpreadFind4, bestSpreadFind4Exp := BestSpreadIterativeFind4(pdist,callList,weightDelta,iterations,riskTolSpline,starts)
		fmt.Println("BestSpreadIterativeFind4() compared",4*4*iterations*starts," spreads and found",bestSpreadFind4,"with an exp. return of",bestSpreadFind4Exp)

		msg = fmt.Sprintf("Assuming the probability distribution for the date %v, the 4-spread with strikes and weights",bestSpreadFind4,"has the highest expected return out of all call options available with %.1f Percent expected return. Owning the underlying asset (%v) has an expected return of %.1f Percent. %s", bestSpreadFind4.calls[0].date, bestSpreadFind4Exp, ticker, long.ExpectedReturn(pdist),riskTolExclusion)
		longSpline = long.ToSpline(min(pdist.x),max(pdist.x))
		bestSpreadFindSpline = bestSpreadFind4.ToSpline(min(pdist.x),max(pdist.x))
		content += bestSpreadFindSpline.MathematicaExport2("Blue",longSpline,"Red",msg,false,folderName,"-bestSpreadFind4",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(pdist,bestSpreadFindSpline,longSpline))
		rpFind,err1 = bestSpreadFind4.riskProfile(pdist)
		if len(rpFind.x) == 0 || len(rpFind.y) == 0 {continue}
		if err1 == nil {
			content += rpFind.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-bestSpreadFind4RiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rpFind,riskTolSpline))
		} else {
			fmt.Println("ERROR in riskProfile:",err1)
		}
		if info{
			elapsed := time.Now().Sub(startTime).Milliseconds()
			elapsedPerSpread := float64(elapsed)/float64(4*4*iterations*starts)
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms - ",elapsedPerSpread,"ms per spread)")
		}
		 */


		//Best N=2-Spread Iterative Find
		/*
		N := 2
		contractTypes := []string{"call","put"}
		if info{
			fmt.Println("Starting gradient decent search in all (N=",N,")-combination spreads and creating wolfram mathematica export codes...")
			startTime = time.Now()
		}
		bestSpreadFindN, bestSpreadFindNExp,startRiskProfile, riskTotalCount,riskMatchCount,totalSpreadsCompared := BestSpreadIterativeFindN(N,pdist,callList,weightDelta,iterations,riskTolSpline,starts,contractTypes,maxSearchLength,share_price)
		if info {
			//totalCount := N*4*iterations*starts
			fmt.Println("BestSpreadIterativeFindN() with N=",N," compared",totalSpreadsCompared," spreads and found",bestSpreadFindN,"with an exp. return of",bestSpreadFindNExp)
			riskTolExclusion = fmt.Sprintf("%.5f Percent (%v out of %v) of spreads were excluded due to the risk profile not matching.\n",100.0*float64(riskTotalCount-riskMatchCount)/float64(riskTotalCount),riskTotalCount-riskMatchCount,riskTotalCount)
			fmt.Println(riskTolExclusion)
		}
		msg = fmt.Sprintf("Assuming the probability distribution for the date %v, the (N=",N,")-spread with strikes and weights",bestSpreadFindN,"has the highest expected return out of all call options available with %.1f Percent expected return. Owning the underlying asset (%v) has an expected return of %.1f Percent. %s", bestSpreadFindN.calls[0].date, bestSpreadFindNExp, ticker, long.ExpectedReturn(pdist),riskTolExclusion)
		longSpline := long.ToSpline(min(pdist.x),max(pdist.x))
		bestSpreadFindSpline := bestSpreadFindN.ToSpline(min(pdist.x),max(pdist.x))
		content += bestSpreadFindSpline.MathematicaExport2("Blue",longSpline,"Red",msg,false,folderName,"-bestSpreadFindN2",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(pdist,bestSpreadFindSpline,longSpline))
		rpFind,err1 := bestSpreadFindN.riskProfile(pdist)
		if len(rpFind.x) == 0 || len(rpFind.y) == 0 {continue}
		if err1 == nil {
			content += rpFind.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-bestSpreadFindN2RiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rpFind,riskTolSpline))
		} else {
			fmt.Println("ERROR in riskProfile:",err1)
		}
		content += startRiskProfile.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-bestStartRiskProfileSpreadFindN2",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rpFind,riskTolSpline))
		if info{
			elapsed := time.Now().Sub(startTime).Milliseconds()
			elapsedPerSpread := float64(elapsed)/float64(totalSpreadsCompared)
			fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms - ",elapsedPerSpread,"ms per spread)")
		}
		 */



		//Best N=4-Spread Iterative Find
		if iterativeFindN {
			N := 4
			contractTypes := []string{"call","put","call","put"}
			if info{
				fmt.Println("Starting gradient decent search in all (N=",N,")-combination spreads and creating wolfram mathematica export codes...")
				startTime = time.Now()
			}
			//bestSpread,bestSpreadExp,startRiskProfile, riskTotalCount,riskMatchCount,totalSpreadsCompared
			bestSpreadFindN, bestSpreadFindNExp,_,riskTotalCount,riskMatchCount,totalSpreadsCompared := BestSpreadIterativeFindN(N,pdist,callList,weightDelta,iterations,riskTolSpline,starts,contractTypes,maxSearchLength,share_price)
			if info {
				//totalCount := N*4*iterations*starts
				fmt.Println("BestSpreadIterativeFindN() with N=",N," compared",totalSpreadsCompared," spreads and found",bestSpreadFindN,"with an exp. return of",bestSpreadFindNExp)
				riskTolExclusion = fmt.Sprintf("%.5f Percent (%v out of %v) of spreads were excluded due to the risk profile not matching.\n",100.0*float64(riskTotalCount-riskMatchCount)/float64(riskTotalCount),riskTotalCount-riskMatchCount,riskTotalCount)
				fmt.Println(riskTolExclusion)
			}
			msg = fmt.Sprintf("Assuming the probability distribution for the date %v, the (N=%v)-spread with strikes and weights %v has the highest expected return out of all call options available with %.1f Percent expected return. Owning the underlying asset (%v) has an expected return of %.1f Percent. %s", bestSpreadFindN.calls[0].date, N, bestSpreadFindN, bestSpreadFindNExp, ticker, long.ExpectedReturn(pdist),riskTolExclusion)
			longSpline := long.ToSpline(min(pdist.x),max(pdist.x))
			bestSpreadFindSpline := bestSpreadFindN.ToSpline(min(pdist.x),max(pdist.x))
			content += bestSpreadFindSpline.MathematicaExport2("Blue",longSpline,"Red",msg,false,folderName,"-bestSpreadFindN4",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(pdist,bestSpreadFindSpline,longSpline))
			rpFind,err1 := bestSpreadFindN.riskProfile(pdist);check(err1)
			if len(rpFind.x) == 0 || len(rpFind.y) == 0 {continue}
			if err1 == nil {
				content += rpFind.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-bestSpreadFindN4RiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rpFind,riskTolSpline))
			} else {
				fmt.Println("ERROR in riskProfile:",err1)
			}
			if info{
				elapsed := time.Now().Sub(startTime).Milliseconds()
				elapsedPerSpread := float64(elapsed)/float64(totalSpreadsCompared)
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms - ",elapsedPerSpread,"ms per spread)\n")
			}


			//console print best spread
			resX := 130
			resY := 35
			fmt.Print(bestSpreadFindN.PrintList())
			fmt.Printf("Expected return: %.2f %%\n",bestSpreadFindNExp)
			fmt.Print(bestSpreadFindN.PrintASCIIPerformance(min(pdist.x),max(pdist.x),resX,resY,true))
			fmt.Println("RiskProfile:")
			fmt.Println(rpFind.PrintASCII(0.0,1.0,resX,resY,true))
			fmt.Println(rpFind)
			RiskProfileIterativeMeanReturn(rpFind)

			// riskSpline
			/*
			if info{
				fmt.Print("Creating wolfram mathematica export code for riskSpline...")
				startTime = time.Now()
			}
			riskSpline, err = bestSpreadFindN.riskProfile(pdist)
			if len(riskSpline.x) == 0 || len(riskSpline.y) == 0 {continue}
			if err == nil {
				content += riskSpline.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",true,folderName,"-riskSplineBestSpreadFindN",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
			} else {
				fmt.Println("ERROR in riskProfile:",err)
			}
			if info {
				fmt.Println("done. (took",time.Now().Sub(startTime).Milliseconds(),"ms)")
				fmt.Println("RiskProfile of best spread found by FindN():")
				fmt.Println(riskSpline.PrintASCII(pdist,120,35))
			}
			 */



		}







		// Brute force legacy

		if brute {


			// all 2-combinations of calls and {0.5,0.75,0.25} weighing in both directions (buy&sell)
			if info{
				fmt.Println("Forming and comparing spreads and creating wolfram mathematica export codes...")
				startTime = time.Now()
			}

			bestSpread,bestSpreadExp,riskMatchCount,totalCount,timeTallys := BestSpread2CombinationsManual(pdist,callList,weights,riskCompare,riskTolSpline,selling)


			riskTolExclusion = ""
			if riskCompare {
				riskTolExclusion = fmt.Sprintf("%.5f Percent (%v out of %v) of spreads were excluded due to the risk profile not matching. %v spreads are feasible. \n",100.0*float64(totalCount-riskMatchCount)/float64(totalCount),totalCount-riskMatchCount,totalCount,riskMatchCount)
				overallBestSpreadTotalCount += totalCount
				fmt.Println(riskTolExclusion)
			}

			if info {
				tDelta := time.Now().Sub(startTime).Milliseconds()
				fmt.Println("done. (took",tDelta/1000.0/60.0,"mins, created and compared",totalCount,"spreads (",float64(tDelta)/float64(totalCount),"ms per spread) and found a spread with exp. return",bestSpreadExp,"%)")
				timeTally1 := timeTallys[0]
				timeTally2 := timeTallys[1]
				fmt.Printf("%.1f %% (%.1f seconds) of time was spent riskProfile(), the remaining %.1f %% (%.1f seconds) were spent in the rest of BestSpread2CombinationManual()\n",timeTally2/(timeTally1+timeTally2)*100,timeTally2/1000.0,timeTally1/(timeTally1+timeTally2)*100,timeTally1/1000.0)
			}

			//fmt.Printf("Assuming the probability distribution for the date %v, the 2-spread with strikes and weights {(%.1f, %.2f),(%.1f, %.2f)} has the highest expected return out of all call options available with %.1f %% expected return. Owning the underlying asset (%v) has an expected return of %.1f %%. %s", bestSpread.calls[0].date, bestSpread.calls[0].base,bestSpread.weights[0],bestSpread.calls[1].base,bestSpread.weights[1], bestSpreadExp, ticker, long.ExpectedReturn(pdist),riskTolExclusion)
			if riskMatchCount == 0 {
				fmt.Println("No spreads matching risk profile. Continue with next prompt.")
				content += riskTolSpline.MathematicaExport("Darker[Red]","No spreads matching risk profile. Continue with next prompt.",false,folderName,"-bestSpreadRiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(riskTolSpline,riskTolSpline))
				continue
			} else {
				bestSpreadCAGR := CAGR(bestSpreadExp,yearsToExpiry)
				if debug {
					fmt.Println("yearsToExpiry=",yearsToExpiry)
					fmt.Println("bestSpreadCAGR=",bestSpreadCAGR)
				}
				longExp := long.ExpectedReturn(pdist)
				longCAGR := CAGR(longExp,yearsToExpiry)
				msg := fmt.Sprintf("Assuming the probability distribution for the date %v, the 2-spread with strikes and weights {(%s,%.1f, %.2f),(%s,%.1f, %.2f)} has the highest expected return out of all call options available with %.1f Percent expected return (%.1f Percent CAGR). Owning the underlying asset (%v) has an expected return of %.1f Percent. (%.1f Percent CAGR) %s", bestSpread.calls[0].date,bestSpread.calls[0].optionType, bestSpread.calls[0].base,bestSpread.weights[0],bestSpread.calls[1].optionType,bestSpread.calls[1].base,bestSpread.weights[1], bestSpreadExp,bestSpreadCAGR, ticker, longExp,longCAGR,riskTolExclusion)
				longSpline := long.ToSpline(min(pdist.x),max(pdist.x))
				bestSpreadSpline := bestSpread.ToSpline(min(pdist.x),max(pdist.x))
				content += bestSpreadSpline.MathematicaExport2("Blue",longSpline,"Red",msg,false,folderName,"-bestSpread",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(pdist,bestSpreadSpline,longSpline))
				rp,err1 := bestSpread.riskProfile(pdist);check(err1)
				if len(rp.x) == 0 || len(rp.y) == 0 {continue}
				if err1 == nil {
					content += rp.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",false,folderName,"-bestSpreadRiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rp,riskTolSpline))
				} else {
					fmt.Println("ERROR in riskProfile:",err)
				}

				//console print best spread
				resX := 130
				resY := 35
				fmt.Print(bestSpread.PrintList())
				fmt.Printf("Expected return: %.2f %%\n",bestSpreadExp)
				fmt.Print(bestSpread.PrintASCIIPerformance(min(pdist.x),max(pdist.x),resX,resY,true))
				fmt.Println("RiskProfile:")
				fmt.Println(rp.PrintASCII(0.0,1.0,resX,resY,true))
				fmt.Println(rp)
				RiskProfileIterativeMeanReturn(rp)



				// Overall update check
				if overallBestSpreadExpCAGR < bestSpreadCAGR {
					overallBestSpreadExpCAGR = bestSpreadCAGR
					//overallBestSpreadExp = bestSpreadExp
					overallBestSpread = bestSpread
					//overallBestSpreadRiskTolExclusion = riskTolExclusion
					//overallBestSpreadTotalCount = totalCount
					overallBestSpreadRiskMatchCount = riskMatchCount
					overallBestSpreadPDist = pdist
					overallMsg = msg
					//overallDate = d
				}


			}

		}

		if info {
			fmt.Println("Done with date", d,"\n\n")
		}

	}


	// Create Overall review and folder and mathematica export
	folderName := ticker + " (Overall)"
	err = os.Mkdir(path+folderName, 0755);check(err)

	if overallBestSpreadRiskMatchCount == 0 {
		fmt.Println("No spreads matching risk profile. Continue with next prompt.")
		content += riskTolSpline.MathematicaExport("Darker[Red]","No spreads matching risk profile. Continue with next prompt.",false,folderName,"-bestSpreadRiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(riskTolSpline,riskTolSpline))
	} else {
		msg := overallMsg //wrong overallBestSpreadTotalCount
		fmt.Println("overallBestSpreadTotalCount=",overallBestSpreadTotalCount)
		longSpline := long.ToSpline(min(overallBestSpreadPDist.x),max(overallBestSpreadPDist.x))
		bestSpreadSpline := overallBestSpread.ToSpline(min(overallBestSpreadPDist.x),max(overallBestSpreadPDist.x))
		content += bestSpreadSpline.MathematicaExport2("Blue",longSpline,"Red",msg,false,folderName,"-bestSpread",mathematicaCompressionLevel,mathematicaImageResolution,PlotRange(overallBestSpreadPDist,bestSpreadSpline,longSpline))
		rp,err1 := overallBestSpread.riskProfile(overallBestSpreadPDist);check(err1)
		content += rp.MathematicaExport2("Blue",riskTolSpline,"Darker[Red]","",false,folderName,"-bestSpreadRiskProfile",mathematicaCompressionLevel,mathematicaImageResolution,riskPlotRange(rp,riskTolSpline))
	}



	//WriteFile("sigmas.nb",mathCodeSigma,"/")

	WriteFile("output.nb",content,"/tmp/"+live+"/"+promptName+"/")

	//Portfolio
	/*
	files, err := ioutil.ReadDir(path+promptSubPath[1:])
	check(err)
	for _, file := range files {
		if strings.Contains(file.Name(),"portfolio_") {
			//sell assets and buy Overall best
			LoadPortfolioJson()
		}
	}
	 */

}



// ------------------------------- to be implemented -------------------------------
func LoadPortfolioJson(){}
func PDistOptimalTransport(){}
func ChangeMeritPortfolio(best spread, transactionHistory int, taxProfile int){}
func ChangePortfolio(newSpread spread){}
func LeastSquare(x,y []float64){}

//careful: only all calls or all puts
func OptionLeastSquare(callList []callfunc) error {
	callType := callList[0].optionType
	for _,c := range callList[1:] {
		if c.optionType != callType {return fmt.Errorf("not same option type (call,put)")}
	}
	return nil
}

func HuellkurveIdeal(){}
func Huellkurve(){}

func PowellWolfe(){}




// ------------------------------- to be sorted -------------------------------



// ------------------------------- testing functions -------------------------------

func testingDumpster() {
	//testing .ToSpline()
	/*
		bestCallSpline := bestcall.ToSpline(min(pdist.x),max(pdist.x))
		tmp,id := bestCallSpline.PrintMathematicaCode(false)
		mathCode += tmp+"\n"
		mathCode += "Export[\"" + folderName + "\\TEST_bestCallSpline.png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
		content += mathCode
		//content += bestCallSpline.MathematicaExport("Blue","",false,folderName,"TEST_bestCallSpline",mathematicaCompressionLevel,mathematicaImageResolution)
	*/

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
	/*
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

	//testing UnionXYCC()
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

	//sell put testing
	/*
		fmt.Println("sell put testing")
		var putTesting callfunc
		for _,opt := range callList{
			if opt.optionType == "put" {
				if opt.base > 50.0{
					putTesting = opt
					break
				}
			} else {
				continue
			}
		}
		fmt.Println("use put:",putTesting)
		sellputSpread := spread{
			num:     1,
			calls:   []callfunc{putTesting},
			weights: []float64{-1.0},
		}
		sellputExp := sellputSpread.ExpectedReturn(pdist)
		fmt.Println("expected return:",sellputExp)
		sellputSpline := sellputSpread.ToSpline(min(pdist.x),max(pdist.x))
		content += putTesting.ToSpline(min(pdist.x),max(pdist.x)).MathematicaExport("Blue","",false,folderName,"-PutTestingToSpline",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
		content += putTesting.ToSpread().ToSpline(min(pdist.x),max(pdist.x)).MathematicaExport("Blue","",false,folderName,"-PutTestingToSpreadToSpline",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
		content += sellputSpline.MathematicaExport("Blue","",false,folderName,"-sellPutTesting",mathematicaCompressionLevel,mathematicaImageResolution,"Automatic")
		fmt.Println("done.")
	*/

	//PosNegRange testing
	/*
		fmt.Println("PosNegRange testing...")
		bestSpreadSpline = bestSpread.ToSpline(min(pdist.x),max(pdist.x))
		//bestSpreadSpline := bestcall.ToSpline(min(pdist.x),max(pdist.x))
		yTest := 120.0
		roots := bestSpreadSpline.NewtonRoots(yTest,0.0001,40)
		fmt.Println("roots: ",roots)
		negTest,_ := bestSpreadSpline.PosNegRange(yTest,0.0001,40)
		fmt.Println("negRange: ",negTest)
		fmt.Println("done.")
	*/

	//Integral testing
	/*
	a:=0.0;b:=1000.0;
	fmt.Println(fmt.Sprintf("pdist.IntegralSpline(%v,%v)=",a,b),pdist.IntegralSpline(a,b))
	*/

	//testing spreads with ASCII Plots
	/*
		call1 := callfunc{
			base:       200,
			cost:       3,
			factor:     1,
			date:       nil,
			optionType: "call",
		}
		put1 := callfunc{
			base:       200,
			cost:       3,
			factor:     -1,
			date:       nil,
			optionType: "put",
		}

		testSpread := spread{
			num:     1,
			calls:   []callfunc{long},
			weights: []float64{1.0},
		}
		fmt.Println(testSpread.PrintList())
		fmt.Println(testSpread.PrintASCIIPerformance(my_spline{
			deg:        1,
			splineType: nil,
			x:          []float64{0.0,500.0},
			y:          nil,
			coeffs:     nil,
			unique:     false,
		},120,35))

		testSpread = spread{
			num:     1,
			calls:   []callfunc{call1,long},
			weights: []float64{-1,0},
		}
		fmt.Println(testSpread.PrintList())
		fmt.Println(testSpread.PrintASCIIPerformance(my_spline{
			deg:        1,
			splineType: nil,
			x:          []float64{0.0,500.0},
			y:          nil,
			coeffs:     nil,
			unique:     false,
		},120,35))

		x := 190.0
		fmt.Printf("put1.At(%.0f)=%v\n",x,put1.At(x))

		os.Exit(123)
	*/

	//Investment testing
	/*
		var investments []Investment
		investments = append(investments,callfunc{
			base:       0,
			cost:       200,
			factor:     1,
			date:       nil,
			optionType: "call",
		})
		investments = append(investments,callfunc{
			base:       360,
			cost:       62.6,
			factor:     100,
			date:       nil,
			optionType: "call",
		})
		investments = append(investments,callfunc{
			base:       10,
			cost:       0.08,
			factor:     100,
			date:       nil,
			optionType: "put",
		})
		investments = append(investments,nonInvested{})
		for _,inv := range investments {
			fmt.Println(inv.ExpectedReturn(pdist))
		}
	*/

}


// ------------------------------- features functions -------------------------------






// ------------------------------- Polygon API specific functions -------------------------------

func PolygonAPISharePrice(ticker string, apiKey string) float64 {
	url := "https://api.polygon.io/v2/aggs/ticker/"+ticker+"/prev?adjusted=true&apiKey="+apiKey
	_,body,err := opt.APIRequest(url,1)
	check(err)
	body = strings.Split(body,"\"c\":")[1]
	body = strings.Split(body,",")[0]

	share_price,err := strconv.ParseFloat(body,64)
	check(err)
	return share_price
}

func MergeRequests(optreqs []opt.OptionURLReq, nMax int) []opt.Option{
	var options []opt.Option
	log := ""
	var msg string
	var options_tmp []opt.Option
	for _,optreq := range optreqs{
		options_tmp, msg = opt.GetOptions(optreq,nMax)
		for _,opt := range options_tmp {
			options = append(options,opt)
		}
		log += msg
	}
	opt.WriteJson("log.json",log)
	return options
}


// ------------------------------- spread specific functions -------------------------------

func (sp spread) PrintList() string {
	output := ""
	var buysell string
	for i,c := range sp.calls {
		if math.Abs(sp.weights[i]/math.Abs(sp.weights[i]) - 1) < 0.01 {buysell = "buy"}
		if math.Abs(sp.weights[i]/math.Abs(sp.weights[i]) + 1 ) < 0.01 {buysell = "sell"}
		output += fmt.Sprintf("%.1f %% : %s %.1f %s %v \n",sp.weights[i]*100,c.optionType,c.base,buysell,sp.weights[i])
	}
	return output
}

/*
func (sp spread) PrintASCIIPerformanceTransposed(pdist my_spline) string {
	output := ""
	minX := min(pdist.x)
	maxX := max(pdist.x)
	spSpline := sp.ToSpline(minX,maxX)
	minY := min(spSpline.y)
	maxY := max(spSpline.y)
	resX := 100
	resY := 100
	var plot [][]string = make([][]string,resX)
	for i := range plot {
		plot[i] = make([]string,resY)
	}
	eps := 0.5*(maxY-minY)/float64(resY)
	output += "Y "
	for j := 0 ; j < resY ; j++ {
		if math.Mod(float64(j),25) == 0 {
			output += fmt.Sprintf("%.0f ",minY+float64(j)*(maxY-minY)/float64(resY))
		} else {
			output += " "
		}
	}
	output += "\n X"
	for i := 0 ; i < resX ; i++ {
		output += "\n"
		output += fmt.Sprintf("%.0f",minX+float64(i)*(maxX-minX)/float64(resX))
		for j := 0 ; j < resY ; j++ {
			if math.Abs(spSpline.At(minX+float64(i)*(maxX-minX)/float64(resX)) - (minY+float64(j)*(maxY-minY)/float64(resY))) < eps {
				plot[i][j] = "*"
				output += "*"
			}else {
				plot[i][j] = " "
				output += " "
			}
		}
	}
	output += "\n"
	return output
}
 */


func (sp spread) PrintASCIIPerformance(minX float64, maxX float64, resX int, resY int, zeroLine bool) string {
	debug := true
	output := ""
	//minX := min(pdist.x)
	//maxX := max(pdist.x)
	spSpline := sp.ToSpline(minX,maxX)
	if debug{
		fmt.Println(spSpline)
	}
	minY := min(spSpline.y)
	maxY := max(spSpline.y)
	if debug {
		fmt.Printf("minX=%.1f maxX=%.1f minY=%.1f maxY=%.1f",minX,maxX,minY,maxY)
	}
	//resX := 120
	//resY := 35
	/*
	var plot [][]string = make([][]string,resX)
	for i := range plot {
		plot[i] = make([]string,resY)
	}
	 */
	eps := 0.5*(maxY-minY)/float64(resY)


	//testing
	//fmt.Println("maxY=",maxY,"\nDigitsBase10(maxY)=",DigitsBase10(maxY))

	var tmp string
	var maxPrintLength int
	for j := 0 ; j < resY ; j++ {
		output += "\n"

		//spacing y axis start
		if maxY-minY < 10 {
			if maxY-minY < 1000{
				tmp = fmt.Sprintf("%.6f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
			} else {
				tmp = fmt.Sprintf("%.2f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
			}
		} else {
			tmp = fmt.Sprintf("%.0f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
		}
		if maxPrintLength < len(tmp) {
			if minY<0{
				maxPrintLength = len(tmp)+1
			} else {
				maxPrintLength = len(tmp)
			}
		}

		output += repeatstr(" ",maxPrintLength-len(tmp)) + tmp + " |"
		if math.Floor(maxY-float64(j)*(maxY-minY)/float64(resY-1)) < 0 {output+=" "}

		for i := 0 ; i < resX ; i++ {
			if math.Abs(  spSpline.At( minX+float64(i+1)*(maxX-minX)/float64(resX) ) - (maxY-float64(j)*(maxY-minY)/float64(resY-1))  ) < eps {
				//plot[i][j] = "*"
				output += "*"
				//break
			}else if math.Abs((maxY-float64(j)*(maxY-minY)/float64(resY-1))) < eps && zeroLine {
				//plot[i][j] = " "
				output += "-"
			} else {
				//plot[i][j] = " "
				output += " "
			}
		}
	}

	/*
	for i := 0 ; i < resX ; i++ {
		for j := 0 ; j < resY ; j++ {
			output += plot[i][j]
		}
	}
	 */

	output+="\n"+repeatstr(" ",DigitsBase10(maxY)+3)+repeatstr("_",resX+4)+"\n"+repeatstr(" ",DigitsBase10(maxY)+4)
	var lastXLength int = 0
	for i := 0 ; i < resX ; i++ {
		if math.Mod(float64(i),10) == 0 || i == resX-1 {
			output = output[0:len(output)-1-lastXLength] + fmt.Sprintf("|%.0f  ",minX+float64(i)*(maxX-minX)/float64(resX-1))
			lastXLength = DigitsBase10(minX+float64(i)*(maxX-minX)/float64(resX))+1
			//fmt.Print("lastXLength=",lastXLength)
		} else {
			output += " "
		}
	}
	output += "\n"

	return output
}

func (ms my_spline) PrintASCII(minX float64, maxX float64, resX int, resY int, zeroLine bool) string {
	debug := false
	if debug{fmt.Println("Print ASCII debug")}
	output := ""
	//minX := min(pdist.x)
	//maxX := max(pdist.x)
	//spSpline := sp.ToSpline(minX,maxX)
	minY := min(ms.y)
	maxY := max(ms.y)
	//resX := 120
	//resY := 35
	/*
	var plot [][]string = make([][]string,resX)
	for i := range plot {
		plot[i] = make([]string,resY)
	}
	 */
	eps := 0.5*(maxY-minY)/float64(resY)


	//testing
	//fmt.Println("maxY=",maxY,"\nDigitsBase10(maxY)=",DigitsBase10(maxY))

	var tmp string
	var maxPrintLength int
	for j := 0 ; j < resY ; j++ {
		output += "\n"

		//spacing y axis start


		if maxY-minY < 10 {
			if maxY-minY < 1000{
				tmp = fmt.Sprintf("%.6f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
			} else {
				tmp = fmt.Sprintf("%.2f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
			}
		} else {
			tmp = fmt.Sprintf("%.0f",maxY-float64(j)*(maxY-minY)/float64(resY-1))
		}
		if maxPrintLength < len(tmp) {
			if minY<0{
				maxPrintLength = len(tmp)+1
			} else {
				maxPrintLength = len(tmp)
			}
		}


		output += repeatstr(" ",maxPrintLength-len(tmp)) + tmp + " |"
		if math.Floor(maxY-float64(j)*(maxY-minY)/float64(resY-1)) < 0 {output+=" "}

		for i := 0 ; i < resX ; i++ {
			if math.Abs(  ms.At( minX+float64(i+1)*(maxX-minX)/float64(resX) ) - (maxY-float64(j)*(maxY-minY)/float64(resY-1))  ) < eps {
				if debug{fmt.Printf("at:%.1f",ms.At( minX+float64(i+1)*(maxX-minX)/float64(resX) ))}
				//plot[i][j] = "*"
				output += "*"
				//break
			}else if math.Abs((maxY-float64(j)*(maxY-minY)/float64(resY-1))) < eps && zeroLine {
				//plot[i][j] = " "
				output += "-"
			} else {
				//plot[i][j] = " "
				output += " "
			}
		}
	}

	/*
		for i := 0 ; i < resX ; i++ {
			for j := 0 ; j < resY ; j++ {
				output += plot[i][j]
			}
		}
	*/

	output+="\n"+repeatstr(" ",maxPrintLength+2)+repeatstr("_",resX+4)+"\n"+repeatstr(" ",maxPrintLength+2)
	var lastXLength int = 0
	//var tmp string
	for i := 0 ; i < resX ; i++ {
		if math.Mod(float64(i),11) == 0 || i == resX-1 {
			if maxX-minX < 10 {
				tmp = fmt.Sprintf("|%.2f",minX+float64(i)*(maxX-minX)/float64(resX-1))
			} else {
				tmp = fmt.Sprintf("|%.1f",minX+float64(i)*(maxX-minX)/float64(resX-1))
			}
			output = output[0:len(output)-1-int(math.Max(0.0,float64(lastXLength-2)))] /*+ repeatstr(" ",resX/24.0)*/ + tmp
			lastXLength = len(tmp)
			//fmt.Print("lastXLength=",lastXLength)
		} else {
			output += " "
		}
	}
	output += "\n"

	return output
}

//returns number of digits of a number including the negative sign
func DigitsBase10 (number float64) int {
	return len(fmt.Sprintf("%.0f",number))
	/*
	if math.Abs(number) < 1 {
		if number < 0 {
			return 2
		} else {
			return 1
		}
	} else if number < 0 {
		return int(1+math.Log(-number)*math.Log10E)+1
	}else{
		return int(1+math.Log(number)*math.Log10E)
	}
	return 1
	 */
}

//returns best spread,bestSpreadExp, riskMatchCount, totalCount
func BestSpread2CombinationsManual(pdist my_spline, callList []callfunc, weights []float64, riskCompare bool, riskTolSpline my_spline, selling bool) (spread,float64,int,int,[]float64) {

	bestSpreadExp := -10000.0
	var bestSpread spread
	var spread_tmp spread
	var riskMatchBool bool
	spreadStart := 0
	timeTally1 := 0.0
	timeTally2 := 0.0
	timeStart := time.Now()
	timeStart1 := time.Now()
	timeStart2 := time.Now()
	var percSteps float64
	if riskCompare{
		percSteps = 0.001
	} else {
		percSteps = 0.001
	}
	spreadCount := 0
	riskMatchCount := 0
	ws := weights
	//ws := []float64{0.0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0}
	//ws := []float64{0.25,0.5,0.75}
	totalCount := (int)(4*len(ws)*len(callList)*(len(callList)-1)/2)
	for i := 0 ; i < len(callList)-1 ; i++ {
		for j := i+1 ; j < len(callList)-1 ; j++{
			for _,w := range ws {

				//buy-buy
				spread_tmp = spread{
					num:     2,
					calls:   []callfunc{callList[i],callList[j]},
					weights: []float64{w,(1.0-w)},
				}
				spreadCount++
				//totalCount++
				if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
					if riskCompare {

						timeTally1 += float64(time.Now().Sub(timeStart1).Microseconds())/1000.0
						timeStart2 = time.Now()
						rp,err := spread_tmp.riskProfile(pdist)
						timeTally2 += float64(time.Now().Sub(timeStart2).Microseconds())/1000.0
						timeStart1 = time.Now()

						rp,err = spread_tmp.riskProfile(pdist)
						if len(rp.x) == 0 || len(rp.y) == 0 {continue}
						if err == nil{
							riskMatchBool = riskMatch(rp,riskTolSpline)
						} else {
							fmt.Println("ERROR in riskProfile:",err)
						}
					} else {riskMatchBool = true}
					if riskMatchBool {riskMatchCount++}
					if riskMatchBool {
						//spreads = append(spreads,spread_tmp)
						bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
						bestSpread = spread_tmp
					}
				}


				if selling {

					//buy-sell
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{w,-(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {

							timeTally1 += float64(time.Now().Sub(timeStart1).Microseconds())/1000.0
							timeStart2 = time.Now()
							rp,err := spread_tmp.riskProfile(pdist)
							timeTally2 += float64(time.Now().Sub(timeStart2).Microseconds())/1000.0
							timeStart1 = time.Now()

							rp,err = spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}

					//sell-buy
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{-w,(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {

							timeTally1 += float64(time.Now().Sub(timeStart1).Microseconds())/1000.0
							timeStart2 = time.Now()
							rp,err := spread_tmp.riskProfile(pdist)
							timeTally2 += float64(time.Now().Sub(timeStart2).Microseconds())/1000.0
							timeStart1 = time.Now()

							rp,err = spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}


					//sell-sell
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{-w,-(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {

							timeTally1 += float64(time.Now().Sub(timeStart1).Microseconds())/1000.0
							timeStart2 = time.Now()
							rp,err := spread_tmp.riskProfile(pdist)
							timeTally2 += float64(time.Now().Sub(timeStart2).Microseconds())/1000.0
							timeStart1 = time.Now()

							rp,err = spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}

				}



				//if math.Mod(float64(spreadCount)/float64(totalCount),percSteps) < percSteps/2000 {
				if math.Mod(float64(spreadCount),percSteps*float64(totalCount)) < 4 {
					elapsed := time.Now().Sub(timeStart).Milliseconds()
					elapsedPerSpread := float64(elapsed)/float64(spreadCount-spreadStart)
					//fmt.Println(fmt.Sprintf("%.1f",100.0*float64(spreadCount)/float64(totalCount)) , "% (took "+fmt.Sprintf("%v milliseconds - %.4f ms per spread",elapsed,elapsedPerSpread)+") - est. time rem.:",fmt.Sprintf("%.2f",float64(totalCount-spreadCount)*elapsedPerSpread/1000/60),"mins")
					fmt.Printf("\r                                                                                                                                                         ")
					fmt.Printf("\r %.1f %% (compared %v spreads, took %v ms - %.4f ms per spread - %.0f spreads/s - best exp. return found: %.1f %%) - est. time rem.: %vm%.0fs",100.0*float64(spreadCount)/float64(totalCount),spreadCount,elapsed,elapsedPerSpread,1000*1.0/elapsedPerSpread,bestSpreadExp,int(float64(totalCount-spreadCount)*elapsedPerSpread/1000/60),math.Mod(float64(int(float64(totalCount-spreadCount)*elapsedPerSpread/1000)),60))
					//timeStart = time.Now()
					//spreadStart = spreadCount
				}

			}
		}
	}
	elapsed := time.Now().Sub(timeStart).Milliseconds()
	elapsedPerSpread := float64(elapsed)/float64(spreadCount-spreadStart)
	fmt.Printf("\r %.1f %% (compared %v spreads, took %v milliseconds - %.4f ms per spread - %.0f spreads/s) - est. time rem.: %vm%.0fs\n",100.0,spreadCount,elapsed,elapsedPerSpread,1000*1.0/elapsedPerSpread,int(float64(totalCount-spreadCount)*elapsedPerSpread/1000/60),math.Mod(float64(int(float64(totalCount-spreadCount)*elapsedPerSpread/1000)),60))

	timeTally1 += float64(time.Now().Sub(timeStart1).Microseconds())/1000.0

	timeTallys := []float64{timeTally1,timeTally2}

	if riskMatchCount == 0 {
		return spread{
			num:     0,
			calls:   nil,
			weights: nil,
		},bestSpreadExp,0,0,timeTallys
	}

	return bestSpread,bestSpreadExp, riskMatchCount, totalCount, timeTallys
}

//Add all combination of calls with weights to baseSpread, filter for riskTol and to bestSpreadExp and return the (new) best spread, its expected return and how many spreads have been created and checked and how many passed the risk tolerance
func BestSpreadInterationBrute(n int, baseSpread spread, bestSpreadExp float64, pdist my_spline, callList []callfunc, weights []float64, riskCompare bool, riskTolSpline my_spline, selling bool) (spread,float64,string,int,int) {


	return spread{
		num:     0,
		calls:   nil,
		weights: nil,
	},0,"",100,10
}

/*
func BestSpreadIterativeFind2(pdist my_spline, optionsList []callfunc, weightDelta float64, iterations int, riskTolSpline my_spline, starts int) (spread,float64) {

	debug := true

	var callList []callfunc
	var putList []callfunc
	//separation into call and put
	for _,o := range optionsList {
		if o.optionType == "call" || o.optionType == "long" {
			callList = append(callList,o)
		} else if o.optionType == "put" || o.optionType == "short" {
			putList = append(putList,o)
		}
	}

	//sort callList by strike
	var callStrikes []float64
	for i := range callList {
		callStrikes = append(callStrikes,callList[i].base)
	}
	callStrikesIdxs := sortAndReturnIdxsFloat64(callStrikes)
	var callListSorted []callfunc
	for i := range callList {
		callListSorted = append(callListSorted,callList[callStrikesIdxs[i]])
	}

	//sort putList by strike
	var putStrikes []float64
	for i := range putList {
		putStrikes = append(putStrikes,putList[i].base)
	}
	putStrikesIdxs := sortAndReturnIdxsFloat64(putStrikes)
	var putListSorted []callfunc
	for i := range putList {
		putListSorted = append(putListSorted,putList[putStrikesIdxs[i]])
	}

	//sort.Float64s(weights)

	fmt.Println("callListSorted=",callListSorted)
	fmt.Println("putListSorted=",putListSorted)
	fmt.Println("weightDelta=",weightDelta)

	var nc,np int
	var bestSpread spread
	var bestSpreadExp float64

	bestSpread = spread{
		num:     2,
		calls:   []callfunc{callList[np],putList[nc]},
		weights: []float64{0.5,0.5},
	}
	bestSpreadExp = bestSpread.ExpectedReturn(pdist)

	var startSpread spread
	var startSpreadExp float64


	for starts >= 0 {
		if debug{fmt.Println("starts=",starts)}

		time.Sleep(10*time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		nc = rand.Intn(len(callList)-1)
		np = rand.Intn(len(putList)-1)
		//nw := rand.Intn(len(weights))
		//fmt.Println(nc,np,nw)
		currentIdxs := []int{nc,np}

		if debug {fmt.Println("currentIdxs=",currentIdxs)}

		searchLength := 1

		startSpread = spread{
			num:     2,
			calls:   []callfunc{callList[nc],putList[np]},
			weights: []float64{0.5,0.5},
		}
		startSpreadExp = startSpread.ExpectedReturn(pdist)

		var tmp spread
		var tmpExp float64
		var tmpRiskProfile my_spline
		var err error
			//var tmpBest spread
			//var tmpBestExp float64
		startSpreadExpBefore := -1000.0


		for iterations > 0  {
			startSpreadExpBefore = startSpreadExp
				//tmpBestExp = bestSpreadExp
				//tmpBest = bestSpread

			//o1 +
			if currentIdxs[0]+searchLength < len(callList) {
				tmp = spread{
					num:     startSpread.num,
					calls:   []callfunc{callList[currentIdxs[0]+searchLength],putList[currentIdxs[1]]},
					weights: startSpread.weights,
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					currentIdxs[0] += searchLength
					//if debug {fmt.Println(startSpreadExp)}
				}
			}

			//o1 -
			if currentIdxs[0]-searchLength >= 0 {
				tmp = spread{
					num:     startSpread.num,
					calls:   []callfunc{callList[currentIdxs[0]-searchLength],putList[currentIdxs[1]]},
					weights: startSpread.weights,
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					currentIdxs[0] -= searchLength
					//if debug {fmt.Println(startSpreadExp)}
				}
			}

			//o2 +1
			if currentIdxs[1]+searchLength < len(putList) {
				tmp = spread{
					num:     startSpread.num,
					calls:   []callfunc{callList[currentIdxs[0]],putList[currentIdxs[1]+searchLength]},
					weights: startSpread.weights,
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					currentIdxs[1] += searchLength
					//if debug {fmt.Println(startSpreadExp)}
				}
			}

			//o2 -1
			if currentIdxs[1]-searchLength >= 0 {
				tmp = spread{
					num:     startSpread.num,
					calls:   []callfunc{callList[currentIdxs[0]],putList[currentIdxs[1]-searchLength]},
					weights: startSpread.weights,
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					currentIdxs[1] -= searchLength
					//if debug {fmt.Println(startSpreadExp)}
				}
			}


			//w1 +1
			if startSpread.weights[0]+float64(searchLength)*weightDelta <= 1.0 && startSpread.weights[1]-float64(searchLength)*weightDelta >= 0 {
				tmp = spread{
					num:     startSpread.num,
					calls:   startSpread.calls,
					weights: []float64{startSpread.weights[0]+float64(searchLength)*weightDelta,startSpread.weights[1]-float64(searchLength)*weightDelta},
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					//if debug {fmt.Println(startSpreadExp)}
				}
			}

			//w2 +1
			if startSpread.weights[1]+float64(searchLength)*weightDelta <= 1.0 && startSpread.weights[0]-float64(searchLength)*weightDelta >= 0 {
				tmp = spread{
					num:     startSpread.num,
					calls:   startSpread.calls,
					weights: []float64{startSpread.weights[0]-float64(searchLength)*weightDelta,startSpread.weights[1]+float64(searchLength)*weightDelta},
				}
				tmpExp = tmp.ExpectedReturn(pdist)
				tmpRiskProfile,err = tmp.riskProfile(pdist)
				check(err)
				if tmpExp > startSpreadExp && riskMatch(tmpRiskProfile,riskTolSpline){
					startSpread = tmp
					startSpreadExp = tmpExp
					//if debug {fmt.Println(startSpreadExp)}
				}
			}



			if startSpreadExp == startSpreadExpBefore {
				if searchLength < 100{
					searchLength++
				} else {
					break
				}

				//fmt.Println("searchLength=",searchLength)
			} else {
				searchLength = 1
			}

			iterations--
		}

		 if startSpreadExp > bestSpreadExp {
		 	bestSpread = startSpread
		 	bestSpreadExp = startSpreadExp
		 	if debug {fmt.Println("final:",bestSpreadExp)}
		 }

		starts--
	}




	return bestSpread,bestSpreadExp

}
 */

/*
func BestSpreadIterativeFind4(pdist my_spline, optionsList []callfunc, weightDelta float64, iterations int, riskTolSpline my_spline, starts int) (spread,float64) {

	debug := true
	selling := true

	var callList []callfunc
	var putList []callfunc
	//separation into call and put
	for _,o := range optionsList {
		if o.optionType == "call" || o.optionType == "long" {
			callList = append(callList,o)
		} else if o.optionType == "put" || o.optionType == "short" {
			putList = append(putList,o)
		}
	}

	//sort callList by strike
	var callStrikes []float64
	for i := range callList {
		callStrikes = append(callStrikes,callList[i].base)
	}
	callStrikesIdxs := sortAndReturnIdxsFloat64(callStrikes)
	var callListSorted []callfunc
	for i := range callList {
		callListSorted = append(callListSorted,callList[callStrikesIdxs[i]])
	}

	//sort putList by strike
	var putStrikes []float64
	for i := range putList {
		putStrikes = append(putStrikes,putList[i].base)
	}
	putStrikesIdxs := sortAndReturnIdxsFloat64(putStrikes)
	var putListSorted []callfunc
	for i := range putList {
		putListSorted = append(putListSorted,putList[putStrikesIdxs[i]])
	}

	//sort.Float64s(weights)

	fmt.Println("callListSorted=",callListSorted)
	fmt.Println("putListSorted=",putListSorted)
	fmt.Println("weightDelta=",weightDelta)

	var nc1,np1,nc2,np2 int
	var bestSpread spread
	var bestSpreadExp float64

	bestSpread = spread{
		num:     4,
		calls:   []callfunc{callList[nc1],putList[np1],callList[nc2],putList[np2]},
		weights: []float64{0.25,0.25,0.25,0.25},
	}
	bestSpreadExp = bestSpread.ExpectedReturn(pdist)

	var startSpread spread
	var startSpreadExp float64


	for starts >= 0 {
		if debug{fmt.Println("starts=",starts)}

		time.Sleep(10*time.Millisecond)
		rand.Seed(time.Now().UnixNano())
		nc1 = rand.Intn(len(callList)-1)
		np1 = rand.Intn(len(putList)-1)
		nc2 = rand.Intn(len(callList)-1)
		np2 = rand.Intn(len(putList)-1)
		//nw := rand.Intn(len(weights))
		//fmt.Println(nc,np,nw)
		currentIdxs := []int{nc1,np1,nc2,np2}

		if debug {fmt.Println("currentIdxs=",currentIdxs)}

		searchLength := 1

		startSpread = spread{
			num:     4,
			calls:   []callfunc{callList[nc1],putList[np1],callList[nc2],putList[np2]},
			weights: []float64{0.25,0.25,0.25,0.25},
		}
		startSpreadExp = startSpread.ExpectedReturn(pdist)


		//var tmp spread
		//var tmpExp float64
		//var tmpRiskProfile my_spline
		//var err error
		//var tmpCalls []callfunc
		//var tmpWeightsAlterable []bool
		//var tmpWeightsAlterableCount int

		var alterIndex int
		var alterDirection int

		startSpreadExpBefore := -1000.0


		for iterations > 0 {
			startSpreadExpBefore = startSpreadExp
				//tmpBestExp = bestSpreadExp
				//tmpBest = bestSpread

			//o1 +
			alterIndex = 0
			alterDirection = 1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)


			//o1 -
			alterIndex = 0
			alterDirection = -1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)

			//o2 +
			//for o2 it should be putList
			alterIndex = 1
			alterDirection = 1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)


			//o2 -
			alterIndex = 1
			alterDirection = -1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)


			//o3 +
			alterIndex = 2
			alterDirection = 1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)

			//o3 -
			alterIndex = 2
			alterDirection = -1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)

			//o4 +
			alterIndex = 3
			alterDirection = 1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)

			//o4 -
			alterIndex = 3
			alterDirection = -1
			startSpread,startSpreadExp,currentIdxs,_,_ = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)


			//w1 +
			alterIndex = 0
			alterDirection = 1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w1 -
			alterIndex = 0
			alterDirection = -1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w2 +
			alterIndex = 1
			alterDirection = 1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w2 -
			alterIndex = 1
			alterDirection = -1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w3 +
			alterIndex = 2
			alterDirection = 1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w3 -
			alterIndex = 2
			alterDirection = -1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w4 +
			alterIndex = 3
			alterDirection = 1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)

			//w4 -
			alterIndex = 3
			alterDirection = -1
			startSpread, startSpreadExp,_,_ = alterWeights(alterIndex, alterDirection, weightDelta, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, selling, debug)


			if (startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore < 0.001 {
				if searchLength < 100{
					searchLength++
					if debug {fmt.Println("searchLength=",searchLength)}
				} else {
					searchLength = 1
					break
				}

				//fmt.Println("searchLength=",searchLength)
			}

			iterations--
			if math.Mod(float64(iterations),1000.0) < 1{
				if debug {
					fmt.Println("iterations=",iterations)
					fmt.Println("(startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore=",(startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore)
				}
			}
		}

		if startSpreadExp > bestSpreadExp {
			bestSpread = startSpread
			bestSpreadExp = startSpreadExp
			if debug {fmt.Println("final:",bestSpreadExp)}
		}

		starts--
	}




	return bestSpread,bestSpreadExp

}
 */


//return bestSpread,bestSpreadExp,startRiskProfile, riskTotalCount,riskMatchCount,totalSpreadsCompared
func BestSpreadIterativeFindN(N int, pdist my_spline, optionsList []callfunc, weightDelta float64, iterations int, riskTolSpline my_spline, starts int, contractTypes []string, maxSearchLength int, share_price float64) (spread,float64,my_spline,int,int,int) {

	debug := false
	selling := true
	info := true

	startStarts := starts
	//startIterations := iterations

	riskTotalCount := 0
	riskMatchCount := 0

	totalSpreadsCompared := 0

	if len(contractTypes) != N {
		fmt.Println("Error in BestSpreadIterativeFindN: len(contractTypes) != N !")
		os.Exit(0)
	}

	var callList []callfunc
	var putList []callfunc
	//separation into call and put
	for _,o := range optionsList {
		if o.optionType == "call" || o.optionType == "long" {
			callList = append(callList,o)
		} else if o.optionType == "put" || o.optionType == "short" {
			putList = append(putList,o)
		}
	}

	//sort callList by strike
	var callStrikes []float64
	for i := range callList {
		callStrikes = append(callStrikes,callList[i].base)
	}
	callStrikesIdxs := sortAndReturnIdxsFloat64(callStrikes)
	var callListSorted []callfunc
	for i := range callList {
		callListSorted = append(callListSorted,callList[callStrikesIdxs[i]])
	}

	//sort putList by strike
	var putStrikes []float64
	for i := range putList {
		putStrikes = append(putStrikes,putList[i].base)
	}
	putStrikesIdxs := sortAndReturnIdxsFloat64(putStrikes)
	var putListSorted []callfunc
	for i := range putList {
		putListSorted = append(putListSorted,putList[putStrikesIdxs[i]])
	}

	//sort.Float64s(weights)

	/*
	fmt.Println("callListSorted=",callListSorted)
	fmt.Println("putListSorted=",putListSorted)
	fmt.Println("weightDelta=",weightDelta)
	 */

	//var nc1,np1,nc2,np2 int

	var bestSpread spread
	var bestSpreadExp float64
	var currentIdxs []int = make([]int,N)

	bestSpread = spread{
		num:     N,
		calls:   []callfunc{},
		weights: []float64{},
	}


	//random start
	/*
	var ran int
	for i := 0 ; i <= N-1 ; i++{
		time.Sleep(10*time.Nanosecond)
		rand.Seed(time.Now().UnixNano())
		switch contractTypes[i] {
		case "call":
			ran = rand.Intn(len(callList)-1)
			currentIdxs[i] = ran
			bestSpread.calls = append(bestSpread.calls,callList[ran])
		case "put":
			ran = rand.Intn(len(putList)-1)
			currentIdxs[i] = ran
			bestSpread.calls = append(bestSpread.calls,putList[ran])
		}

		bestSpread.weights = append(bestSpread.weights,1.0/float64(N))
	}
	 */


	//long and short start
	/*
	var long_index int = -1
	var short_index int = -1
	for i,c := range callList {
		if c.optionType == "long" {
			long_index = i
		}
	}
	for i,p := range putList {
		if p.optionType == "short" {
			short_index = i
		}
	}
	if long_index == -1 || short_index == -1 {
		fmt.Println("Didn't fing long or short in callList and putList.")
		os.Exit(3)
	}
	for i := 0 ; i <= N-1 ; i++ {
		if math.Mod(float64(i),2)==0{
			bestSpread.calls = append(bestSpread.calls, callList[long_index])
		} else {
			bestSpread.calls = append(bestSpread.calls, putList[short_index])
		}
	}
	 */


	// start at long (first) and short (last)
	for i := 0 ; i <= N-1 ; i++ {
		if contractTypes[i] == "call" {
			bestSpread.calls = append(bestSpread.calls, callList[0])
		} else if contractTypes[i] == "put" {
			bestSpread.calls = append(bestSpread.calls, putList[len(putList)-1])
		}
		bestSpread.weights = append(bestSpread.weights,1.0/float64(N))
	}


	/*
	if debug {
		fmt.Println("Initial bestSpread:",bestSpread)
		fmt.Println("Perform alterWeights()")
		bestSpread,bestSpreadExp,_,_ = alterWeights(0,1,0.1,1,bestSpread,pdist,riskTolSpline,bestSpreadExp,false,true)
		fmt.Println("bestSpread now:",bestSpread)
		os.Exit(1)
	}
	 */



	bestSpreadExp = bestSpread.ExpectedReturn(pdist)

	var startSpread spread
	var startSpreadExp float64
	var startRiskProfile my_spline
	var err error

	//progress printing variables
	/*
	var timeStart time.Time = time.Now()
	var spreadStart int = 0
	percSteps := 0.1
	 */


	for starts >= 0 {
		if debug{fmt.Println("starts=",starts)}

		searchLength := 1

		/*
		time.Sleep(10*time.Nanosecond)
		rand.Seed(time.Now().UnixNano())
		nc1 = rand.Intn(len(callList)-1)
		np1 = rand.Intn(len(putList)-1)
		nc2 = rand.Intn(len(callList)-1)
		np2 = rand.Intn(len(putList)-1)
		//nw := rand.Intn(len(weights))
		//fmt.Println(nc,np,nw)
		currentIdxs = []int{nc1,np1,nc2,np2}

		startSpread = spread{
				num:     4,
				calls:   []callfunc{callList[nc1],putList[np1],callList[nc2],putList[np2]},
				weights: []float64{0.25,0.25,0.25,0.25},
			}
			startSpreadExp = startSpread.ExpectedReturn(pdist)
		 */


		//setup a new starting spread
		startSpread = spread{
			num:     N,
			calls:   []callfunc{},
			weights: []float64{},
		}


		if starts == startStarts {
			// first and last
			for i := 0 ; i <= N-1 ; i++ {
				if contractTypes[i] == "call" {
					startSpread.calls = append(startSpread.calls, callList[0])
					currentIdxs[i] = 0
				} else if contractTypes[i] == "put" {
					startSpread.calls = append(startSpread.calls, putList[len(putList)-1])
					currentIdxs[i] = len(putList)-1
				}
				startSpread.weights = append(startSpread.weights,1.0/float64(N))
			}
		} else {
			/*
			riskMatching := false
			for !riskMatching {

			 */

				//random
				var ran int
				for i := 0 ; i <= N-1 ; i++{
					time.Sleep(10*time.Nanosecond)
					rand.Seed(time.Now().UnixNano())
					switch contractTypes[i] {
					case "call":
						ran = rand.Intn(len(callList)-1)
						currentIdxs[i] = ran
						startSpread.calls = append(startSpread.calls,callList[ran])
					case "put":
						ran = rand.Intn(len(putList)-1)
						currentIdxs[i] = ran
						startSpread.calls = append(startSpread.calls,putList[ran])
					}

					startSpread.weights = append(startSpread.weights,1.0/float64(N))
				}

				/*
				startSpreadRiskProfile,err := startSpread.riskProfile(pdist)
				check(err)
				riskMatching = riskMatch(startSpreadRiskProfile,riskTolSpline)
			}
				 */

		}


		startSpreadExp = startSpread.ExpectedReturn(pdist)



		if debug {fmt.Println("currentIdxs=",currentIdxs)}

		/*
			var tmp spread
			var tmpExp float64
			var tmpRiskProfile my_spline
			var err error
			var tmpCalls []callfunc
			var tmpWeightsAlterable []bool
				var tmpWeightsAlterableCount int
		*/
		//var alterIndex int
		var alterDirection int

		startSpreadExpBefore := -1000.0
		startRiskProfile,err = startSpread.riskProfile(pdist)
		check(err)

		var tmpRiskTotalCount int
		var tmpRiskMatchCount int

		var tmpBest spread = startSpread
		var tmpBestExp float64 = startSpreadExp
		var tmpBestRiskProfile my_spline = startRiskProfile

		var tmpBestCurrentIdxs []int = currentIdxs
		var tmp spread
		var tmpCurrentIdxs []int

		for iterations > 0 /*&& math.Abs((bestSpreadExp-bestSpreadExpBefore)/bestSpreadExpBefore) > 0.01*/ {
			if info && math.Mod(float64(iterations),10) < 1 {
				//bestSpread.PrintList()
				fmt.Printf("\rCompared %v spreads, current best expected return found: %.3f %%",totalSpreadsCompared,bestSpreadExp)
			}
			startSpreadExpBefore = startSpreadExp
			/*
				tmpBestExp = bestSpreadExp
				tmpBest = bestSpread
			*/

			totalSpreadsCompared += 4*N


			//use tmp and go direction of best change, not first positive change
			for alterIndex := 0 ; alterIndex < N ; alterIndex++ {

				//o +
				alterDirection = 1
				//startSpread,startSpreadExp, startRiskProfile, currentIdxs,tmpRiskTotalCount,tmpRiskMatchCount = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)
				tmp, tmpCurrentIdxs = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)
				tmpBest, tmpBestExp, tmpBestRiskProfile, tmpBestCurrentIdxs, riskTotalCount, riskMatchCount, err = betterSpread(tmpBest,tmp,tmpBestCurrentIdxs,tmpCurrentIdxs,riskTolSpline,pdist)
				check(err)
				riskTotalCount += tmpRiskTotalCount
				riskMatchCount += tmpRiskMatchCount

				//o -
				alterDirection = -1
				//startSpread,startSpreadExp, startRiskProfile, currentIdxs,tmpRiskTotalCount,tmpRiskMatchCount = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)
				tmp, tmpCurrentIdxs = alterCall(alterIndex, alterDirection, callList, putList, currentIdxs, searchLength, startSpread, pdist, riskTolSpline, startSpreadExp, debug)
				//fmt.Println(startSpread.calls)
				//fmt.Println(tmpBest.calls)
				tmpBest, tmpBestExp, tmpBestRiskProfile, tmpBestCurrentIdxs, riskTotalCount, riskMatchCount, err = betterSpread(tmpBest,tmp,tmpBestCurrentIdxs,tmpCurrentIdxs,riskTolSpline,pdist)
				check(err)
				riskTotalCount += tmpRiskTotalCount
				riskMatchCount += tmpRiskMatchCount

				//w +
				alterDirection = 1
				//startSpread, startSpreadExp, startRiskProfile, tmpRiskTotalCount,tmpRiskMatchCount = alterWeights(alterIndex, alterDirection, weightDelta, 1, startSpread, pdist, riskTolSpline, startSpreadExp, startRiskProfile, selling, debug)
				tmp = alterWeights(alterIndex, alterDirection, weightDelta, 1, startSpread, pdist, riskTolSpline, startSpreadExp, startRiskProfile, selling, debug)
				tmpBest, tmpBestExp, tmpBestRiskProfile, tmpBestCurrentIdxs, tmpRiskTotalCount, tmpRiskMatchCount, err = betterSpread(tmpBest,tmp,tmpBestCurrentIdxs,tmpCurrentIdxs,riskTolSpline,pdist)
				check(err)
				riskTotalCount += tmpRiskTotalCount
				riskMatchCount += tmpRiskMatchCount

				//w -
				alterDirection = -1
				//startSpread, startSpreadExp, startRiskProfile, tmpRiskTotalCount,tmpRiskMatchCount = alterWeights(alterIndex, alterDirection, weightDelta, 1, startSpread, pdist, riskTolSpline, startSpreadExp, startRiskProfile, selling, debug)
				tmp = alterWeights(alterIndex, alterDirection, weightDelta, 1, startSpread, pdist, riskTolSpline, startSpreadExp, startRiskProfile, selling, debug)
				tmpBest, tmpBestExp, tmpBestRiskProfile, tmpBestCurrentIdxs, riskTotalCount, riskMatchCount, err = betterSpread(tmpBest,tmp,tmpBestCurrentIdxs,tmpCurrentIdxs,riskTolSpline,pdist)
				check(err)
				riskTotalCount += tmpRiskTotalCount
				riskMatchCount += tmpRiskMatchCount

				//o+ w-


				//implement combination of directions ; in other words: diagonal movement
				//also make tmp and select best tmp instead of first?
				//don't already check ExpReturn and riskTol, just return arrived
				//Gridsearch?
				//Powell-Wolff?

				/*
				for alterIndex2 := 0 ; alterIndex2 < N ; alterIndex2++ {
					//o1+ o2+
				}
				 */





				startSpread = tmpBest
				startSpreadExp = tmpBestExp
				startRiskProfile = tmpBestRiskProfile

			}


			//if debug {fmt.Printf("(startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore=%.4f\n",(startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore)}
			if (startSpreadExp - startSpreadExpBefore)/startSpreadExpBefore < 0.00000001 {
				if searchLength < maxSearchLength {
					searchLength++
					//if debug {fmt.Println("searchLength=",searchLength)}
				} else {
					searchLength = 1
					break
				}

				//fmt.Println("searchLength=",searchLength)
			} else {
				searchLength = 1
			}


			//progress print
			/*
			spreadCount := 4*N*(startIterations-iterations)*(startStarts-starts)
			totalCount := 4*N*startIterations*startStarts
			if math.Mod(float64(spreadCount),percSteps*float64(totalCount)) < 4 {
				elapsed := time.Now().Sub(timeStart).Nanoseconds()/1000.0
				elapsedPerSpread := float64(elapsed)/float64(spreadCount-spreadStart)
				fmt.Println(fmt.Sprintf("%.1f",100.0*float64(spreadCount)/float64(totalCount)) , "% (took "+fmt.Sprintf("%v milliseconds - %.4f ms per spread",elapsed,elapsedPerSpread)+") - est. time rem.:",fmt.Sprintf("%.1f",float64(totalCount-spreadCount)*elapsedPerSpread/1000),"seconds")
				timeStart = time.Now()
				spreadStart = spreadCount
			}
			 */


			iterations--
		}

		if startSpreadExp > bestSpreadExp {
			bestSpread = startSpread
			bestSpreadExp = startSpreadExp
			if debug {fmt.Println("final:",bestSpreadExp)}
		}

		starts--
	}
	if info {
		//bestSpread.PrintList()
		fmt.Printf("\rCompared %v spreads, current best expected return found: %.3f %%",totalSpreadsCompared,bestSpreadExp)
	}
	if info{fmt.Println("")}



	/*
		separate into callList and putList and sort them by strike
		(this way you won't find mispricings and might get stuck due to a mispricing; go perhaps 5 steps in the direction and compare them all)
	*/

	return bestSpread,bestSpreadExp,startRiskProfile, riskTotalCount,riskMatchCount,totalSpreadsCompared

}

//for diagonal movement, divided into quadrants (through alterDirectionOption and alterDirectionWeight)
/*
func alterCall2Weights(alterIndex1 int, alterDirectionOption1 int, alterIndex2 int, alterDirectionOption2 int, alterDirectionWeight int, callList []callfunc, putList []callfunc , currentIdxs []int, searchLength int, startSpread spread, pdist my_spline, riskTolSpline my_spline, startSpreadExp float64, debug bool, weightDelta float64, selling bool) (spread,float64,my_spline,[]int,int,int) {
	var searchLengthWeight int

	tmpOptionType := startSpread.calls[alterIndex].optionType
	var optList []callfunc
	//var startRiskProfile my_spline
	var tmp spread
	//var tmpExp float64
	if tmpOptionType == "call" || tmpOptionType == "long" {
		optList = callList
	} else if tmpOptionType == "put" || tmpOptionType == "short" {
		optList = putList
	} else {
		fmt.Println("Error in alterCall(): OptionType neither call,long,put, nor short!")
		os.Exit(420)
	}

	opt1Bounds := currentIdxs[alterIndex1]+alterDirectionOption1*searchLength < len(optList) && currentIdxs[alterIndex1]+alterDirectionOption1*searchLength >= 0
	opt2Bounds := currentIdxs[alterIndex2]+alterDirectionOption2*searchLength < len(optList) && currentIdxs[alterIndex2]+alterDirectionOption2*searchLength >= 0
	//wgt1Bounds := currentIdxs[alterIndex1]+alterDirectionOption1*searchLength < len(optList) && currentIdxs[alterIndex1]+alterDirectionOption1*searchLength >= 0

	if  {

	}


	for searchLengthCall := 0 ; searchLengthCall < searchLength ; searchLengthCall++ {
		searchLengthWeight = searchLength-searchLengthCall

	}

	fmt.Println(searchLengthWeight)
	return spread{},0,my_spline{},nil,0,0
}
 */


func alterCall(alterIndex int, alterDirection int, callList []callfunc, putList []callfunc , currentIdxs []int, searchLength int, startSpread spread, pdist my_spline, riskTolSpline my_spline, startSpreadExp float64, debug bool) /*(spread,float64,my_spline,[]int,int,int)*/(spread,[]int) {
	/*
	riskMatchCount := 0
	riskTotalCount := 0
	 */
	tmpOptionType := startSpread.calls[alterIndex].optionType
	var optList []callfunc
	//var startRiskProfile my_spline
	var tmp spread
	//var tmpExp float64
	if tmpOptionType == "call" || tmpOptionType == "long" {
		optList = callList
	} else if tmpOptionType == "put" || tmpOptionType == "short" {
		optList = putList
	} else {
		fmt.Println("Error in alterCall(): OptionType neither call,long,put, nor short!")
		os.Exit(420)
	}
	//if debug {fmt.Println("currentIdxs[alterIndex]+alterDirection*searchLength=",currentIdxs[alterIndex]+alterDirection*searchLength)}
	if currentIdxs[alterIndex]+alterDirection*searchLength < len(optList) && currentIdxs[alterIndex]+alterDirection*searchLength >= 0 {

		tmpCalls := startSpread.calls
		tmpCalls[alterIndex] = optList[currentIdxs[alterIndex]+alterDirection*searchLength]
		tmp = spread{
			num:     startSpread.num,
			calls:   tmpCalls,
			weights: startSpread.weights,
		}

		currentIdxs[alterIndex] += alterDirection*searchLength

		/*
		//tmpExp = tmp.ExpectedReturn(pdist)
		if tmpExp > startSpreadExp{
			tmpRiskProfile,err := tmp.riskProfile(pdist)
			check(err)
			riskTotalCount++
			if riskMatch(tmpRiskProfile,riskTolSpline) {
				riskMatchCount++
				startSpread = tmp
				startSpreadExp = tmpExp
				currentIdxs[alterIndex] += alterDirection*searchLength
				startRiskProfile = tmpRiskProfile
				if debug {fmt.Println(startSpreadExp)}
			}
		}
		 */
	} else {
		return startSpread,currentIdxs
	}
	//return startSpread,startSpreadExp,startRiskProfile,currentIdxs,riskTotalCount,riskMatchCount
	return tmp,currentIdxs
}

func alterWeights(alterIndex int, alterDirection int, weightDelta float64, searchLength int, startSpread spread, pdist my_spline, riskTolSpline my_spline, startSpreadExp float64, startRiskProfile my_spline, selling bool, debug bool) /*(spread,float64,my_spline,int,int)*/ spread {
	/*
	riskMatchCount := 0
	riskTotalCount := 0
	 */
	var tmp spread
	if math.Abs( startSpread.weights[alterIndex] + float64(alterDirection*searchLength)*weightDelta ) <= 1.0 {
		if !selling && startSpread.weights[alterIndex] + float64(alterDirection*searchLength)*weightDelta < 0.0 {
			return startSpread
		}
		tmpWeights := startSpread.weights
		tmpWeights[alterIndex] = tmpWeights[alterIndex]+float64(alterDirection*searchLength)*weightDelta

		SumAbs := 0.0
		for i := range tmpWeights{
			SumAbs += math.Abs(tmpWeights[i])
		}
		for i := range tmpWeights {
			//if i != alterIndex {
				tmpWeights[i] = tmpWeights[i] / SumAbs
			//}
		}
		tmp = spread{
			num:     startSpread.num,
			calls:   startSpread.calls,
			weights: tmpWeights,
		}

		/*
		tmpExp := tmp.ExpectedReturn(pdist)
		if tmpExp > startSpreadExp{
			riskTotalCount++
			tmpRiskProfile,err := tmp.riskProfile(pdist)
			check(err)
			if riskMatch(tmpRiskProfile,riskTolSpline) {
				riskMatchCount++
				startSpread = tmp
				startSpreadExp = tmpExp
				startRiskProfile = tmpRiskProfile
				if debug {fmt.Println(startSpreadExp)}
			}
		}
		 */
	} else {
		return startSpread
	}
	//return startSpread,startSpreadExp,startRiskProfile,riskTotalCount,riskMatchCount
	return tmp
}

func betterSpread(s1 spread, s2 spread, currentIdx1 []int, currentIdx2 []int, riskTolSpline my_spline, pdist my_spline) (spread,float64,my_spline,[]int,int,int,error) {
	debug := false
	s1Exp := s1.ExpectedReturn(pdist)
	s2Exp := s2.ExpectedReturn(pdist)
	if debug{
		fmt.Println("s1Exp=",s1Exp)
		fmt.Println("s2Exp=",s2Exp)
	}
	riskMatchCount := 0
	riskTotalCount := 0

	var s1RiskProfile my_spline
	var s2RiskProfile my_spline
	var err error

	if s1Exp > s2Exp {
		s1RiskProfile, err = s1.riskProfile(pdist)
		check(err)
		riskTotalCount++
		if riskMatch(s1RiskProfile, riskTolSpline) {
			riskMatchCount++
			return s1,s1Exp,s1RiskProfile,currentIdx1,riskTotalCount,riskMatchCount,nil
			//if debug {fmt.Println(startSpreadExp)}
		} else {
			s2RiskProfile, err = s2.riskProfile(pdist)
			check(err)
			riskTotalCount++
			if riskMatch(s2RiskProfile, riskTolSpline) {
				riskMatchCount++
				return s2,s2Exp,s2RiskProfile,currentIdx2,riskTotalCount,riskMatchCount,nil
				//if debug {fmt.Println(startSpreadExp)}
			}
		}
	} else {
		s2RiskProfile, err = s2.riskProfile(pdist)
		check(err)
		riskTotalCount++
		if riskMatch(s2RiskProfile, riskTolSpline) {
			riskMatchCount++
			return s2,s2Exp,s2RiskProfile,currentIdx2,riskTotalCount,riskMatchCount,nil
			//if debug {fmt.Println(startSpreadExp)}
		} else {
			//tmpRiskProfile, err = s1.riskProfile(pdist)
			check(err)
			riskTotalCount++
			if riskMatch(s1RiskProfile, riskTolSpline) {
				riskMatchCount++
				return s1,s1Exp,s1RiskProfile,currentIdx1,riskTotalCount,riskMatchCount,nil
				//if debug {fmt.Println(startSpreadExp)}
			}
		}
	}
	if len(s1RiskProfile.x) == 0 {
		s1RiskProfile, err = s1.riskProfile(pdist)
	}
	if debug{fmt.Println("warning: neither spread in betterSpread() satisfies riskTol!")}
	return s1,-1000.0,s1RiskProfile,currentIdx1,riskTotalCount,riskMatchCount,nil
}


//building all spreads first and only then comparing is very expensive w.r.t. RAM usage and should be avoided. Compare while building the spreads already and check for risk tolerance and discard bad spreads
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

func(sp spread) Add (c callfunc, w float64) spread {
	var weights_new []float64
	for i := range sp.weights {
		weights_new = append(weights_new,sp.weights[i]*(1-w))
	}
	weights_new = append(weights_new,w)

	return spread{
		num:     sp.num+1,
		calls: append(sp.calls, c),
		weights: weights_new,
	}
}

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

func (sp spread) riskProfile(pdist my_spline) (my_spline,error) {
	debug := false
	if debug {fmt.Println("riskProfile debug:")}
	spreadPerf := sp.ToSpline(min(pdist.x),max(pdist.x))
	var ys []float64
	var probs []float64

	n := 10
	tolYPerc := 0.01 //basically in NewtonRoot()
	var probsMap map[float64]float64
	probsMap = make(map[float64]float64)
	start := min(spreadPerf.y)
	end := max(spreadPerf.y)
	dy := (end-start)/float64(n)
	if math.Abs(dy) < 0.00000001 {
		return NewSpline([]string{"1","2"},[]float64{min(pdist.x),max(pdist.x)},[]float64{min(spreadPerf.y),max(spreadPerf.y)}), nil
	}
	if debug {
		fmt.Println("dy: ",dy)
		fmt.Println("min(spreadPerf.y) :",min(spreadPerf.y) )
		fmt.Println( "max(spreadPerf.y) :",max(spreadPerf.y) )
		//os.Exit(0)
	}

	for y := start  /*-2*math.Abs(dy*min(spreadPerf.y))*/ ; y <= end ; y += dy /*dy + dy*(-math.Pow (y - (start + end)/2, 2) + math.Pow ((start + end)/2, 2))/math.Pow ((start + end)/2, 2)*/ { //normally dy but I'm trying to concentrate more towards start and end
		if debug {fmt.Println("y: ",y)}
		neg,_ := spreadPerf.PosNegRange(y,tolYPerc,30/*4*n*/)
		if debug {fmt.Println(neg)}
		if len(neg) == 0 {
			probs = append(probs,0.0)
			probsMap[y] = 0.0
			ys = append(ys, y)
		}

		prob_tmp := 0.0
		for i := range neg {
			prob_tmp += pdist.IntegralSpline(neg[i][0],neg[i][1])
		}
		//discard decreasing probs - not an ideal solution! kinda hacky, messy
		if len(probs) == 0 || prob_tmp > probs[len(probs)-1] {
			probs = append(probs,prob_tmp)
			probsMap[y] = prob_tmp
			ys = append(ys, y)
		}
		if debug {
			fmt.Println("probs: ",probs)
			fmt.Println("ys: ",ys)
		}
	}

	//sort probs and accordingly change ys
	//ideally probs should already be sorted
	/*
		idx := sortAndReturnIdxsFloat64(probs)
		var ys_tmp []float64
		for i := range ys {
			ys_tmp = append(ys_tmp,ys[idx[i]])
		}
		ys = ys_tmp
	*/


	//sort ys and accordingly change probs
	/*
		sortAndReturnIdxsFloat64(ys)
		fmt.Println("(probs,ys):")
		var probs2 []float64
		probs2 = make([]float64,len(ys))
		for i := range probs {
			fmt.Println("(",probsMap[ys[i]],",",ys[i],")")
			probs2[i] = probsMap[ys[i]]
		}
		probs = probs2
	*/


	//splinetype = []string{"3","2","=Sl","=Cv","E0Sl"}
	splinetype := []string{"1","2"}
	if len(probs) != len(ys){
		fmt.Println("len(probs)!=len(ys)!!")
		os.Exit(0)
	}
	riskSpline := NewSpline(splinetype,probs,ys)

	if len(riskSpline.x) == 0 || len(riskSpline.y) == 0 {
		return my_spline{}, fmt.Errorf("x or y empty")
	}
	return riskSpline,nil
}






// ------------------------------- my spline specific functions -------------------------------

func constSpline(c float64, xrange []float64) my_spline{
	return my_spline{
		deg:        0,
		splineType: nil,
		x:          []float64{xrange[0],xrange[1]},
		y:          []float64{c,c},
		coeffs:     []float64{c},
		unique:     false,
	}
}

func riskMatch(riskProfile my_spline, riskTol my_spline) bool {
	dx := 0.001
	for x := 0.0 ; x <= 1.0 ; x += dx {
		if riskProfile.At(x) < riskTol.At(x) {return false}
	}
	return true
}

func riskPlotRange(ms1 my_spline, ms2 my_spline) string {
	return fmt.Sprintf("{{0,1},{%.1f,%.1f}}",math.Min(ms1.At(0.0),ms2.At(0.0)),math.Max(ms1.At(1.0),ms2.At(1.0)))
}

func PlotRange(pdist my_spline, ms1 my_spline, ms2 my_spline) string {
	return fmt.Sprintf("{{%.1f,%.1f},{%.1f,%.1f}}",min(pdist.x),max(pdist.x),math.Min(ms1.At(min(pdist.x)),ms2.At(min(pdist.x))),math.Max(ms1.At(max(pdist.x)),ms2.At(max(pdist.x))))
}

func RiskProfileIterativeMeanReturn(riskProfile my_spline) {
	fmt.Println("RiskProfileIterativeMeanReturn:")
	/*
		//norming
		normedRP := riskProfile.Add(constSpline(-min(riskProfile.y),[]float64{min(riskProfile.x),max(riskProfile.x)}))
		normedRP = normedRP.Factor(1.0/(max(normedRP.y)))
		fmt.Println("riskProfile normed")
		fmt.Println(normedRP.PrintASCII(0.0,1.0,130,35))
		fmt.Println("riskProfile Integrated")
		fmt.Println(normedRP.Integrate().PrintASCII(0.0,1.0,130,35))
	*/

	runs := 1000000
	iterations := 20
	fmt.Println("iterations=",iterations)
	fmt.Println("(avg-ing out) runs=",runs)
	var returns []float64
	/*
		var tmpOld float64
		var tmp float6
		tmpOld = 0.0
		tmp = 1.0
	*/

	//tolPerc := 0.0001

	//var tmp []float64 = []float64{0.0,1.0}
	var cagrs []float64

	var tmpstr string
	for i:=0;i<runs;i++{
		iteration := 0
		tmp := []float64{0.0,1.0}
		for /*math.Abs(tmp[len(tmp)-2]-tmp[len(tmp)-1])/tmp[len(tmp)-2] > tolPerc && tmp[len(tmp)-1] > tolPerc &&*/ iteration <= iterations {
			randReturn := riskProfile.At(rand.Float64())
			//fmt.Println("rand=",randReturn)
			tmp = append(tmp,tmp[len(tmp)-1] * (1.0+randReturn/100.0) )
			//fmt.Println("tmp=",tmp[len(tmp)-1],"=",tmp[len(tmp)-2],"*(1+",randReturn/100.0,")")
			iteration++
		}
		//fmt.Println("final:",tmp[len(tmp)-1])
		returns = append(returns,tmp[len(tmp)-1]-1)
		cagrs = append(cagrs,math.Pow(tmp[len(tmp)-1],1.0/float64(iteration))-1.0)

		//Print override testing
		if math.Mod(float64(i),1000000)<1 {
			tmpstr += fmt.Sprintf("\r%s",repeatstr(" ",len(tmpstr)))
			sort.Float64s(returns)
			//fmt.Println("returns=",returns)
			tmpstr = "\r"
			tmpstr += fmt.Sprintf("Mean[returns]=%.2f %% ",returns[len(returns)/2]*100)
			tmpstr += fmt.Sprintf("Worst[returns]=%.2f %% ",returns[0]*100)
			tmpstr += fmt.Sprintf("Best[returns]=%.2f %% ",returns[len(returns)-1]*100)
			sort.Float64s(cagrs)
			tmpstr += fmt.Sprintf("Mean[cagrs]=%.2f %% ",cagrs[len(cagrs)/2]*100)
			tmpstr += fmt.Sprintf("Worst[cagrs]=%.2f %% ",cagrs[0]*100)
			tmpstr += fmt.Sprintf("Best[cagrs]=%.2f %% ",cagrs[len(cagrs)-1]*100)
			fmt.Print(tmpstr)
		}


	}
	fmt.Printf("\r%s\n",repeatstr(" ",len(tmpstr)))
	sort.Float64s(returns)
	//fmt.Println("returns=",returns)
	fmt.Printf("Mean[returns]=%.2f %% ",returns[len(returns)/2]*100)
	fmt.Printf("Worst[returns]=%.2f %% ",returns[0]*100)
	fmt.Printf("Best[returns]=%.2f %% \n",returns[len(returns)-1]*100)

	sort.Float64s(cagrs)
	fmt.Printf("Mean[cagrs]=%.2f %% ",cagrs[len(cagrs)/2]*100)
	fmt.Printf("Worst[cagrs]=%.2f %% ",cagrs[0]*100)
	fmt.Printf("Best[cagrs]=%.2f %% \n",cagrs[len(cagrs)-1]*100)

	resX := 130-1
	xs := []float64{}
	ys := []float64{}
	for i := 0 ; i <= resX ; i++ {
		idx := int(float64(len(cagrs)-1)/float64(resX))*i
		xs = append(xs,float64(idx))
		ys = append(ys,cagrs[idx]*100)
		//fmt.Printf("idx= %v  -  cagrs[idx]= %.2f\n",idx,cagrs[idx])
	}
	splineType := []string{"1","2"}
	IterativeRiskProfile := NewSpline(splineType,xs,ys)
	fmt.Println("IterativeRiskProfile CAGR:")
	fmt.Println(IterativeRiskProfile.PrintASCII(0,float64(len(cagrs)),130,35,true))

	//fmt.Println(MathematicaXYPlot(xs,ys))

}

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
	/*
	if !ms.unique{
		fmt.Println("Caution: Solution not unique!")
	}
	 */
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
	//where is the 3 coming from? it cannot be deg bc it's deg-3, right?! Use var oldthree and check.
	oldthree := deg-1
	x_curv := make([][]float64,len(x))
	for i:=0 ; i < len(x_curv) ; i++ {
		x_curv[i] = make([]float64,deg)
		for j:=0 ; j < len(x_curv[0])-2 ; j++ {
			//fmt.Println(factorial(3-j)/factorial(deg-3-j),"xi^",deg-3-j)
			x_curv[i][j] = float64(factorial(oldthree-j)/factorial(deg-oldthree-j))*math.Pow(x[i],float64(deg-oldthree-j))
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
		M.SetRow(cur,x_var[i],deg*i,y[i])
		cur++
	}

	//x_var right
	for i := 0 ; i < len(x_var)-1 ; i++ {
		//fmt.Println("fct val cond. right")
		M.SetRow(cur,x_var[i+1],deg*i,y[i+1])
		cur++
	}

	//=Sl
	if contains(constraints,"=Sl"){
		for i := 0 ; i < len(x_slope)-2 ; i++ {
			//fmt.Println("=Sl")
			row := floatlist_cat(x_slope[i+1],floatlist_negation_compwise(x_slope[i+1]))
			//fmt.Println(row)
			M.SetRow(cur,row,deg*i,0)
			cur++
		}
	}

	//0Sl
	if contains(constraints,"0Sl"){
		for i := 0 ; i < len(x_slope)-2 ; i++ {
			//fmt.Println("0Sl")
			M.SetRow(cur,x_slope[i],deg*i,0)
			cur++
		}
	}

	//=Cv
	if contains(constraints,"=Cv"){
		for i := 1 ; i < len(x_curv)-1 ; i++ {
			//fmt.Println("=Cv")
			row := floatlist_cat(x_curv[i],floatlist_negation_compwise(x_curv[i]))
			//fmt.Println(row)
			M.SetRow(cur,row,deg*(i-1),0)
			cur++
		}
	}

	//0Cv
	if contains(constraints,"0Cv"){
		for i := 1 ; i < len(x_curv)-1 ; i++ {
			//fmt.Println("0Cv")
			//row := floatlist_cat(x_curv[i],floatlist_negation_compwise(x_curv[i]))
			//M.AddRow(cur,row,4*(i-1),0)
			M.SetRow(cur,x_curv[i],deg*(i-1),0)
			cur++
		}
	}

	//E0Sl
	if contains(constraints,"E0Sl") {
		//fmt.Println("E0Sl")
		//first
		M.SetRow(cur,x_slope[0],deg*0,0)
		cur++
		//last
		M.SetRow(cur,x_slope[len(x_slope)-1],deg*(len(x_slope)-2),0)
		cur++
	}

	//E0Cv
	if contains(constraints,"E0Cv"){
		//fmt.Println("E0Cv")
		//first
		M.SetRow(cur,x_curv[0],deg*(len(x_curv)-2),0)
		cur++
		//last
		M.SetRow(cur,x_curv[len(x_curv)-1],deg*(len(x_curv)-2),0)
		cur++
	}

	//EQSl
	if contains(constraints,"EQSl") {
		//fmt.Println("EQSl")
		//first
		M.SetRow(cur,x_slope[0],deg*0,(y[1]-y[0])/(x[1]-x[0]))
		cur++
		//last
		M.SetRow(cur,x_slope[len(x_slope)-1],deg*(len(x_slope)-2),(y[len(y)-1]-y[len(y)-2])/(x[len(y)-1]-x[len(y)-2]))
		cur++
	}


	return M,nil
}

func (ms my_spline) PrintMathematicaCode(points bool, color string, plotRange string) (string,string){

	//either input or make dependent on ranges (ms.x,ms.y) AND minimal distance between entries
	//careful with too low Accu!! it can cause the plots to be wrong and finding this as a cause isn't easy
	/*
	xAccu := 4
	yAccu := 10
	xAccuStr := "%."+string(xAccu)+"f"
	yAccuStr := "%."+string(yAccu)+"f"
	 */

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
	if points {
		result += fmt.Sprintf("xy%v:=ListPlot[Transpose[{x%v, y%v}], PlotStyle -> {AbsolutePointSize[8]},ImageSize -> Large, PlotRange -> %s, AxesOrigin -> {0, 0}];\n",id,id,id,plotRange)
	}

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
		result += fmt.Sprintf("%.4f",ms.x[i/(ms.deg+1)])
		result += fmt.Sprint("<=x<=")
		result += fmt.Sprintf("%.4f",ms.x[i/(ms.deg+1)+1])
		result += fmt.Sprint("}")
		if i<(ms.deg+1)*(len(ms.x)-1)-(ms.deg+1) {
			result += fmt.Sprint(",")
		}
	}
	result += fmt.Sprint("}];\n")



	result += fmt.Sprintf("fctplot%v := Plot[fct%v[x]",id,id)
	result += ",{x,"
	result += fmt.Sprintf("%.4f",ms.x[0])
	result += fmt.Sprint(",")
	result += fmt.Sprintf("%.4f",ms.x[len(ms.x)-1])
	result += fmt.Sprint("},ImageSize->Large, PlotRange -> "+plotRange+", PlotStyle -> "+color+", AxesOrigin -> {0, 0}];\n")

	//Show
	if points {
		result += fmt.Sprintf("s%v:=Show[fctplot%v, xy%v];\n\n",id,id,id)
	} else {
		result += fmt.Sprintf("s%v:=Show[fctplot%v];\n\n",id,id)
	}

	return result,id
}

func (ms my_spline) MathematicaExport(color string, msg string, points bool, folderName string, fileName string, mathematicaCompressionLevel string, mathematicaImageResolution string, plotRange string) string {
	if len(ms.x) == 0 {return ""}
	content := ""
	tmp,id := ms.PrintMathematicaCode(points,color,plotRange)
	content += tmp+"\n"
	if msg == "" {
		content += "Export[\"" + folderName + "\\"+fileName+".png\"," + fmt.Sprintf("Show[s%v]",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
	} else {
		content += fmt.Sprintf("msg1 := Text[\""+msg+"\"];\n")
		content += "Export[\"" + folderName + "\\"+fileName+".png\"," + fmt.Sprintf("{msg1,Show[s%v]}",id) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
	}
	return content
}

func (ms my_spline) MathematicaExport2(color1 string, ms2 my_spline, color2 string, msg string, points bool, folderName string, fileName string, mathematicaCompressionLevel string, mathematicaImageResolution string, plotRange string) string {
	if len(ms.x) == 0 {return ""}
	content := ""
	tmp,id := ms.PrintMathematicaCode(points,color1,plotRange)
	tmp2,id2 := ms2.PrintMathematicaCode(points,color2,plotRange)
	content += tmp+"\n"
	content += tmp2+"\n"
	if msg == "" {
		content += "Export[\"" + folderName + "\\"+fileName+".png\"," + fmt.Sprintf("Show[s%v,s%v]",id,id2) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
	} else {
		content += fmt.Sprintf("msg1 := Text[\""+msg+"\"];\n")
		content += "Export[\"" + folderName + "\\"+fileName+".png\"," + fmt.Sprintf("{msg1,Show[s%v,s%v]}",id,id2) + ", \"CompressionLevel\" -> "+mathematicaCompressionLevel+", \n ImageResolution -> "+mathematicaImageResolution+"];\n"
	}
	return content
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

	//which spline(Nr) is relevent for x?
	for i := 0 ; i < len(ms.x) ; i++ {
		if i+1 < len(ms.x){
			if x >= ms.x[i] && x <= ms.x[i+1]{
				splineNr = i
				break
			}
		} else {
			splineNr = i
		}
	}

	//splineNr := len(ms.x)-1

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

/*
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
 */

//returns int_{a}^{b} ms.At(x) dx
func (ms my_spline) IntegralSpline(a,b float64) float64 {
	integral := 0.0
	i1 := 0
	i2 := len(ms.x)-1
	for i1 < len(ms.x) && ms.x[i1] < a {i1++};
	for i2 > 0 &&  ms.x[i2] > b {i2--};
	if i1 < 1 {i1=1}
	if i1 >= len(ms.x) {i1=len(ms.x)-1}
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
//unionizes two splines to have the same x intervals, appropriately adjust y arrays as well as the coefficients in both splines and returns the two, now compatible splines, compatible for e.g. addition, multiplication, etc.
func UnionXYCC (ms1, ms2 my_spline) (my_spline , my_spline) {

	if len(ms1.x) == 0 || len(ms2.x) == 0 {
		return ms1,ms2
	}

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

	if len(ms1.x) == 0 || len(ms2.x) == 0 {
		return my_spline{}
	}

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
	if len(ms.x) == 0 {return my_spline{}}
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


//For degrees <=3, should be replaced by pq or cubic root formula
//finds roots (y=0) of ms, starting at xo with a tolerance of 0<tol. For other y's it doesn't find roots but where ms is y.
//implement Derive() and calculate it once instead of using D() multiple times
func (ms my_spline) NewtonRoot(x0 float64, y float64, tolYPerc float64) (float64,error) {
	debug := false
	derivative := ms.Derive()
	if debug{
		fmt.Println("calculated derivative.")
		fmt.Println("at xn=",x0," the derivitive is ",derivative.At(x0))
	}
	xn := x0
	skip := 20
	span := max(ms.y)-min(ms.y)
	tolY := tolYPerc * span
	for math.Abs(y-ms.At(xn)) > tolY {
		skip--
		if skip <= 0 || xn > max(ms.x) || xn < min (ms.x) {return 0,fmt.Errorf("newton couldn't find root")}
		//fmt.Println("old xn: ",xn," , old yn: ", ms.At(xn), " , D(xn)=",ms.D(xn))

		xn = math.Min(max(ms.x),math.Max(min(ms.x),  xn+(y-ms.At(xn))/derivative.At(xn)  ))

		if debug{
			fmt.Println(xn,":",ms.At(xn) , " , difference to ",y," is ",ms.At(xn)-y)
		}
	}
	return xn,nil
}

// divide range into n pieces, start NewtonRoot on each and collect all roots, gathering with tolerance tol
func (ms my_spline) NewtonRoots (y float64, tolYPerc float64, n int) []float64 {
	debug := false
	if len(ms.x) == 0 {return []float64{}}
	tolXPerc := 0.01
	span := ms.x[len(ms.x)-1]-ms.x[0]
	dx := span / float64(n)
	var intersections []float64
	for x := ms.x[0] ; x < ms.x[len(ms.x)-1] ; x += dx {
		root,err := ms.NewtonRoot(x,y,tolYPerc)
		if err != nil {
			if debug{
				fmt.Println("error: ",err)
				fmt.Println("continue with next starting point")
			}
			continue
		}
		if !containsFloat(intersections,root,tolXPerc*span){
			if debug {
				fmt.Println("NewtonRoots: root:",root," since ms.At(root) = ",ms.At(root)," ~= y = ",y, "." +
					"\nThe abs. val. of the difference is ",math.Abs(ms.At(root)-y)," < ",tolYPerc*(max(ms.y)-min(ms.y)))
			}
			intersections = append(intersections,root)
		}
	}
	return intersections
}

func DoubleFloatUnionTol(ar [][]float64, tol float64) [][]float64 {
	debug := false
	var result [][]float64
	//result = make([][]float64,0)
	for i := 0 ; i < len(ar)-1 ; i++{
		for j := range result {
			if ar[i][0]-result[j][0] < tol && ar[i][1]-result[j][1] < tol {
				if debug{
					fmt.Println("ar[",i,"][0]=",ar[i][0], " , result [",j,"][0]=",result[j][0] , " (dif: ",math.Abs(ar[i][0]-result[j][0]),"<tol=",tol , " AND ar[",i,"][1]=",ar[i][1], " ,  result [",j,"][1]=",result[j][1]," , (dif: ",math.Abs(ar[i][0]-result[j][0]),"<tol=",tol,")")
				}
				i++
				break
			}
		}
		if debug {
			fmt.Println("adding in DoubleFloatUnionTol")
		}
		result = append(result, ar[i])
	}
	return result
}


// returns regions where ms is less than y and where ms is greater than y
func (ms my_spline) PosNegRange (y float64, tolYPerc float64, n int) ([][]float64 , [][]float64) {
	debug := false
	//tolXPerc := 0.01
	//span := max(ms.x)-min(ms.x)
	//tolX := tolXPerc*span
	roots := ms.NewtonRoots(y,tolYPerc,n)

	//replace start and end separation
	roots = append(roots,min(ms.x))
	roots = append(roots,max(ms.x))
	sort.Float64s(roots)

	if len(roots) == 0 {return [][]float64{},[][]float64{}}
	var neg [][]float64
	var pos [][]float64

	//start
	/*
	if ms.At((min(ms.x)+roots[0])/2) < y  {
		if debug{ fmt.Println("ms.At((min(ms.x)+roots[0])/2)=",ms.At((min(ms.x)+roots[0])/2) ,"<y=",y,"." +
			"Therefore [min(ms.x),roots[0]]=[",min(ms.x),",",roots[0],"] is added to neg.")}
		neg = append(neg, []float64{min(ms.x),roots[0]})
	} else {
		if debug{ fmt.Println("ms.At((min(ms.x)+roots[0])/2)=",ms.At((min(ms.x)+roots[0])/2) ,">y=",y,"." +
			"Therefore [min(ms.x),roots[0]]=[",min(ms.x),",",roots[0],"] is added to pos.")}
		pos = append(pos, []float64{min(ms.x),roots[0]})
	}
	 */

	//all middle parts
	for i := range roots[0:len(roots)-1] {
		if ms.At((roots[i+1]+roots[i])/2 ) < y  {
			if debug{ fmt.Println("For i=",i," ms.At((roots[i+1]+roots[i])/2)=",ms.At((roots[i+1]+roots[i])/2) ,"<",y,"." +
				"Therefore [roots[i+1],roots[i]]=[",roots[i],",",roots[i+1],"] is added to neg.")}
			neg = append(neg, []float64{roots[i],roots[i+1]})
		} else {
			if debug{ fmt.Println("For i=",i," ms.At((roots[i+1]+roots[i])/2)=",ms.At((roots[i+1]+roots[i])/2) ,">",y,"." +
				"Therefore [roots[i+1],roots[i]]=[",roots[i],",",roots[i+1],"] is added to pos.")}
			pos = append(pos, []float64{roots[i],roots[i+1]})
		}
	}

	//end
	/*
	if ms.At((max(ms.x)+roots[len(roots)-1])/2 ) < y  {
		if debug{ fmt.Println("ms.At((max(ms.x)+roots[len(roots)-1])/2)=",ms.At((max(ms.x)+roots[len(roots)-1])/2) ,"<",y,"." +
			"Therefore [roots[len(roots)-1],max(ms.x)]=[",roots[len(roots)-1],",",max(ms.x),"] is added to neg.")}
		neg = append(neg, []float64{roots[len(roots)-1],max(ms.x)})
	} else {
		if debug{ fmt.Println("ms.At((max(ms.x)+roots[len(roots)-1])/2)=",ms.At((max(ms.x)+roots[len(roots)-1])/2) ,">",y,"." +
			"Therefore [roots[len(roots)-1],max(ms.x)]=[",roots[len(roots)-1],",",max(ms.x),"] is added to pos.")}
		pos = append(pos, []float64{roots[len(roots)-1],max(ms.x)})
	}
	 */

	//var neg_tmp,pos_tmp [][]float64

	//remove duplicates
	/*
	for i := 0 ; i < len(neg)-1 ; i++{
		for j := 0 ; j < i ; j++ {
			if neg[i][0]-neg[j][0]<tolX && neg[i][1]-neg[j][1]<tolX {
				i++
				break
			} else {
				neg_tmp = append(neg_tmp, neg[i])
			}
		}
	}
	neg = neg_tmp
	 */

	/*
	fmt.Println("neg before union: ",neg)
	fmt.Println("pos before union: ",pos)
	neg = DoubleFloatUnionTol(neg,tolX)
	pos = DoubleFloatUnionTol(pos,tolX)
	fmt.Println("neg after union: ",neg)
	fmt.Println("pos after union: ",pos)
	 */

	//remove small interval
	/*
	tolX := 0.01*(max(ms.x)+min(ms.x))/2
	var neg_tmp,pos_tmp [][]float64
	for i := range neg {
		if math.Abs(neg[i][0]-neg[i][1]) < tolX{
			continue
		}
		neg_tmp = append(neg_tmp,neg[i])
	}
	for i := range pos {
		if math.Abs(pos[i][0]-pos[i][1]) < tolX{
			continue
		}
		pos_tmp = append(neg_tmp,pos[i])
	}
	neg = neg_tmp
	pos = pos_tmp
	 */

	//merge intervals with tolerance
	//make into separate function
	/*
		if len(neg) == 0 || len(pos) == 0 {return neg,pos}
		tolPerc := 0.01
		span := max(ms.y)-min(ms.y)
		tol := tolPerc * span

		neg_tmp = make([][]float64,0)
		var left, right float64
		for i := 0 ; i < len(neg) ; i++ {
			left = neg[i][0]
			right = neg[i][1]
			for i+1 < len(neg) && math.Abs(neg[i][1]-neg[i+1][0]) < tol {
				right = neg[i+1][1]
				i++
			}
			neg_tmp = append(neg_tmp,[]float64{left,right})
		}
		neg = neg_tmp

		pos_tmp = make([][]float64,0)
		for i := 0 ; i < len(pos) ; i++ {
			left = pos[i][0]
			right = pos[i][1]
			for i+1 < len(pos) && math.Abs(pos[i][1]-pos[i+1][0]) < tol {
				right = pos[i+1][1]
				i++
			}
			pos_tmp = append(pos_tmp,[]float64{left,right})
		}
		pos = pos_tmp
	*/
	neg = MergeLeftRight(neg,0.01)
	pos = MergeLeftRight(pos,0.01)



	return neg,pos
}

//e.g. MergeLeftRight([[100,200],[200.1,300]],0.1/(300-100))=[[100,300]]
func MergeLeftRight (a [][]float64, tolPerc float64) [][]float64 {
	if len(a) == 0 {return a}
	span := a[0][0]-a[len(a)-1][1]
	tol := tolPerc * span

	var tmp [][]float64
	var left, right float64
	for i := 0 ; i < len(a) ; i++ {
		left = a[i][0]
		right = a[i][1]
		for i+1 < len(a) && math.Abs(a[i][1]-a[i+1][0]) <= tol {
			right = a[i+1][1]
			i++
		}
		tmp = append(tmp,[]float64{left,right})
	}
	return tmp
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

func getSplineCoeffsDegree (ms my_spline, splineNr int, deg int) float64 {
	return ms.coeffs[(splineNr+1)*(ms.deg+1)-1-deg]
}

// ! only up to degree 4 currently but extendable (see pdf)
func (ms my_spline) Inversion() my_spline {
	var coeffs []float64
	for i:=0; i < len(ms.x)-1 ; i++ {
		for d := ms.deg ; d >= 0 ; d-- {
			switch d {
			//as in coeffs, the degree per polynomial is decreasing for increasing index, go one polynomial further and subtract d
			case 0: coeffs = append(coeffs, //A_0=a_0
				getSplineCoeffsDegree(ms,i,0))
			case 1: coeffs = append(coeffs,//A_1=a_1^-1
				math.Pow(getSplineCoeffsDegree(ms,i,1),-1)	)
			case 2: coeffs = append(coeffs,//A_2=-a_1^-3a_2
				-math.Pow(getSplineCoeffsDegree(ms,i,1),-3)*getSplineCoeffsDegree(ms,i,2))
			case 3: coeffs = append(coeffs,//A_3=a_1^-5(2a_2^2-a_1a_3)
				math.Pow(getSplineCoeffsDegree(ms,i,1),-5)*(2*math.Pow(getSplineCoeffsDegree(ms,i,2),2)-getSplineCoeffsDegree(ms,i,1)*getSplineCoeffsDegree(ms,i,3))	)
			case 4: coeffs = append(coeffs,//A_4=a_1^-7(5a_1a_2a_3-a_1^2a_4-5a_2^3)
				math.Pow(getSplineCoeffsDegree(ms,i,1),-7)*(
					5*getSplineCoeffsDegree(ms,i,1)*getSplineCoeffsDegree(ms,i,2)*getSplineCoeffsDegree(ms,i,3)-
					math.Pow(getSplineCoeffsDegree(ms,i,1),2)*getSplineCoeffsDegree(ms,i,4)-5*math.Pow(getSplineCoeffsDegree(ms,i,2),3))	)
			}
		}
	}
	inversion := my_spline{
		deg:        ms.deg,
		splineType: ms.splineType,
		x:          ms.y,
		y:          ms.x,
		coeffs:     coeffs,
		unique:     false,
	}
	fmt.Println("inversion: deg=",inversion.deg,", len(x)=",len(inversion.x),", len(y)=", len(inversion.y),", len(coeffs)=",len(coeffs))
	return inversion
}






// ------------------------------- call specific functions -------------------------------

func (call callfunc) ToSpread() spread {
	return spread{
		num:     1,
		calls:   []callfunc{call},
		weights: []float64{1},
	}
}

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
	code += "xy:=ListPlot[Transpose[{x, y}], PlotStyle -> {AbsolutePointSize[8]},ImageSize -> Large, PlotRange -> Automatic, Joined -> True];\n"
	code += "Show[xy];"
	return code
}

func (call callfunc) At (x float64) float64{
	return call_gain_perc(x, call)
}

func call_breakeven_ground(call callfunc) float64{
	return call.base+call.cost*1/call.factor
}

func call_breakeven_base(call callfunc, curbase float64) float64{
	return call.base*call.factor*curbase/(call.factor*curbase-call.cost)
}

func call_gain_perc(x float64, call callfunc) float64{
	//return math.Max(-1,x/(call.cost/call.factor)-call.base/(call.cost/call.factor)-1)*100
	//return math.Max(-1,((1.0/call.factor*(x-call.base))/call.cost)-1)*100
	//100*Max[-1,math.Abs(1.0/call.factor)*((x-call.base)/(call.cost*call.factor)))-1]
	return math.Max(-1,math.Abs(1.0/call.factor)*((x-call.base)/(call.cost*call.factor))-1)*100
}

func (call callfunc) PrintMathematicaCode(lr float64) string{
	//fmt.Println("Mathematica Code to visualize call option value\n\n")
	code := ""
	//code += fmt.Sprintln("call:=Plot[100*Max[-1,(x/(",call.cost/call.factor,")-(",call.base/(call.cost/call.factor),")-1)],{x,0,",lr,"},ImageSize->Large, PlotRange->Automatic];")
	code += fmt.Sprintf("call:=Plot[100*Max[-1,%.4f*((x-%.3f)/%.10f)-1],{x,0,%.3f},ImageSize->Large, PlotRange->Automatic];\n",math.Abs(1.0/call.factor),call.base,call.cost*call.factor,lr)
	code += fmt.Sprintln("Show[call]")
	return code
}

//doesn't do it properly for puts
func PrintMathematicaCode(calls []callfunc, share_price float64, callsColor string, longColor string, includeShort bool) string {
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
		//code += fmt.Sprint("call"+strconv.Itoa(i)+":=Plot[100*Max[-1,(x/(",call.cost/call.factor,")-(",call.base/(call.cost/call.factor),")-1)],{x,0,xmax},ImageSize->Large, PlotRange->Automatic, PlotStyle -> "+callsColor+", AxesOrigin -> {0, 0}];\n")
		code += fmt.Sprintf("call%v:=Plot[100*Max[-1,%.4f*((x-%.3f)/%.10f)-1],{x,0,xmax},ImageSize->Large, PlotRange->Automatic, PlotStyle -> %s, AxesOrigin -> {0, 0}];\n",strconv.Itoa(i),math.Abs(1.0/call.factor),call.base,call.cost/call.factor,callsColor )
	}

	code += "long := Plot[100*(x - "+fmt.Sprintf("%.2f",share_price)+")/"+fmt.Sprintf("%.2f",share_price)+", {x, 0, xmax}, PlotStyle -> "+longColor+", AxesOrigin -> {0, 0}];\n"
	if includeShort {
		code += "short := Plot[100*(-x + "+fmt.Sprintf("%.2f",share_price)+")/"+fmt.Sprintf("%.2f",share_price)+", {x, 0, xmax}, PlotStyle -> "+longColor+", AxesOrigin -> {0, 0}];\n"
	}

	for i := range calls {
		if i==0{
			code += "s:=Show[{call"+strconv.Itoa(i)
			continue
		}
		code += fmt.Sprintln(",call"+strconv.Itoa(i))
	}
	code += ",long"
	if includeShort {
		code += ",short"
	}
	code += "}]\n"
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

//full_integral call*pdist
func (call callfunc) ExpectedReturn (pdist my_spline) float64{
	debug := false
	var ref float64
	var dx float64
	//result := pdist.SplineMultiply(call.ToSpline(min(pdist.x),max(pdist.x))).Integral(min(pdist.x),max(pdist.x),dx)
	callSpline := call.ToSpline(min(pdist.x),max(pdist.x))
	if len(callSpline.x) == 0 || len(callSpline.y) == 0 {
		fmt.Println("unexpected error in callfunc.ExpectedReturn()")
		return -10000.0
	}
	result := pdist.SplineMultiply(callSpline).FullIntegralSpline()
	if debug{
		fmt.Print("debug in call.ExpectedReturn ")
		dx = 0.001
		ref = call.ExpectedReturnDX(pdist,dx)
		diff := math.Abs(ref - result)
		fmt.Println("math.Abs(ref - result)=",diff)
		if diff > 0.01*math.Abs(ref) {
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
	if call.optionType == "put" {
		return my_spline{
			deg:        1,
			splineType: []string{"3","2","=Sl","=Cv","EQSl"},//not really
			x:          []float64{a,call.base,b},
			y:          []float64{call.At(a),-100,-100},
			//coeffs:     []float64{call.factor/call.cost*100,-100-100*call.base*call.factor/call.cost,0,-100},
			//100*Max[-1,  %.4f*((x-%.3f)/%.10f)-1   ,math.Abs(1.0/call.factor)  ,  call.base  ,  call.cost*call.factor
			//				math.Abs(1.0/call.factor)*((x-call.base)/(call.cost*call.factor)))-1
			//				 = math.Abs(1.0/call.factor)/(call.cost*call.factor) *x - math.Abs(1.0/call.factor)*call.base/(call.cost*call.factor)
			//coeffs:     []float64{call.factor*call.cost*100,-100-100*call.base*call.factor*call.cost   ,0,-100},
			//coeffs:     []float64{math.Abs(1.0/call.factor)/(call.cost*call.factor)*100 , -100-100*call.base*math.Abs(1.0/call.factor)/(call.cost*call.factor)   ,0,-100},
			coeffs:     []float64{call.factor/math.Abs(call.factor)*1.0/(call.cost)*100 , -100-100*call.factor/math.Abs(call.factor)*call.base*1.0/(call.cost)   ,0,-100},
			unique:     false,
		}
	}
	return my_spline{
			deg:        1,
			splineType: []string{"3","2","=Sl","=Cv","EQSl"},//not really
			x:          []float64{a,call.base,b},
			y:          []float64{-100,-100,call.At(b)},
			//coeffs:     []float64{0,-100,call.factor/call.cost*100,-100-100*call.base*call.factor/call.cost},
			//coeffs:     []float64{0,-100,call.factor*call.cost*100,-100-100*call.base*call.factor*call.cost},
			//coeffs:     []float64{0,-100  ,  math.Abs(1.0/call.factor)/(call.cost*call.factor)*100,-100-100*call.base*math.Abs(1.0/call.factor)/(call.cost*call.factor)},
			coeffs:     []float64{0,-100  ,  call.factor/math.Abs(call.factor)*1.0/(call.cost)*100,-100-100*call.base*call.factor/math.Abs(call.factor)*1.0/(call.cost)},
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
		if optt.Contract_type == "call" {
			callList = append(callList, callfunc{
				base:   float64(optt.Strike_price),
				cost:   optt.Close,
				factor: float64(optt.Shares_per_contract),//1.0/float64(optt.Shares_per_contract),//1.0,
				date:   dateInt,
			})
		} else if optt.Contract_type == "put"{
			callList = append(callList, callfunc{
				base:   float64(optt.Strike_price),
				cost:   optt.Close,
				factor: -float64(optt.Shares_per_contract),//-1.0,//-1.0/float64(optt.Shares_per_contract),
				date:   dateInt,
			})
		}

		volumes = append(volumes, optt.Volume)
	}
	var interListVol []float64
	for i,call := range callList {
		for v := 0 ; v < volumes[i] ; v++ {
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

func FindCloestPDist(d string, pdistDates []string, pdistSplines map[string]my_spline, debug bool) my_spline{
	ymd := strings.Split(d,"-")
	y,err := strconv.ParseInt(ymd[0],10,64);check(err)
	m,err := strconv.ParseInt(ymd[1],10,64);check(err)

	var pdist my_spline
	var pdDy int64 = -1
	var yBestDist int = 1000 //just sufficiently large to start
	if m > 6 {
		for i,pdD := range pdistDates {
			pdDy,err = strconv.ParseInt(strings.Split(pdD,"-")[0],10,64)
			check(err)
			if int(math.Abs(float64(pdDy - (y+1)))) < yBestDist {
				yBestDist = int(math.Abs(float64(pdDy - (y+1))))
				pdist = pdistSplines[pdistDates[i]]
			}
		}
	} else {
		for i,pdD := range pdistDates {
			pdDy,err = strconv.ParseInt(strings.Split(pdD,"-")[0],10,64);check(err)
			if int(math.Abs(float64(pdDy - (y)))) < yBestDist {
				yBestDist = int(math.Abs(float64(pdDy - (y))))
				pdist = pdistSplines[pdistDates[i]]
			}
		}
	}
	if len(pdist.x) == 0 {
		fmt.Println("no year match in pdistDates.")
		os.Exit(69)
	}

	if debug {
		fmt.Println("pdDy=",pdDy)
		fmt.Println("y=",y)
		fmt.Println("yBestDist=",yBestDist)
		//os.Exit(120)
	}

	return pdist
}

func CAGR(returnPerc float64,years float64) float64 {
	return (math.Pow(1.0+returnPerc/100.0,1.0/years)-1.0)*100.0
}


// returns if date is in dateRange
func DateInRange(dateRange []string, dateStr string) bool {
	dateLowStr := dateRange[0]
	dateHighStr := dateRange[1]
	dateLow, err := time.Parse("2006-01-02",dateLowStr)
	check(err)
	dateHigh, err := time.Parse("2006-01-02",dateHighStr)
	check(err)
	date, err := time.Parse("2006-01-02",dateStr)
	check(err)
	return date.Before(dateHigh) && date.After(dateLow)
}


//sorting and keeping index changes
type sortable struct {
	nums []float64
	idxs []int
}
func (s sortable) Len() int           { return len(s.nums) }
func (s sortable) Less(i, j int) bool { return s.nums[i] < s.nums[j] }
func (s sortable) Swap(i, j int) {
	s.nums[i], s.nums[j] = s.nums[j], s.nums[i]
	s.idxs[i], s.idxs[j] = s.idxs[j], s.idxs[i]
}
func sortAndReturnIdxsFloat64(nums []float64) []int {
	idxs := make([]int, len(nums))
	for i := range idxs {
		idxs[i] = i
	}
	sort.Sort(sortable{nums, idxs})
	return idxs
}


//saves pulled options as well as the date and the contract type in a json at location os.Getwd()/pathExt/filename.json
func SaveOptionsJson(pathExt string, filename string, pulledDate time.Time, options []opt.Option, Contract_type []string){
	//.UTC().Format("2006-01-02T15:04:05Z07:00")

	data := map[string]interface{}{
		"pulledDate":    pulledDate.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"options":   options,
		"Contract_type":	Contract_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	jsonDataStr := fmt.Sprintf("%s\n",jsonData)
	jsonDataStr = strings.Replace(jsonDataStr,",",",\n",-1)

	path, err := os.Getwd()
	check(err)
	os.Mkdir(path+pathExt,0755)

	WriteFile(filename+".json",jsonDataStr,pathExt)

}

func LoadOptionsJson(path string, filename string) (time.Time,[]string,[]opt.Option) {
	jsonFile, err := os.Open(path+filename)
	check(err)

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	//ticker := result["ticker"].(string)

	pulledDate_tmp := result["pulledDate"].(string)
	//pulledDate_tmp = strings.Split(pulledDate_tmp,".")[0]
	//fmt.Println("pulledDate_tmp=",pulledDate_tmp)
	//layout := "2023-09-14T15:28:06.702131+02:00"
	//layout := "2023-09-15T13:41:45Z"
	pulledDate, err := time.Parse("2006-01-02T15:04:05Z07:00", pulledDate_tmp)
	check(err)
	//fmt.Println("pulledDate=",pulledDate)

	/*
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	t, err := time.Parse(layout, str)
	 */


	var Contract_type []string

	for _,ct := range result["Contract_type"].([]interface{}) {
		Contract_type = append(Contract_type,ct.(string))
	}

	var options []opt.Option

	options_tmp := result["options"]
	for _,o := range options_tmp.([]interface{}) {
		oCasted := o.(map[string]interface{})
		opt_tmp := opt.Option{}

		//fmt.Println(o)

		/*
			Cfi string
			Contract_type string
			Exerciese_style string
			Expiration_date string
			Primaty_exchange string
			Shares_per_contract int
			Strike_price float64
			Ticker string
			Underlying_ticker string
			Volume int
			Vw float64
			Open float64
			Close float64
			High float64
			Low float64
			T int
			N int
		 */

		opt_tmp.Cfi 				= oCasted["Cfi"].(string)
		opt_tmp.Contract_type 		= oCasted["Contract_type"].(string)
		opt_tmp.Exerciese_style 	= oCasted["Exerciese_style"].(string)
		opt_tmp.Expiration_date		= oCasted["Expiration_date"].(string)
		opt_tmp.Primaty_exchange 	= oCasted["Primaty_exchange"].(string)
		opt_tmp.Shares_per_contract = int(oCasted["Shares_per_contract"].(float64))
		opt_tmp.Strike_price 		= oCasted["Strike_price"].(float64)
		opt_tmp.Ticker 				= oCasted["Ticker"].(string)
		opt_tmp.Underlying_ticker 	= oCasted["Underlying_ticker"].(string)
		opt_tmp.Volume 				= int(oCasted["Volume"].(float64))
		opt_tmp.Vw 					= oCasted["Vw"].(float64)
		opt_tmp.Open 				= oCasted["Open"].(float64)
		opt_tmp.Close 				= oCasted["Close"].(float64)
		opt_tmp.High 				= oCasted["High"].(float64)
		opt_tmp.Low 				= oCasted["Low"].(float64)
		opt_tmp.T 					= int(oCasted["T"].(float64))
		opt_tmp.N 					= int(oCasted["N"].(float64))

		//fmt.Println(opt_tmp)
		options = append(options,opt_tmp)

	}

	return pulledDate,Contract_type,options
}

func SavePromptJson(ticker string ,pdistDates []string, pdistX map[string][]float64, pdistY map[string][]float64, StrikeRange []int ,Contract_type string) {
	data := map[string]interface{}{
		"ticker":    ticker,
		"pdistDates":   pdistDates,
		"pdistX": pdistX,
		"pdistY": pdistY,
		"StrikeRange": StrikeRange,
		"Contract_type": Contract_type,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	jsonDataStr := fmt.Sprintf("%s\n",jsonData)
	jsonDataStr = strings.Replace(jsonDataStr,",",",\n",-1)

	currentTime := time.Now()
	live := currentTime.Format("2006-01-02")

	filename := "prompt2_"+ticker+" ("+live+").json"
	path, err := os.Getwd()
	check(err)
	os.Mkdir(path+"\\prompts",0755)
	pathExt := "\\prompts\\"

	WriteFile(filename,jsonDataStr,pathExt)

}

func SavePromptJsonEasy(ticker string ,pdistDates []string, pdistX map[string][]float64, pdistY map[string][]float64, StrikeRange []int ,Contract_type string){

	pdistDatesXY := make(map[string][][]float64,len(pdistDates))
	for _,d := range pdistDates {
		pdistDatesXY[d] = make([][]float64,len(pdistX[d]))
		for i,_ := range pdistX[d] {
			pdistDatesXY[d][i] = make([]float64,2)
			pdistDatesXY[d][i][0] = pdistX[d][i]
			pdistDatesXY[d][i][1] = pdistY[d][i]
		}
	}


	data := map[string]interface{}{
		"ticker":    ticker,
		"pdistDatesXY":   pdistDatesXY,
		"StrikeRange": StrikeRange,
		"Contract_type": Contract_type,
	}


	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	jsonDataStr := fmt.Sprintf("%s\n",jsonData)
	jsonDataStr = strings.Replace(jsonDataStr,"],[","],\n[",-1)
	jsonDataStr = strings.Replace(jsonDataStr,",[",",[\n",-1)
	jsonDataStr = strings.Replace(jsonDataStr,"[[","[\n[",-1)
	jsonDataStr = strings.Replace(jsonDataStr,"]]","]\n]",-1)
	jsonDataStr = strings.Replace(jsonDataStr,",\"",",\n\"",-1)
	jsonDataStr = strings.Replace(jsonDataStr,"{\"","{\n\"",-1)
	jsonDataStr = strings.Replace(jsonDataStr,"}","\n}",-1)



	currentTime := time.Now()
	live := currentTime.Format("2006-01-02")
	filename := "promptEasy_"+ticker+" ("+live+").json"
	path, err := os.Getwd()
	check(err)
	os.Mkdir(path+"\\prompts",0755)
	pathExt := "\\prompts\\"

	WriteFile(filename,jsonDataStr,pathExt)

}

func SavePromptJsonOld(ticker string ,pdistDates []string, pdistX map[string][]float64, pdistY map[string][]float64, StrikeRange []int ,Contract_type string){

	content := "ticker="+ticker+";\n"
	/*
	content += "[ticker={"+ticker+"};" + "pdistDates={"
	for i,d := range pdistDates {
		content += d
		if i<len(pdistDates)-1{ content += ","}
	}
	content += ";"
	 */

	content += "pdistDatesXY={\n"
	for _,d := range pdistDates {
		content += "[\n" + d + ",\n{"
		X := pdistX[d]
		Y := pdistY[d]
		for i,_ := range X {
			content += fmt.Sprintf("(%.1f,%.1f)\n",X[i],Y[i])
			if i<len(X)-1{content+=","}
		}
		content += "},\n"
		content += "];\n"
	}
	content += "};\n"

	content += "StrikeRange=["+fmt.Sprint(StrikeRange)+"];\n"
	content += "ContractType=["+fmt.Sprintf(Contract_type)+"];\n"

	/*
	content += ";"

	content += "pdistY={"
	for i,d := range pdistDates {
		content += "[" + d + ",{"
		Y := pdistY[d]
		for i,y := range Y {
			content += fmt.Sprintf("%.2f",y)
			if i<len(Y)-1{content+=","}
		}
		content += "}"
		if i<len(pdistDates)-1{content+=","}
	}
	 */
	 //tobecontinued

	currentTime := time.Now()
	live := currentTime.Format("2006-01-02")
	filename := "prompt_"+ticker+" ("+live+").json"
	path, err := os.Getwd()
	check(err)
	os.Mkdir(path+"\\prompts",0755)
	pathExt := "\\prompts\\"

	WriteFile(filename,content,pathExt)
}

func LoadPromptEasy(path string, filename string) (string, []string, map[string][]float64, map[string][]float64, []int, []string, []string, []float64, []float64){

	jsonFile, err := os.Open(path+filename)
	check(err)

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	ticker := result["ticker"].(string)

	var DateRange []string
	for _,d := range result["DateRange"].([]interface{}){
		DateRange = append(DateRange,d.(string))
	}


	/* 		-- cast map[string][]interface to map[string][][]float64 --		 */
	// Initialize an empty map with string keys and [][]float64 values
	var pdistDatesXY = make(map[string][][]float64)
	// Loop through data and cast each value to [][]float64
	for k, v := range result["pdistDatesXY"].(map[string]interface{}) {
		// Use type assertion to attempt to cast v to []interface{}
		if outer, ok := v.([]interface{}); ok {
			// Initialize a 2D float64 slice
			twoD := make([][]float64, len(outer))
			for i, inner := range outer {
				// Use type assertion to attempt to cast inner to []interface{}
				if innerCasted, ok := inner.([]interface{}); ok {
					// Initialize a float64 slice
					floats := make([]float64, len(innerCasted))
					for j, val := range innerCasted {
						// Use type assertion to cast val to float64
						if floatVal, ok := val.(float64); ok {
							floats[j] = floatVal
						} else {
							fmt.Printf("Inner value at index %d is not of type float64\n", j)
						}
					}
					twoD[i] = floats
				} else {
					fmt.Printf("Value at index %d is not of type []interface{}\n", i)
				}
			}
			pdistDatesXY[k] = twoD
		} else {
			fmt.Printf("Value for key '%s' is not of type []interface{}\n", k)
		}
	}


	// old debugging
	/*
	fmt.Printf("pdistDatesXY datatype: %T",result["pdistDatesXY"])
	pdistDatesXYMap := result["pdistDatesXY"].(map[string]interface{})
	var pdistDatesXY map[string][][]float64
	for d,x := range pdistDatesXYMap {
			for i, y := range x.([]interface{}){
				pdistDatesXY[d] = make([][]float64,len(x.([]interface{})))
				fmt.Printf("y type is %T\n",y)
				fmt.Println("d=",d,"i=",i)
				for j, z := range y.([]interface{}){
					pdistDatesXY[d][i] = make([]float64,len(y.([]interface{})))
					z = z.(float64)
					fmt.Printf("z type is %T\n",z)
					pdistDatesXY[d][i][j] = z.(float64)
				}
			}
	}
	//pdistDatesXY := pdistDatesXYMap.(map[string][][]float64)
	*/

	var StrikeRange = make([]int,2)
	castedSR := result["StrikeRange"].([]interface{})
	StrikeRange[0] = int(castedSR[0].(float64))
	StrikeRange[1] = int(castedSR[1].(float64))

	var Contract_type []string
	castedCT := result["Contract_type"].([]interface{})
	for _,c := range castedCT {
		Contract_type = append(Contract_type,c.(string))
	}

	//Convert pdistDatesXY into pdistDates, pdistX, pdistY
	var pdistDates = make ([]string,0)
	var pdistX = make(map[string][]float64)
	var pdistY = make(map[string][]float64)
	for d, xy := range pdistDatesXY {
		pdistDates = append(pdistDates,d)
		pdistX[d] = make([]float64,0/*len(pdistDatesXY[d])*/)
		pdistY[d] = make([]float64,0/*len(pdistDatesXY[d])*/)
		for _, x := range xy {
			pdistX[d] = append(pdistX[d],x[0])
			pdistY[d] = append(pdistY[d],x[1])
		}
	}

	riskTol := result["riskTol"].([]interface{})

	var riskTolX = make([]float64,len(riskTol))
	var riskTolY = make([]float64,len(riskTol))
	for i := range riskTol {
		riskTolX[i] = riskTol[i].([]interface{})[0].(float64)
		riskTolY[i] = riskTol[i].([]interface{})[1].(float64)
	}

	return ticker,pdistDates,pdistX,pdistY,StrikeRange,DateRange,Contract_type,riskTolX,riskTolY
}

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

//returns a bool whether item is contained in list with tolerance eps
func containsFloat(list []float64, item float64, eps float64) bool {
	//eps := 0.01
	for _,l := range list {
		if math.Abs(l-item)<eps {
			return true
		}
	}
	return false
}

//piecewise addition for two []float64
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
		for i:=0;i</*3*/len(dateStr);i++ {
			tmp,_ := strconv.Atoi(dateStr[i])
			dateInt = append(dateInt,tmp)
		}

		if len(optionsMap[optt.Expiration_date])>0 {
			optionsMap[optt.Expiration_date] = append(optionsMap[optt.Expiration_date],optt)
			if optt.Contract_type == "call"{
				callListMap[optt.Expiration_date] = append(callListMap[optt.Expiration_date],callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Close,
					factor: 1.0,//float64(optt.Shares_per_contract),//1.0/flhoat64(optt.Sares_per_contract),
					date:   dateInt,
					optionType: optt.Contract_type,
				})
			} else if optt.Contract_type == "put"{
				callListMap[optt.Expiration_date] = append(callListMap[optt.Expiration_date],callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Close,
					factor: -1.0,//-float64(optt.Shares_per_contract),//-1.0/float64(optt.Shares_per_contract),
					date:   dateInt,
					optionType: optt.Contract_type,
				})
			}
		} else {
			optionsDates = append(optionsDates,optt.Expiration_date)
			tmp := make([]opt.Option,1)
			tmp[0] = optt
			optionsMap[optt.Expiration_date] = tmp
			tmpp := make([]callfunc,len(addToAll))
			for i := range addToAll {
				tmpp[i] = addToAll[i]
			}
			if optt.Contract_type == "call" {
				tmpp = append(tmpp,callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Close,
					factor: float64(optt.Shares_per_contract),//1.0,//1.0/float64(optt.Shares_per_contract),
					date:   dateInt,
					optionType: optt.Contract_type,
				})
			} else if optt.Contract_type == "put" {
				tmpp = append(tmpp,callfunc{
					base:   float64(optt.Strike_price),
					cost:   optt.Close,
					factor: -float64(optt.Shares_per_contract),//-1.0,//-1.0/float64(optt.Shares_per_contract),
					date:   dateInt,
					optionType: optt.Contract_type,
				})
			}

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





// ------------------------------- old functions -------------------------------
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



/*
	bestSpreadExp := -1000.0
	var bestSpread spread
	var spread_tmp spread
	var riskMatchBool bool
	spreadStart := 0
	timeStart := time.Now()
	var percSteps float64
	if riskCompare{
		percSteps = 0.1
	} else {
		percSteps = 0.1
	}
	spreadCount := 0
	riskMatchCount := 0
	ws := []float64{0.0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0}
	//ws := []float64{0.25,0.5,0.75}
	totalCount := (int)(4*len(ws)*len(callList)*(len(callList)-1)/2)
	for i := 0 ; i < len(callList)-1 ; i++ {
		for j := i+1 ; j < len(callList)-1 ; j++{
			for _,w := range ws {

				//buy-buy
				spread_tmp = spread{
					num:     2,
					calls:   []callfunc{callList[i],callList[j]},
					weights: []float64{w,(1.0-w)},
				}
				spreadCount++
				//totalCount++
				if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
					if riskCompare {
						rp,err := spread_tmp.riskProfile(pdist)
						if len(rp.x) == 0 || len(rp.y) == 0 {continue}
						if err == nil{
							riskMatchBool = riskMatch(rp,riskTolSpline)
						} else {
							fmt.Println("ERROR in riskProfile:",err)
						}
					} else {riskMatchBool = true}
					if riskMatchBool {riskMatchCount++}
					if riskMatchBool {
						//spreads = append(spreads,spread_tmp)
						bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
						bestSpread = spread_tmp
					}
				}


				if selling {

					//buy-sell
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{w,-(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {
							rp,err := spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}

					//sell-buy
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{-w,(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {
							rp,err := spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}


					//sell-sell
					spread_tmp = spread{
						num:     2,
						calls:   []callfunc{callList[i],callList[j]},
						weights: []float64{-w,-(1.0-w)},
					}
					spreadCount++
					//totalCount++
					if spread_tmp.ExpectedReturn(pdist) > bestSpreadExp {
						if riskCompare {
							rp,err := spread_tmp.riskProfile(pdist)
							if len(rp.x) == 0 || len(rp.y) == 0 {continue}
							if err == nil{
								riskMatchBool = riskMatch(rp,riskTolSpline)
							} else {
								fmt.Println("ERROR in riskProfile:",err)
							}
						} else {riskMatchBool = true}
						if riskMatchBool {riskMatchCount++}
						if riskMatchBool {
							//spreads = append(spreads,spread_tmp)
							bestSpreadExp = spread_tmp.ExpectedReturn(pdist)
							bestSpread = spread_tmp
						}
					}

				}



				//if math.Mod(float64(spreadCount)/float64(totalCount),percSteps) < percSteps/2000 {
				if math.Mod(float64(spreadCount),percSteps*float64(totalCount)) < 4 {
					elapsed := time.Now().Sub(timeStart).Milliseconds()
					elapsedPerSpread := float64(elapsed)/float64(spreadCount-spreadStart)
					fmt.Println(fmt.Sprintf("%.1f",100.0*float64(spreadCount)/float64(totalCount)) , "% (took "+fmt.Sprintf("%v",elapsed)+" milliseconds - ",elapsedPerSpread," per spread)")
					timeStart = time.Now()
					spreadStart = spreadCount
				}

			}
		}
	}
	if riskMatchCount == 0 {
		continue
	}

	riskTolExclusion := ""
	if riskCompare {
		riskTolExclusion = fmt.Sprintf("%.2f Percent (%v out of %v) of spreads were excluded due to the risk profile not matching.\n",100.0*float64(totalCount-riskMatchCount)/float64(totalCount),totalCount-riskMatchCount,totalCount)
		fmt.Println(riskTolExclusion)
	}
*/