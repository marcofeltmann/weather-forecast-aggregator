How to call the API (w/o key):
```sh
curl 'https://api.open-meteo.com/v1/forecast?latitude=42.6493934&longitude=-8.8201753&start_date=2024-11-09&end_date=2024-11-09&daily=temperature_2m_max'
```

How the data returns:
```json
{"latitude":42.5625,"longitude":-8.8125,"generationtime_ms":0.033020973205566406,"utc_offset_seconds":0,"timezone":"GMT","timezone_abbreviation":"GMT","elevation":2.0,"daily_units":{"time":"iso8601","temperature_2m_max":"°C"},"daily":{"time":["2024-11-09"],"temperature_2m_max":[20.6]}}
```
