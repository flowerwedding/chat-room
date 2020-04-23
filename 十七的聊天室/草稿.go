package main

/*
var upgrader = websocket.Upgrader{ // 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main(){
	_ = http.ListenAndServe(":8080", nil)
}

func init() {
	http.HandleFunc("/v/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn){
			for{
				mType,msg,_:=conn.ReadMessage()
				_ = conn.WriteMessage(mType, msg)
			}
		}(conn)
	})
}*/