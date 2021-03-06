package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Stats explains the stats functionality
func Stats(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: stats / class"
	embed.Description = "`(stats|class) [text] [num]` gives stats for a specific amount of skills, alongside a class randomly chosen from the adjectives and nouns added by server members combined."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[text]",
			Value:  "The text to get the stats for. No text will give your own stats instead.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "The number of skills to print (Default is 4 for stats, 0 for class).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`adjective`, `nouns`, `skills`, `toggle`",
		},
	}
	return embed
}

// Adjectives explains the adjective functionality
func Adjectives(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: adj / adjective / adjectives"
	embed.Description = "`(adj|adjective|adjectives) [add|remove] [adj]` lets you add/remove adjectives from the stats feature, or see the list of current adjectives if no word is given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[add|remove]",
			Value:  "State add / remove to add / remove the adjective (Default: add).",
			Inline: true,
		},
		{
			Name:   "[adj]",
			Value:  "The word to add/remove. No word lets you see the full list.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`stats`, `nouns`, `skills`, `toggle`",
		},
	}
	return embed
}

// Nouns explains the noun functionality
func Nouns(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: noun / nouns"
	embed.Description = "`(noun|nouns) [add|remove] [noun]` lets you add/remove nouns from the stats feature, or see the list of current nouns if no word is given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[add|remove]",
			Value:  "State add / remove to add / remove the noun (Default: add).",
			Inline: true,
		},
		{
			Name:   "[noun]",
			Value:  "The word to add/remove. No word lets you see the full list.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`stats`, `adjective`, `skills`, `toggle`",
		},
	}
	return embed
}

// Skills explains the skills functionality
func Skills(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: skill / skills"
	embed.Description = "`(skill|skills) [add|remove] [skill]` lets you add/remove skills from the stats feature, or see the list of current skills if no word is given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[add|remove]",
			Value:  "State add / remove to add / remove the skill (Default: add).",
			Inline: true,
		},
		{
			Name:   "[skill]",
			Value:  "The word to add/remove. No word lets you see the full list.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`stats`, `adjective`, `nouns`, `toggle`",
		},
	}
	return embed
}
