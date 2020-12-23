{
  "$schema": "https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.4.json",
  "version": "2.1.0",
  "runs": [
    {{- $run_first := true }}
    {{- range $report_index, $report := . }}
    {{- if and $report.Valid (not (eq $report.Message "This resource kind is not supported by kubesec")) -}}
      {{- if $run_first -}}
        {{- $run_first = false -}}
      {{ else -}}
        ,
      {{- end }}
    {
      "tool": {
        "driver": {
          "name": "Kubesec",
          "fullName": "Kubesec Kubernetes Resource Security Policy Validator",
          "rules": [
        {{- $rule_first := true }}
          {{- range .Rules }}
            {{- if $rule_first -}}
              {{- $rule_first = false -}}
            {{ else -}}
              ,
            {{- end }}
            {
              "id": "{{ .ID }}",
              "shortDescription": {
                "text": "{{ .Reason }}"
              },
              "messageStrings": {
                "selector": {
                  "text": {{ escapeString .Selector | printf "%q" }}
                }
              },
              "properties": {
                "points": "{{ .Points }}"
              }
            }
          {{- end -}}
          ]
        }
      },
      "results": [
      {{- $result_first := true }}
      {{- range $result_index, $res := joinSlices .Scoring.Advise .Scoring.Critical -}}
        {{- if $result_first -}}
          {{- $result_first = false -}}
        {{ else -}}
          ,
        {{- end }}
        {
          "ruleId": "{{ $res.ID }}",
          "level": "warning",
          "message": {
            "text": {{ endWithPeriod $res.Reason | printf "%q" }},
            "properties": {
              "score": "{{ $res.Points }}",
              "selector": {{ escapeString $res.Selector | printf "%q" }}
            }
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "{{ $report.FileName }}"
                }
              }
            }
          ]
        }
      {{- end -}}
      ],
      "columnKind": "utf16CodeUnits"
    }
  {{- end -}}
  {{- end }}
  ]
}
