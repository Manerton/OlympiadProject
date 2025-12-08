package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"sync"
	"time"
)

type VerifyApplicationsStrategy struct {
	timeout time.Duration
}

func NewVerifyApplicationsStrategy(timeout time.Duration) *VerifyApplicationsStrategy {
	return &VerifyApplicationsStrategy{timeout: timeout}
}

func (s *VerifyApplicationsStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {
	fmt.Printf("=== НАЧАЛО ОБРАБОТКИ ЗАПРОСА VerifyApplications ===\n")
	fmt.Printf("Время начала: %s\n", time.Now().Format(time.RFC3339))

	// ДОБАВЛЕНО: отладочный вывод конфигурации targets
	fmt.Println("\n=== КОНФИГУРАЦИЯ TARGETS ===")
	for i, target := range targets {
		fmt.Printf("Target[%d]: URL='%s', Fields=%v\n", i, target.URL, target.Fields)
	}

	if len(targets) < 4 {
		fmt.Printf("ОШИБКА: ожидается 4 target, получено %d\n", len(targets))
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("expected 4 targets, got %d", len(targets)),
		}, nil
	}

	client := &http.Client{Timeout: s.timeout}

	// ========== 1. Читаем тело оригинального запроса ==========
	fmt.Println("\n=== ЭТАП 1: Чтение входящего запроса ===")

	// Сохраняем тело запроса для отладки
	bodyBytes, _ := io.ReadAll(origReq.Body)
	fmt.Printf("Тело запроса (raw): %s\n", string(bodyBytes))

	// Восстанавливаем тело для декодирования
	origReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var incomingBody struct {
		ID   string `json:"id"`
		Role int    `json:"role"`
	}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&incomingBody); err != nil {
		fmt.Printf("ОШИБКА декодирования тела запроса: %v\n", err)
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "invalid request body",
		}, err
	}

	fmt.Printf("Парсинг успешен: ID='%s', Role=%d\n", incomingBody.ID, incomingBody.Role)

	// ========== 2. Получаем список школ ==========
	fmt.Println("\n=== ЭТАП 2: Получение списка школ ===")
	var schoolIDs []string

	if incomingBody.Role == 0 {
		fmt.Println("ROLE = 0: запрашиваем школы района")
		// ПРОВЕРЯЕМ что targets[0] существует
		if len(targets) > 0 {
			fmt.Printf("Target[0].URL: '%s'\n", targets[0].URL)
			if targets[0].URL == "" {
				fmt.Println("ОШИБКА: targets[0].URL пустой!")
				return &responcetypes.ApiResponse{
					Status:     "error",
					StatusCode: http.StatusInternalServerError,
					Error:      "district service URL is empty",
				}, nil
			}

			// Исправляем URL: добавляем ID как path parameter
			districtURL := targets[0].URL
			// Убедимся, что URL заканчивается на /, иначе добавим
			if len(districtURL) > 0 && districtURL[len(districtURL)-1] != '/' {
				districtURL += "/"
			}
			districtURL += incomingBody.ID

			fmt.Printf("Формируем URL запроса к district service: %s\n", districtURL)

			req1, err := http.NewRequest("GET", districtURL, nil)
			if err != nil {
				fmt.Printf("ОШИБКА создания запроса: %v\n", err)
				return nil, err
			}
			req1.Header = origReq.Header.Clone()

			// ДОБАВЛЕНО: отладочный вывод заголовков
			fmt.Printf("Заголовки запроса: %v\n", req1.Header)

			fmt.Println("Отправка запроса к district service...")
			resp1, err := client.Do(req1)
			if err != nil {
				fmt.Printf("ОШИБКА вызова district service: %v\n", err)
				return &responcetypes.ApiResponse{
					Status:     "error",
					StatusCode: http.StatusBadGateway,
					Error:      fmt.Sprintf("failed to call district service: %v", err),
				}, err
			}
			defer resp1.Body.Close()

			fmt.Printf("Ответ от district service: статус %d\n", resp1.StatusCode)

			if resp1.StatusCode >= 400 {
				body, _ := io.ReadAll(resp1.Body)
				fmt.Printf("ОШИБКА от district service: %s\n", string(body))
				return &responcetypes.ApiResponse{
					Status:     "error",
					StatusCode: resp1.StatusCode,
					Error:      "district service returned error",
				}, nil
			}

			var districtResp struct {
				Data struct {
					Schools []struct {
						ID string `json:"id"`
					} `json:"schools"`
				} `json:"data"`
				Status string `json:"status"`
			}

			respBody1, _ := io.ReadAll(resp1.Body)
			fmt.Printf("Тело ответа от district service: %s\n", string(respBody1))

			if err := json.Unmarshal(respBody1, &districtResp); err != nil {
				fmt.Printf("ОШИБКА парсинга ответа от district service: %v\n", err)
				// Попробуем альтернативную структуру
				var altDistrictResp struct {
					Schools []struct {
						ID string `json:"id"`
					} `json:"data"`
				}
				if err := json.Unmarshal(respBody1, &altDistrictResp); err != nil {
					fmt.Printf("ОШИБКА парсинга альтернативной структуры: %v\n", err)
					return &responcetypes.ApiResponse{
						Status:     "error",
						StatusCode: http.StatusInternalServerError,
						Error:      "failed to parse district response",
					}, err
				}

				fmt.Printf("Получено школ от района (альтернативная структура): %d\n", len(altDistrictResp.Schools))
				for i, school := range altDistrictResp.Schools {
					schoolIDs = append(schoolIDs, school.ID)
					fmt.Printf("  Школа %d: ID=%s\n", i+1, school.ID)
				}
			} else {
				fmt.Printf("Получено школ от района (стандартная структура): %d\n", len(districtResp.Data.Schools))
				for i, school := range districtResp.Data.Schools {
					schoolIDs = append(schoolIDs, school.ID)
					fmt.Printf("  Школа %d: ID=%s\n", i+1, school.ID)
				}
			}
		} else {
			fmt.Println("ОШИБКА: targets[0] не существует")
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Error:      "district target not configured",
			}, nil
		}
	} else {
		fmt.Printf("ROLE != 0: используем ID как школу: %s\n", incomingBody.ID)
		schoolIDs = []string{incomingBody.ID}
	}

	fmt.Printf("Итого schoolIDs: %v\n", schoolIDs)

	if len(schoolIDs) == 0 {
		fmt.Println("Нет школ для обработки, возвращаем пустой массив")
		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Data:       json.RawMessage("[]"),
		}, nil
	}

	// ========== 3. Получаем заявки для школ ==========
	fmt.Println("\n=== ЭТАП 3: Получение заявок для школ ===")

	// ПРОВЕРЯЕМ что targets[1] существует
	if len(targets) > 1 {
		fmt.Printf("Target[1].URL: '%s'\n", targets[1].URL)

		if targets[1].URL == "" {
			fmt.Println("ОШИБКА: targets[1].URL пустой!")
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Error:      "application service URL is empty",
			}, nil
		}

		appPayload := map[string]any{
			"ids": schoolIDs,
		}
		payloadBytes, _ := json.Marshal(appPayload)
		fmt.Printf("Тело запроса к application service: %s\n", string(payloadBytes))

		req2, err := http.NewRequest("POST", targets[1].URL, bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Printf("ОШИБКА создания запроса: %v\n", err)
			return nil, err
		}
		req2.Header = origReq.Header.Clone()
		req2.Header.Set("Content-Type", "application/json")

		fmt.Printf("Заголовки запроса: %v\n", req2.Header)

		fmt.Println("Отправка запроса к application service...")
		resp2, err := client.Do(req2)
		if err != nil {
			fmt.Printf("ОШИБКА вызова application service: %v\n", err)
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusBadGateway,
				Error:      fmt.Sprintf("failed to call application-service: %v", err),
			}, err
		}
		defer resp2.Body.Close()

		fmt.Printf("Ответ от application service: статус %d\n", resp2.StatusCode)

		if resp2.StatusCode >= 400 {
			body, _ := io.ReadAll(resp2.Body)
			fmt.Printf("ОШИБКА от application service: %s\n", string(body))
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: resp2.StatusCode,
				Error:      "application service returned error",
			}, nil
		}

		// ОБНОВЛЕНО: структура для парсинга ответа от application service
		var appResponse struct {
			Data []struct {
				ID                 string `json:"id"`
				UserID             string `json:"userId"`
				SchoolID           string `json:"schoolId"`
				EventID            string `json:"eventId"`
				Profile            string `json:"profile"`
				ClassParticipation int    `json:"class_participation"`
				Status             int    `json:"status"`
				Reason             int    `json:"reason"`
				Code               string `json:"code"`
				SubmittedAt        string `json:"submittedAt"`
				UpdatedAt          string `json:"updatedAt"`
			} `json:"data"`
			Status     string `json:"status"`
			StatusCode int    `json:"status_code"`
		}

		respBody2, _ := io.ReadAll(resp2.Body)
		fmt.Printf("Тело ответа от application service (сырое): %s\n", string(respBody2))

		if err := json.Unmarshal(respBody2, &appResponse); err != nil {
			fmt.Printf("ОШИБКА парсинга ответа от application service: %v\n", err)

			// Попробуем альтернативную структуру (просто массив)
			var applications []struct {
				ID                 string `json:"id"`
				UserID             string `json:"userId"`
				SchoolID           string `json:"schoolId"`
				EventID            string `json:"eventId"`
				Profile            string `json:"profile"`
				ClassParticipation int    `json:"class_participation"`
				Status             int    `json:"status"`
				Reason             int    `json:"reason"`
				Code               string `json:"code"`
				SubmittedAt        string `json:"submittedAt"`
				UpdatedAt          string `json:"updatedAt"`
			}

			if err := json.Unmarshal(respBody2, &applications); err != nil {
				fmt.Printf("ОШИБКА парсинга как чистого массива: %v\n", err)
				return &responcetypes.ApiResponse{
					Status:     "error",
					StatusCode: http.StatusInternalServerError,
					Error:      "failed to parse applications response",
				}, err
			}

			appResponse.Data = applications
			fmt.Printf("Успешно распарсено как чистый массив, заявок: %d\n", len(appResponse.Data))
		} else {
			fmt.Printf("Успешно распарсено как объект с data полем, заявок: %d\n", len(appResponse.Data))
		}

		applications := appResponse.Data

		fmt.Printf("Получено заявок: %d\n", len(applications))
		for i, app := range applications {
			fmt.Printf("  Заявка %d: ID=%s, UserID=%s, EventID=%s, Status=%d, Profile=%s\n",
				i+1, app.ID, app.UserID, app.EventID, app.Status, app.Profile)
		}

		if len(applications) == 0 {
			fmt.Println("Нет заявок для обработки, возвращаем пустой массив")
			return &responcetypes.ApiResponse{
				Status:     "success",
				StatusCode: http.StatusOK,
				Data:       json.RawMessage("[]"),
			}, nil
		}

		// ========== 4. Собираем уникальные UserID и EventID ==========
		fmt.Println("\n=== ЭТАП 4: Сбор уникальных ID ===")
		userIDSet := make(map[string]bool)
		eventIDSet := make(map[string]bool)

		for _, app := range applications {
			userIDSet[app.UserID] = true
			eventIDSet[app.EventID] = true
		}

		userIDs := make([]string, 0, len(userIDSet))
		eventIDs := make([]string, 0, len(eventIDSet))

		for id := range userIDSet {
			userIDs = append(userIDs, id)
		}

		for id := range eventIDSet {
			eventIDs = append(eventIDs, id)
		}

		fmt.Printf("Уникальных UserID: %d -> %v\n", len(userIDs), userIDs)
		fmt.Printf("Уникальных EventID: %d -> %v\n", len(eventIDs), eventIDs)

		// ========== 5. Параллельно запрашиваем пользователей и события ==========
		fmt.Println("\n=== ЭТАП 5: Параллельные запросы пользователей и событий ===")

		// Проверяем остальные targets
		if len(targets) < 4 {
			fmt.Printf("ОШИБКА: недостаточно targets для запросов пользователей и событий\n")
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Error:      "not enough targets configured",
			}, nil
		}

		// ОБНОВЛЕНО: структура для парсинга ответа от user service
		var userResponse struct {
			Data []struct {
				ID            string `json:"id"`
				Email         string `json:"email"`
				FirstName     string `json:"firstname"`
				Surname       string `json:"surname"`
				Patronymic    string `json:"patronymic"`
				PhoneNumber   string `json:"phone_number"`
				BirthDate     string `json:"birthdate"`
				Gender        int    `json:"gender"`
				Role          int    `json:"role"`
				Activated     bool   `json:"activated"`
				ParticipantID string `json:"participant_id"`
				Disability    int    `json:"disability"`
				SchoolID      string `json:"school_id"`
				Citizenship   int    `json:"citizenship"`
				ClassNumber   int    `json:"class_number"`
			} `json:"data"`
			Status string `json:"status"`
		}

		// ОБНОВЛЕНО: структура для парсинга ответа от event service
		var eventResponse struct {
			Data []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
			Status string `json:"status"`
		}

		var wg sync.WaitGroup
		var userErr, eventErr error

		// Запрос пользователей
		fmt.Printf("URL user service: %s\n", targets[2].URL)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if targets[2].URL == "" {
				userErr = fmt.Errorf("user service URL is empty")
				fmt.Println("ОШИБКА: targets[2].URL пустой!")
				return
			}

			userPayload := map[string]any{"ids": userIDs}
			payloadBytes, _ := json.Marshal(userPayload)
			fmt.Printf("Тело запроса к user service: %s\n", string(payloadBytes))

			req3, err := http.NewRequest("POST", targets[2].URL, bytes.NewReader(payloadBytes))
			if err != nil {
				userErr = err
				fmt.Printf("ОШИБКА создания запроса к user service: %v\n", err)
				return
			}
			req3.Header = origReq.Header.Clone()
			req3.Header.Set("Content-Type", "application/json")

			fmt.Println("Отправка запроса к user service...")
			resp3, err := client.Do(req3)
			if err != nil {
				userErr = err
				fmt.Printf("ОШИБКА вызова user service: %v\n", err)
				return
			}
			defer resp3.Body.Close()

			fmt.Printf("Ответ от user service: статус %d\n", resp3.StatusCode)

			if resp3.StatusCode >= 400 {
				body, _ := io.ReadAll(resp3.Body)
				fmt.Printf("ОШИБКА от user service: %s\n", string(body))
				userErr = fmt.Errorf("user service returned status %d", resp3.StatusCode)
				return
			}

			respBody3, _ := io.ReadAll(resp3.Body)
			fmt.Printf("Тело ответа от user service (сырое): %s\n", string(respBody3))

			if err := json.Unmarshal(respBody3, &userResponse); err != nil {
				userErr = err
				fmt.Printf("ОШИБКА парсинга ответа от user service как объекта с data: %v\n", err)

				// Попробуем как чистый массив
				var users []struct {
					ID            string `json:"id"`
					Email         string `json:"email"`
					FirstName     string `json:"firstname"`
					Surname       string `json:"surname"`
					Patronymic    string `json:"patronymic"`
					PhoneNumber   string `json:"phone_number"`
					BirthDate     string `json:"birthdate"`
					Gender        int    `json:"gender"`
					Role          int    `json:"role"`
					Activated     bool   `json:"activated"`
					ParticipantID string `json:"participant_id"`
					Disability    int    `json:"disability"`
					SchoolID      string `json:"school_id"`
					Citizenship   int    `json:"citizenship"`
					ClassNumber   int    `json:"class_number"`
				}

				if err := json.Unmarshal(respBody3, &users); err != nil {
					userErr = fmt.Errorf("failed to parse user response as array: %v", err)
					return
				}

				userResponse.Data = users
				fmt.Printf("Успешно распарсено как чистый массив пользователей, количество: %d\n", len(userResponse.Data))
			} else {
				fmt.Printf("Успешно распарсено как объект с data, пользователей: %d\n", len(userResponse.Data))
			}
		}()

		// Запрос событий
		fmt.Printf("URL event service: %s\n", targets[3].URL)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if targets[3].URL == "" {
				eventErr = fmt.Errorf("event service URL is empty")
				fmt.Println("ОШИБКА: targets[3].URL пустой!")
				return
			}

			eventPayload := map[string]any{"ids": eventIDs}
			payloadBytes, _ := json.Marshal(eventPayload)
			fmt.Printf("Тело запроса к event service: %s\n", string(payloadBytes))

			req4, err := http.NewRequest("POST", targets[3].URL, bytes.NewReader(payloadBytes))
			if err != nil {
				eventErr = err
				fmt.Printf("ОШИБКА создания запроса к event service: %v\n", err)
				return
			}
			req4.Header = origReq.Header.Clone()
			req4.Header.Set("Content-Type", "application/json")

			fmt.Println("Отправка запроса к event service...")
			resp4, err := client.Do(req4)
			if err != nil {
				eventErr = err
				fmt.Printf("ОШИБКА вызова event service: %v\n", err)
				return
			}
			defer resp4.Body.Close()

			fmt.Printf("Ответ от event service: статус %d\n", resp4.StatusCode)

			if resp4.StatusCode >= 400 {
				body, _ := io.ReadAll(resp4.Body)
				fmt.Printf("ОШИБКА от event service: %s\n", string(body))
				eventErr = fmt.Errorf("event service returned status %d", resp4.StatusCode)
				return
			}

			respBody4, _ := io.ReadAll(resp4.Body)
			fmt.Printf("Тело ответа от event service (сырое): %s\n", string(respBody4))

			if err := json.Unmarshal(respBody4, &eventResponse); err != nil {
				eventErr = err
				fmt.Printf("ОШИБКА парсинга ответа от event service как объекта с data: %v\n", err)

				// Попробуем как чистый массив
				var events []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}

				if err := json.Unmarshal(respBody4, &events); err != nil {
					eventErr = fmt.Errorf("failed to parse event response as array: %v", err)
					return
				}

				eventResponse.Data = events
				fmt.Printf("Успешно распарсено как чистый массив событий, количество: %d\n", len(eventResponse.Data))
			} else {
				fmt.Printf("Успешно распарсено как объект с data, событий: %d\n", len(eventResponse.Data))
			}
		}()

		fmt.Println("Ожидание завершения параллельных запросов...")
		wg.Wait()
		fmt.Println("Параллельные запросы завершены")

		if userErr != nil {
			fmt.Printf("ОШИБКА при получении пользователей: %v\n", userErr)
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusBadGateway,
				Error:      fmt.Sprintf("failed to get users: %v", userErr),
			}, nil
		}

		if eventErr != nil {
			fmt.Printf("ОШИБКА при получении событий: %v\n", eventErr)
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusBadGateway,
				Error:      fmt.Sprintf("failed to get events: %v", eventErr),
			}, nil
		}

		// ========== 6. Создаем мапы для быстрого поиска ==========
		fmt.Println("\n=== ЭТАП 6: Создание мап для поиска ===")

		userMap := make(map[string]struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			FirstName     string `json:"firstname"`
			Surname       string `json:"surname"`
			Patronymic    string `json:"patronymic"`
			PhoneNumber   string `json:"phone_number"`
			BirthDate     string `json:"birthdate"`
			Gender        int    `json:"gender"`
			Role          int    `json:"role"`
			Activated     bool   `json:"activated"`
			ParticipantID string `json:"participant_id"`
			Disability    int    `json:"disability"`
			SchoolID      string `json:"school_id"`
			Citizenship   int    `json:"citizenship"`
			ClassNumber   int    `json:"class_number"`
		})

		for _, user := range userResponse.Data {
			userMap[user.ID] = user
			fmt.Printf("  Пользователь добавлен в мапу: ID=%s, ФИО=%s %s %s\n",
				user.ID, user.Surname, user.FirstName, user.Patronymic)
		}

		fmt.Printf("Размер userMap: %d\n", len(userMap))

		eventMap := make(map[string]string)
		for _, event := range eventResponse.Data {
			eventMap[event.ID] = event.Name
			fmt.Printf("  Событие добавлено в мапу: ID=%s, Name=%s\n", event.ID, event.Name)
		}

		fmt.Printf("Размер eventMap: %d\n", len(eventMap))

		// ========== 7. Формируем итоговый ответ ==========
		fmt.Println("\n=== ЭТАП 7: Формирование итогового ответа ===")

		type ResultItem struct {
			ID           string `json:"id"`
			OlympiadName string `json:"olympiadName"`
			Profile      string `json:"profile"`
			Category     int    `json:"category"`
			Status       int    `json:"status"`
			Surname      string `json:"surname"`
			FirstName    string `json:"firstName"`
			Patronymic   string `json:"patronymic"`
			Birthdate    string `json:"birthdate"`
			Gender       int    `json:"gender"`
			ClassNumber  int    `json:"classNumber"`
			Citizenship  int    `json:"citizenship"`
			Disability   int    `json:"disability"`
		}

		results := make([]ResultItem, 0, len(applications))
		skippedCount := 0

		for i, app := range applications {
			user, userExists := userMap[app.UserID]
			eventName, eventExists := eventMap[app.EventID]

			if !userExists || !eventExists {
				skippedCount++
				fmt.Printf("  Пропуск заявки %d (ID=%s): userExists=%v, eventExists=%v\n",
					i+1, app.ID, userExists, eventExists)
				continue
			}

			fmt.Printf("  Обработка заявки %d (ID=%s): User=%s %s, Event=%s\n",
				i+1, app.ID, user.Surname, user.FirstName, eventName)

			results = append(results, ResultItem{
				ID:           app.ID,
				OlympiadName: eventName,
				Profile:      app.Profile,
				Category:     app.ClassParticipation,
				Status:       app.Status,
				Surname:      user.Surname,
				FirstName:    user.FirstName,
				Patronymic:   user.Patronymic,
				Birthdate:    user.BirthDate,
				Gender:       user.Gender,
				ClassNumber:  user.ClassNumber,
				Citizenship:  user.Citizenship,
				Disability:   user.Disability,
			})
		}

		fmt.Printf("Обработано заявок: %d, пропущено: %d, итого в ответе: %d\n",
			len(applications), skippedCount, len(results))

		responseData, err := json.Marshal(results)
		if err != nil {
			fmt.Printf("ОШИБКА маршалинга ответа: %v\n", err)
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Error:      "failed to marshal response",
			}, err
		}

		fmt.Printf("Итоговый ответ (сырой JSON): %s\n", string(responseData))

		// Проверяем, что JSON валиден
		var checkResult []map[string]interface{}
		if err := json.Unmarshal(responseData, &checkResult); err != nil {
			fmt.Printf("ОШИБКА: сгенерированный JSON невалиден: %v\n", err)
		} else {
			fmt.Printf("JSON валиден, элементов: %d\n", len(checkResult))
		}

		fmt.Println("\n=== КОНЕЦ ОБРАБОТКИ ЗАПРОСА VerifyApplications ===")
		fmt.Printf("Время окончания: %s\n", time.Now().Format(time.RFC3339))

		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Data:       json.RawMessage(responseData),
		}, nil
	} else {
		fmt.Println("ОШИБКА: targets[1] не существует")
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      "application target not configured",
		}, nil
	}
}
