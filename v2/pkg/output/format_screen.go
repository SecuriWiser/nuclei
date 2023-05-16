package output

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/config"
	"github.com/projectdiscovery/nuclei/v2/internal/firebase"

	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

// formatScreen formats the output for showing on screen.
func (w *StandardWriter) formatScreen(output *ResultEvent) []byte {
	builder := &bytes.Buffer{}
	riskData := make(map[string]any)

	if !w.noMetadata {
		if w.timestamp {
			timestamp := output.Timestamp.Format("2006-01-02 15:04:05")
			builder.WriteRune('[')
			builder.WriteString(w.aurora.Cyan(output.Timestamp.Format("2006-01-02 15:04:05")).String())
			builder.WriteString("] ")
			riskData["timestamp"] = timestamp
		}
		builder.WriteRune('[')
		builder.WriteString(w.aurora.BrightGreen(output.TemplateID).String())
		riskData["templateID"] = output.TemplateID

		if output.MatcherName != "" {
			builder.WriteString(":")
			builder.WriteString(w.aurora.BrightGreen(output.MatcherName).Bold().String())
			riskData["matcherName"] = output.MatcherName
		} else if output.ExtractorName != "" {
			builder.WriteString(":")
			builder.WriteString(w.aurora.BrightGreen(output.ExtractorName).Bold().String())
			riskData["extractorName"] = output.ExtractorName
		}

		if w.matcherStatus {
			builder.WriteString("] [")
			if !output.MatcherStatus {
				builder.WriteString(w.aurora.Red("failed").String())
			} else {
				builder.WriteString(w.aurora.Green("matched").String())
			}
		}

		builder.WriteString("] [")
		builder.WriteString(w.aurora.BrightBlue(output.Type).String())
		builder.WriteString("] ")
		riskData["type"] = output.Type

		builder.WriteString("[")
		builder.WriteString(w.severityColors(output.Info.SeverityHolder.Severity))
		builder.WriteString("] ")
		riskData["severity"] = output.Info.SeverityHolder.Severity
	}
	if output.Matched != "" {
		builder.WriteString(output.Matched)
		riskData["matchedHost"] = output.Matched
	} else {
		builder.WriteString(output.Host)
		riskData["matchedHost"] = output.Host
	}

	// If any extractors, write the results
	if len(output.ExtractedResults) > 0 {
		builder.WriteString(" [")

		for i, item := range output.ExtractedResults {
			builder.WriteString(w.aurora.BrightCyan(item).String())

			if i != len(output.ExtractedResults)-1 {
				builder.WriteRune(',')
			}
		}
		builder.WriteString("]")
		riskData["extractedResults"] = strings.Join(output.ExtractedResults, ",")
	}

	if len(output.Lines) > 0 {
		builder.WriteString(" [LN: ")

		for i, line := range output.Lines {
			builder.WriteString(strconv.Itoa(line))

			if i != len(output.Lines)-1 {
				builder.WriteString(",")
			}
		}
		builder.WriteString("]")
	}

	// Write meta if any
	if len(output.Metadata) > 0 {
		builder.WriteString(" [")

		first := true
		for name, value := range output.Metadata {
			if !first {
				builder.WriteRune(',')
			}
			first = false

			builder.WriteString(w.aurora.BrightYellow(name).String())
			builder.WriteRune('=')
			builder.WriteString(w.aurora.BrightYellow(strconv.QuoteToASCII(types.ToString(value))).String())
		}
		builder.WriteString("]")
	}
	riskData["metadata"] = output.Metadata

	builder.WriteString(" [")
	builder.WriteString(strings.Replace(output.Info.Description, "\n", "", 1))
	builder.WriteRune(']')
	riskData["description"] = output.Info.Description
	riskData["time"] = time.Now().Format("2006-01-02 15:04:05")

	coll := firebase.Client.Collection("scanning_dev").Doc("risk-profiles").Collection(config.RiskID)
	_, _, err := coll.Add(context.Background(), riskData)
	if err != nil {
		gologger.Error().Msgf("Add data for %s to firebase error %s\n", config.RiskID, err.Error())
	}

	return builder.Bytes()
}
