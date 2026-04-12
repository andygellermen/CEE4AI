roadmap-v2.md
Leitgedanke dieser Version
Die Roadmap wird um drei wesentliche Entwicklungsachsen erweitert:
	1	flexiblere Fragemodelle statt nur Single-Choice
	2	Test- und Review-Fähigkeit zur Qualitätssicherung
	3	Coaching-Modul für Angst-, Druck- und Belastungs-Navigation als spätere Ausbaustufe
Damit wächst CEE4AI von einem reinen Denkprofil-MVP zu einer mehrdomänigen, adaptiven Engine.

Gesamtphasen im Überblick
	0	Strategische Klärung
	1	Fragenmodell und Datenfundament V2
	2	Review- und Testmodus
	3	Technischer MVP V2
	4	Profil- und Ergebnislogik
	5	Sanfte Gamification
	6	Adaptive Evolution Engine
	7	Coaching-Modul: Angst-/Druck-Analyzer
	8	Profi- und Coaching-Fähigkeit
	9	Plattform, Skalierung und API

Phase 0 – Strategische Klärung
Ziel
Den Produktkern weiter schärfen und die neue Mehrdomänen-Logik sauber verankern.
Zusätzliche Inhalte in V2
	•	Trennung von Engine und Domänenmodulen
	•	Klarheit über Fragearten
	•	Klarheit über Review-Status und Testbetrieb
	•	Einordnung des Coaching-Moduls als fachlich eigene Ausbaustufe
Ergebnis
	•	klares Verständnis: CEE4AI ist die Engine und zugleich das erste Kernmodul
	•	Perspektive für weitere Module ohne Produktverwässerung

Phase 1 – Fragenmodell und Datenfundament V2
Ziel
Das Datenmodell so erweitern, dass verschiedene Fragetypen, Domänen und sensible Inhalte abbildbar werden.
Neue Schwerpunkte
	•	Einführung von domain
	•	Einführung von question_type
	•	Unterstützung von single_choice, multiple_select, scale, reflection, optional später open_text
	•	Einführung von scoring_mode
	•	Hierarchien über parent_id
	•	Markierung sensibler Inhalte über is_sensitive und age_gate
	•	thematische Triggerlogik über trigger_group
	•	optionale körperbezogene Verweise über body_reference
	•	Review-Fähigkeit über review_status
Zielbild
Ein Fragenmodell, das sowohl Denkprofil-Fragen als auch Coaching-Navigationsfragen tragen kann.
Offen
	•	genaue Typdefinitionen finalisieren
	•	CSV- und Sheet-Strukturen anpassen
	•	Review-Werte standardisieren

Phase 2 – Review- und Testmodus
Ziel
Fragen mit echten Menschen prüfen und systematisch verbessern.
Inhalte
	•	strukturierte Markierungen pro Frage
	•	optionale Kommentarfelder
	•	Testerrollen
	•	Statuswechsel wie draft, active, review, flagged, disabled, archived
	•	einfache Moderations- und Review-Queue
Beispielhafte Test-Flags
	•	Frage unklar
	•	mehrere Antworten wirken richtig
	•	keine Antwort passt
	•	zu leicht
	•	zu schwer
	•	sprachlich unklar
	•	fachlich problematisch
Ergebnis
Ein belastbarer Validierungsmodus, bevor Inhalte breit ausgerollt werden.
Warum diese Phase wichtig ist
Ohne diesen Schritt riskierst du früh strukturelle Qualitätsfehler im gesamten Fragenpool.

Phase 3 – Technischer MVP V2
Ziel
Die erste lauffähige Version auf die neue Flexibilität vorbereiten.
Erweiterungen gegenüber MVP V1
	•	Fragen sind einer domain zugeordnet
	•	Antwortlogik berücksichtigt question_type
	•	Bewertung berücksichtigt scoring_mode
	•	Review-Felder sind technisch vorhanden
	•	Hierarchische Fragelogik ist vorbereitet
Zusätzlich sinnvoll
	•	Rollenmodell vorbereiten: Nutzer, Tester, Coach, Admin
	•	API bewusst so benennen, dass spätere Domänen erweiterbar bleiben
Ergebnis
Ein technischer Kern, der nicht nur ein Quiz ausliefert, sondern als Engine tragfähig bleibt.

Phase 4 – Profil- und Ergebnislogik
Ziel
Flexible Ergebnislogiken je nach Domäne aufbauen.
Bereich A – Denkprofil
	•	Denktypen-Scores
	•	Interessenfelder
	•	Entwicklungshinweise
Bereich B – Testmodus
	•	Qualitätsmetriken pro Frage
	•	Häufung von Flags
	•	Review-Prioritäten
Bereich C – Coaching-Navigation
	•	Trigger-Landkarten
	•	Belastungscluster
	•	Körperreaktions-Mapping
	•	vorsichtige Reflexionshinweise
Leitplanke
Keine Vermischung der Ergebnisarten.Denkprofil ist etwas anderes als Review-Auswertung, und beides ist etwas anderes als eine Coaching-orientierte Selbstklärung.

Phase 5 – Sanfte Gamification
Ziel
Motivation aufbauen, ohne Tiefe zu verlieren.
Erweiterung in V2
Auch der Testmodus und spätere Coaching-Pfade können von sanfter UX profitieren, aber jeweils anders:
	•	Profilmodus: Fortschritt, Missionen, Denkpfade
	•	Testmodus: Feedbackfreundlichkeit, Klarheit, Review-Leichtigkeit
	•	Coaching-Modul: Sicherheit, Ruhe, Orientierung, kleine Schritte
Vermeidungsstrategie
Das Coaching-Modul darf niemals durch aggressive Spielmechanik trivialisiert werden.

Phase 6 – Adaptive Evolution Engine
Ziel
Die Fragesteuerung intelligent und domänensensibel ausbauen.
Neue Anforderungen in V2
	•	Pfadwahl anhand domain
	•	Hierarchische Navigation über parent_id
	•	Triggersteuerung über trigger_group
	•	situationsabhängige Folgefragen
	•	unterschiedliche Bewertungslogiken je Fragetyp
Ergebnis
Die Engine kann nicht nur Schwierigkeit anpassen, sondern auch thematische Pfade dynamisch steuern.

Phase 7 – Coaching-Modul: Angst-/Druck-Analyzer
Ziel
Ein eigenständiges Coaching-Modul auf Basis derselben Engine aufbauen.
Sinnvolle Heimat
Dieses Modul gehört in die Ausbaustufe Coaching und kann dort als zusätzliche, adaptive Unterstützung eingesetzt werden.
Typischer Ablauf
	1	Abfrage grober Themenräume
	2	Verzweigung in Subkategorien
	3	Zuordnung möglicher Körperreaktionen
	4	Exploration belastender Gedankenmuster
	5	Identifikation von Triggern und Ressourcen
	6	Formulierung vorsichtiger Reflexions- und Entlastungsimpulse
Potenzielle Anwendungsfelder
	•	Coaching
	•	Workshops zur Angst- und Belastungsklärung
	•	strukturierte Selbsthilfe
	•	Vororientierung vor einem fachlichen Gespräch
Wichtige Leitplanken
	•	keine medizinische oder psychotherapeutische Diagnose
	•	keine definitive Ursachenzuschreibung
	•	sensible Sprache
	•	Schutzmechanismen bei potenziell belastenden Inhalten
Ergebnis
Ein starkes, verantwortungsvoll positioniertes Zusatzmodul mit hoher praktischer Relevanz.

Phase 8 – Profi- und Coaching-Fähigkeit
Ziel
Die Engine für professionelle Nutzerkontexte erweitern.
Neue V2-Schwerpunkte
	•	getrennte Nutzerrollen
	•	Coach-Ansichten
	•	Review-Ansichten
	•	Klienten- oder Teilnehmerbezug
	•	domänenspezifische Reports
Beispielhafte Rollen
	•	Proband/Nutzer
	•	Tester
	•	Coach
	•	Admin
	•	später Profi-Anwender
Ergebnis
CEE4AI wird für professionelle Nutzung strukturell anschlussfähig.

Phase 9 – Plattform, Skalierung und API
Ziel
Die Engine produktreif und integrationsfähig machen.
Erweiterungen
	•	modulare Domänenverwaltung
	•	APIs für unterschiedliche Frontends oder Partner
	•	White-Label-Fähigkeit
	•	Multi-Tenant-Strukturen
	•	feinere Sensitivitäts- und Freigabelogik
Ergebnis
Eine Plattform, die unterschiedliche Anwendungsfelder sauber unter einem methodischen Dach tragen kann.

Neue Querschnittsthemen in V2
A. Sensible Inhalte
Sensible Domänen brauchen:
	•	is_sensitive
	•	age_gate
	•	klare Sprachregeln
	•	domänenspezifische Ergebnistexte
B. Mehrfach richtige Antworten
In ausgewählten Feldern müssen multiple_select und passende scoring_mode-Varianten sauber unterstützt werden.
C. Qualitätssicherung
Reviewability ist kein Zusatz, sondern Kernbestandteil einer verlässlichen Wissens- und Coaching-Engine.

Empfohlene nächste Umsetzungsschritte
Sprint 1 – Dokumente und Schema V2
	•	Dokumentation aktualisieren
	•	Datenmodell V2 formulieren
	•	neue CSV-/Sheet-Strukturen definieren
Sprint 2 – Technischer Kern V2
	•	Go-Modelle erweitern
	•	SQLite-Schema anpassen
	•	Review- und Domänenfelder einführen
Sprint 3 – Testbetrieb
	•	Testmodus aktivieren
	•	echte Rückmeldungen sammeln
	•	Fragenqualität iterativ verbessern
Sprint 4 – Domänenpfade
	•	Profildomäne stabilisieren
	•	Coaching-Domäne vorbereiten

Definition von „fertig“ in V2
V2-Kern ist fertig, wenn
	•	Fragen Domänen kennen
	•	verschiedene Fragetypen technisch modelliert sind
	•	Review-Status gespeichert werden können
	•	Test-Flags aufgenommen werden können
	•	hierarchische Folgefragen vorbereitet sind
	•	sensible Inhalte markiert werden können
Coaching-Ausbaustufe ist fertig, wenn
	•	Themenpfade adaptiv verzweigen
	•	Trigger- und Körperbezüge abbildbar sind
	•	Ergebnistexte sicher und nicht-diagnostisch formuliert sind
	•	professionelle Nutzung methodisch abgesichert ist

Strategische Zusammenfassung
Mit dieser Roadmap wächst CEE4AI nicht chaotisch, sondern kontrolliert in Richtung einer methodisch starken, adaptiven, reviewbaren und professionell anschlussfähigen Engine.
Genau dadurch kann das Projekt später mehrere ernsthafte Anwendungsräume tragen, ohne seinen Kern zu verlieren.
