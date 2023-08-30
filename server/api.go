package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

// Item represents an item for sale
type Frag struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// Path to parent
	ParentImg string `json:"parentImg"`
	// Paths to images associated with the item
	FragImg string  `json:"fragImg"`
	Price   float64 `json:"price"`
}

var frags = []Frag{
	{
		ID:        0,
		Name:      "Frag 1",
		ParentImg: "/static/test-images/parent.jpeg",
		FragImg:   "/static/test-images/frag.jpeg",
		Price:     1000.00,
	},
	{
		ID:        1,
		Name:      "Frag 2",
		ParentImg: "/static/test-images/pic1.jpeg",
		FragImg:   "/static/test-images/pic2.jpeg",
		Price:     100.00,
	},
	{
		ID:        2,
		Name:      "Frag 3",
		ParentImg: "/static/test-images/pic3.jpeg",
		FragImg:   "/static/test-images/pic4.jpeg",
		Price:     10.00,
	},
}

func makeApiHandler() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/forsale/frags", forsaleFragsHandler)
	r.HandleFunc("/forsale/frags/buy", forsaleFragsBuyHandler)
	return r
}

func forsaleFragsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(frags)
}

var (
	emailAddr     = os.Getenv("EMAIL_ADDR")
	emailPwd      = os.Getenv("EMAIL_PWD")
	emailSrvr     = os.Getenv("EMAIL_SRVR")
	emailSrvrPort = os.Getenv("EMAIL_SRVR_PORT")
)
var emailAuth = smtp.PlainAuth(
	"", emailAddr, emailPwd, emailSrvr,
)

func forsaleFragsBuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	idStr := r.URL.Query().Get("frag_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id >= uint64(len(frags)) {
		idStr = ""
	}
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	frag := frags[id]
	createMsg := func(to string) []byte {
		return []byte(
			"To: " + to + "\r\n" +
				"Subject: You've got $$$\r\n" +
				"\r\n" +
				fmt.Sprintf("A purchase of $%.2f was made!!!\r\n", frag.Price),
		)
	}
	addr := emailSrvr + ":" + emailSrvrPort
	to := "9726793337@txt.att.net"
	if err := smtp.SendMail(
		addr, emailAuth, emailAddr,
		[]string{to}, createMsg(to),
	); err != nil {
		log.Printf("error sending to %s: %v", to, err)
	}
}
