package utils

import (
	"fmt"
	"strings"

	"github.com/istvzsig/eve-gate-watch/internal/dscan"
	"github.com/istvzsig/eve-gate-watch/internal/pvp"
)

func FormatISK(value float64) string {
	switch {
	case value >= 1_000_000_000:
		return fmt.Sprintf("%.1fB ISK", value/1_000_000_000)
	case value >= 1_000_000:
		return fmt.Sprintf("%.1fM ISK", value/1_000_000)
	case value >= 1_000:
		return fmt.Sprintf("%.1fK ISK", value/1_000)
	default:
		return fmt.Sprintf("%.0f ISK", value)
	}
}

func FormatDScan(result dscan.DScanResult) string {
	return fmt.Sprintf(
		"```text\n"+
			"%s SYSTEM INTELLIGENCE\n"+
			"--------------------------------\n"+
			"Region             %s\n"+
			"Constellation      %s\n"+
			"Security           %.1f %s\n"+
			"Planets            %d\n"+
			"Stargates          %d\n"+
			"Stations           %d\n"+
			"\n"+
			"LAST HOUR\n"+
			"--------------------------------\n"+
			"Ship jumps         %d\n"+
			"Ship kills         %d\n"+
			"Pod kills          %d\n"+
			"NPC kills          %d\n"+
			"```",
		strings.ToUpper(result.SystemName),
		result.RegionName,
		result.ConstellationName,
		result.SecurityStatus,
		securityClass(result.SecurityStatus),
		result.Planets,
		result.Stargates,
		result.Stations,
		result.ShipJumps,
		result.ShipKills,
		result.PodKills,
		result.NPCKills,
	)
}

func FormatPVP(result pvp.Result) string {
	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"```text\n"+
			"%s PVP ACTIVITY\n"+
			"--------------------------------\n"+
			"%-12s %-18s %-12s %s\n"+
			"--------------------------------\n",
		strings.ToUpper(result.SystemName),
		"ID",
		"SHIP",
		"VALUE",
		"OCCURRED",
	)

	if len(result.Kills) == 0 {
		builder.WriteString("No recent PvP activity.\n")
	} else {
		for _, kill := range result.Kills {
			fmt.Fprintf(
				&builder,
				"%-12d %-18s %-12s %s\n",
				kill.ID,
				kill.ShipName,
				FormatISK(kill.Value),
				kill.OccurredAt.Local().Format("15:04"),
			)
		}
	}

	builder.WriteString("```")

	return builder.String()
}

func securityClass(value float64) string {
	switch {
	case value >= 0.5:
		return "HIGH SEC"
	case value > 0.0:
		return "LOW SEC"
	default:
		return "NULL SEC"
	}
}
