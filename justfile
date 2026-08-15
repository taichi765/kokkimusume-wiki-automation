console_log:
    az monitor log-analytics query \
        --workspace 'a773bac6-9977-4dc9-9faa-76d0f48927bc' \
        --analytics-query "ContainerAppConsoleLogs_CL\
        | where RevisionName_s == 'app--0000007'\
        | where time_t >= todatetime('2026/08/15 16:10:00')\
        | project Log_s"\
        |jq -r "[ \
            .[].Log_s \
            | fromjson \
            | select(.source.file | startswith(\"/build/discord-bot/\"))\
        ]"

system_log:
    az monitor log-analytics query \
        --workspace 'a773bac6-9977-4dc9-9faa-76d0f48927bc' \
        --analytics-query "ContainerAppSystemLogs_CL\
        | where RevisionName_s == 'app--0000007' and time_t >= todatetime('2026/08/15 16:10:00')\
        | project TimeGenerated, time_t, Log_s, Level, Type_s, Reason_s\
        | order by time_t desc "