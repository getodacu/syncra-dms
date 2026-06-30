CREATE TABLE "json_recipes" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "title" varchar(160) NOT NULL,
  "description" text,
  "json" jsonb NOT NULL,
  "counter" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_json_recipes_counter" CHECK ("counter" >= 0)
);

CREATE INDEX "idx_json_recipes_created_id" ON "json_recipes" ("created_at", "id");

INSERT INTO "json_recipes" (
  "id",
  "title",
  "description",
  "json",
  "counter",
  "created_at",
  "updated_at"
) VALUES
(
  '00000000-0000-4000-8000-000000000101',
  'Invoice',
  'Romanian invoice extraction fields.',
  $json$
{
    "type": "object",
    "additionalProperties": false,
    "properties": {
      "furnizor_nume": {
        "description": "Denumirea furnizorului",
        "type": "string"
      },
      "furnizor_cui": {
        "description": "CUI/CIF furnizor",
        "type": "string"
      },
      "client_nume": {
        "description": "Denumirea clientului",
        "type": "string"
      },
      "client_cui": {
        "description": "CUI/CIF client",
        "type": "string"
      },
      "numar_factura": {
        "description": "Numărul facturii, extras ca string",
        "type": "string"
      },
      "data_emitere": {
        "description": "Data emiterii facturii în format YYYY-MM-DD",
        "type": "string"
      },
      "data_scadenta": {
        "description": "Data scadentă în format YYYY-MM-DD, dacă există",
        "type": "string"
      },
      "moneda": {
        "description": "Moneda facturii, de exemplu RON, EUR, USD",
        "type": "string"
      },
      "subtotal_fara_tva": {
        "description": "Valoarea totală fără TVA",
        "type": "number"
      },
      "total_tva": {
        "description": "Valoarea totală TVA",
        "type": "number"
      },
      "total_cu_tva": {
        "description": "Valoarea totală cu TVA",
        "type": "number"
      },
      "linii": {
        "description": "Liniile facturii",
        "type": "array",
        "items": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "descriere": {
              "description": "Descriere produs sau serviciu",
              "type": "string"
            },
            "cantitate": {
              "description": "Cantitatea",
              "type": "number"
            },
            "pret_unitar": {
              "description": "Preț unitar fără TVA sau preț unitar detectat",
              "type": "number"
            },
            "cota_tva": {
              "description": "Cota TVA",
              "type": "number"
            },
            "valoare_totala": {
              "description": "Valoarea totală a liniei",
              "type": "number"
            }
          },
          "required": [
            "descriere",
            "cantitate",
            "pret_unitar",
            "cota_tva",
            "valoare_totala"
          ]
        }
      },
      "observatii": {
        "description": "Observații sau mențiuni de pe factură",
        "type": "string"
      }
    },
    "required": [
      "furnizor_nume",
      "furnizor_cui",
      "client_nume",
      "client_cui",
      "numar_factura",
      "data_emitere",
      "data_scadenta",
      "moneda",
      "subtotal_fara_tva",
      "total_tva",
      "total_cu_tva",
      "linii",
      "observatii"
    ]
  }
$json$::jsonb,
  0,
  now(),
  now()
),
(
  '00000000-0000-4000-8000-000000000102',
  'Bon fiscal',
  'Romanian fiscal receipt extraction fields.',
  $json$
{
    "type": "object",
    "additionalProperties": false,
    "properties": {
      "comerciant_nume": {
        "description": "Denumirea comerciantului",
        "type": "string"
      },
      "comerciant_cif": {
        "description": "CIF/CUI comerciant",
        "type": "string"
      },
      "comerciant_adresa": {
        "description": "Adresa comerciantului sau a punctului de lucru",
        "type": "string"
      },
      "numar_bon": {
        "description": "Numărul bonului fiscal, extras ca string",
        "type": "string"
      },
      "data_emitere": {
        "description": "Data emiterii în format YYYY-MM-DD",
        "type": "string"
      },
      "ora_emitere": {
        "description": "Ora emiterii în format HH:MM:SS, dacă apare",
        "type": "string"
      },
      "moneda": {
        "description": "Moneda documentului, de obicei RON",
        "type": "string"
      },
      "linii": {
        "description": "Produsele sau serviciile cumpărate",
        "type": "array",
        "items": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "descriere": {
              "description": "Denumirea produsului sau serviciului",
              "type": "string"
            },
            "cantitate": {
              "description": "Cantitatea",
              "type": "number"
            },
            "pret_unitar": {
              "description": "Prețul unitar",
              "type": "number"
            },
            "valoare_totala": {
              "description": "Valoarea totală a liniei",
              "type": "number"
            },
            "cota_tva": {
              "description": "Cota TVA, dacă poate fi identificată",
              "type": "number"
            }
          },
          "required": [
            "descriere",
            "cantitate",
            "pret_unitar",
            "valoare_totala",
            "cota_tva"
          ]
        }
      },
      "total_cu_tva": {
        "description": "Totalul final al bonului fiscal",
        "type": "number"
      },
      "total_tva": {
        "description": "Total TVA, dacă apare sau poate fi calculat",
        "type": "number"
      },
      "metoda_plata": {
        "description": "Metoda de plată: numerar, card, voucher, tichet etc.",
        "type": "string"
      },
      "suma_platita": {
        "description": "Suma plătită",
        "type": "number"
      },
      "rest": {
        "description": "Restul acordat clientului, dacă apare",
        "type": "number"
      },
      "observatii": {
        "description": "Alte mențiuni de pe bon",
        "type": "string"
      }
    },
    "required": [
      "comerciant_nume",
      "comerciant_cif",
      "comerciant_adresa",
      "numar_bon",
      "data_emitere",
      "ora_emitere",
      "moneda",
      "linii",
      "total_cu_tva",
      "total_tva",
      "metoda_plata",
      "suma_platita",
      "rest",
      "observatii"
    ]
  }
$json$::jsonb,
  0,
  now(),
  now()
),
(
  '00000000-0000-4000-8000-000000000103',
  'Carte de identitate',
  'Romanian identity card extraction fields.',
  $json$
{
    "type": "object",
    "properties": {
      "Nume": {
        "type": "string",
        "description": ""
      },
      "Prenume": {
        "type": "string",
        "description": ""
      },
      "Sex": {
        "type": "string",
        "description": "Sexul biologic: M for male, F for female",
        "enum": [
          "M",
          "F"
        ]
      },
      "CNP": {
        "type": "number",
        "description": "Cod numeric personal format din 13 cifre"
      },
      "Data nasterii": {
        "type": "string",
        "description": "Data nasterii in format DD.MM.YYYY"
      },
      "Nr. document": {
        "type": "string",
        "description": "Numarul documentului"
      }
    },
    "required": [
      "Nume",
      "Sex",
      "CNP",
      "Data nasterii",
      "Nr. document",
      "Prenume"
    ]
  }
$json$::jsonb,
  0,
  now(),
  now()
)
ON CONFLICT ("id") DO NOTHING;
