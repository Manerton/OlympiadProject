<?php

namespace App\Services;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;
use Psr\Http\Message\ResponseInterface;
class ApiService
{
    public $defaultHeaders = [
        'Content-Type' => 'application/json',
        'Accept' => 'application/json',
    ];
    public function get(string $url, array $params = [], array $headers = []): array
    {
        $client = new Client();

        try {
            $response = $client->request('GET', $url, [
                'query' => $params,  // Параметры GET-запроса
                'headers' => array_merge($this->defaultHeaders, $headers), // Объединяем с дефолтными заголовками
                'http_errors' => false // Чтобы не выбрасывать исключения при HTTP ошибках
            ]);

            return $this->handleResponse($response);

        } catch (RequestException $e) {
            return [
                'error' => true,
                'message' => 'API request failed',
                'details' => $e->getMessage(),
                'status' => $e->getCode() ?: 500
            ];
        }
    }
    public function post(string $url, array $data = [], array $headers = []): array
    {
        $client = new Client();

        try {
            $response = $client->request('POST', $url, [
                'json' => $data,  // Автоматически кодирует в JSON и устанавливает Content-Type
                'headers' => array_merge($this->defaultHeaders, $headers),
                'http_errors' => false
            ]);

            return $this->handleResponse($response);

        } catch (RequestException $e) {
            return [
                'error' => true,
                'message' => 'API POST request failed',
                'details' => $e->getMessage(),
                'status' => $e->getCode() ?: 500
            ];
        }
    }
    protected function handleResponse(ResponseInterface $response): array
    {
        $content = $response->getBody()->getContents();
        $data = json_decode($content, true) ?? $content;

        return [
            'success' => $response->getStatusCode() >= 200 && $response->getStatusCode() < 300,
            'status' => $response->getStatusCode(),
            'data' => $data,
            'headers' => $response->getHeaders()
        ];
    }
}
