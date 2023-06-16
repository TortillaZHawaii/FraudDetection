## Alerts notification API (aka Alerts Reader)

Z wielu powodów nie można odczytywać bezpośrednio z przeglądarki wiadomości zapisywanych na topicach Kafki.
Głównym powodem jest ograniczony dostęp do niższych warstw sieciowych w API przeglądarek.
Z tego powodu powstał `alerts-notification-api`.
Jest to proxy pomiędzy Kafką, a Dashboardem.

### Sposób działania

Ze względu na wymaganie równoczesnego połączenia się wielu klientów, a także czytania alertów z Kafki zdecydowano się na skorzystanie z biblioteki `asyncio`.

Zamiast `python-kafka`, która działa sekwencyjnie, skorzystano z klienta `aiokafka`, który umożliwia asynchroniczność odczytu i nie blokuje działania w trakcie oczekiwania na kolejne wiadomości.

Asynchronicznie do odczytu działa serwer Web Socket, domyślnie nasłuchujący na [ws://localhost:12000](). 
Sama biblioteka `ws` dodaje funkcję `broadcast`, która umożliwia wysłanie wiadomości wszystkim połączonym klientom, w tym przypadku po odczytaniu wiadomości z Kafki.
