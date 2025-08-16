<?php
namespace app\services;

use yii\httpclient\Client;
use yii\base\Component;
use Yii;

class ApiService extends Component
{
    /**
     * @var string Базовый URL API
     */
    public $baseUrl;

    /**
     * @var int Таймаут запроса в секундах
     */
    public $timeout = 30;

    /**
     * @var Client HTTP клиент
     */
    protected $client;

    public function init()
    {
        parent::init();
        $this->client = new Client([
            'baseUrl' => $this->baseUrl,
            'requestConfig' => [
                'format' => Client::FORMAT_JSON,
                'options' => [
                    'timeout' => $this->timeout
                ]
            ],
            'responseConfig' => [
                'format' => Client::FORMAT_JSON
            ]
        ]);
    }

    /**
     * Отправка GET запроса
     * @param string $url
     * @param array $params
     * @return array
     * @throws \yii\httpclient\Exception
     */
    protected function get($url, $params = [])
    {
        $response = $this->client->get($url, $params)->send();

        if (!$response->isOk) {
            throw new \Exception('API request failed: ' . $response->statusCode);
        }

        return $response->data;
    }

    /**
     * Отправка POST запроса
     * @param string $url
     * @param array $data
     * @return array
     * @throws \yii\httpclient\Exception
     */
    protected function post($url, $data = [])
    {
        $response = $this->client->post($url, $data)->send();

        if (!$response->isOk) {
            throw new \Exception('API request failed: ' . $response->statusCode);
        }

        return $response->data;
    }

    /**
     * Отправка PUT запроса
     * @param string $url
     * @param array $data
     * @return array
     * @throws \yii\httpclient\Exception
     */
    protected function put($url, $data = [])
    {
        $response = $this->client->put($url, $data)->send();

        if (!$response->isOk) {
            throw new \Exception('API request failed: ' . $response->statusCode);
        }

        return $response->data;
    }

    /**
     * Отправка DELETE запроса
     * @param string $url
     * @return array
     * @throws \yii\httpclient\Exception
     */
    protected function delete($url)
    {
        $response = $this->client->delete($url)->send();

        if (!$response->isOk) {
            throw new \Exception('API request failed: ' . $response->statusCode);
        }

        return $response->data;
    }
}