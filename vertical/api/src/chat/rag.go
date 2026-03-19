package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"orbita-api/db"

	"go.uber.org/zap"
)

const systemPrompt = `Sei un assistente AI per Orbita, una piattaforma per calcolare i benefici dell'automazione RPA.
Rispondi sempre in italiano. Sii preciso e conciso.

HAI ACCESSO A TRE FONTI DI DATI:

1. **Processi** (schede di valutazione business): contengono anagrafica (nome, codice, area, proponente, responsabile, dipartimento), parametri operativi (periodicita, tecnologia, sistemi coinvolti, costi, tempi, tassi errore), risultati calcolati (ROI, risparmi, ore risparmiate, impact score), bot collegati (linked_bots) e documento caricato.

2. **Orchestrator UiPath** (dati reali di esecuzione): esecuzioni bot (job con stato Successful/Faulted/Running), code di lavoro (queue items), schedulazioni (trigger con cron expression e prossima esecuzione). Ogni processo puo avere piu bot collegati (es. FL003_DSP, FL003_PRF sono bot dello stesso processo FL003).

3. **Documenti** (PPTX, DOCX, PDF) caricati per ogni processo: contengono la descrizione operativa dettagliata del processo (passi, sistemi, eccezioni, flusso).

REGOLE FONDAMENTALI:

A) PARTI SEMPRE DAI PROCESSI. Ogni risposta deve essere organizzata per processo. I bot collegati (linked_bots) NON sono processi separati, sono componenti del processo padre. Se l'utente chiede una tabella, le righe devono essere i PROCESSI, non i singoli bot o schedulazioni.

B) CODICE PROCESSO: il codice del processo si ricava dal nome dei bot collegati (es. bot "FL003_ControlloGiacenzeFiorano_DSP" -> codice processo "FL003"). Quando mostri tabelle, includi sempre il codice processo.

C) SCHEDULAZIONI: ogni processo puo avere piu bot, ognuno con la sua schedulazione. Quando l'utente chiede "schedulazione" o "orario di partenza", mostra SOLO l'orario/frequenza della schedulazione PRINCIPALE (tipicamente il primo trigger DSP o il trigger con frequenza piu bassa). NON elencare tutti i trigger separatamente a meno che non venga chiesto esplicitamente.

D) DESCRIZIONI: quando l'utente chiede una "descrizione" o "descrizione sintetica" di un processo, usa il contenuto del documento caricato (sezione "Documenti rilevanti" nel contesto) per dare una descrizione operativa in 1-2 frasi. NON ripetere solo il campo "Descrizione" dell'anagrafica. Se il documento non e disponibile, usa la descrizione dall'anagrafica.

E) DOCUMENTI: ogni processo puo avere un documento caricato. Quando l'utente chiede informazioni sui documenti o file, descrivi TUTTI i documenti disponibili, non solo uno.

F) FORMATTAZIONE: usa Markdown. Per tabelle usa il formato: | Col1 | Col2 |\n|------|------|\n| val | val |. Quando l'utente chiede una tabella, scegli le colonne piu pertinenti alla domanda. Tieni le tabelle compatte e leggibili.

G) RISPOSTE DIRETTE: rispondi subito con i dati richiesti. Non chiedere chiarimenti se puoi dedurre cosa serve dal contesto. Se l'utente chiede "tabella delle schedulazioni e descrizione", dai subito una tabella con processo, schedulazione e descrizione sintetica.

H) BOT E PROCESSI: i codici come FL001, FL002, FL003, FL003b, FL006, FL008 sono codici processo. I nomi come "FL001-IWP_DSP", "FL003_ControlloGiacenzeFiorano_DSP" sono nomi bot. Ogni bot appartiene a un processo. Raggruppa sempre per processo.`

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

	// 1. Search document chunks (most relevant to question)
	if azure.embeddingDeployment != "" {
		embedding, err := azure.CreateEmbedding(question)
		if err != nil {
			zap.S().Warnw("Failed to create embedding for RAG", "error", err)
		} else {
			chunkService := ChunkService{}
			chunks, err := chunkService.SearchSimilar(companyID, embedding, 10)
			if err != nil {
				zap.S().Warnw("Failed to search chunks", "error", err)
			} else if len(chunks) > 0 {
				var docParts []string
				for _, chunk := range chunks {
					docParts = append(docParts, fmt.Sprintf("[Documento: %s]\n%s", chunk.FileName, chunk.Content))
				}
				contextParts = append(contextParts, "## Contenuto documenti rilevanti\n"+strings.Join(docParts, "\n\n"))
			}
		}

		// Also get first chunk of ALL documents (for overview/descriptions)
		allFirstChunks := queryAllDocumentFirstChunks(companyID)
		if allFirstChunks != "" {
			contextParts = append(contextParts, "## Riepilogo di tutti i documenti caricati\n"+allFirstChunks)
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

func queryAllDocumentFirstChunks(companyID int) string {
	database := db.DB()
	ctx := context.Background()

	rows, err := database.C.Query(ctx, `SELECT DISTINCT ON (file_name) file_name, content
		FROM chat_document_chunks 
		WHERE company_id = $1 AND chunk_index = 0
		ORDER BY file_name, chunk_index`, companyID)
	if err != nil {
		return ""
	}
	defer rows.Close()

	var parts []string
	for rows.Next() {
		var fileName, content string
		if err := rows.Scan(&fileName, &content); err != nil {
			continue
		}
		summary := content
		if len(summary) > 500 {
			summary = summary[:500] + "..."
		}
		parts = append(parts, fmt.Sprintf("- **%s**: %s", fileName, summary))
	}
	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("Documenti indicizzati (%d):\n%s", len(parts), strings.Join(parts, "\n"))
}

func queryProcessData(companyID int) string {
	database := db.DB()

	type ProcessRow struct {
		ProcessName  string  `db:"process_name"`
		Status       string  `db:"status"`
		Data         string  `db:"data"`
		Results      string  `db:"results"`
		DocumentName *string `db:"document_name"`
	}

	sql := `SELECT 
		process_name,
		status,
		COALESCE(data::text, '{}') as data,
		COALESCE(results::text, '{}') as results,
		document_name
	FROM processes 
	WHERE company_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC 
	LIMIT 50`

	rows, err := database.C.Query(context.Background(), sql, companyID)
	if err != nil {
		zap.S().Warnw("Failed to query processes for RAG", "error", err)
		return ""
	}
	defer rows.Close()

	schedMap := queryScheduleMap(companyID)

	var lines []string
	for rows.Next() {
		var p ProcessRow
		if err := rows.Scan(&p.ProcessName, &p.Status, &p.Data, &p.Results, &p.DocumentName); err != nil {
			continue
		}

		var data map[string]interface{}
		var results map[string]interface{}
		json.Unmarshal([]byte(p.Data), &data)
		json.Unmarshal([]byte(p.Results), &results)

		line := fmt.Sprintf("### %s\n  Stato: %s", p.ProcessName, p.Status)

		// Anagrafica
		addField := func(label, key string) {
			if v, ok := data[key]; ok && v != nil && fmt.Sprintf("%v", v) != "" {
				line += fmt.Sprintf("\n  %s: %v", label, v)
			}
		}
		addField("Area", "area")
		addField("Proponente", "proposer")
		addField("Responsabile", "responsibleManager")
		addField("Dipartimento", "department")
		addField("Tipo processo", "processType")
		addField("Periodicita", "periodicity")
		addField("Sistemi coinvolti", "systemsInvolved")
		addField("Tecnologia", "technology")
		addField("Altra tecnologia", "technologyOther")
		addField("Descrizione", "processDescription")

		// Costi e tempi
		addNumField := func(label, key, unit string) {
			if v, ok := data[key]; ok && v != nil {
				if f, ok := v.(float64); ok && f > 0 {
					line += fmt.Sprintf("\n  %s: %.0f%s", label, f, unit)
				}
			}
		}
		addNumField("Costo implementazione", "implementationCost", "€")
		addNumField("Costo formazione", "trainingCost", "€")
		addNumField("Costo manutenzione annuo", "maintenanceCost", "€/anno")
		addNumField("Costo orario", "hourlyCost", "€/h")
		addNumField("Tempo per attivita", "timePerActivity", " min")
		addNumField("Attivita al giorno", "activitiesPerDay", "")
		addNumField("Giorni lavorativi/anno", "workingDaysPerYear", "")
		addNumField("Tasso errore attuale", "currentErrorRate", "%")
		addNumField("Tasso errore post-RPA", "postErrorRate", "%")
		addNumField("Costo per errore", "errorCost", "€")
		addNumField("Riduzione tempo", "timeReductionFactor", "%")

		// Risultati
		addResField := func(label, key, unit string) {
			if v, ok := results[key]; ok && v != nil {
				if f, ok := v.(float64); ok && f != 0 {
					line += fmt.Sprintf("\n  %s: %.1f%s", label, f, unit)
				}
			}
		}
		addResField("ROI", "roi", "%")
		addResField("Risparmio annuale", "annualSavings", "€")
		addResField("Risparmio operativo", "operationalSavings", "€")
		addResField("Risparmio riduzione errori", "errorReductionSavings", "€")
		addResField("Beneficio produttivita", "productivityBenefit", "€")
		addResField("Ore risparmiate/mese", "hoursSavedMonthly", " h")
		addResField("Ore risparmiate/anno", "hoursSavedAnnually", " h")
		addResField("Impact Score", "impactScore", "")
		if v, ok := results["breakEvenMonths"]; ok && v != nil {
			if f, ok := v.(float64); ok && f > 0 {
				line += fmt.Sprintf("\n  Break-even: %.0f mesi", f)
			}
		}

		// Documento
		if p.DocumentName != nil && *p.DocumentName != "" {
			line += fmt.Sprintf("\n  Documento caricato: %s", *p.DocumentName)
		}

		// Bot collegati e schedulazioni
		linkedBots := fmt.Sprintf("%v", data["linkedBots"])
		if linkedBots != "" && linkedBots != "<nil>" && linkedBots != "null" && linkedBots != "[]" {
			line += fmt.Sprintf("\n  Bot collegati: %s", linkedBots)
			linkedBotsJSON, _ := json.Marshal(data["linkedBots"])
			botSchedules := matchBotSchedules(string(linkedBotsJSON), schedMap)
			if botSchedules != "" {
				line += fmt.Sprintf("\n  Schedulazioni Orchestrator: %s", botSchedules)
			}
		}
		if v, ok := data["botNotes"]; ok && v != nil && fmt.Sprintf("%v", v) != "" {
			line += fmt.Sprintf("\n  Ruolo bot: %v", v)
		}

		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return "Nessun processo trovato per questa azienda."
	}
	return fmt.Sprintf("Processi aziendali (%d):\n%s", len(lines), strings.Join(lines, "\n\n"))
}

type scheduleInfo struct {
	Name        string
	ReleaseName string
	Enabled     bool
	CronSummary string
	NextOcc     string
}

func queryScheduleMap(companyID int) map[string][]scheduleInfo {
	database := db.DB()
	ctx := context.Background()
	result := make(map[string][]scheduleInfo)

	rows, err := database.C.Query(ctx, `SELECT name, COALESCE(release_name, ''), enabled, 
		COALESCE(cron_summary, COALESCE(cron_expression, '')), 
		COALESCE(TO_CHAR(next_occurrence, 'YYYY-MM-DD HH24:MI:SS'), '')
		FROM orchestrator_schedules WHERE company_id = $1`, companyID)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var s scheduleInfo
		if err := rows.Scan(&s.Name, &s.ReleaseName, &s.Enabled, &s.CronSummary, &s.NextOcc); err != nil {
			continue
		}
		key := strings.ToLower(s.ReleaseName)
		result[key] = append(result[key], s)
	}
	return result
}

func matchBotSchedules(linkedBotsJSON string, schedMap map[string][]scheduleInfo) string {
	linkedBotsJSON = strings.TrimSpace(linkedBotsJSON)
	if linkedBotsJSON == "" || linkedBotsJSON == "null" || linkedBotsJSON == "[]" {
		return ""
	}

	var botNames []string
	if strings.HasPrefix(linkedBotsJSON, "[") {
		if err := json.Unmarshal([]byte(linkedBotsJSON), &botNames); err != nil {
			return ""
		}
	} else {
		botNames = []string{linkedBotsJSON}
	}

	var matched []string
	seen := make(map[string]bool)
	for _, bot := range botNames {
		botLower := strings.ToLower(strings.TrimSpace(bot))
		for key, schedules := range schedMap {
			if strings.Contains(key, botLower) || strings.Contains(botLower, key) || key == botLower {
				for _, s := range schedules {
					label := s.Name
					if seen[label] {
						continue
					}
					seen[label] = true
					stato := "attivo"
					if !s.Enabled {
						stato = "disattivato"
					}
					entry := fmt.Sprintf("%s (%s, %s", label, stato, s.CronSummary)
					if s.NextOcc != "" {
						entry += fmt.Sprintf(", prossima: %s", s.NextOcc)
					}
					entry += ")"
					matched = append(matched, entry)
				}
			}
		}
	}
	if len(matched) == 0 {
		return ""
	}
	return strings.Join(matched, "; ")
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
