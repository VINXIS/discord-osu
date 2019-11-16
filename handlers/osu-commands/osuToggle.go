package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// OsuToggle toggles beatmap/user messages on/off
func OsuToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if m.Author.ID != server.OwnerID {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin && len(m.Mentions) >= 1 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
			return
		}
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.OsuToggle = !serverData.OsuToggle

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if serverData.OsuToggle {
		s.ChannelMessageSend(m.ChannelID, "Enabled map/user links O_o")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Disabled map/user links O_o")
	}
	return
}