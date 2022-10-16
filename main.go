package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"
)

type Status struct {
	Water       int    `json:"water"`
	Wind        int    `json:"wind"`
	StatusWater string `json:"statuswater"`
	StatusWind  string `json:"statuswind"`
}

type Data struct {
	Status `json:"status"`
}

func updateData() {

	for {
		var data = Data{Status: Status{}}
		Min := 1
		Max := 20

		data.Status.Water = rand.Intn(Max - Min + 1)

		data.Status.Wind = rand.Intn(Max - Min + 1)

		switch {
		case data.Status.Water <= 5:
			data.Status.StatusWater = "Aman"
		case data.Status.Water <= 8:
			data.Status.StatusWater = "Siaga"
		default:
			data.Status.StatusWater = "Bahaya"
		}

		switch {
		case data.Status.Wind <= 6:
			data.Status.StatusWind = "Aman"
		case data.Status.Wind <= 15:
			data.Status.StatusWind = "Siaga"
		default:
			data.Status.StatusWind = "Bahaya"
		}

		b, err := json.MarshalIndent(&data, "", " ")

		if err != nil {
			log.Fatalln("error while marshalling json data  =>", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)

		if err != nil {
			log.Fatalln("error while writing value to data.json file  =>", err.Error())
		}
		fmt.Println("wait 10 seconds")
		time.Sleep(time.Second * 10)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file := path.Join("views", "index.html")
		tpl, _ := template.ParseFiles(file)

		var data = Data{Status: Status{}}

		b, err := ioutil.ReadFile("data.json")

		if err != nil {
			fmt.Fprint(w, "error while reading data.json file  =>", err.Error())
			return
		}

		err = json.Unmarshal(b, &data)

		err = tpl.ExecuteTemplate(w, "index.html", data.Status)
	})

	http.ListenAndServe(":8000", nil)
}
