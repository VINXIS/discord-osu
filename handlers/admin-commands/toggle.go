package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	tools "../../tools"
	gencommands "../general-commands"
	"github.com/bwmarrin/discordgo"
)

// Toggle toggles server options on/off
func Toggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	// Set new information in server data
	serverData.Time = time.Now()
	flagged := false
	if strings.Contains(m.Content, "-a") || strings.Contains(m.Content, "-announce") {
		if serverData.AnnounceChannel == m.ChannelID {
			serverData.AnnounceChannel = ""
		} else {
			serverData.AnnounceChannel = m.ChannelID
		}
		flagged = true
	}
	if strings.Contains(m.Content, "-d") || strings.Contains(m.Content, "-daily") {
		serverData.Daily = !serverData.Daily
		flagged = true
	}
	if strings.Contains(m.Content, "-os") || strings.Contains(m.Content, "-osu") {
		serverData.OsuToggle = !serverData.OsuToggle
		flagged = true
	}
	if strings.Contains(m.Content, "-s") || strings.Contains(m.Content, "-stats") {
		serverData.AllowAnyoneStats = !serverData.AllowAnyoneStats
		flagged = true
	}
	if strings.Contains(m.Content, "-v") || strings.Contains(m.Content, "-vibe") {
		serverData.Vibe = !serverData.Vibe
		flagged = true
	}
	if !flagged {
		s.ChannelMessageSend(m.ChannelID, "No flags given! Please use one of the flags listed in `help toggle`!")
		return
	}

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	go gencommands.ServerInfo(s, m)
}
