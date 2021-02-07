package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"server-side-mirror/servers/gateway/handlers"
	"server-side-mirror/servers/gateway/indexes"
	"server-side-mirror/servers/gateway/models/users"
	"server-side-mirror/servers/gateway/sessions"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// IndexHandler does something
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from API Gateway"))
}

// Director is a function wrapper
type Director func(r *http.Request)

// CustomDirector does load balancing using the round-robin method
func CustomDirector(targets []*url.URL, ctx *handlers.HandlerContext) Director {
	var counter int32
	counter = 0

	return func(r *http.Request) {
		state := &handlers.SessionState{}
		_, err := sessions.GetState(r, ctx.Key, ctx.SessionStore, state)
		if err != nil {
			r.Header.Del("X-User")
			log.Println("Error getting User from GetState")
			log.Println(err)
			return
		}

		userJSON, _ := json.Marshal(state.User)
		userString := string(userJSON)

		targ := targets[counter%int32(len(targets))]
		atomic.AddInt32(&counter, 1)
		r.Header.Add("X-User", userString)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
	}
}

func getAllUrls(addresses string) []*url.URL {

	urlSlice := strings.Split(addresses, ",")
	var urls []*url.URL
	for _, u := range urlSlice {
		url := url.URL{Scheme: "http", Host: u}
		urls = append(urls, &url)
	}
	return urls
}

//main is the main entry point for the server
func main() {
	address := os.Getenv("ADDR")
	// Default address the server should listen on
	if len(address) == 0 {
		address = ":443"
	}
	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	sessionkey := os.Getenv("SESSIONKEY")
	redisaddr := os.Getenv("REDISADDR")
	dsn := os.Getenv("DSN")

	// create redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisaddr, // use default Addr
	})
	redisStore := sessions.NewRedisStore(redisClient, 0)
	dsn = fmt.Sprintf("root:%s@tcp("+dsn+")/users", os.Getenv("MYSQL_ROOT_PASSWORD"))
	userStore := users.NewMysqlStore(dsn)

	messagesaddr := os.Getenv("MESSAGEADDR")
	summaryaddr := os.Getenv("SUMMARYADDR")
	messagingUrls := getAllUrls(messagesaddr)

	conns := make(map[int64]*websocket.Conn)
	socketStore := handlers.NewNotify(conns, &sync.Mutex{})
	IDs := make(map[string]int64)
	indexedUsers := indexes.NewTrie(IDs, &sync.Mutex{})
	userStore.IndexUsers(indexedUsers)
	ctx := handlers.NewHandlerContext(sessionkey, userStore, *indexedUsers, redisStore, *socketStore)

	// proxies
	messagesProxy := &httputil.ReverseProxy{Director: CustomDirector(messagingUrls, ctx)}
	summaryProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: summaryaddr})
	// starting a new mux session
	// mux := http.NewServeMux()

	mux := mux.NewRouter()

	mux.HandleFunc("/", IndexHandler)

	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/users/{userID}", ctx.UpdateUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/mine", ctx.SpecificSessionsHandler)
	mux.HandleFunc("/v1/user/search", ctx.SearchHandler)
	mux.Handle("/v1/summary", summaryProxy)
	mux.HandleFunc("/v1/ws", ctx.WebSocketConnectionHandler)
	mux.Handle("/v1/channels/{channelID}/members", messagesProxy)
	mux.Handle("/v1/channels/{channelID}", messagesProxy)
	mux.Handle("/v1/channels", messagesProxy)
	mux.Handle("/v1/messages/{messageID}", messagesProxy)

	wrappedMux := handlers.NewLogger(mux)

	// logging server location or errors
	log.Printf("server is listening %s...", address)
	// config := &tls.Config{
	// 	MinVersion: tls.VersionTLS12,
	// 	MaxVersion: tls.VersionTLS13,
	// 	CipherSuites: []uint16{
	// 		TLS_AES_128_GCM_SHA256,
	// 		TLS_AES_256_GCM_SHA384,
	// 		TLS_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// 	PreferServerCipherSuites: true,
	// 	ClientSessionCache:       tls.NewLRUClientSessionCache(128),
	// }

	// server := &http.Server{Addr: address, Handler: nil, TLSConfig: config}
	// log.Fatal(server.ListenAndServeTLS(tlsCertPath, tlsKeyPath))

	log.Fatal(http.ListenAndServeTLS(address, tlsCertPath, tlsKeyPath, wrappedMux))

	// server := &http.Server{Addr: ":4000", Handler: nil, TLSConfig: config}
	// http2.ConfigureServer(server, nil)

	// log.Printf("Staring webserver ...")
	// go http.ListenAndServe(":3000", nil)
	// server.ListenAndServeTLS(TLS_PUBLIC_KEY, TLS_PRIVATE_KEY)
	/* To host server:
	- change path until in folder with main.go in it
	- 'go install main.go'
	- navigate to 'go' bin folder and run main.exe
	*/
}
