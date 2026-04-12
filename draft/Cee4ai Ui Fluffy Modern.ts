export default function CEE4AIUI() {
  const answers = [
    { id: 1, label: "18" },
    { id: 2, label: "24" },
    { id: 3, label: "32" },
    { id: 4, label: "34" },
  ];

  const chips = [
    "Analytisch-logisch",
    "Kreativ-assoziativ",
    "Systemisch-vernetzt",
    "Intuitiv-empathisch",
    "Reflektiv-bewusst",
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-white via-slate-50 to-violet-50 text-slate-800">
      <div className="mx-auto max-w-6xl px-6 py-10">
        <div className="mb-8 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <p className="mb-2 inline-flex rounded-full border border-violet-200 bg-white/80 px-3 py-1 text-xs tracking-wide text-violet-700 shadow-sm backdrop-blur">
              Cognitive Evolution Engine for Adaptive Intelligence
            </p>
            <h1 className="text-3xl font-semibold tracking-tight md:text-4xl">
              Sanfte, moderne Fragenreise mit Profilbildung
            </h1>
            <p className="mt-3 max-w-3xl text-sm leading-6 text-slate-600 md:text-base">
              Kein klassischer IQ-Test. Eher eine bewusste Reise durch Denkstile,
              Mustererkennung, Intuition und Selbstreflexion.
            </p>
          </div>

          <div className="w-full max-w-sm rounded-3xl border border-white/70 bg-white/80 p-4 shadow-xl shadow-slate-200/60 backdrop-blur">
            <div className="mb-2 flex items-center justify-between text-sm text-slate-500">
              <span>Fortschritt</span>
              <span>3 / 10</span>
            </div>
            <div className="h-3 overflow-hidden rounded-full bg-slate-100">
              <div className="h-full w-[30%] rounded-full bg-slate-800" />
            </div>
            <p className="mt-3 text-xs text-slate-500">
              Jede Antwort stärkt dein aktuelles Denkprofil.
            </p>
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[1.35fr_0.65fr]">
          <section className="rounded-[28px] border border-white/70 bg-white/85 p-6 shadow-2xl shadow-slate-200/50 backdrop-blur">
            <div className="mb-6 flex flex-wrap items-center gap-2">
              <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">
                Kategorie: Allgemeinwissen
              </span>
              <span className="rounded-full bg-violet-100 px-3 py-1 text-xs font-medium text-violet-700">
                Schwierigkeit: 1
              </span>
              <span className="rounded-full bg-emerald-100 px-3 py-1 text-xs font-medium text-emerald-700">
                Denktyp: analytisch-logisch
              </span>
            </div>

            <div className="mb-8">
              <p className="mb-2 text-sm uppercase tracking-[0.18em] text-slate-400">
                Frage 3
              </p>
              <h2 className="text-2xl font-semibold leading-tight md:text-3xl">
                Welche Zahl setzt die Reihe sinnvoll fort: 2, 4, 8, 16, ?
              </h2>
              <p className="mt-4 max-w-2xl text-sm leading-6 text-slate-500">
                Antworte intuitiv, aber aufmerksam. Dein Profil entwickelt sich nicht
                nur durch richtig oder falsch, sondern auch durch Muster in deinen Entscheidungen.
              </p>
            </div>

            <div className="grid gap-3">
              {answers.map((answer) => (
                <button
                  key={answer.id}
                  className="group rounded-2xl border border-slate-200 bg-white px-4 py-4 text-left transition hover:-translate-y-0.5 hover:border-slate-400 hover:shadow-lg"
                >
                  <div className="flex items-center gap-4">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-slate-100 text-sm font-semibold text-slate-700 transition group-hover:bg-slate-800 group-hover:text-white">
                      {String.fromCharCode(64 + answer.id)}
                    </div>
                    <div>
                      <p className="text-base font-medium text-slate-800">{answer.label}</p>
                      <p className="text-sm text-slate-500">Antwort auswählen</p>
                    </div>
                  </div>
                </button>
              ))}
            </div>

            <div className="mt-8 flex flex-wrap items-center gap-3">
              <button className="rounded-2xl bg-slate-900 px-5 py-3 text-sm font-medium text-white shadow-lg shadow-slate-300 transition hover:-translate-y-0.5">
                Antwort speichern
              </button>
              <button className="rounded-2xl border border-slate-200 bg-white px-5 py-3 text-sm font-medium text-slate-700 transition hover:border-slate-400">
                Frage überspringen
              </button>
            </div>
          </section>

          <aside className="space-y-6">
            <div className="rounded-[28px] border border-white/70 bg-white/85 p-5 shadow-xl shadow-slate-200/50 backdrop-blur">
              <h3 className="text-lg font-semibold">Denktypen im Blick</h3>
              <p className="mt-2 text-sm leading-6 text-slate-500">
                Dein Profil ist nicht starr. Es entwickelt sich mit jeder beantworteten Frage.
              </p>
              <div className="mt-4 flex flex-wrap gap-2">
                {chips.map((chip) => (
                  <span
                    key={chip}
                    className="rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs text-slate-600"
                  >
                    {chip}
                  </span>
                ))}
              </div>
            </div>

            <div className="rounded-[28px] border border-violet-100 bg-gradient-to-br from-violet-50 to-white p-5 shadow-xl shadow-violet-100/50">
              <h3 className="text-lg font-semibold text-slate-800">Evolutionärer Hinweis</h3>
              <p className="mt-3 text-sm leading-6 text-slate-600">
                Du bewegst dich aktuell stabil auf einer analytischen Ebene. Der nächste Entwicklungsschritt
                liegt darin, Muster nicht nur logisch, sondern auch vernetzt zu betrachten.
              </p>
            </div>

            <div className="rounded-[28px] border border-emerald-100 bg-gradient-to-br from-emerald-50 to-white p-5 shadow-xl shadow-emerald-100/50">
              <h3 className="text-lg font-semibold text-slate-800">Sanfte UX-Richtung</h3>
              <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-600">
                <li>• viel Weißraum und ruhige Flächen</li>
                <li>• weiche Rundungen und dezente Schatten</li>
                <li>• Fortschritt statt Leistungsdruck</li>
                <li>• Ergebnis als Profil, nicht als Urteil</li>
              </ul>
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
}
