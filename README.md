# Fraud Detection

Rozwiązanie składa się z wielu modułów.

![architecture](docs/diagram.png)

Projekt można uruchomić korzystając z polecenia `docker compose up`. Należy chwilę odczekać, gdyż Kafka się długo uruchamia i dopiero po jej uruchomieniu rozwiązanie zaczyła działać stabilnie.

## [Generator](generator/README.md)
Tworzy losowe transakcje które są zapisywane na topic Kafki `transactions` w formacie JSON.
Parametry generatora można zmieniać poprzez parametry w pliku `docker-compose.yaml`.

## [Flink](frauddetection/README.md)
Zarządza i uruchamia algorytmy sprawdzające oszustwa. Czyta z topicu Kafki `transactions`, a zapisuje wyniki do topicu `alerts`.

Jest kilka algorytmów sprawdzających wiarygodność transakcji.
- `ExpiredCardDetector` - sprawdza, czy transakcja nie odbyła się po skończeniu ważności karty,
- `OverLimitDetector` - sprawdza, czy transakcja mieści się w limicie,
- `SmallThenLargeDetector` - sprawdza, czy w okienku minutowym nie dopuszczono się małej transakcji (poniżej 20 zł), po czym wykonano dużą transakcję (powyżej 500 zł),
- `NormalDistributionDetector` - oblicza parametry rozkładu normalnego: średnią oraz wariancję, a na ich podstawie odrzuca transakcje które są oddalone od średniej o więcej niż odchylenie standardowe,
- `LocationDetector` - który sprawdza dla ustalonego "session window" czy w trakcie okna transakcje nie są nadmiarowo oddalone od średniej lokalizacji transakcji w oknie.

Algorytmy zapisane są w folderze `frauddetection/src/java/spendreport/detectors`.

## [Alerts Reader](alerts-notification-api-py/README.md)
Zczytuje alerty z topicu Kafki `alerts` i przekazuje je klientom nasłuchującym na Websockecie.

## [Dashboard](dashboard/README.md)
Łączy się z Alerts Reader po Websockecie i na żywo wyświetla alerty oraz zlicza je na wykresie w okienkach 30 sekundowych.

## [Kafka i Kafdrop](kafka/README.md)
Kafka ma dwa topici umożliwiające komunikację. Działanie Kafki możemy podejrzeć korzystając z interfejsu graficznego Kafdrop.