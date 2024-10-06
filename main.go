package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/json"
	"io"
	"github.com/bwmarrin/discordgo"
)

var (
	Token     = os.Getenv("TOKEN")
	ChannelID0 = os.Getenv("CHANNEL_ID_0")
	ChannelID1 = os.Getenv("CHANNEL_ID_1")
)

func main() {
	apiURL := os.Getenv("API_URL")

    // Check if the environment variable is set
    if apiURL == "" {
        fmt.Println("API_URL error tidak disetting")
        return
    }

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error membuat sesi Discord:", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error membuka koneksi:", err)
		return
	}

	fmt.Println("Bot berjalan. Tekan CTRL-C untuk keluar.")
	go renameChannel(dg, ChannelID0, "0", apiURL)
	go renameChannel(dg, ChannelID1, "1", apiURL)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func renameChannel(s *discordgo.Session, channelID string, gameID string, apiURL string) {
	for {
		// Panggil fungsi terpisah untuk mendapatkan status game
		updateChannelName(s, channelID, gameID, apiURL)
		time.Sleep(5 * time.Minute)
	}
}

func updateChannelName(s *discordgo.Session, channelID string, gameID string, apiURL string) {
	url := fmt.Sprintf("%s/status?games=%s", apiURL, gameID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error mengambil data dari API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error membaca body respons:", err)
		return
	}

	var data struct {
		Online       bool
		Players      struct {
			Online int
		}
		ServerStatus string
	}
	
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	newName := fmt.Sprintf("Online: %d", data.Players.Online)
	if gameID == "0" {
		newName = fmt.Sprintf("Minecraft : %d players", data.Players.Online)
	} else if gameID == "1" {
		newName = fmt.Sprintf("Ragnarok : %d players", data.Players.Online)
	}
	_, err = s.ChannelEdit(channelID, &discordgo.ChannelEdit{
		Name: newName,
	})
	if err != nil {
		fmt.Println("Error mengubah nama saluran:", err)
	} else {
		fmt.Printf("Nama saluran diubah menjadi: %s\n", newName)
	}
}
