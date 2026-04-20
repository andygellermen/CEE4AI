# credential-trust-ladder-v3-final.md

## Zweck
Dieses Dokument definiert die Zertifikats- und Vertrauensarchitektur der CPE-Produktfamilie in der V3-Endfassung.

Es beantwortet:
- Welche Zertifikatsstufen gibt es?
- Welche Regeln müssen erfüllt sein?
- Wann ist ein Zertifikat motivierend und wann wirklich belastbar?
- Wie wird Meaning-/Journey-Zertifizierung sauber eingebunden?
- Wie entsteht langfristig Außenwirkung und Marktrelevanz?

---

## Leitprinzip
Ein Zertifikat ist zunächst **keine Aura**, sondern eine **verifizierbare Aussage unter definierten Regeln**.

Vertrauen entsteht aus:
1. sauberer Messlogik
2. nachvollziehbarer Evidenz
3. Governance
4. technischer Verifizierbarkeit
5. Reputation und Partnernetzwerk

---

## Die Trust Ladder
### Level 0 – Participation Acknowledgement
**Charakter:** Teilnahme- und Aktivitätsnachweis  
**Zweck:** Motivation, Einstieg, geringe Hürde

### Level 1 – Snapshot Certificate
**Charakter:** frühes Trend- oder Vorprofil  
**Zweck:** Appetizer, Rückkehrofferte, erste Orientierung

### Level 2 – Assessed Profile Certificate
**Charakter:** qualifizierte Profilbescheinigung  
**Zweck:** belastbarere interne oder externe Orientierung

### Level 3 – Extended Profile & Pathway Certificate
**Charakter:** vertieftes Fähigkeits-, Persönlichkeits- oder Entwicklungsprofil  
**Zweck:** stärkere Nutzung für Karriere, Transformation und Entwicklung

### Level 4 – Verified Certification
**Charakter:** verifizierte und stärker abgesicherte Zertifizierung  
**Zweck:** ernsthafte Außenwirkung, höhere Glaubwürdigkeit

### Level 5 – Partner-Endorsed / Industry-Linked Credential
**Charakter:** externer Vertrauenshebel durch Partner oder Ökosystem  
**Zweck:** echte Marktrelevanz und Türöffner-Funktion

---

## Zertifikatsobjekte
- `assessment_program`
- `assessment_track`
- `certificate_rule`
- `certificate_template`
- `credential_assertion`
- `credential_evidence`
- `credential_verification`
- `credential_revocation`

---

## Mindestregeln pro Zertifikatsstufe
### Snapshot
- `minimum_completion`: niedrig
- `minimum_confidence`: niedrig
- `identity_verification`: nein
- `proctoring`: nein

### Assessed Profile
- `minimum_completion`: mittel
- `minimum_confidence`: mittel
- `identity_verification`: optional
- `proctoring`: nein

### Extended Profile
- `minimum_completion`: hoch
- `minimum_confidence`: mittel bis hoch
- `identity_verification`: optional bis empfohlen
- `proctoring`: optional je Track

### Verified Certification
- `minimum_completion`: hoch
- `minimum_confidence`: hoch
- `identity_verification`: ja
- `proctoring`: ja oder gleichwertiger Integritätsmechanismus

### Partner-Endorsed
- alle Punkte von Verified
- zusätzliche Partner- und Auditlogik

---

## Meaning-/Journey-Zertifikate
### Sinnvolle Formen
- Meaning Journey Completion
- Reflective Development Track Completion
- Guided Spiritual Pathway Completion
- Contemplative Practice Completion

### Charakter
Diese Zertifikate bestätigen:
- Teilnahme
- Durchlaufen eines strukturierten Pfades
- dokumentierte Reflexions- oder Entwicklungsarbeit

Sie bestätigen **nicht** objektive spirituelle Höherstufung.

### Zulässige Sprache
- Journey
- Completion
- Reflective Track
- Guided Pathway

### Nicht empfohlene Sprache
- höhere Bewusstseinsstufe bestätigt
- spirituelle Reife objektiv zertifiziert
- metaphysische Autorität bescheinigt

---

## Confidence- und Integrity-Logik
Ein Zertifikat darf nie nur auf Punkten beruhen.

Zusätzlich sollten einfließen:
- Vollständigkeit
- Antwortkonsistenz
- Session-Integrität
- Zeitmuster
- Reviewfreiheit kritischer Fragen
- Reife des zugrunde liegenden Inhaltsmoduls

### Empfohlene Ergebnisfelder
- `result_confidence`
- `assessment_integrity_score`
- `profile_depth`
- `ruleset_version`
- `issued_under_ruleset_version`

---

## Zertifikatstexte: Pflichtbausteine
Jedes Zertifikat braucht klar sichtbare Metadaten:
- Zertifikatsstufe
- ausstellende Stelle
- Ausstellungsdatum
- Regelversion
- ggf. Ablaufdatum
- Verifikationsmöglichkeit
- Hinweis auf Art des Nachweises

### Beispielhafte Klarstellungen
- „Trendprofil / vorläufige Einschätzung“
- „qualifizierter Profilnachweis“
- „verifizierte Zertifizierung unter Regelwerk Version X“
- „strukturierter Reflexions- und Entwicklungsweg erfolgreich abgeschlossen“

---

## Digitale Verifizierbarkeit
Das Zertifikatssystem sollte auf verifizierbare digitale Credentials ausgerichtet sein.

### Zielbild
- maschinenverifizierbare Assertions
- trennbare Rolle von Aussteller, Inhaber und Prüfer
- nachprüfbare Evidenz
- Widerruf / Sperrung möglich

### Strategische Richtung
- Open Badges / CLR / Verifiable Credentials als Zielkompatibilität
- anfangs pragmatische interne Struktur erlaubt

---

## Revocation und Re-Issuance
### Gründe für Entzug
- Regelverletzung
- Integritätsprobleme
- fehlerhafte Auswertung
- methodische Revision
- Missbrauch

### Gründe für Neu-Ausstellung
- aktualisierte Regelversion
- nachträgliche Identitätsprüfung
- ergänzte Evidenz
- korrigierte Bewertung

---

## Ehrfurcht und Außenwirkung – die nüchterne Wahrheit
Wahre Schwere entsteht nicht durch Design allein, sondern durch:
- methodische Güte
- verlässliche Regeln
- technische Nachprüfbarkeit
- Partner und Reputation
- Sichtbarkeit im Markt

### Deshalb zweigleisig denken
#### 1. Technische Vertrauensarchitektur
Regeln, Verifikation, Audit, Versionierung.

#### 2. Markt- und Reputationsarchitektur
Partnerschaften, Lobbying, Sichtbarkeit, Anschlussfähigkeit.

---

## Zertifikatskategorien in der Produktfamilie
### CEE4AI Cognitive Profile
- Snapshot
- Assessed Cognitive Profile
- Extended Cognitive Pathway Certificate

### CEE Personality Profile
- Personality Snapshot
- Personality Pattern Certificate
- Extended Personality Development Certificate

### CEE Coaching Analyzer
- Completion- und Reflexionsnachweise
- vorsichtig im Wording
- keine klinische Zertifizierung suggerieren

### CEE Pathways
- Pathway Readiness
- Transformation Track Completion
- Career / Learning Direction Certificate

### CEE Meaning & Spiritual Depth
- Meaning Journey Completion
- Reflective Development Track Completion
- Guided Spiritual Pathway Completion

---

## Abschluss
Die Trust Ladder V3 Final schützt das Produkt vor Überversprechen und schafft zugleich einen strategischen Weg zu echter Außenwirkung.  
Sie verbindet Motivation, Governance, Technik, Meaning-Sensibilität und Marktvertrauen zu einer aufbaubaren Zertifikatsarchitektur.

