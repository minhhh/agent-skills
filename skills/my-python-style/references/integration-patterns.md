# External Integrations

## GSheets: Non-Destructive "Shifting" Updates
When updating spreadsheets that act as ledgers (where new data is prepended), avoid overwriting blocks. Use the `insertRange` API request to physically shift existing cells down before writing new data. This preserves formatting and data boundaries.