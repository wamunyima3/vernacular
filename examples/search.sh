#!/bin/sh
# Search the Bemba dictionary for "amenshi", unchecked words included.
curl --fail-with-body --silent \
  --header "X-API-Key: $VERNACULAR_KEY" \
  "$VERNACULAR_URL/v1/entries/search?language=bem&q=amenshi&min_status=machine"
