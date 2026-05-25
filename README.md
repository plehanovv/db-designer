# DB Designer VKR

Subsystem prototype for automated database structure design from a natural language domain description.

## Run Go service

```powershell
go run -buildvcs=false ./cmd/server
```

Open `http://localhost:8080`.

Optional environment variables:

- `PORT` changes the Go HTTP port.
- `NLP_SERVICE_URL` changes the NLP endpoint. Default: `http://localhost:8000/analyze`.

## Run NLP service

From `D:\Python\nlp-service`:

```powershell
py -3.11 -m venv .venv
.\.venv\Scripts\python.exe -m pip install -r requirements.txt
.\.venv\Scripts\python.exe -m spacy download en_core_web_sm
.\.venv\Scripts\python.exe -m spacy download ru_core_news_sm
.\.venv\Scripts\python.exe -m uvicorn main:app --reload --port 8000
```

The Go service works without the NLP service by using a local rule-based fallback, but spaCy improves tokenization, lemmas and part-of-speech tags.
