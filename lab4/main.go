package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("--- Symulacja Równoległego Przetwarzania Zamówień ---")

	// Inicjalizacja kanałów
	// Używamy buforowanych kanałów, aby uniknąć blokowania, jeśli workerzy
	// chwilowo nie nadążają lub generator jest szybszy.
	ordersChan := make(chan Order, numOrders)     // Kanał dla przychodzących zamówień
	resultsChan := make(chan ProcessResult, numOrders) // Kanał dla wyników przetwarzania

	// Inicjalizacja WaitGroup do synchronizacji gorutyn
	var wg sync.WaitGroup

	// 1. Uruchomienie workerów (pulę workerów)
	fmt.Printf("Uruchamianie %d workerów...\n", numWorkers)
	wg.Add(numWorkers) // Oczekujemy na zakończenie numWorkers gorutyn workerów
	for i := 1; i <= numWorkers; i++ {
		go processOrderWorker(i, ordersChan, resultsChan, &wg)
	}

	// 2. Uruchomienie kolektora wyników w osobnej gorutynie
	fmt.Println("Uruchamianie kolektora wyników...")
	wg.Add(1) // Oczekujemy na zakończenie gorutyny kolektora
	go resultsCollector(resultsChan, &wg)

	// 3. Uruchomienie generatora zamówień w osobnej gorutynie
	fmt.Println("Uruchamianie generatora zamówień...")
	wg.Add(1) // Oczekujemy na zakończenie gorutyny generatora
	go orderGenerator(ordersChan, &wg)

	// Czekanie na zakończenie wszystkich gorutyn *związanych z generowaniem i przetwarzaniem*
	// Musimy poczekać, aż generator skończy i zamknie ordersChan,
	// a workerzy przetworzą wszystko z ordersChan i zakończą pracę.
	// Jednakże, bezpośrednie czekanie na wg.Wait() tutaj spowodowałoby deadlock,
	// bo kolektor też jest w tej samej WaitGroup i czeka na zamknięcie resultsChan,
	// które zamykamy dopiero *po* zakończeniu workerów.

	// Zamiast tego, stworzymy osobną logikę czekania:
	// - Uruchomimy gorutynę, która poczeka na generator i workerów, a *następnie* zamknie kanał wyników.
	// - Główna gorutyna (main) poczeka tylko na kolektora.

	// Gorutyna zarządzająca zamknięciem kanału wyników
	go func() {
		// Czekamy, aż generator *i wszyscy workerzy* zakończą pracę.
		// Musimy śledzić ich osobno lub użyć licznika.
		// Prostsze podejście: rozdzielmy WaitGroupy lub poczekajmy na generatora, a potem na workerów.

		// Podejście z jedną WG - wymaga ostrożności:
		// Stworzymy tymczasową wg tylko dla generatora i workerów
		var wgProducers sync.WaitGroup
		wgProducers.Add(1 + numWorkers) // 1 generator + numWorkers workerów

		// Uruchom generatora (ponownie, ale teraz z wgProducers)
		// UWAGA: To jest koncepcyjny błąd w tym przepływie. Powinniśmy użyć jednej WG
		// i zamknąć resultsChan *po* tym jak wg.Wait() dla workerów się zakończy.
		// Poprawmy logikę w main:

		// --- Poprawiona Logika main ---
		// 1. Generator startuje (i zamknie ordersChan po skończeniu)
		// 2. Workerzy startują (i zakończą się po przetworzeniu wszystkiego z ordersChan)
		// 3. Kolektor startuje (i zakończy się po przetworzeniu wszystkiego z resultsChan)

		// Potrzebujemy sposobu, by wiedzieć, kiedy *wszyscy* workerzy skończyli,
		// aby móc bezpiecznie zamknąć `resultsChan`.

		// Użyjemy osobnej WaitGroup dla workerów
		var wgWorkers sync.WaitGroup
		wgWorkers.Add(numWorkers)

		// Zmodyfikowane uruchomienie workerów
		fmt.Printf("Uruchamianie %d workerów...\n", numWorkers)
		for i := 1; i <= numWorkers; i++ {
			// Przekazujemy wgWorkers do workerów
			go func(workerID int) {
				defer wgWorkers.Done() // Worker sygnalizuje zakończenie
				processOrderWorker(workerID, ordersChan, resultsChan, nil) // Przekazujemy nil, bo główna wg nie śledzi workerów bezpośrednio
			}(i)
		}

		// Uruchomienie generatora
		var wgGenerator sync.WaitGroup
		wgGenerator.Add(1)
		fmt.Println("Uruchamianie generatora zamówień...")
		go func() {
			defer wgGenerator.Done()
			orderGenerator(ordersChan, nil) // Przekazujemy nil, bo generator ma swoją WG
		}()


		// Uruchomienie kolektora
		var wgCollector sync.WaitGroup
		wgCollector.Add(1)
		fmt.Println("Uruchamianie kolektora wyników...")
		go func() {
			defer wgCollector.Done()
			resultsCollector(resultsChan, nil) // Przekazujemy nil
		}()


		// Teraz główna gorutyna koordynuje zamknięcie kanałów
		// 1. Poczekaj na zakończenie generatora
		wgGenerator.Wait()
		fmt.Println("Main: Generator zakończył pracę.")
		// ordersChan zostanie zamknięty przez generatora

		// 2. Poczekaj na zakończenie *wszystkich* workerów
		wgWorkers.Wait()
		fmt.Println("Main: Wszyscy workerzy zakończyli pracę.")

		// 3. *Po* zakończeniu workerów, zamknij kanał wyników
		// To sygnał dla kolektora, że nie będzie więcej wyników
		close(resultsChan)
		fmt.Println("Main: Kanał wyników zamknięty.")

		// 4. Poczekaj na zakończenie kolektora
		wgCollector.Wait()
		fmt.Println("Main: Kolektor zakończył pracę.")

	}() // Koniec gorutyny zarządzającej

	// Główna gorutyna nie musi już czekać tutaj na wg,
	// ponieważ gorutyna zarządzająca koordynuje całe zakończenie.
	// Ale żeby program główny nie zakończył się przedwcześnie,
	// musimy poczekać na zakończenie tej gorutyny koordynującej.
	// Zamiast tego, możemy po prostu poczekać na kolektora w głównej funkcji.

	// --- Ostateczna, uproszczona logika main ---
	fmt.Println("--- Symulacja Równoległego Przetwarzania Zamówień (Finalna Wersja) ---")

	ordersChanFinal := make(chan Order, numOrders)
	resultsChanFinal := make(chan ProcessResult, numOrders)
	var wgFinal sync.WaitGroup

	// Uruchomienie kolektora (musi być pierwszy, by nasłuchiwać)
	wgFinal.Add(1)
	go func() {
		defer wgFinal.Done()
		resultsCollector(resultsChanFinal, nil) // nil, bo używamy wgFinal w main
	}()

	// Uruchomienie workerów
	wgFinal.Add(numWorkers) // Dodajemy workerów do głównej WG
	for i := 1; i <= numWorkers; i++ {
		go func(workerID int) {
			defer wgFinal.Done() // Worker sygnalizuje zakończenie w głównej WG
			processOrderWorker(workerID, ordersChanFinal, resultsChanFinal, nil) // nil, bo używamy wgFinal w main
		}(i)
	}

	// Uruchomienie generatora
	wgFinal.Add(1) // Dodajemy generatora do głównej WG
	go func() {
		defer wgFinal.Done() // Generator sygnalizuje zakończenie w głównej WG
		orderGenerator(ordersChanFinal, nil) // nil, bo używamy wgFinal w main
	}()

	// Musimy poczekać, aż *generator i workerzy* skończą, zanim zamkniemy resultsChan.
	// Najprościej jest uruchomić osobną gorutynę, która poczeka na zakończenie
	// generatora i workerów (czyli na `numWorkers + 1` sygnałów Done()), a potem zamknie resultsChan.

	go func() {
		// Stworzymy nową WG tylko do śledzenia generatora i workerów
		var wgProducersAndWorkers sync.WaitGroup
		wgProducersAndWorkers.Add(1 + numWorkers) // 1 generator + numWorkers

		// Re-implementacja uruchomienia z tą nową WG
		ordersChanTemp := make(chan Order, numOrders)
		resultsChanTemp := make(chan ProcessResult, numOrders)

		// Kolektor (bez zmian, nasłuchuje na resultsChanTemp)
		var wgCollectorOnly sync.WaitGroup
		wgCollectorOnly.Add(1)
		go func(){
			defer wgCollectorOnly.Done()
			resultsCollector(resultsChanTemp, nil)
		}()

		// Workerzy (używają wgProducersAndWorkers)
		for i := 1; i <= numWorkers; i++ {
			go func(workerID int) {
				defer wgProducersAndWorkers.Done() // Sygnalizuj zakończenie w tej specyficznej WG
				processOrderWorker(workerID, ordersChanTemp, resultsChanTemp, nil)
			}(i)
		}

		// Generator (używa wgProducersAndWorkers)
		go func() {
			defer wgProducersAndWorkers.Done() // Sygnalizuj zakończenie w tej specyficznej WG
			orderGenerator(ordersChanTemp, nil)
		}()

		// Ta gorutyna czeka, aż generator i workerzy skończą
		wgProducersAndWorkers.Wait()
		// Dopiero teraz bezpiecznie zamykamy kanał wyników
		close(resultsChanTemp)
		fmt.Println("Koordynator: Generator i workerzy zakończyli. Zamykanie kanału wyników.")

		// Teraz czekamy tylko na zakończenie kolektora
		wgCollectorOnly.Wait()
		fmt.Println("Koordynator: Kolektor zakończył. Koniec symulacji.")

	}() // Koniec gorutyny koordynującej


	// Główna gorutyna musi poczekać na zakończenie symulacji.
	// W tym przypadku, ponieważ cała logika koordynacji jest w powyższej gorutynie,
	// główna funkcja mogłaby się zakończyć. Aby temu zapobiec,
	// można użyć np. kanału do sygnalizacji zakończenia całej symulacji.

	// Prostsze rozwiązanie (wracając do jednej WG, ale z poprawną koordynacją zamknięcia):

	fmt.Println("\n--- Symulacja Równoległego Przetwarzania Zamówień (Poprawiona Koordynacja) ---")

	ordersChanCoord := make(chan Order, numOrders)
	resultsChanCoord := make(chan ProcessResult, numOrders)
	var wgCoord sync.WaitGroup

	// 1. Uruchom workerów (dodajemy do wgCoord)
	fmt.Printf("Uruchamianie %d workerów...\n", numWorkers)
	wgCoord.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go func(workerID int) {
			defer wgCoord.Done() // Worker sygnalizuje zakończenie
			processOrderWorker(workerID, ordersChanCoord, resultsChanCoord, nil)
		}(i)
	}

	// 2. Uruchom generatora (w osobnej gorutynie, ale *nie* dodajemy go do wgCoord)
	//    Generator zamknie ordersChanCoord, gdy skończy.
	fmt.Println("Uruchamianie generatora zamówień...")
	go orderGenerator(ordersChanCoord, nil) // Sam zarządza swoim cyklem życia i zamknięciem kanału

	// 3. Uruchom gorutynę, która poczeka na *zakończenie wszystkich workerów*
	//    a następnie zamknie kanał wyników.
	go func() {
		wgCoord.Wait() // Czekaj aż licznik wgCoord (czyli wszyscy workerzy) dojdzie do zera
		close(resultsChanCoord) // Bezpiecznie zamknij kanał wyników
		fmt.Println("Koordynator zamknięcia: Wszyscy workerzy zakończyli. Kanał wyników zamknięty.")
	}()

	// 4. Uruchom kolektora (w głównej gorutynie lub osobnej).
	//    Kolektor będzie działał dopóki resultsChanCoord nie zostanie zamknięty.
	//    Jeśli uruchomimy go w głównej gorutynie, program poczeka na jego zakończenie.
	fmt.Println("Uruchamianie kolektora wyników (w głównej gorutynie)...")
	resultsCollector(resultsChanCoord, nil) // Wywołanie blokujące, czeka na zamknięcie resultsChanCoord

	fmt.Println("\n--- Koniec Symulacji ---")
}
