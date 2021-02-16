package main

import (
	"../google_protocol_buffer/pb/protobuf"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pmylund/go-cache"
	"hash/crc32"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"
	"hash/fnv"
)


//REFRENCES--------------------------------
//https://stackoverflow.com/questions/28400340/how-support-concurrent-connections-with-a-udp-server-using-go
//https://stackoverflow.com/questions/27625787/measuring-memory-usage-of-executable-run-using-golang
//https://stackoverflow.com/questions/31879817/golang-os-exec-realtime-memory-usage
//https://golangcode.com/print-the-current-memory-usage/

//-----------------------------------------
//CONTROL PANEL----------------------------
var debug = true
var MEMORY_LIMIT = 51200
var MULTI_CORE_MODE = true
var MAP_SIZE_MB = 70
var CACHE = 10
var CACHE_LIFESPAN=5// how long should the cache persist

//-----------------------------------------

//MAP_AND_CACHE----------------------------
var storage =make(map[string][]byte);
var mutex sync.Mutex;
// Create a cache with a default expiration time of 5 seconds, and which
// purges expired items every 1 seconds
var message_cache = cache.New(5*time.Second, 1*time.Nanosecond)

//-----------------------------------------

//HELPER FUNCTIONS-------------------------
/**
 * @Description: checks if port is valid
 * @param port
 * @return bool
 */
func check_if_port_is_valid(port string) bool {
	casted_port, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("CASTING ERROR [port valid]")
	}
	if (casted_port > 65535) || (casted_port < 0) {
		return false
	} else {
		return true
	}
}
/**
 * @Description:  returns the ieee checksum
 * @param a
 * @param b
 * @return uint64
 */
func calculate_checksum(a []byte, b []byte) uint64 {
	var concat_byte_arr = append(a, b...)
	var check_sum = uint64(crc32.ChecksumIEEE(concat_byte_arr))
	return check_sum
}
/**
 * @Description: converts the input key into a hash
 * @param s
 * @return uint32
 */
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
//-----------------------------------------
// MESSAGE PARSING FUNC--------------------
/**
 * @Description: peels back the message layer and searches the cache for a reacent message else calls message_handles
 * @param whole_message
 * @return []byte
 * @return bool
 */
func message_broker(whole_message []byte)([]byte,bool){
	cast_whole_req:=&protobuf.Msg{
		MessageID: nil,
		Payload:   nil,
		CheckSum:  0,
	}

	error := proto.Unmarshal(whole_message, cast_whole_req)
	if error != nil {
		fmt.Printf("\nUNPACK ERROR[W2] %+v\n", error)
	}
	if(cast_whole_req.CheckSum!=calculate_checksum(cast_whole_req.MessageID,cast_whole_req.Payload)){
		fmt.Printf("\n[CHECKSUM WRONG NO RESPONSE SENT] %+v\n", error)
		return nil,false;
	}
	data_cached,found:=check_cache(cast_whole_req.MessageID);
	if(found){
		res_message:=&protobuf.Msg{
			MessageID: cast_whole_req.MessageID,
			Payload:   data_cached,
			CheckSum:  calculate_checksum(cast_whole_req.MessageID,data_cached),
		}
		payload, err := proto.Marshal(res_message)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		} else {
			fmt.Printf("payload generated\n")
		}
		return payload,true

	} else {
		uncached_data:=message_handler(cast_whole_req.Payload);
		cache_data(cast_whole_req.MessageID,uncached_data);
		res_message:=&protobuf.Msg{
			MessageID: cast_whole_req.MessageID,
			Payload: uncached_data,
			CheckSum:  calculate_checksum(cast_whole_req.MessageID,uncached_data),
		}
		payload, err := proto.Marshal(res_message)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		} else {
			fmt.Printf("payload generated\n")
		}
		return payload,true
	}


}
/**
 * @Description: peals the seconday message layer and performs server functions returns the genarated payload
 * @param message
 * @return []byte
 */
func message_handler(message []byte) []byte {
	cast_req := &protobuf.KVRequest{
		Command: 0,
		Key:     nil,
		Value:   nil,
		Version: nil,
	}
	error := proto.Unmarshal(message, cast_req)
	if error != nil {
		fmt.Printf("\nUNPACK ERROR %+v\n", error)
	}

	//* 0x01 - Put: This is a put operation.
	//* 0x02 - Get: This is a get operation.
	//* 0x03 - Remove: This is a remove operation.
	//* 0x04 - Shutdown: shuts-down the node (used for testing and management).
	//* 0x05 - Wipeout: deletes all keys stored in the node (used for testing).
	//* 0x06 - IsAlive: does nothing but replies with success if the node is alive.
	//* 0x07 - GetPID: the node is expected to reply with the processID of the Go process
	//* 0x08 - GetMembershipCount:(This will be used later, for now you are expected to return 1.)
	switch cast_req.GetCommand() {
	case 1:
		{
			if debug {
				fmt.Println("PUT")
			}
			return put(cast_req.GetKey(), cast_req.GetValue(), cast_req.GetVersion())
		}
	case 2:
		{
			if debug {
				fmt.Println("GET")
			}
			return get(cast_req.Key);
		}
	case 3:
		{
			if debug {
				fmt.Println("REMOVE")
			}
			return remove(cast_req.Key);
		}
	case 4:
		{
			if debug {
				fmt.Println("SHUTDOWN")
			}
			fmt.Println("[TERMINATING]....")
			shutdown();
		}
	case 5:
		{
			if debug {
				fmt.Println("WIPEOUT")
			}
			return wipeout();
		}
	case 6:
		{
			if debug {
				fmt.Println("IsALIVE")
			}
			return is_alive();
		}
	case 7:
		{
			if debug {
				fmt.Println("GetPID")
			}
			return getpid();
		}
	case 8:
		{
			if debug {
				fmt.Println("GetMembershipCount")
			}
			return getmemcount();
		}
	default:
		{
			if debug {
				fmt.Println("INVALID COMMAND")
			}
			return message;
		}
	}
	return is_alive();
}
/**
 * @Description:  check if the data is in the cache
 * @param message_id
 * @return []byte
 * @return bool
 */
func check_cache(message_id []byte)([]byte,bool){
	response, found := message_cache.Get(string(message_id))

	if ! found {
		return nil, false;
	} else {
		str := fmt.Sprintf("%v", response)
		var cached_data=[]byte(str);
		return cached_data,true;
	}
}
/**
 * @Description: put data in the cache
 * @param message_id
 * @param data
 * @return bool
 */
func cache_data(message_id []byte,data []byte)(bool){
	if get_mem_usage() > uint64(MEMORY_LIMIT-20) {
		return true;
		fmt.Println("\n[MEMORY WARNING]\n")
		message_cache.Flush();
	}
	if get_mem_usage() > uint64(MEMORY_LIMIT/2) {
		fmt.Println("\n[MEMORY WARNING 50%]\n")
		message_cache.DeleteExpired();
	}

	message_cache.Set(string(message_id), string(data), cache.DefaultExpiration);
	return true;
}
//-----------------------------------------
//DATABASE FUNCTIONS-----------------------

/**
 * @Description:Puts some value (and corresponding version)
 *into the store. The value (and version) can be later retrieved using the key.
 * @param key
 * @param value
 * @param version
 * @return []byte
 */
func put(key []byte, value []byte, version int32) []byte {
	var error_code uint32=0;
	error_code=0;
	if(int(len(value))>10000) {
		error_code=7;

		var err_code *uint32;
		err_code=new(uint32);
		err_code=&error_code;

		var _pid *int32;
		_pid=new(int32);
		procress_id:=int32(syscall.Getpid())
		_pid=&procress_id
		payload:=&protobuf.KVResponse{
			ErrCode: err_code,
			Value:   value,
			Pid:_pid,
		}
		out, err := proto.Marshal(payload)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		} else {
			if(debug){
				//fmt.Printf("NEW BUFFER-> %+v",out);
			}
		}

		return out;
	}

	mutex.Lock()//<<<<<<<<<<<<<<<MAP LOCK

	if(bToMb(uint64(unsafe.Sizeof(storage)))<uint64(MAP_SIZE_MB)){

		storage[string(key)] = value // adding the value
		error_code=0;
		if debug {
			fmt.Println("PUT",string(key),",<->",string(value));
		}
	} else{
		if debug {
			fmt.Println("ERROR PUTTING");
		}
		error_code=2;
	}
	mutex.Unlock()//<<<<<<<<<<<<<<<MAP UNLOCK

	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Value:   value,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}

	return out;
}
/**
 * @Description:Returns the value and version that is associated with the key. If there is no such key in your store, the store should return error (not found).
 * @param key
 * @return []byte
 */
func get(key []byte) []byte{
	mutex.Lock();
	value, found := storage[string(key)];
	mutex.Unlock();
	var error_code uint32;
	if(found) {
		error_code=0;
	} else {
		error_code=1;
	}

	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Value:   value,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}

	return out;
}
/**
 * @Description:Removes the value that is associated with the key.
 * @param key
 * @return []byte
 */
func remove(key []byte) []byte{
	mutex.Lock();
	value, found := storage[string(key)];
	if(found){
		delete(storage, string(key));
	}
	mutex.Unlock();
	if(found) {
		var error_code uint32;
		error_code=0;
		var err_code *uint32;
		err_code=new(uint32);
		err_code=&error_code;

		var _pid *int32;
		_pid=new(int32);
		procress_id:=int32(syscall.Getpid())
		_pid=&procress_id
		payload:=&protobuf.KVResponse{
			ErrCode: err_code,
			Value:   value,
			Pid:_pid,
		}
		out, err := proto.Marshal(payload)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		} else {
			if(debug){
				//fmt.Printf("NEW BUFFER-> %+v",out);
			}
		}

		return out;

	} else {
		var error_code uint32;
		error_code=1;
		var err_code *uint32;
		err_code=new(uint32);
		err_code=&error_code;

		var _pid *int32;
		_pid=new(int32);
		procress_id:=int32(syscall.Getpid())
		_pid=&procress_id
		payload:=&protobuf.KVResponse{
			ErrCode: err_code,
			Value:   value,
			Pid:_pid,
		}
		out, err := proto.Marshal(payload)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		} else {
			if(debug){
				//fmt.Printf("NEW BUFFER-> %+v",out);
			}
		}

		return out;
	}
}
/**
 * @Description: calls os.shutdown
 */
func shutdown(){
	os.Exit(555);
}
/**
 * @Description: clears the database
 * @return []byte
 */
func wipeout() []byte{
	mutex.Lock()//<<<<<<<<<<<<<<<MAP LOCK
	for k := range storage {
		delete(storage, k)
	}
	mutex.Unlock()//<<<<<<<<<<<<<<<MAP UNLOCK
	var error_code uint32;
	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}

	return out;
}
/**
 * @Description: response indicating server is alive
 * @return []byte
 */
func is_alive()[]byte{
	fmt.Println("CLIENT ASKED IF SERVER ALIVE");

	var error_code uint32;
	error_code=0;

	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}
	return out;
}
/**
 * @Description: gets the current procressID
 * @return []byte
 */
func getpid()[]byte{
	var error_code uint32;
	error_code=0;

	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}
	return out;
}
/**
 * @Description: returns number of members
 * @return []byte
 */
func getmemcount() []byte{
	var error_code uint32;
	error_code=0;

	var err_code *uint32;
	err_code=new(uint32);
	err_code=&error_code;

	var _pid *int32;
	_pid=new(int32);
	procress_id:=int32(syscall.Getpid())
	_pid=&procress_id
	payload:=&protobuf.KVResponse{
		ErrCode: err_code,
		Pid:_pid,
	}
	out, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	} else {
		if(debug){
			//fmt.Printf("NEW BUFFER-> %+v",out);
		}
	}
	return out;
}
//-----------------------------------------

//UDP SERVER FUNC--------------------------
/**
 * @Description:  go routine to serve a client
 * @param connection
 * @param conduit
 * @param thread_num
 */
func UDP_daemon(connection *net.UDPConn, conduit chan int, thread_num int) {
	if debug {
		fmt.Println("THREAD SPAWNED->", thread_num)
	}
	buffer := make([]byte, 65535)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for err != nil {
		fmt.Println("listener failed - ", err)
		conduit <- 2
	}
	n, remoteAddr, err = connection.ReadFromUDP(buffer)
	conduit <- 1
	response_val,valid:=message_broker(buffer[:n]);
	if(valid){
		n, err = connection.WriteToUDP(response_val, remoteAddr)
		if err != nil {
			fmt.Println("[FAIL] err serving ->", remoteAddr,"error is ",err)
			fmt.Println("\n[IMP ERROR] message size was ",len(response_val));
		}
		if debug {

			fmt.Println("THREAD EXIT->", thread_num, "\n")
		}
	} else{
		fmt.Println("WRONG CHECKSUM");
	}


}
/**
 * @Description: genrates several go routine depending on the memory
 * @param _port
 * @return func()
 */
func spawn_UDP_daemon(_port string) func() {
	port_num, err := strconv.Atoi(_port)
	if err != nil {
		fmt.Printf("PORT CASTING ERROR [EXIT]")
	}
	address := net.UDPAddr{
		Port: port_num,
		IP:   net.IP{127, 0, 0, 1}, // local ip address
	}
	return func() {
		var thread_num int = 0
		connection, err := net.ListenUDP("udp", &address)
		err=connection.SetWriteBuffer(20000);
		if(err!=nil){
			fmt.Printf("[WRITE SETTING ERROR]%+v",err);
		}
		err=connection.SetReadBuffer(20000);
		if(err!=nil){
			fmt.Printf("[READ SETTING ERROR]%+v",err);
		}
		if err != nil {
			panic(err)
		}
		conduit := make(chan int)
		if MULTI_CORE_MODE {
			for i := 0; i < runtime.NumCPU(); i++ {
				go UDP_daemon(connection, conduit, thread_num)
				thread_num++
			}
		} else {
			go UDP_daemon(connection, conduit, thread_num);
			thread_num++;
		}
		if !debug {
			fmt.Print("\033[s")
		}

		for {
			if get_mem_usage() < uint64(MEMORY_LIMIT) {
				switch <-conduit {
				case 1:
					go UDP_daemon(connection, conduit, thread_num)
					thread_num++
				case 2:
					if debug {
						fmt.Println("THREAD NUMBER ", thread_num, "EXITED ON ERROR")
					}
					thread_num++
				}
				if !debug {
					fmt.Print("\033[u\033[K")
					fmt.Println("NUMBER OF ACTIVE THREADS", thread_num)
				}
			} else {
				fmt.Printf("\n[HALT] MEMORY FULL calling garbage collector-> %+v\n",get_mem_usage())
				time.Sleep(1 * time.Second)
				message_cache.Flush();
				message_cache.DeleteExpired();
				runtime.GC();
			}
		}

	}
}
/**
 * @Description: inialises the server based on port number
 * @param args
 */
func init_server(args string) {
	spawner := spawn_UDP_daemon(args)
	if MULTI_CORE_MODE {
		fmt.Println("[MULTICORE MODE] [", runtime.NumCPU(), "] SPANNERS AT PORT [",args,"] SYS MEM LIMIT [",MEMORY_LIMIT,"]");
	}
	spawner()
}
/**
 * @Description: for debug to view packet b4 sending
 * @param arr
 */
func double_check(arr []byte){
	cast_whole_req:=&protobuf.Msg{
		MessageID: nil,
		Payload:   nil,
		CheckSum:  0,
	}

	error := proto.Unmarshal(arr, cast_whole_req)
	if error != nil {
		fmt.Printf("\nUNPACK ERROR[W2] %+v\n", error)
	}
	if(cast_whole_req.CheckSum!=calculate_checksum(cast_whole_req.MessageID,cast_whole_req.Payload)){
		fmt.Printf("\n[CHECKSUM WRONG NO RESPONSE SENT] %+v\n", error)
	}

	cast_req := &protobuf.KVResponse{
		ErrCode:          nil,
		Value:            nil,
		Pid:              nil,
	}
	error = proto.Unmarshal(cast_whole_req.Payload, cast_req)
	if error != nil {
		fmt.Printf("\nUNPACK ERROR %+v\n", error)
	} else {
		//fmt.Println("\n XXXXX values areXXXXX \n");
		//fmt.Println(cast_req.ErrCode);
		//fmt.Println(cast_req.Value);
		//fmt.Println(cast_req.Pid);
	}



}

//copied form
//https://golang.org/pkg/runtime/#MemStats
/**
 * @Description: retunrs current sys mem usage
 * @return uint64
 */
func get_mem_usage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if debug {
		//fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
		//fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
		//fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	}
	return bToMb(m.Sys)
}
/**
 * @Description: convert bytes to mb
 * @param b
 * @return uint64
 */
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//-----------------------------------------
//MAIN FUNCTION
//-----------------------------------------
func main() {

	argsWithProg := os.Args
	if len(argsWithProg) != 2 {
		fmt.Printf("INVALID NUMBER OF ARGUMENTS, EXITTING....\n")
	} else {
		if !check_if_port_is_valid(argsWithProg[1]) {
			fmt.Printf("PORT NOT VALID,EXITTING...\n")
		} else {
			fmt.Printf("------------------\n")
			init_server(argsWithProg[1]) // initaliaze;
		}
	}

}