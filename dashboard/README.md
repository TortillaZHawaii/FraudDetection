## Dashboard

Moduł projektu odpowiedzialny za prezentowanie użytkownikowi końcowemu w wygodny sposób alertów na żywo.
Napisany w języku Typescript przy wykorzystaniu m.in. Next.js, React, web socket oraz pakietu komponentów tremor.

### Sposób działania

Po wejściu na stronę (domyślnie [http://localhost:3000]) komponent `alertTable.tsx` łączy się przy użyciu Web Socket z `alerts-notification-api-py` i nasłuchuje na nadchodzące alerty. Po otrzymaniu alertu, informacje o alercie są dodawane do tabeli, a także zliczane na wykresie w równych przedziałach czasowych.
