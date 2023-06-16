## Fraud detection

Kluczowy moduł aplikacji odpowiedzialny za przetwarzanie danych w czasie rzeczywistym.
Wykorzystuje on bibliotekę Flink oraz jest napisany w Java 11.
Zczytuje dane (source) z topiku `transactions` Kafki, a następnie przekazuje je do detektorów. 
Same detektory działają równolegle i przekazują wykryte anomalie zarówno do topiku `alerts` Kafki, a także do procesu `AlertSink`, który loguje na bieżąco alerty.

Do przekazywania informacji wykorzystywane są `DataStream` w przypadku alertów oraz `KeyedStream` w przypadku transakcji.
Część detektorów, m.in. `NormalDistributionDetector` korzysta ze stanu przechowywanego przez Flink.
`KeyedStream` jest wymagany, aby wymusić przetwarzanie przez ten sam proces kolejnych transakcji z tej samej karty.
W przeciwnym wypadku przy rozproszonym przetwarzaniu moglibyśmy mieć problemy z niespójnym stanem, brudnymi odczytami, czy zbędnym i czasochłonnym przenoszeniem stanu między procesami.
Opis przepływu danych i konfiguracja połączeń z Kafką są zapisane w `FraudDetectionJob.java`.
Zależności między taskami we Flinku przedstawione są poniżej:
![]().
<!-- TODO: wstawic screenshot z Flinka -->

### Detektory
<!-- TODO: dodac opis detektorów -->

### Dockerfile
Flink umożliwia przekazywanie Job w formacie JAR na dwa sposoby. 
<!-- TODO: dodac link -->
Pierwszy z nich pozwala na dołączanie JAR w runtime Flinka, a drugi poprzez dołączenie ich do odpowiedniego folderu.
W projekcie wykorzystywana jest opcja druga.

W pliku `Dockerfile` możemy zauważyć, że projekt jest najpierw kompilowany przy użyciu Maven, a następnie stworzony JAR jest kopiowany do obrazu Flinka w katalogu `/opt/flink/usrlib/`.

Następnie w `docker-compose.yml` stworzone zostały dwa rodzaje kontenerów na podstawie jednego obrazu.
Są to `jobmanager` oraz `taskmanager`.
Do uruchomienia `jobmanager` wykorzystywane jest polecenie Flink `standalone-job` w którym przekazujemy klasę opisującą Job, czyli w naszym przypadku `--job-classname spendreport.FraudDetection`.
Z drugiej strony mamy `taskmanager`, który będzie wykonywał pojedyncze taski. Uruchamiany jest na tym samym obrazie co `jobmanager`, tylko z poleceniem `taskmanager`.
Dodatkowo możemy zwiększyć liczbę kontenerów `taskmanager` poprzez zwiększenie parametru `scale` (domyślnie: 1).

W zmiennych środowiskowych ustawiane są wspólnie `jobmanager.rpc.address` wzkazujący adres jobmanagera oraz `parallelism.default` mówiący o ilości działających równolegle tasków tego samego rodzaju.
Warto zauważyć, że `jobmanager.rpc.address` ma ustawiony adres jaki przypisze mu `docker-compose`, czyli w tym przypadku skoro nazwaliśmy go `jobmanager`, to `jobmanager.rpc.address: jobmanager`.
W `taskmanager` ustalany jest także `taskmanager.numberOfTaskSlots`.