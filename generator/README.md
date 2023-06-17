## Generator
Generuje transakcje i zapisuje je na topicu `transactions` w formacie JSON.

CardOwner

    Id: Numer identyfikacyjny właściciela karty (liczba całkowita).
    FirstName: Imię właściciela karty (łańcuch znaków).
    LastName: Nazwisko właściciela karty (łańcuch znaków).

Card

    CardNumber: Numer karty (łańcuch znaków).
    CardType: Typ karty (łańcuch znaków).
    ExpMonth: Miesiąc ważności karty (liczba całkowita z przedziału od 1 do 12).
    ExpYear: Rok ważności karty (liczba całkowita z przedziału od 2023 do 2030).
    CVV: Kod CVV karty (łańcuch znaków).

CardTransaction

    Amount: Kwota transakcji (liczba całkowita).
    LimitLeft: Pozostały limit na karcie po transakcji (liczba całkowita).
    Currency: Waluta transakcji (łańcuch znaków).
    Latitude: Szerokość geograficzna transakcji (liczba zmiennoprzecinkowa z przedziału od 52.1 do 52.4).
    Longitude: Długość geograficzna transakcji (liczba zmiennoprzecinkowa z przedziału od 20.8 do 21.2).
    Card: Obiekt karty (Card).
    Owner: Obiekt właściciela karty (CardOwner).
    UTC: Czas transakcji w formacie czasu UTC (time.Time).

### Generowanie danych

Wartości dla pól są generowane przy użyciu biblioteki gofakeit.
Pole Latitude i Longitude są generowane w prostokącie obejmującym Warszawę.

### Parametry
Generator udostępnia kilka parametrów zmieniających sposób generowania danych. Są one zczytywane ze zmiennych środowiskowych:
- `SEED` - ziarno generatora liczb pseudolosowych (domyślnie 0, czyli losowe),
- `CARD_OWNER_COUNT` - liczba właścicieli kart,
- `CARD_COUNT` - łączna liczba kart (właściciele przypisywani są do kart w sposób równomierny),
- `WAIT_TIME_MS` - czas oczekiwania pomiędzy generowaniem kolejnych transakcji (w milisekundach).
Generator także pozwala na zmianę parametrów Kafki:
- `KAFKA_BROKER` - adres brokera Kafki,
- `KAFKA_TOPIC` - nazwa topicu do którego wysyłane są transakcje.
