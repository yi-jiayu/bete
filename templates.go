package bete

const (
	templateArrivalSummary = `{{ if .Stop.Description }}<strong>{{ .Stop.Description }} ({{ .Stop.ID }})</strong>
{{ else }}<strong>{{ .Stop.ID }}</strong>
{{ end }}{{ with .Stop.RoadName }}{{ . }}
{{ end }}<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
{{- range (.Services | filterByService .Filter | sortByService) }}
{{ $fst := until $.Time .NextBus.EstimatedArrival -}}
{{ $snd := until $.Time .NextBus2.EstimatedArrival -}}
{{ $thd := until $.Time .NextBus3.EstimatedArrival -}}
| {{ printf "%-4v" .ServiceNo }} | {{ printf "%3v" $fst }} | {{ printf "%3v" $snd }} | {{ printf "%3v" $thd }} |
{{- end }}
</pre>
{{ with .Filter }}Filtered by services: {{ join . ", " }}
{{ end }}<em>Last updated on {{ .Time | inSGT }}</em>`
	templateArrivalDetails = `{{ if .Stop.Description }}<strong>{{ .Stop.Description }} ({{ .Stop.ID }})</strong>
{{ else }}<strong>{{ .Stop.ID }}</strong>
{{ end }}{{ with .Stop.RoadName }}{{ . }}
{{ end }}<pre>
Svc   Eta  Sea  Typ  Fea
---   ---  ---  ---  ---
{{- range (.Services | filterByService .Filter | arrivingBuses | sortByArrival | take 10) }}
{{ $eta := until $.Time .EstimatedArrival -}}
{{ $sea := .Load -}}
{{ $typ := .Type -}}
{{ $fea := .Feature -}}
{{ printf "%-4v" .ServiceNo }} {{ printf "%4v" $eta }} {{ printf "%4v" $sea }} {{ printf "%4v" $typ }} {{ printf "%4v" $fea -}}
{{ end }}
</pre>
{{ with .Filter }}Filtered by services: {{ join . ", " }}
{{ end }}<em>Last updated on {{ .Time | inSGT }}</em>`
)
