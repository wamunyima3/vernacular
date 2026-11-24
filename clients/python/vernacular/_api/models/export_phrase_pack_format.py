from enum import Enum


class ExportPhrasePackFormat(str, Enum):
    JSON = "json"
    PO = "po"
    XLIFF = "xliff"

    def __str__(self) -> str:
        return str(self.value)
