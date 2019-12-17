package gencommands

import (
	"net/http"

	config "../../config"
	"github.com/bwmarrin/discordgo"
)

// Crab sends a crab rave image
func Crab(s *discordgo.Session, m *discordgo.MessageCreate) {
	response, err := http.Get(config.Conf.Crab)
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "crab.gif",
			Reader: response.Body,
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
