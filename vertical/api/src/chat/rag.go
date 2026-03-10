package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"benefit-calculator-api/db"

	"go.uber.org/zap"
)

const systemPrompt = `Sei un assistente AI per il Benefit Calculator, una piattaforma per calcolare i benefici dell'automazione RPA.
Rispondi sempre in italiano. Sii preciso e conciso.

Hai accesso a DUE fonti di dati DISTINTE:

1. **Proposte di automazione** (tabella processes): sono le schede di valutazione business dei processi, con ROI, risparmi, costi, stato di avanzamento. Sono le proposte inserite manualmente dagli utenti.

2. **Esecuzioni bot Orchestrator** (tabelle orchestrator_job_executions, orchestrator_queue_items, orchestrator_schedules): sono i dati REALI delle esecuzioni dei robot UiPath. Contengono job eseguiti, stati (Successful/Faulted/Running), code di lavoro, errori, tempi di esecuzione, e le schedulazioni (trigger) configurate con cron expression e prossima esecuzione. Questi dati vengono sincronizzati dall'Orchestrator UiPath.

IMPORTANTE: Ogni proposta (processo) può avere un campo linked_bots con i nomi dei bot Orchestrator collegati. Un singolo processo/assessment può generare MOLTI bot con codici diversi (es. un processo può avere bot come ProcA_DSP, ProcA_PRF, ProcB_DSP, ecc.). Quando un bot appare nel campo linked_bots di un processo, NON trattarlo come un processo separato: fa parte di quel processo. Se l'utente chiede "di che parla FL002", cerca PRIMA se FL002 appare come linked_bot di qualche processo e rispondi che fa parte di quel processo. Usa la tabella orchestrator_process_queue_map per sapere quale bot è collegato a quale coda. Quando l'utente chiede di un "processo che ha girato" o "ultimo bot eseguito" o "errori", cerca SEMPRE nei dati Orchestrator (esecuzioni/code), NON nelle proposte. Quando chiede di ROI, costi, risparmi, valutazione, cerca nelle proposte. Quando l'utente chiede di schedulazioni, trigger, cron, "quando gira", "prossima esecuzione", cerca nei dati schedule dell'Orchestrator.

3. **Documenti** (file PDF, PowerPoint e Word) caricati dall'azienda per ogni processo.

Quando citi dati, indica i numeri esatti. Quando citi documenti, menziona il nome del file.
Se non hai informazioni sufficienti per rispondere, dillo chiaramente.`

const dbSchemaPrompt = `Schema del database disponibile:
- processes: id, process_name, status ('To Valuate','Analysis','Ongoing','Production'), created_at, data (JSONB con: area, proposer, responsible_manager, department, process_type, periodicity, technology, systems_involved, implementation_cost, training_cost, maintenance_cost, hourly_cost, time_per_activity, activities_per_day, working_days_per_year, current_error_rate, post_error_rate, error_cost, productivity_factor, time_reduction_factor, data_quality_score, audit_score, customer_experience_score, error_reduction_score, standardization_score, scalability_score, linked_bots), results (JSONB con: operational_savings, error_reduction_savings, productivity_benefit, annual_savings, roi, break_even_months, hours_saved_monthly, hours_saved_annually, impact_score). NOTA: linked_bots e un array di nomi bot Orchestrator collegati a questo processo.
- orgs_companies: id, name, parent_id
- orgs_areas: id, company_id, name
- orchestrator_job_executions: id, company_id, connector_id, process_name, state ('Successful','Faulted','Stopped','Running','Pending'), source_type, host_machine, start_time, end_time, info, folder_name
- orchestrator_queue_items: id, company_id, connector_id, queue_name, status ('Successful','Failed','New','InProgress','Retried','Abandoned'), priority, reference, processing_exception_type, error_message, start_processing, end_processing, retry_number
- orchestrator_queue_definitions: id, company_id, connector_id, name, max_retries
- orchestrator_process_queue_map: id, company_id, connector_id, process_name, queue_name, auto_detected (collegamento tra bot e code)
- orchestrator_schedules: id, company_id, connector_id, external_schedule_id, name, enabled (bool), release_name (nome processo UiPath), package_name, cron_expression (es. "0 0/5 * 1/1 * ? *"), cron_summary (es. "Every 5 minutes"), next_occurrence (timestamp prossima esecuzione), timezone_id, timezone_iana, start_strategy, folder_name, input_arguments (JSONB)`

func BuildRAGContext(azure *AzureClient, companyID int, question string) string {
	var contextParts []string

	// 1. Search document chunks
	if azure.embeddingDeployment != "" {
		embedding, err := azure.CreateEmbedding(question)
		if err != nil {
			zap.S().Warnw("Failed to create embedding for RAG", "error", err)
		} else {
			chunkService := ChunkService{}
			chunks, err := chunkService.SearchSimilar(companyID, embedding, 5)
			if err != nil {
				zap.S().Warnw("Failed to search chunks", "error", err)
			} else if len(chunks) > 0 {
				var docParts []string
				for _, chunk := range chunks {
					docParts = append(docParts, fmt.Sprintf("[Documento: %s]\n%s", chunk.FileName, chunk.Content))
				}
				contextParts = append(contextParts, "## Documenti rilevanti\n"+strings.Join(docParts, "\n\n"))
			}
		}
	}

	// 2. Query DB for relevant process data
	dbContext := queryProcessData(companyID)
	if dbContext != "" {
		contextParts = append(contextParts, "## Dati processi\n"+dbContext)
	}

	// 3. Query orchestrator data
	orchContext := queryOrchestratorData(companyID)
	if orchContext != "" {
		contextParts = append(contextParts, "## Dati Orchestrator\n"+orchContext)
	}

	if len(contextParts) == 0 {
		return ""
	}
	return strings.Join(contextParts, "\n\n---\n\n")
}

func queryProcessData(companyID int) string {
	db := db.DB()

	type ProcessSummary struct {
		ProcessName string  `db:"process_name"`
		Status      string  `db:"status"`
		Area        string  `db:"area"`
		Technology  string  `db:"technology"`
		Description string  `db:"description"`
		LinkedBots  string  `db:"linked_bots"`
		BotNotes    string  `db:"bot_notes"`
		ROI         float64 `db:"roi"`
		Savings     float64 `db:"savings"`
	}

	sql := `SELECT 
		process_name,
		status,
		COALESCE(data->>'area', '') as area,
		COALESCE(data->>'technology', '') as technology,
		COALESCE(data->>'processDescription', '') as description,
		COALESCE(data->>'linkedBots', '') as linked_bots,
		COALESCE(data->>'botNotes', '') as bot_notes,
		COALESCE((results->>'roi')::numeric, 0) as roi,
		COALESCE((results->>'annual_savings')::numeric, 0) as savings
	FROM processes 
	WHERE company_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC 
	LIMIT 50`

	rows, err := db.C.Query(context.Background(), sql, companyID)
	if err != nil {
		zap.S().Warnw("Failed to query processes for RAG", "error", err)
		return ""
	}
	defer rows.Close()

	var lines []string
	for rows.Next() {
		var p ProcessSummary
		if err := rows.Scan(&p.ProcessName, &p.Status, &p.Area, &p.Technology, &p.Description, &p.LinkedBots, &p.BotNotes, &p.ROI, &p.Savings); err != nil {
			continue
		}
		line := fmt.Sprintf("- %s | Area: %s | Stato: %s | Tecnologia: %s | ROI: %.1f%% | Risparmio: €%.0f/anno",
			p.ProcessName, p.Area, p.Status, p.Technology, p.ROI, p.Savings)
		if p.Description != "" {
			line += fmt.Sprintf("\n  Descrizione: %s", p.Description)
		}
		if p.LinkedBots != "" && p.LinkedBots != "null" && p.LinkedBots != "[]" {
			line += fmt.Sprintf("\n  Bot collegati: %s", p.LinkedBots)
		}
		if p.BotNotes != "" {
			line += fmt.Sprintf("\n  Ruolo bot: %s", p.BotNotes)
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return "Nessun processo trovato per questa azienda."
	}
	return fmt.Sprintf("Processi aziendali (%d):\n%s", len(lines), strings.Join(lines, "\n"))
}

func queryOrchestratorData(companyID int) string {
	database := db.DB()
	ctx := context.Background()
	var parts []string

	// Job stats
	type JobStats struct {
		Total      int `db:"total"`
		Successful int `db:"successful"`
		Faulted    int `db:"faulted"`
		Stopped    int `db:"stopped"`
		Running    int `db:"running"`
	}
	var js JobStats
	err := database.C.QueryRow(ctx, `SELECT 
		COUNT(*) as total,
		COUNT(*) FILTER (WHERE state = 'Successful') as successful,
		COUNT(*) FILTER (WHERE state = 'Faulted') as faulted,
		COUNT(*) FILTER (WHERE state = 'Stopped') as stopped,
		COUNT(*) FILTER (WHERE state = 'Running') as running
		FROM orchestrator_job_executions WHERE company_id = $1`, companyID).
		Scan(&js.Total, &js.Successful, &js.Faulted, &js.Stopped, &js.Running)
	if err == nil && js.Total > 0 {
		parts = append(parts, fmt.Sprintf("Esecuzioni bot: %d totali, %d completate, %d in errore, %d fermate, %d in esecuzione",
			js.Total, js.Successful, js.Faulted, js.Stopped, js.Running))
	}

	// Recent jobs (any state) with timestamps
	recentRows, err := database.C.Query(ctx, `SELECT process_name, state, host_machine, start_time, end_time
		FROM orchestrator_job_executions 
		WHERE company_id = $1 
		ORDER BY start_time DESC LIMIT 10`, companyID)
	if err == nil {
		defer recentRows.Close()
		var recentLines []string
		for recentRows.Next() {
			var name, state string
			var machine *string
			var startTime, endTime *time.Time
			if err := recentRows.Scan(&name, &state, &machine, &startTime, &endTime); err == nil {
				st, et := "-", "-"
				if startTime != nil {
					st = startTime.Format("2006-01-02 15:04:05")
				}
				if endTime != nil {
					et = endTime.Format("2006-01-02 15:04:05")
				}
				m := "-"
				if machine != nil {
					m = *machine
				}
				recentLines = append(recentLines, fmt.Sprintf("  - %s | Stato: %s | Macchina: %s | Inizio: %s | Fine: %s", name, state, m, st, et))
			}
		}
		if len(recentLines) > 0 {
			parts = append(parts, "Ultime 10 esecuzioni bot (ordine cronologico decrescente):\n"+strings.Join(recentLines, "\n"))
		}
	}

	// Recent failed jobs
	failedRows, err := database.C.Query(ctx, `SELECT process_name, info, start_time 
		FROM orchestrator_job_executions 
		WHERE company_id = $1 AND state = 'Faulted' 
		ORDER BY start_time DESC LIMIT 10`, companyID)
	if err == nil {
		defer failedRows.Close()
		var failedLines []string
		for failedRows.Next() {
			var name string
			var info *string
			var startTime *time.Time
			if err := failedRows.Scan(&name, &info, &startTime); err == nil {
				ts := "-"
				if startTime != nil {
					ts = startTime.Format("2006-01-02 15:04:05")
				}
				infoStr := ""
				if info != nil {
					infoStr = *info
				}
				failedLines = append(failedLines, fmt.Sprintf("  - %s | %s | Info: %s", name, ts, infoStr))
			}
		}
		if len(failedLines) > 0 {
			parts = append(parts, "Ultimi job falliti:\n"+strings.Join(failedLines, "\n"))
		}
	}

	// Queue stats
	type QueueStats struct {
		Total      int `db:"total"`
		Successful int `db:"successful"`
		Failed     int `db:"failed"`
	}
	var qs QueueStats
	err = database.C.QueryRow(ctx, `SELECT 
		COUNT(*) as total,
		COUNT(*) FILTER (WHERE status = 'Successful') as successful,
		COUNT(*) FILTER (WHERE status = 'Failed') as failed
		FROM orchestrator_queue_items WHERE company_id = $1`, companyID).
		Scan(&qs.Total, &qs.Successful, &qs.Failed)
	if err == nil && qs.Total > 0 {
		parts = append(parts, fmt.Sprintf("Elementi code: %d totali, %d completati, %d falliti",
			qs.Total, qs.Successful, qs.Failed))
	}

	// Queue breakdown by name
	queueRows, err := database.C.Query(ctx, `SELECT queue_name, COUNT(*) as cnt,
		COUNT(*) FILTER (WHERE status = 'Successful') as ok,
		COUNT(*) FILTER (WHERE status = 'Failed') as ko
		FROM orchestrator_queue_items WHERE company_id = $1
		GROUP BY queue_name ORDER BY cnt DESC`, companyID)
	if err == nil {
		defer queueRows.Close()
		var queueLines []string
		for queueRows.Next() {
			var name string
			var cnt, ok, ko int
			if err := queueRows.Scan(&name, &cnt, &ok, &ko); err == nil {
				queueLines = append(queueLines, fmt.Sprintf("  - %s: %d totali, %d ok, %d falliti", name, cnt, ok, ko))
			}
		}
		if len(queueLines) > 0 {
			parts = append(parts, "Dettaglio per coda:\n"+strings.Join(queueLines, "\n"))
		}
	}

	// Schedule stats
	type ScheduleStats struct {
		Total    int `db:"total"`
		Enabled  int `db:"enabled"`
		Disabled int `db:"disabled"`
	}
	var ss ScheduleStats
	err = database.C.QueryRow(ctx, `SELECT 
		COUNT(*) as total,
		COUNT(*) FILTER (WHERE enabled = true) as enabled,
		COUNT(*) FILTER (WHERE enabled = false) as disabled
		FROM orchestrator_schedules WHERE company_id = $1`, companyID).
		Scan(&ss.Total, &ss.Enabled, &ss.Disabled)
	if err == nil && ss.Total > 0 {
		parts = append(parts, fmt.Sprintf("Schedulazioni: %d totali, %d attive, %d disattivate",
			ss.Total, ss.Enabled, ss.Disabled))
	}

	// Schedule details
	schedRows, err := database.C.Query(ctx, `SELECT name, release_name, enabled, cron_summary, cron_expression, next_occurrence, timezone_iana, folder_name
		FROM orchestrator_schedules 
		WHERE company_id = $1 
		ORDER BY enabled DESC, name ASC`, companyID)
	if err == nil {
		defer schedRows.Close()
		var schedLines []string
		for schedRows.Next() {
			var name string
			var releaseName, cronSummary, cronExpr, tzIana, folderName *string
			var enabled bool
			var nextOcc *time.Time
			if err := schedRows.Scan(&name, &releaseName, &enabled, &cronSummary, &cronExpr, &nextOcc, &tzIana, &folderName); err == nil {
				stato := "attivo"
				if !enabled {
					stato = "disattivato"
				}
				proc := "-"
				if releaseName != nil {
					proc = *releaseName
				}
				cron := "-"
				if cronSummary != nil {
					cron = *cronSummary
				} else if cronExpr != nil {
					cron = *cronExpr
				}
				next := "-"
				if nextOcc != nil {
					next = nextOcc.Format("2006-01-02 15:04:05")
				}
				tz := ""
				if tzIana != nil {
					tz = " (" + *tzIana + ")"
				}
				folder := ""
				if folderName != nil {
					folder = " | Cartella: " + *folderName
				}
				schedLines = append(schedLines, fmt.Sprintf("  - %s | Processo: %s | %s | Frequenza: %s | Prossima: %s%s%s", name, proc, stato, cron, next, tz, folder))
			}
		}
		if len(schedLines) > 0 {
			parts = append(parts, "Dettaglio schedulazioni:\n"+strings.Join(schedLines, "\n"))
		}
	}

	// Recent queue errors
	errRows, err := database.C.Query(ctx, `SELECT queue_name, processing_exception_type, error_message, start_processing
		FROM orchestrator_queue_items 
		WHERE company_id = $1 AND status = 'Failed' AND error_message IS NOT NULL
		ORDER BY start_processing DESC LIMIT 10`, companyID)
	if err == nil {
		defer errRows.Close()
		var errLines []string
		for errRows.Next() {
			var qName string
			var excType, errMsg, startProc *string
			if err := errRows.Scan(&qName, &excType, &errMsg, &startProc); err == nil {
				et := ""
				if excType != nil {
					et = *excType
				}
				em := ""
				if errMsg != nil {
					em = *errMsg
				}
				errLines = append(errLines, fmt.Sprintf("  - %s | %s | %s", qName, et, em))
			}
		}
		if len(errLines) > 0 {
			parts = append(parts, "Ultimi errori code:\n"+strings.Join(errLines, "\n"))
		}
	}

	// Process-Queue mappings
	mapRows, err := database.C.Query(ctx, `SELECT process_name, queue_name, auto_detected
		FROM orchestrator_process_queue_map WHERE company_id = $1 ORDER BY process_name`, companyID)
	if err == nil {
		defer mapRows.Close()
		var mapLines []string
		for mapRows.Next() {
			var proc, queue string
			var auto bool
			if err := mapRows.Scan(&proc, &queue, &auto); err == nil {
				src := "manuale"
				if auto {
					src = "auto"
				}
				mapLines = append(mapLines, fmt.Sprintf("  - %s → %s (%s)", proc, queue, src))
			}
		}
		if len(mapLines) > 0 {
			parts = append(parts, "Mapping processo-coda:\n"+strings.Join(mapLines, "\n"))
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "\n\n")
}

func BuildMessages(azure *AzureClient, companyID int, conversationHistory []Message, userMessage string) []chatMessage {
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "system", Content: dbSchemaPrompt},
	}

	// RAG context
	ragContext := BuildRAGContext(azure, companyID, userMessage)
	if ragContext != "" {
		messages = append(messages, chatMessage{
			Role:    "system",
			Content: "Contesto recuperato per rispondere alla domanda:\n\n" + ragContext,
		})
	}

	// Conversation history (last 20 messages)
	for _, msg := range conversationHistory {
		messages = append(messages, chatMessage{Role: msg.Role, Content: msg.Content})
	}

	// Current user message
	messages = append(messages, chatMessage{Role: "user", Content: userMessage})

	return messages
}

func GenerateTitle(azure *AzureClient, userMessage string) string {
	messages := []chatMessage{
		{Role: "system", Content: "Genera un titolo breve (max 6 parole) in italiano per questa conversazione. Rispondi SOLO con il titolo, senza virgolette."},
		{Role: "user", Content: userMessage},
	}
	title, err := azure.ChatCompletion(messages)
	if err != nil {
		zap.S().Warnw("Failed to generate title", "error", err)
		if len(userMessage) > 50 {
			return userMessage[:50] + "..."
		}
		return userMessage
	}
	return strings.TrimSpace(title)
}
