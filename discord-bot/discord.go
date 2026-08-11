package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/taichi765/kokkimusume-wiki-automation/common"
)

// /new
func newCharaModalSlashCommand(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	slog.Debug("received slash command /new")
	// FIXME: 日本のサーバで実行されるとは限らないかも? UTC+9と指定したほうが良いかも知らん
	today := time.Now().Local().Format("2006/01/02")

	areaOptions := make([]discord.StringSelectMenuOption, len(common.ValidAreas))
	for i, a := range common.ValidAreas {
		areaOptions[i] = discord.NewStringSelectMenuOption(a, a)
	}
	fmt.Println(areaOptions)

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

	if err := e.DeferCreateMessage(true); err != nil {
		return err
	}

	if a.client == nil {
		return fmt.Errorf("client was not initialized")
	}

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
		return e.CreateMessage(discord.NewMessageCreate().WithContentf("the value of date-input was invalid: %v", err))
	}

	chara := common.CharacterData{
		Name:                nameInput.Value,
		Area:                areaInput.Values[0],
		FirstAppearenceDate: dateInput.Value,
	}
	go a.handleUpdateCsv(chara, e.Message.ChannelID)
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

func (a *App) handleUpdateCsv(chara common.CharacterData, channelId snowflake.ID) {
	if err := updateCsv(chara, a.envVars.githubAppId, a.envVars.githubInstallationId); err != nil {
		_, err := a.client.Rest.CreateMessage(channelId, discord.NewMessageCreate().WithContentf("failed to update CSV file on GitHub: %v", err))
		if err != nil {
			slog.Error("failed to send message", slog.Any("err", err))
		}
		return
	}
	_, err := a.client.Rest.CreateMessage(channelId, discord.NewMessageCreate().WithContent("CSVファイルの更新に成功しました"))
	if err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}
