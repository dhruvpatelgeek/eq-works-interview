package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)





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
}
/**
 * @Description: converts the input key into a hash
 * @param s
 * @return uint32
 */
func hash(s string) uint32 {

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
}
/**
 * @Description: peals the seconday message layer and performs server functions returns the genarated payload
 * @param message
 * @return []byte
 */
func message_handler(message []byte) []byte {
}
/**
 * @Description:  check if the data is in the cache
 * @param message_id
 * @return []byte
 * @return bool
 */
func check_cache(message_id []byte)([]byte,bool){
}
/**
 * @Description: put data in the cache
 * @param message_id
 * @param data
 * @return bool
 */
func cache_data(message_id []byte,data []byte)(bool){

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
}
/**
 * @Description:Returns the value and version that is associated with the key. If there is no such key in your store, the store should return error (not found).
 * @param key
 * @return []byte
 */
func get(key []byte) []byte{
}
/**
 * @Description:Removes the value that is associated with the key.
 * @param key
 * @return []byte
 */
func remove(key []byte) []byte{
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
}
/**
 * @Description: response indicating server is alive
 * @return []byte
 */
func is_alive()[]byte{
}
/**
 * @Description: gets the current procressID
 * @return []byte
 */
func getpid()[]byte{
}
/**
 * @Description: returns number of members
 * @return []byte
 */
func getmemcount() []byte{
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
}
/**
 * @Description: genrates several go routine depending on the memory
 * @param _port
 * @return func()
 */
func spawn_UDP_daemon(_port string) func() {
}
/**
 * @Description: inialises the server based on port number
 * @param args
 */
func init_server(args string) {
}

/**
 * @Description: retunrs current sys mem usage
 * @return uint64
 */
func get_mem_usage() uint64 {
}

//-----------------------------------------
//MAIN FUNCTION
//-----------------------------------------
func main() {

}
