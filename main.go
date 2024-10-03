package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token     = os.Getenv("TOKEN")
	ChannelID = os.Getenv("CHANNEL_ID")
)

func main() {
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
	go renameChannel(dg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func renameChannel(s *discordgo.Session) {
	for {
		newName := fmt.Sprintf("Waktu: %s", time.Now().Format("15:04:05"))
		_, err := s.ChannelEdit(ChannelID, &discordgo.ChannelEdit{
			Name: newName,
		})
		if err != nil {
			fmt.Println("Error mengubah nama saluran:", err)
		} else {
			fmt.Printf("Nama saluran diubah menjadi: %s\n", newName)
		}
		time.Sleep(5 * time.Minute)
	}
}
