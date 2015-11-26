package main 
import(
"fmt"
"net/http"
"github.com/julienschmidt/httprouter"
"log"
"strconv"
"bytes"

)
var maps map[int]string
func Putvalue(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k:= p.ByName("key_id")
	v := p.ByName("value")
	key,err:= strconv.Atoi(k)
	if err!=nil {
		panic(err)
	}else{
	maps[key]= v
	w.WriteHeader(200)
}
}
func GetOneData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k:= p.ByName("key_id")
	key, err := strconv.Atoi(k)
	if err!=nil {
		panic(err)
	}else{
	value, right:= maps[key]
	if right {
		jsonresponse := `{
                "key" : "` + strconv.Itoa(key) + `",
                "value" : "` + value + `"
            }`
          fmt.Fprintf(w, "%s\n", jsonresponse)
	}else{
		fmt.Printf("Not Valid Key")
	}
	}
}
func GetAllData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var buf bytes.Buffer
	buf.WriteString("{\n[\n")
	for key, value := range maps {
		   string := `{
                        "key" : "` + strconv.Itoa(key) + `",
                        "value" : "` + value + `"
                    },` + "\n"
                    buf.WriteString(string)
	}
	resp := buf.String()
	resp = resp[:len(resp)-2]
	resp = resp + "\n]\n}"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", resp)

}
func main() {
	maps= make(map[int]string)
	router := httprouter.New()
	router.PUT("/keys/:key_id/:value", Putvalue)
	router.GET("/keys", GetAllData)
	router.GET("/keys/:key_id", GetOneData)
	log.Fatal(http.ListenAndServe(":3001", router))
	
}