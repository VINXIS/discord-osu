package osucommands

import (
	"regexp"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	"github.com/bwmarrin/discordgo"
)

// ScorePost posts your score in a single line format
func ScorePost(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData, postType string) {
	mapRegex, _ := regexp.Compile(`(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	scorePostRegex, _ := regexp.Compile(`sc?(orepost)?\s+(\S+)`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	mod2Regex, _ := regexp.Compile(`\+(\S+)`)
	scoreRegex, _ := regexp.Compile(`\*\*(([0-9]|,)+)\*\* `)

	var beatmap osuapi.Beatmap
	var username string
	var user osuapi.User
	mods := "NM"
	parsedMods := osuapi.Mods(0)
	scoreVal := int64(0)
	if postType == "scorePost" {
		if !scorePostRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, "You did not give a username / map / mods / anything! See `help sc` for more details.")
			return
		}
		username = scorePostRegex.FindStringSubmatch(m.Content)[2]
		if modRegex.MatchString(username) {
			mods = strings.ToUpper(modRegex.FindStringSubmatch(username)[1])
			if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
				mods += "DT"
			}
			parsedMods = osuapi.ParseMods(mods)

			username = strings.TrimSpace(strings.Replace(username, modRegex.FindStringSubmatch(username)[0], "", 1))
		}
		// Get the map
		var submatches []string
		if mapRegex.MatchString(m.Content) {
			submatches = mapRegex.FindStringSubmatch(m.Content)
		} else {
			// Get prev messages
			messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
				return
			}

			// Look for a valid beatmap ID
			for _, msg := range messages {
				if len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
					if mapRegex.MatchString(msg.Embeds[0].URL) {
						submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].URL)
						break
					} else if mapRegex.MatchString(msg.Embeds[0].Author.URL) {
						submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)
						break
					}
				} else if mapRegex.MatchString(msg.Content) {
					submatches = mapRegex.FindStringSubmatch(msg.Content)
					break
				}
			}
		}

		// Check if found
		if len(submatches) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}

		// Get the map
		nomod := osuapi.Mods(0)
		switch submatches[3] {
		case "s":
			beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
		case "b":
			beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
		case "beatmaps":
			beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
		case "beatmapsets":
			if len(submatches[7]) > 0 {
				beatmap = osutools.BeatmapParse(submatches[7], "map", &nomod)
			} else {
				beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
			}
		}
		if beatmap.BeatmapID == 0 {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		} else if beatmap.Approved < 1 {
			s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
			return
		}
		username = strings.TrimSpace(strings.Replace(username, submatches[0], "", -1))

		// Get user
		for _, player := range cache {
			if username != "" {
				if username == player.Osu.Username {
					user = player.Osu
					break
				}
			} else if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				user = player.Osu
				break
			}
		}

		// Check if user even exists
		if user.UserID == 0 {
			if username == "" {
				s.ChannelMessageSend(m.ChannelID, "No user mentioned in message/linked to your account! Please use `set` or `link` to link an osu! account to you, or name a user to obtain their recent score of!")
			}
			test, err := OsuAPI.GetUser(osuapi.GetUserOpts{
				Username: username,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "User "+username+" may not exist! Are you sure you replaced spaces with `_`?")
				return
			}
			user = *test
		}
	} else {
		nomod := osuapi.Mods(0)
		beatmap = osutools.BeatmapParse(mapRegex.FindStringSubmatch(m.Embeds[0].URL)[4], "map", &nomod)

		username = m.Embeds[0].Author.Name
		test, err := OsuAPI.GetUser(osuapi.GetUserOpts{
			Username: username,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "User "+username+" may not exist! Are you sure you replaced spaces with `_`?")
			return
		}
		user = *test

		mods = mod2Regex.FindStringSubmatch(m.Embeds[0].Description)[1]
		if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
			mods += "DT"
		}
		parsedMods = osuapi.ParseMods(mods)

		scoreText := strings.Replace(scoreRegex.FindStringSubmatch(m.Embeds[0].Description)[1], ",", "", -1)
		scoreVal, _ = strconv.ParseInt(scoreText, 10, 64)
	}

	// API call
	var score osuapi.Score
	if postType == "recent" {
		scoreOpts := osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Limit:  50,
		}
		scores, err := OsuAPI.GetUserRecent(scoreOpts)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
		if len(scores) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Could not create a scorepost for the score above! Can't create scoreposts for unfinished scores currently.")
			return
		}
		for _, recentScore := range scores {
			if recentScore.Score.Score == scoreVal {
				score = recentScore.Score
				break
			}
		}

	} else if postType == "recentBest" {
		scoreOpts := osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Limit:  100,
		}
		scores, err := OsuAPI.GetUserBest(scoreOpts)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
		if len(scores) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Could not create a scorepost for the score above! Can't create scoreposts for unfinished scores currently.")
			return
		}
		for _, recentScore := range scores {
			if recentScore.Score.Score == scoreVal {
				score = recentScore.Score
				break
			}
		}

	} else {
		scoreOpts := osuapi.GetScoresOpts{
			BeatmapID: beatmap.BeatmapID,
			UserID:    user.UserID,
		}
		if parsedMods != 0 {
			scoreOpts.Mods = &parsedMods
		}
		scores, err := OsuAPI.GetScores(scoreOpts)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
		if len(scores) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Could not create a scorepost for the score above! Can't create scoreposts for unfinished scores currently.")
			return
		}
		score = scores[0].Score
	}

	diffMods := 338 & score.Mods
	if diffMods&256 != 0 && diffMods&64 != 0 { // Remove DTHT
		diffMods -= 320
	}
	if diffMods&2 != 0 && diffMods&16 != 0 { // Remove EZHR
		diffMods -= 18
	}
	beatmap = osutools.BeatmapParse(strconv.Itoa(beatmap.BeatmapID), "map", &diffMods)

	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	acc := strconv.FormatFloat(accCalc, 'f', 2, 64) + "%"

	text := user.Username + " | " +
		beatmap.Artist + " - " + beatmap.Title +
		" [" + beatmap.DiffName + "] +"

	modText := strings.Replace(score.Mods.String(), "DTNC", "NC", -1)
	newModText := ""
	for i := range modText {
		newModText += string(modText[i])
		if i > 0 && (i-1)%2 == 0 && i != len(modText)-1 {
			newModText += ","
		}
	}
	text += newModText +
		" (" + acc + ")" +
		" (" + strings.Replace(beatmap.Creator, "_", `\_`, -1) + " | " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + "★) "

	text = strings.Replace(text, " +NM", "", -1)

	if score.MaxCombo == beatmap.MaxCombo {
		if accCalc == 100.0 {
			text += "SS "
		} else {
			text += "FC "
		}
	} else {
		if score.CountMiss == 0 {
			text += strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x "
		} else {
			text += strconv.Itoa(score.CountMiss) + "m " + strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x "
		}
	}

	leaderboard, _ := OsuAPI.GetScores(osuapi.GetScoresOpts{
		BeatmapID: beatmap.BeatmapID,
		Limit:     50,
	})
	for i, mapScore := range leaderboard {
		if score.UserID == mapScore.UserID && score.Score == mapScore.Score.Score {
			text += "#" + strconv.Itoa(i+1) + " "
			break
		}
	}

	ppValues := make(chan string, 1)
	go osutools.PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), modText, ppValues)
	ppVal, _ := strconv.ParseFloat(<-ppValues, 64)
	if beatmap.Approved == osuapi.StatusLoved {
		text += "LOVED | " + strconv.FormatFloat(ppVal, 'f', 0, 64) + "pp if ranked | "
	} else if beatmap.Approved == osuapi.StatusQualified {
		text += "QUALIFIED | " + strconv.FormatFloat(ppVal, 'f', 0, 64) + "pp if ranked | "
	} else {
		text += "| " + strconv.FormatFloat(ppVal, 'f', 0, 64) + "pp | "
	}

	if score.Mods&256 != 0 || score.Mods&64 != 0 {
		text += "xxxx cv. UR"
	} else {
		text += "xxxx UR"
	}

	s.ChannelMessageSend(m.ChannelID, text)

}