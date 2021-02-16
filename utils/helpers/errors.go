package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/birukbelay/item/utils/global"
)

// runtime.FuncForPC(pc).Name() :-- main.main
// fn  :--  /tmp/sandbox819000997/prog.go
// line : -- 32
// err :-- is the error

var re = regexp.MustCompile(`^.*\.(.*)$`)

func Trim(input string) string{
	res1 := strings.Trim(input, "C:/Users/b/go/src/github.com/birukbelay/F")
	return res1
}

//retrieveCallInfo shows which error handler was called
func retrieveCallInfo(status int) {
	pc, _, _, _ := runtime.Caller(1)

	funcName := runtime.FuncForPC(pc).Name()

	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash



	log.Printf("       %d \n  Err_handler_func: %s %s          ", status, Trim(funcName), funcName[lastDot :])
	//-----------------------------------------------------------------------

	pc1 := make([]uintptr, 15)
	n := runtime.Callers(2, pc1)
	frames := runtime.CallersFrames(pc1[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
}

func PrintTopUi1(funcName, name interface{}){
	//....,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
	fmt.Println("        |")
	fmt.Println("        |")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++|")
	fmt.Printf("                                          %v  %v                                                  ",funcName, name)
}

func PrintTopUiWithStatusCode( statusCode int){
	//....,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
	fmt.Println("        |")
	fmt.Println("        |")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++|")
	fmt.Printf("|                                         %v                                                               |\n",
		string(statusCode))
}

func PrintUi(sign string, amount int){
	fmt.Println(sign)

	var sum string
	//var i int
	for i := 1; i < amount; i++ {
		//fmt.Println(sum,"|")

		sum+=sign

	}
	fmt.Println(sum,"|")

}

func LogValue(name string,errs interface{}){

	PrintTopUi1("Log", name)

	pcc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pcc).Name()

	function_name := re.ReplaceAllString(funcName, "$1")

	log.Printf(":=:-> func_name:- %s \n:: [LogMessage]:- %v \n\n", function_name, errs)
	log.Printf(" %s",Trim(funcName))


	PrintUi(".", 40)

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	retrieveCallInfo(500)

	PrintUi("=", 40)

}

// The D/C bn  HandleErr & RenderResponse is that RenderResponse shows all the error as json
// Their main Dc is that they use dt methods to Get The function Name



//HandleErr only displays custom Header Errors
func HandleErr(w http.ResponseWriter, errs interface{}, header string, statusCode int){

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ResponseH", header)

	http.Error(w, global.ErrorText(header), statusCode)

	PrintTopUi1("HandleErr", "")
//------------------
	pcc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pcc).Name()
	//-------
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash
	log.Printf(":=:->FunctionName - .%s :> [ErrorMessage]:- %v  \n", funcName[lastDot+1:], errs)
	log.Printf("%s  \n",Trim(funcName))
	PrintUi("-", 40)

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	retrieveCallInfo(statusCode)

	PrintUi("=", 40)

}

//RenderResponse returns the full output /error message as Json
// outPut is the full message Returned; it could be the
func RenderResponse(w http.ResponseWriter, output interface{}, header string, statusCode int){

	jsonResponse, err := json.MarshalIndent(output, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ResponseH", header)

	w.WriteHeader(statusCode)
	_,_ = w.Write(jsonResponse)
	//http.Error(w, ErrorText(header), statusCode)


	//	=======================================- Logging -========================================
	//                                  For production The logging will be removed
	PrintTopUi1("renderResponse", statusCode)
//----------------------
	pcc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pcc).Name()

	package_name := re.ReplaceAllString(funcName, "$1")
	log.Printf("++++ function_name:- %s \n:> [OUTPUT]:- %v  \n\n", package_name, output)
	log.Printf(" %s ", Trim(funcName))
	fmt.Println("40...")

	PrintUi(".", 40)

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s:%d\n", frame.File, frame.Line)
	retrieveCallInfo(statusCode)
	fmt.Println("40...")

	PrintUi("=", 40)


}

// Being deprecated
////////////////////////////////////////////////////////////////////////////////////////////////////////////


// ErrSingle ...
func ErrSingle(w http.ResponseWriter, err error, errorCode int) error {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(errorCode), errorCode)

		//show error location
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf(">>>>[error] in \n %s \n [%s: %d]\n %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		retrieveCallInfo(errorCode)


		return err
	}
	return nil
}


// ErrMultiple ...
func ErrMultiple(w http.ResponseWriter, errs []error, errorCode int) []error {

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(errorCode), errorCode)

		//show error location
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf(">>>>[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, errs)
		retrieveCallInfo(errorCode)

		return errs
	}
	return nil

}

//Deprecated ..................
//JsonResponse returns the json indent of the response error
func JsonResponse(w http.ResponseWriter, response interface{}, header string, statusCode int) {

	jsonResponse, err := json.MarshalIndent(response, "", "\t\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ResponseH", header)

	w.WriteHeader(statusCode)
	_,_ = w.Write(jsonResponse)
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}


