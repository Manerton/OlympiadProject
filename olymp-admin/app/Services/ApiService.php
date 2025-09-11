<?php

namespace App\Services;

use GuzzleHttp\Client;
use GuzzleHttp\Cookie\CookieJar;
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
                'query' => $params,
                'headers' => array_merge($this->defaultHeaders, $headers),
                'http_errors' => false
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
    public function post(string $url, array $data = [], array $headers = [], bool $withCredentials = false): array
    {
        $client = new Client();

        try {
            $options = [
                'json' => $data,
                'headers' => array_merge($this->defaultHeaders, $headers),
                'http_errors' => false
            ];

            if ($withCredentials) {
                $cookies = [];
                foreach ($_COOKIE as $name => $value) {
                    $cookies[$name] = $value;
                }

                $options['cookies'] = CookieJar::fromArray($cookies, parse_url($url, PHP_URL_HOST));
            }
            $response = $client->request('POST', $url, $options);

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
