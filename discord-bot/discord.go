package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/taichi765/kokkimusume-wiki-automation/common"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "new",
		Description: "append new character to CSV file on github",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleJapanese: "新しい国旗娘をGitHub上のCSVファイルに追加する",
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	},
	discord.SlashCommandCreate{
		Name:        "version",
		Description: "show bot's version",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleJapanese: "ボットのバージョンを表示する",
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	},
	discord.SlashCommandCreate{
		Name:        "help",
		Description: "print help",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleJapanese: "ボットの使い方を表示する",
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	},
}

type commandInfo struct {
	name string
	desc string
}

func helpSlashCommand(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	infos := listAllCommands()

	b := &strings.Builder{}
	for _, info := range infos {
		fmt.Fprintf(b, "%s: %s\n", info.name, info.desc)
	}

	return e.CreateMessage(discord.NewMessageCreate().WithContent(b.String()))
}

func listAllCommands() []commandInfo {
	infos := make([]commandInfo, len(commands))
	for i, cmd := range commands {
		switch c := cmd.(type) {
		case discord.SlashCommandCreate:
			infos[i] = commandInfo{
				name: c.Name,
				desc: c.Description,
			}
		default:
			panic("this type of command is not supported")
		}
	}
	return infos
}

func versionSlashCommand(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.NewMessageCreate().WithContentf("version: %s\ncommit hash: %s", version, commitHash))
}

// /new
func newCharaModalSlashCommand(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	slog.Debug("received slash command /new")
	// FIXME: 日本のサーバで実行されるとは限らないかも? UTC+9と指定したほうが良いかも知らん
	today := time.Now().Local().Format("2006/01/02")

	areaOptions := make([]discord.StringSelectMenuOption, len(common.ValidAreas))
	for i, a := range common.ValidAreas {
		areaOptions[i] = discord.NewStringSelectMenuOption(a, a)
	}

	return e.Modal(
		discord.NewModalCreate("/modals/new", "新しい国旗娘を追加",
			discord.NewLabel("名前",
				discord.NewTextInput("/modals/new/name-input", discord.TextInputStyleShort).WithRequired(true),
			),
			discord.NewLabel("地域",
				discord.NewStringSelectMenu("/modals/new/area-input", "placeholder",
					areaOptions...,
				).WithRequired(true),
			),
			discord.NewLabel("初出",
				discord.NewTextInput("/modals/new/date-input", discord.TextInputStyleShort).WithRequired(true).WithValue(today),
			),
		),
	)
}

// /modals/new
func (a *App) onNewCharaModalSubmitted(e *handler.ModalEvent) error {
	slog.Debug("received modal event for /modals/new")

	if err := e.Respond(
		discord.InteractionResponseTypeDeferredCreateMessage,
		discord.NewMessageCreate().
			WithEphemeral(true).
			WithContent("waiting for inputs..."),
	); err != nil {
		return err
	}
	/*if err := e.DeferCreateMessage(false); err != nil {
	return err
	}*/

	nameInput, ok := e.Data.TextInput("/modals/new/name-input")
	if !ok {
		return fmt.Errorf("can't get name-input data from modal event")
	}
	areaInput, ok := e.Data.StringSelectMenu("/modals/new/area-input")
	if !ok {
		return fmt.Errorf("can't get area-input data from modal event")
	}
	if len(areaInput.Values) != 1 {
		return fmt.Errorf("number of the values of area-input should be exactly one")
	}
	dateInput, ok := e.Data.TextInput("/modals/new/date-input")
	if !ok {
		return fmt.Errorf("can't get date-input data from modal event")
	}

	if err := validateDateInput(dateInput.Value); err != nil {
		_, err := e.CreateFollowupMessage(
			discord.NewMessageCreate().WithContentf("the value of date-input was invalid: %v", err),
		)
		return err
	}

	chara := common.CharacterData{
		Name:                nameInput.Value,
		Area:                areaInput.Values[0],
		FirstAppearenceDate: dateInput.Value,
	}
	go a.handleUpdateCsv(chara, e.CreateFollowupMessage)
	return nil
}

// 形式が不正または日付が明日以降のときエラーを返す
func validateDateInput(val string) error {
	inputDate, err := time.Parse("2006/01/02", val)
	if err != nil {
		return err
	}
	inputYear, inputMonth, inputDay := inputDate.Date()
	nowYear, nowMonth, nowDay := time.Now().Local().Date()
	if inputYear > nowYear {
		return fmt.Errorf("input date was future")
	}
	if inputYear < nowYear {
		return nil
	}
	if inputMonth > nowMonth {
		return fmt.Errorf("input date was future")
	}
	if inputMonth < nowMonth {
		return nil
	}
	if inputDay > nowDay {
		return fmt.Errorf("input date was future")
	}
	return nil
}

func (a *App) handleUpdateCsv(chara common.CharacterData, createFollowupMessage func(discord.MessageCreate, ...rest.RequestOpt) (*discord.Message, error)) {
	client, err := newClient(a.envVars.githubAppId, a.envVars.githubInstallationId, a.envVars.githubPrivateKey)
	if err != nil {
		createFollowupMessage(discord.NewMessageCreate().WithContentf("failed to build github client: %v", err))
	}

	if err := updateCsv(chara, client); err != nil {
		if _, err := createFollowupMessage(discord.NewMessageCreate().WithContentf("failed to update CSV file on GitHub: %v", err)); err != nil {
			slog.Error("failed to send message", slog.Any("err", err))
		}
		return
	}

	if _, err := createFollowupMessage(
		discord.NewMessageCreate().WithContentf("%sをCSVファイルに追加しました", chara.Name),
	); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}
