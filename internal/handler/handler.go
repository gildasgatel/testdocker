package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testdocker/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

type Handler interface {
	Home(c *gin.Context)
	WS(c *gin.Context)
}

type Handle struct {
	dataCSV map[string]string
}

func New() Handler {
	return &Handle{
		dataCSV: make(map[string]string),
	}
}

func (h *Handle) Home(c *gin.Context) {
	h.loadCSV()
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
func (h *Handle) WS(c *gin.Context) {
	h.handleWebSocket(c.Writer, c.Request)
}

func (h *Handle) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		broadcastMessage(msg)

		resp := h.handleResp(msg)
		broadcastMessage(resp)
	}
}

func broadcastMessage(msg []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
			client.Close()
			delete(clients, client)
		}
	}
}

func (h *Handle) handleResp(msg []byte) []byte {
	m := string(msg)
	sl := strings.Split(m, ":")
	if len(sl) < 2 {
		return msg
	}
	key := generalizeKey(sl[1])
	resp := []byte("Bot: ")

	rep, dist := utils.GetBestMatch(key)
	if dist > 4 {
		resp = append(resp, []byte("Désolé je n'ai pas compris la question.")...)
		return resp
	}
	value := h.dataCSV[rep]
	resp = append(resp, []byte(value)...)

	return resp
}

func (h *Handle) loadCSV() {
	file, err := os.Open("./data.csv")
	if err != nil {
		log.Println("Erreur lors de l'ouverture du fichier:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Erreur lors de la lecture du fichier CSV:", err)
		return
	}

	for _, row := range records {
		key := generalizeKey(row[0])
		utils.DataKey = append(utils.DataKey, key)
		value := row[1]
		h.dataCSV[key] = value
	}
}

func generalizeKey(value string) string {
	value = strings.ToLower(value)
	return strings.ReplaceAll(value, " ", "")
}
