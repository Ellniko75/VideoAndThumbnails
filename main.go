package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type user struct {
	Name    string
	Age     string
	Surname string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Headers", "*")
}

func FirstEndpoint(w http.ResponseWriter, r *http.Request) {
	videoname := r.URL.Query()["videoname"][0]
	log.Println(videoname)
	strVideoPath := fmt.Sprint("./Video/", videoname, ".mp4")
	fmt.Print(strVideoPath)
	http.ServeFile(w, r, strVideoPath)
}

func saveScreenshots(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	w.Write(b)
	//userToSend := user{Name: "juan", Age: "20", Surname: "elmatador"}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(userToSend)
}
func saveVideo(w http.ResponseWriter, r *http.Request) {
	//type video struct {
	//	VideoArr  []int8
	//	VideoName string
	//}
	enableCors(&w)
	v, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	//decode the json
	//var videoData video
	//err234 := json.Unmarshal(v, &videoData)
	//if err234 != nil {
	//	log.Println(err234, "Error aca xd")
	//}

	erro := os.WriteFile("xd.mp4", v, 0o777)
	if erro != nil {
		log.Println("ERROR ON WRITING", erro)
	}
	log.Println("File created")
	//Creates the thumbnails for that video, one each 4 seconds
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprint("xd", ".mp4"), "-vf", "fps=1/4", "%04d.png")
	cmd.Run()

	//userToSend := user{Name: "juan", Age: "20", Surname: "elmatador"}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(userToSend)
}
func saveVideoAndThumbnail(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	v, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if len(v) == 0 {
		return
	}

	finalVideoLength := len(v) - 20
	videoPortion := v[:finalVideoLength]
	videoname := getNameFromArray(v)

	erro := os.WriteFile(fmt.Sprint("Video/", videoname, ".mp4"), videoPortion, 0o777)
	if erro != nil {
		log.Println("ERROR ON WRITING", erro)
	}

	//Creates the thumbnails for that video, one each 4 seconds
	cmd := exec.Command("ffmpeg", "-i", fmt.Sprint("Video/", videoname, ".mp4"), "-vf", "fps=1/4", "./Video/%04d.png")
	cmd.Run()

}

func main() {
	http.HandleFunc("/", FirstEndpoint)
	http.HandleFunc("/postvideo", saveScreenshots)
	http.HandleFunc("/uploadVideo", saveVideo)
	http.HandleFunc("/uploadVideoAndThumbnail", saveVideoAndThumbnail)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// gets the name from the array and returns it in string format
func getNameFromArray(arr []byte) string {
	//this is where the video part of the slice ends
	whereVideoSectionEnds := len(arr) - 20
	//here starts the name
	namePortion := arr[whereVideoSectionEnds:]

	//remove padding of blank space (ascii code 32) from the end, until you hit something that is not a 32
	for i := len(namePortion) - 1; i >= 0; i-- {
		if namePortion[i] != 32 {
			return string(namePortion[:i+1])

		}
	}
	return ""

}
