"""Search the Bemba dictionary for "amenshi", unchecked words included. Python 3.9 or later."""
import json
import os
import urllib.parse
import urllib.error
import urllib.request

query = urllib.parse.urlencode({"language": "bem", "q": "amenshi", "min_status": "machine"})
request = urllib.request.Request(
    f"{os.environ['VERNACULAR_URL']}/v1/entries/search?{query}",
    headers={"X-API-Key": os.environ["VERNACULAR_KEY"]},
)
try:
    with urllib.request.urlopen(request) as response:
        hits = json.load(response)["data"]
except urllib.error.HTTPError as error:
    body = json.load(error)
    raise SystemExit(f"{error.code} {body['code']}: {body['message']}") from error

for hit in hits:
    entry = hit["entry"]
    meanings = [f"{s['gloss']} ({s['trust']['status']})" for s in entry["senses"]]
    print(f"{entry['headword']}: {', '.join(meanings)}")
