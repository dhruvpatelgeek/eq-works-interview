package main

import (
	"../../google_protocol_buffer/pb/protobuf"
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	guuid "github.com/google/uuid"
	"github.com/pmylund/go-cache"
	"hash/crc32"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)
//MESSAGE LIMIT CACHE--------------------
const RATE_LIMIT = 5*time.Second
const REQUEST_LIMIT=10;
// you will have no more than 10 requests in 5 seconds

const PURGE_TIME=1*time.Second;
const UPLOAD_RATE=5*time.Second;

//----------------------------------------//
var messageCache = cache.New(RATE_LIMIT,PURGE_TIME)
// stack to store requests
var counterStack =make(map[string]string);
//mutex to protect the map
var mapMutex sync.Mutex;
//database port
var database_port="3000"
type counters struct {
	sync.Mutex
	view  int
	click int
}

var (
	c = counters{}
	content = []string{"sports", "entertainment", "business", "education"}
)
//----------------------------------------


func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := content[rand.Intn(len(content))]

	c.Lock()
	c.view++
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}

	addToStack(data);

}

func addToStack(data string) {
dt := time.Now()
//Format MM-DD-YYYY hh:mm:ss
currTime :=dt.Format("01-02-2006 15:04")
key:=data+":"+ currTime;
// if you want to use JSON you can parse this as JSON
// for now i am just going to push as string
views:= strconv.Itoa((c.view))
clicks:= strconv.Itoa((c.click))
value:="{views:"+views+",clicks:"+clicks+"}";
fmt.Println("value->",value);
mapMutex.Lock()
counterStack[key]=value;
mapMutex.Unlock()
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}
}

func isAllowed() bool {
	if(messageCache.ItemCount()<REQUEST_LIMIT) {
		dt := time.Now()
		//Format MM-DD-YYYY hh:mm:ss
		currTime :=dt.Format("01-02-2006 15:04")
		messageCache.Add(currTime,"dummy",cache.DefaultExpiration);
		return true;
	} else {
		return false
	}
}



func uploadCounters() error {
	for{

		time.Sleep(UPLOAD_RATE)
		mapMutex.Lock()
		if(len(counterStack)>1) {
			for k, v := range counterStack {
				fmt.Printf("[SENDING] key[%s] value[%s]\n", k, v)
				sendRequest(k,v);
				delete(counterStack, k);
			}
		}
		mapMutex.Unlock()
	}
}
func startServer() func(port string){
	go uploadCounters()
	return func (port string) {
		http.HandleFunc("/", welcomeHandler)
		http.HandleFunc("/view/", viewHandler)
		http.HandleFunc("/stats/", statsHandler)
		addr:=":"+port;
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
//COMUNICATION WITH satellite------------------------------------
func checkIfPortIsValid(port string) bool {
	castedPort, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("CASTING ERROR [port valid]")
	}
	if (castedPort > 65535) || (castedPort < 0) {
		return false
	} else {
		return true
	}
}
func genUUID() string {
	id := guuid.New()
	//fmt.Printf("UUID GENERATED-> %s\n", id.String())
	return id.String()
}
func sendRequest(key string,value string){
	message,messageId:=generatePayload([]byte(key),[]byte(value));
	server_address:="127.0.0.1"+":"+database_port;
	firePayload(message,server_address,messageId);
}
func calculate_checksum(a []byte,b []byte) uint64{
	var concat_byte_arr=append(a,b...);
	var check_sum=uint64(crc32.ChecksumIEEE(concat_byte_arr))
	return check_sum;
}
func generatePayload(key []byte,value []byte) ([]byte,string){
	//* 0x01 - Put: This is a put operation.
	//* 0x02 - Get: This is a get operation.
	//* 0x03 - Remove: This is a remove operation.
	//* 0x04 - Shutdown: shuts-down the node (used for testing and management).
	//* 0x05 - Wipeout: deletes all keys stored in the node (used for testing).
	//* 0x06 - IsAlive: does nothing but replies with success if the node is alive.
	//* 0x07 - GetPID: the node is expected to reply with the processID of the Go process
	//* 0x08 - GetMembershipCount:(This will be used later)
	
	payload:=&protobuf.KVRequest{
		Command: 0x01,
		Key:     key,
		Value:   value,
	}
	shell, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	message_id:=genUUID()
	checksum:=calculate_checksum([]byte(message_id),shell);

	casing:=&protobuf.Msg{
		MessageID:[]byte(message_id),
		Payload:   shell,
		CheckSum: checksum,
	}

	casted_casing, err := proto.Marshal(casing)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	return casted_casing,message_id
}
func firePayload(shell []byte,server_ip string,message_id string){
	fire(shell,server_ip,0,100,message_id)
}
func fire(payload []byte,address string,itr int,timeout int64,message_id string){
	fmt.Printf("RETRYING REQUEST [%d]--------------------------\n",itr);
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	//writing message->
	fmt.Fprintf(conn,string(payload));
	fmt.Printf("packet-written");

	//reading RESPONSE-<
	response_payload :=  make([]byte, 10000)
	var byte_ctr int;
	timeoutDuration := time.Duration(timeout) * time.Millisecond;
	err = conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	byte_ctr, err = bufio.NewReader(conn).Read(response_payload)
	if e,ok := err.(net.Error); ok && e.Timeout() {
		fmt.Printf("\nTIMEOUT after %v\n",timeoutDuration)
		conn.Close();
		if(itr<3){
			fire(payload,address,itr+1,timeout+100,message_id);
		}
		return
	}else if err != nil {
		fmt.Println("\n[ERROR] reading\n",err);
		conn.Close();
		return
	} else {
		fmt.Printf("packet-received:\n bytes=%d \nfrom=%s\n",byte_ctr,address);
		fmt.Printf("\n>>>>>>>REPLY FORM SERVER[SUCCESS]\n")
		cast_response:=&protobuf.Msg{
			MessageID: nil,
			Payload:   nil,
			CheckSum:  0,
		}
		error:=proto.Unmarshal(response_payload[:byte_ctr],cast_response);
		if(error!=nil) {
			fmt.Printf("\n[e2e3]UNPACK ERROR %+v\n",error)
		}
		local_checksum:=calculate_checksum(cast_response.GetMessageID(),cast_response.GetPayload());
		var flag=false;
		if(local_checksum!=cast_response.GetCheckSum()){
			fmt.Printf("\n[CHECKSUM WRONG]\n")
			flag=true;
		}
		if(flag){
			if(itr<3) {
				fmt.Printf("\n[RETRYING]\n")
				fire(payload,address,itr+1,timeout+100,message_id);
			}
		} else{
			server_response:=cast_response.GetPayload();
			res_struct:=&protobuf.KVResponse{
				ErrCode: nil,
				Value:   nil,
				Pid:     nil,
			}
			error=proto.Unmarshal(server_response,res_struct);
			if(res_struct.GetErrCode()==0){
				fmt.Printf("\nVALLUE WRITTEN SUCESSFULLY\n----------");
				fmt.Printf("value written is \"%+v\"",string(res_struct.GetValue()));
			} else {
				fmt.Printf("\n[SERVER ERROR]satellite internal server error\n----------");
				fmt.Printf("\n[RETRYING]\n")
				fire(payload,address,itr+1,timeout+100,message_id);
			}
		}
		conn.Close();
		return
	}
	return
}
//---------------------------------------------------------------
func main() {

	argsWithProg := os.Args
	if len(argsWithProg) != 3 {
		fmt.Printf("formmat go run main.go [server port] [database port]\n")
	} else {
		if !checkIfPortIsValid(argsWithProg[1]) {
			fmt.Printf("[server port] NOT VALID,EXITTING...\n")
		} else if !checkIfPortIsValid(argsWithProg[2]){
			fmt.Printf("[database port] NOT VALID,EXITTING...\n")
		} else {
			fmt.Printf("STARTING SERVER------------------\n")
			server:=startServer()
			server(argsWithProg[1])
			database_port=argsWithProg[2]
		}
	}
}
